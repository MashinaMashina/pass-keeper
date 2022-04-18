package app

import (
	"fmt"
	"os"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/pkg/clibell"
	"pass-keeper/pkg/cliselect"
)

func FindOne(dto DTO, parameters ...storage.Param) (accesstype.Access, error) {
	rows, err := dto.Storage.List(parameters...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("not found rows")
	}

	if len(rows) > 1 {
		names := make([]string, 0, len(rows))

		for _, row := range rows {
			names = append(names, row.Name())
		}

		var res string

		promt := cliselect.Select{
			Label:     "Multiple options available",
			Items:     names,
			IsVimMode: true,
			Stdout:    clibell.Instance(dto.Stdout),
			Mode:      cliselect.ModeSimple,
		}

		if dto.Config.String("main.mode") == "interactive" {
			promt.Mode = cliselect.ModeInteractive
		}

		// При указании promt.Stdin выбор перестает быть интерактивным
		if dto.Stdin != os.Stdin {
			promt.Stdin = dto.Stdin
		}

		_, res, err = promt.Run()
		if err != nil {
			return nil, err
		}

		return FindOne(dto, append(parameters, params.NewEq("name", res))...)
	}

	return rows[0], nil
}
