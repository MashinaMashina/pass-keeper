package scanner

import "strings"

func (ls *linkScanner) getReplaceConfig() [][2]string {
	// putty.replace не map[string]string, потому что нам важен порядок элементов
	if ls.Config.Slice("putty.replace") == nil {
		ls.Config.Set("putty.replace", []string{
			"PuTTY — ",
			"PuTTY ",
		})
	}

	slice := ls.Config.Slice("putty.replace")
	var res [][2]string
	var parts []string
	for _, v := range slice {
		parts = strings.SplitN(v, ":", 2)

		if len(parts) != 2 {
			parts = append(parts, "")
		}

		res = append(res, [2]string{
			parts[0],
			parts[1],
		})
	}

	return res
}

func (ls *linkScanner) cleanFilename(name string) string {
	replace := ls.getReplaceConfig()

	for _, v := range replace {
		name = strings.Replace(name, v[0], v[1], -1)
	}

	return name
}
