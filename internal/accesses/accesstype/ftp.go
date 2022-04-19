package accesstype

type ftp struct {
	access
}

func NewFTP() Access {
	a := new()
	a.SetType("ftp")
	a.SetPort(21)

	return &ftp{a}
}
