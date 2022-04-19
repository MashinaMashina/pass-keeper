package accesstype

var Types = map[string]func() Access{
	"ssh": NewSSH,
	"ftp": NewFTP,
}
