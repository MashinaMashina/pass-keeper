package scanner

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/pkg/filesystem"
	"path/filepath"
	"strings"
)

func (ls *linkScanner) action(c *cli.Context) error {
	path, err := filepath.Abs(c.Args().First())

	if err != nil {
		return err
	}

	fmt.Fprintln(ls.Stdout, "Scanning", path)

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
		files, err = filesystem.ReadDir(path)
		if err != nil {
			return err
		}
	} else {
		files = append(files, f)
	}

	var access accesstype.Access
	for _, file := range files {
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".lnk") {
			continue
		}

		fmt.Fprintln(ls.Stdout, fmt.Sprintf("Scan \"%s\"", file.Name()))

		access, err = ls.sshAccessByLnkFile(file)

		if err != nil {
			fmt.Fprintln(ls.Stdout, "Error with parsing .lnk:", err)
			continue
		}

		if existRow, err := ls.Storage.FindExists(access); err == nil {
			access.SetID(existRow.ID())
		}

		err = ls.Storage.Save(access)
		if err != nil {
			fmt.Fprintln(ls.Stdout, "Error with saving access:", err)
			continue
		}
	}

	return nil
}
