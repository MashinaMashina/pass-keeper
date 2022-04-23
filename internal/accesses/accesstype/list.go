package accesstype

var Types = map[string]func() Access{
	"ssh":     NewSSH,
	"ftp":     NewFTP,
	"unknown": NewUnknown,
}

func NewFTP() Access {
	a := NewUnknown()
	a.SetType("ftp")
	a.SetPort(21)

	return a
}

func NewSSH() Access {
	a := NewUnknown()
	a.SetType("ssh")
	a.SetPort(22)

	return a
}
