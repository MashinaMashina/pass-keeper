package parselnk

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

func (lp *linkParser) action(c *cli.Context) error {
	path, err := filepath.Abs(c.String("path"))

	if err != nil {
		return err
	}

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

		fmt.Fprintln(lp.Stdout, fmt.Sprintf("Scan \"%s\"", file.Name()))

		access, err = lp.sshAccessByLnkFile(file)

		if err != nil {
			fmt.Fprintln(lp.Stdout, "Error with parsing .lnk:", err)
			continue
		}

		if existRow, err := lp.Storage.FindExists(access); err == nil {
			access.SetID(existRow.ID())
		}

		err = lp.Storage.Save(access)
		if err != nil {
			fmt.Fprintln(lp.Stdout, "Error with saving access:", err)
			continue
		}
	}

	return nil
}
