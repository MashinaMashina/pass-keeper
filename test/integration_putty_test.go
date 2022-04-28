package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/app"
	"path/filepath"
	"strings"
	"testing"
)

func TestPuttyParselnk(t *testing.T) {
	cases := make([]appTestCase, 0, 2)

	folder := FolderTest() + "/putty-lnk/folder"

	var err error
	folder, err = filepath.Abs(folder)

	if err != nil {
		t.Error(err)
		return
	}

	folder = strings.ReplaceAll(folder, "\\", "/")

	access1 := accesstype.NewSSH()
	access1.SetName("site1.ru")
	access1.SetHost("site-host1.ru")
	access1.SetLogin("login1")
	access1.SetPassword("qwerty123")
	access1.SetGroup(folder)

	access2 := accesstype.NewSSH()
	access2.SetName("site2.ru")
	access2.SetHost("127.0.0.1")
	access2.Params().Set("session_name", "sess name")
	access2.SetGroup(folder)

	cases = append(cases, appTestCase{
		Args:    []string{"putty", "import", folder},
		Comment: "добавляем доступы",
		Check: func(t *testing.T, dto app.DTO, output string) {
			except := fmt.Sprintf("Scanning %s\n"+
				"Scan \"PuTTY %s.lnk\"\n"+
				"Scan \"PuTTY %s.lnk\"\n"+
				"Scan \"PuTTY.lnk\"\n"+
				"Error with parsing .lnk: link has no data\n",
				folder, access1.Name(), access2.Name())

			if assert.Equal(t, except, output) {
				rows, err := dto.Storage.List()
				if err != nil {
					t.Error(err)
					return
				}

				if len(rows) < 2 {
					t.Error(fmt.Sprintf("invalid scan folder; all links not parsed. parsed: %d", len(rows)))
					return
				}

				if len(rows) > 2 {
					t.Error(fmt.Sprintf("invalid scan folder; too many results. parsed: %d", len(rows)))
					return
				}

				access1.SetID(rows[0].ID())
				access2.SetID(rows[1].ID())

				equalAccess(t, access1, rows[0])
				equalAccess(t, access2, rows[1])
			}
		},
	})

	cases = append(cases, appTestCase{
		Args:    []string{"putty", "import", folder + "/PuTTY site1.ru"},
		Comment: "обновляем доступ",
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			access1.SetLogin("login2")
			dto.Storage.Save(access1)
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			rows, err := dto.Storage.List(params.NewEq("name", access1.Name()))
			if err != nil {
				t.Error(err)
				return
			}

			if len(rows) == 0 {
				t.Error("not enough rows")
				return
			}
			if len(rows) > 1 {
				t.Error("too many rows")
				return
			}

			// putty scan должен был обновить строку и установить логин из lnk файла
			access1.SetLogin("login1")

			equalAccess(t, access1, rows[0])
		},
	})

	appTestOutput(t, cases)
}
