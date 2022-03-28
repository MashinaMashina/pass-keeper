package parselnk

import "strings"

func (lp *linkParser) getReplaceConfig() [][2]string {
	if lp.replaceConfig == nil {
		parts := strings.Split(lp.puttyConfig.Get("lnk.replace"), "|")

		lp.replaceConfig = make([][2]string, 0)
		for _, part := range parts {
			p := strings.Split(part, ":")

			lp.replaceConfig = append(lp.replaceConfig, [2]string{
				p[0],
				p[1],
			})
		}
	}

	return lp.replaceConfig
}

func (lp *linkParser) cleanFilename(name string) string {
	replace := lp.getReplaceConfig()

	for _, r := range replace {
		name = strings.Replace(name, r[0], r[1], -1)
	}

	return name
}
