package parselnk

import "strings"

func (lp *linkParser) getReplaceConfig() [][2]string {
	// putty.replace не map[string]string, потому что нам важен порядок элементов
	if lp.Config.Slice("putty.replace") == nil {
		lp.Config.Set("putty.replace", []string{
			"PuTTY — ",
			"PuTTY ",
		})
	}

	slice := lp.Config.Slice("putty.replace")
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

func (lp *linkParser) cleanFilename(name string) string {
	replace := lp.getReplaceConfig()

	for _, v := range replace {
		name = strings.Replace(name, v[0], v[1], -1)
	}

	return name
}
