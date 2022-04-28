package scanner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"pass-keeper/internal/accesses/accesstype"
	"testing"
)

func TestParseName(t *testing.T) {
	s := &scanner{}
	s.parseName("ssh\n" +
		"some-project.com:22\n" +
		"login\n" +
		"qwerty123\n" +
		"" +
		"ssh\n" +
		"ssh.some-project.com\n" +
		"anonim")

	assert.Equal(t, "some-project", s.name)

	s.parseName("ssh\n" +
		"some.com:22\n" +
		"login\n" +
		"qwerty123\n" +
		"" +
		"ssh\n" +
		"ssh.some-project.com\n" +
		"anonim")

	assert.Equal(t, "", s.name)
}

func TestParseTokens(t *testing.T) {
	in, inWriter, err := os.Pipe()
	if err != nil {
		return
	}
	defer inWriter.Close()

	outReader, out, err := os.Pipe()
	if err != nil {
		return
	}

	fmt.Fprint(inWriter, "\n\n\n\n\n\n\n\n")
	fmt.Fprint(inWriter, "incorrecttype\nftp\n\n\n\n\n\n\n\n\n\n")

	s := &scanner{}
	s.Stdin = in
	s.Stdout = out
	s.group = "group"

	err = s.parseText("ssh\nhost1.com\nlogin1\npasswd1$%\n\nType: ft3p\n" +
		"Login: bob\nPassword: pol1010\nHost: ftp.host1.com")
	if err != nil {
		t.Error(err)
		return
	}

	out.Close()

	bytes, err := ioutil.ReadAll(outReader)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "Token: \"ssh\"\n"+
		"Enter token type (default t): Token: \"host1.com\"\n"+
		"Enter token type (default h): Token: \"login1\"\n"+
		"Enter token type (default l): Token: \"passwd1$%\"\n"+
		"Enter token type (default p): Token: \"Type: ft3p\"\n"+
		"Enter token type (default g): Token: \"Type:\"\n"+
		"Enter token type (default c): Token: \"ft3p\"\n"+
		"Enter token type (default t): Creating new access\n"+
		"Invalid type ft3p. Enter a type instead: Invalid type ft3p. Enter a type instead: "+
		"Invalid type incorrecttype. Enter a type instead: Token: \"Login: bob\"\n"+
		"Enter token type (default g): Token: \"Login:\"\n"+
		"Enter token type (default c): Token: \"bob\"\n"+
		"Enter token type (default l): Token: \"Password: pol1010\"\n"+
		"Enter token type (default g): Token: \"Password:\"\n"+
		"Enter token type (default c): Token: \"pol1010\"\n"+
		"Enter token type (default p): Token: \"Host: ftp.host1.com\"\n"+
		"Enter token type (default g): Token: \"Host:\"\n"+
		"Enter token type (default c): Token: \"ftp.host1.com\"\n"+
		"Enter token type (default h): ", string(bytes))

	if len(s.accesses) != 2 {
		t.Error("not all parsed")
		return
	}

	access1 := accesstype.NewSSH()
	access1.SetName("host1")
	access1.SetGroup("group")
	access1.SetHost("host1.com")
	access1.SetLogin("login1")
	access1.SetPassword("passwd1$%")

	access2 := accesstype.NewFTP()
	access2.SetName("host1")
	access2.SetGroup("group")
	access2.SetHost("ftp.host1.com")
	access2.SetLogin("bob")
	access2.SetPassword("pol1010")

	equalAccess(t, s.accesses[0], access1)
	equalAccess(t, s.accesses[1], access2)
}

func equalAccess(t *testing.T, expect, real accesstype.Access) bool {
	if expect.Name() != real.Name() ||
		expect.Host() != real.Host() ||
		expect.Port() != real.Port() ||
		expect.Login() != real.Login() ||
		expect.Password() != real.Password() ||
		expect.Group() != real.Group() ||
		expect.Params().Value("session_name") != real.Params().Value("session_name") {
		t.Error(fmt.Sprintf("accesses not equal\nexpected: %+v,\nbut real %+v", expect, real))

		return false
	}

	return true
}
