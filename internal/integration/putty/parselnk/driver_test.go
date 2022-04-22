package parselnk

import (
	"fmt"
	"pass-keeper/internal/accesses/accesstype"
	"testing"
)

func TestParseArgs(t *testing.T) {
	name := "name"
	group := "any/group"

	access := accesstype.NewSSH()
	access.SetName(name)

	p := &linkParser{}

	access.SetLogin("root")
	access.SetHost("sub.domain.net")
	access.SetPassword("aas@sdd$%^")
	access.SetGroup("any/group")

	resAccess, err := p.accessByArguments("-ssh root@sub.domain.net -pw aas@sdd$%^", name, group)
	if err != nil {
		t.Error(err)
		return
	}

	equalAccess(t, access, resAccess)

	access.SetPort(1010)

	resAccess, err = p.accessByArguments("-ssh root@sub.domain.net -P 1010 -pw aas@sdd$%^", name, group)
	if err != nil {
		t.Error(err)
		return
	}

	equalAccess(t, access, resAccess)

	access.Params().Set("session_name", "session name")

	resAccess, err = p.accessByArguments(`-ssh root@sub.domain.net -P 1010 -pw aas@sdd$%^ --load "session name"`, name, group)
	if err != nil {
		t.Error(err)
		return
	}

	equalAccess(t, access, resAccess)
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
