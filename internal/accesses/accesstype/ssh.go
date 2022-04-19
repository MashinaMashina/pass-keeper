package accesstype

type ssh struct {
	access
}

func NewSSH() Access {
	a := new()
	a.SetType("ssh")
	a.SetPort(22)

	return &ssh{a}
}
