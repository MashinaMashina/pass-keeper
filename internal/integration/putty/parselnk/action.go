package parselnk

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/pkg/filesystem"
	"path/filepath"
	"strings"
)

func (lp *linkParser) action(c *cli.Context) error {
	folder, err := filepath.Abs(c.String("path"))

	if err != nil {
		return err
	}

	files, err := filesystem.ReadDir(folder)
	if err != nil {
		return err
	}

	var access accesstype.Access
	for _, file := range files {
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".lnk") {
			continue
		}

		fmt.Println("Scan", file.Name())

		access, err = lp.sshAccessByLnk(file)

		fmt.Println(lp.puttyConfig.Get("lnk.replace"))

		if err != nil {
			fmt.Println("Error with parsing .lnk:", err)
			continue
		}

		err = lp.storage.Save(access)
		if err != nil {
			fmt.Println("Error with saving access:", err)
			continue
		}
	}

	return nil
}
