package accesstype

type ssh struct {
	access
}

func NewSSH() Access {
	return &ssh{access{typo: "ssh", port: 22}}
}
