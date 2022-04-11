package puttyrun

import (
	"github.com/stretchr/testify/assert"
	"pass-keeper/internal/accesses/accesstype"
	"testing"
)

func TestArgs(t *testing.T) {
	t.Log("test putty run arguments")

	access := accesstype.NewSSH()

	access.SetName("any")
	access.SetHost("host1.com")
	access.SetPort(8000)
	access.SetLogin("login")
	access.SetPassword("qwerty")

	p := &puttyRun{}

	assert.Equal(t, []string{"-ssh", "login@host1.com", "-P", "8000", "-pw", "qwerty"}, p.argsFromAccess(access))

	access.SetSession("sess name")

	assert.Equal(t, []string{"-ssh", "login@host1.com", "-P", "8000", "-pw", "qwerty", "--load", "sess name"}, p.argsFromAccess(access))
}
