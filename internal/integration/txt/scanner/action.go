package scanner

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/pkg/filesystem"
	"path/filepath"
	"regexp"
	"strings"
)

var tokenTypes = map[string]string{
	"t": "type",
	"h": "host",
	"l": "login",
	"p": "password",
	"g": "group of tokens",
	"c": "comments",
}

var typeNames = map[string]string{
	"type":     "t",
	"тип":      "t",
	"login":    "l",
	"user":     "l",
	"логин":    "l",
	"password": "p",
	"passwd":   "p",
	"pass":     "p",
	"пароль":   "p",
	"host":     "h",
	"хост":     "h",
}

func (s *scanner) action(c *cli.Context) error {
	path, err := filepath.Abs(c.Args().First())

	if err != nil {
		return err
	}

	sort := "thlpgc"
	types := ""
	for _, k := range sort {
		if val, ok := tokenTypes[string(k)]; ok {
			types += fmt.Sprintf("%s - %s, ", string(k), val)
		}
	}

	types = types[:len(types)-2]

	f, err := filesystem.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = filesystem.Stat(path + ".lnk")
		}

		if err != nil {
			return err
		}
	}

	var files []filesystem.File
	if f.IsDir() {
		s.group = filesystem.PrettyPath(path)
		files, err = filesystem.ReadDir(path)
		if err != nil {
			return err
		}
	} else {
		s.group = filesystem.PrettyPath(f.Dir())
		files = append(files, f)
	}

	fmt.Fprintln(s.Stdout, "Group name: "+s.group)
	fmt.Fprintln(s.Stdout, "Token types: "+types)

	var (
		bytes []byte
	)

	for _, file := range files {
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".txt") {
			continue
		}

		fmt.Fprintln(s.Stdout, fmt.Sprintf("Scan \"%s\"", file.Name()))

		bytes, err = ioutil.ReadFile(file.FullPath())
		if err != nil {
			return err
		}

		err = s.parseText(string(bytes))
		if err != nil {
			return err
		}
	}

	for _, access := range s.accesses {
		if existRow, err := s.Storage.FindExists(access); err == nil {
			access.SetID(existRow.ID())
		}

		if access.ID() > 0 {
			fmt.Fprintln(s.Stdout, "Update access ", access.Type(), access.Login(), access.Host())
		} else {
			fmt.Fprintln(s.Stdout, "Add access ", access.Type(), access.Login(), access.Host())
		}

		err = s.Storage.Save(access)
		if err != nil {
			fmt.Fprintln(s.Stdout, "Error with saving access:", err)
			continue
		}
	}

	return nil
}

/*
 * finding name from the most frequent word
 */
func (s *scanner) parseName(c string) {
	r := regexp.MustCompile(`[^a-zA-Z0-9а-яёА-ЯЁ\-]`)

	res := make(map[string]int16)
	for _, word := range r.Split(c, -1) {
		word = strings.ToLower(strings.TrimSpace(word))

		if word == "" ||
			word == "ftp" ||
			word == "ssh" ||
			word == "com" ||
			word == "net" ||
			word == "ru" ||
			word == "biz" {
			continue
		}

		if _, exists := typeNames[word]; exists {
			continue
		}

		if _, exists := res[word]; exists {
			res[word] += 1
		} else {
			res[word] = 1
		}
	}

	var oldCnt int16 = 1
	var name string
	for word, cnt := range res {
		if cnt > oldCnt {
			name = word
			oldCnt = cnt
		}
	}

	if name != "" {
		s.name = name
	}
}

func (s *scanner) parseText(str string) error {
	s.resetAccess()
	s.parseName(str)
	s.nextType = "t"

	lines := strings.Split(str, "\n")

	err := s.parseTokens(lines)
	if err != nil {
		return err
	}

	s.resetAccess()

	return nil
}

func (s *scanner) parseTokens(tokens []string) error {
	var typo string
	for _, token := range tokens {
		token = strings.TrimSpace(token)

		if token == "" {
			continue
		}

		def := s.nextType

		if strings.Index(token, " ") > -1 {
			def = "g"
		} else if s.comment != "" {
			if v, ok := typeNames[s.trimNameToken(s.comment)]; ok {
				def = v
			}
		}

		text := "Enter token type"
		if def != "" {
			text += " (default " + def + ")"
		}
		text += ": "

		typo = ""
		fmt.Fprintln(s.Stdout, fmt.Sprintf(`Token: "%s"`, token))

		for {
			fmt.Fprint(s.Stdout, text)
			fmt.Fscanln(s.Stdin, &typo)

			if typo == "" {
				typo = def
			}

			if typo != "" {
				break
			}
		}

		if _, exists := s.filled[typo]; exists {
			s.resetAccess()
			fmt.Fprintln(s.Stdout, "Creating new access")

			s.nextType = "t"
			if typo == "t" {
				s.nextType = "h"
			}
		}

		s.nextType = ""
		s.comment = ""
		if typo != "g" && typo != "c" {
			s.filled[typo] = struct{}{}
		}

		switch typo {
		case "t":
			for {
				if _, ok := accesstype.Types[token]; !ok {
					fmt.Fprint(s.Stdout, "Invalid type ", token, ". Enter a type instead: ")
					fmt.Fscanln(s.Stdin, &token)
				} else {
					break
				}
			}

			s.access.SetType(token)
			s.nextType = "h"
		case "h":
			s.access.SetHost(token)
			s.nextType = "l"
		case "l":
			s.access.SetLogin(token)
			s.nextType = "p"
		case "p":
			s.access.SetPassword(token)
		case "g":
			tokens := strings.Split(token, " ")

			if _, ok := typeNames[s.trimNameToken(tokens[0])]; ok {
				s.nextType = "c"
			}

			err := s.parseTokens(tokens)
			if err != nil {
				return err
			}
		case "c":
			s.comment = token
		default:
			return errors.New("incorrect type" + typo)
		}
	}

	return nil
}

func (s *scanner) trimNameToken(str string) string {
	return strings.ToLower(strings.Trim(str, ":=- "))
}

func (s *scanner) resetAccess() {
	if s.access != nil {
		if creator, ok := accesstype.Types[s.access.Type()]; ok {
			for {
				if s.name != "" {
					break
				}
				fmt.Fprintln(s.Stdout, "Enter name for", s.access.Type(), s.access.Login(), s.access.Host())
				fmt.Fscanln(s.Stdin, &s.name)
			}

			access := creator()

			access.SetGroup(s.group)
			access.SetName(s.name)

			if s.access.Host() != "" {
				access.SetHost(s.access.Host())
			}
			if s.access.Port() != 0 {
				access.SetPort(s.access.Port())
			}
			if s.access.Login() != "" {
				access.SetLogin(s.access.Login())
			}
			if s.access.Password() != "" {
				access.SetPassword(s.access.Password())
			}

			s.accesses = append(s.accesses, access)
		} else {
			fmt.Fprintln(s.Stdout, "Invalid type "+s.access.Type())
		}
	}

	s.access = accesstype.NewUnknown()
	s.filled = make(map[string]struct{})
}
