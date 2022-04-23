package runner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"pass-keeper/internal/accesses/accesstype"
	"testing"
)

func TestArgs(t *testing.T) {
	t.Log("test filezilla run arguments")

	p := &filezillaRun{}

	access := accesstype.NewSSH()
	access.SetName("any")
	access.SetHost("host1.com")
	access.SetPort(8000)
	access.SetLogin("login")
	access.SetPassword("qwerty")

	u := fmt.Sprintf("sftp://%s:%s@%s:%d", url.QueryEscape(access.Login()),
		url.QueryEscape(access.Password()), url.QueryEscape(access.Host()), access.Port())

	assert.Equal(t, []string{u}, p.argsFromAccess(access))

	access2 := accesstype.NewFTP()
	access2.SetName("any")
	access2.SetHost("1.1.1.1")
	access2.SetPort(22)
	access2.SetLogin("bob")
	access2.SetPassword("12@3$$%%^^")

	u = fmt.Sprintf("ftp://%s:%s@%s:%d", url.QueryEscape(access2.Login()),
		url.QueryEscape(access2.Password()), url.QueryEscape(access2.Host()), access2.Port())

	assert.Equal(t, []string{u}, p.argsFromAccess(access2))
}
