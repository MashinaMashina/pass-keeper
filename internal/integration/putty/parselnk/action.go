package parselnk

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses"
	"pass-keeper/pkg/filesystem"
	"path/filepath"
	"strings"
)

func (lp *linkParser) cliAction(c *cli.Context) error {
	folder, err := filepath.Abs(c.String("folder"))

	if err != nil {
		return err
	}

	files, err := filesystem.ReadDir(folder)
	if err != nil {
		return err
	}

	var access accesses.Access
	for _, file := range files {
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".lnk") {
			continue
		}

		fmt.Println("Scan", file.Name())

		access, err = sshAccessByLnk(file)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = lp.storage.Add(access)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}
