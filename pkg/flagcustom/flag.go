package flagcustom

import "errors"

var ErrEscaping = errors.New("invalid escaping")

func ParseFlags(in string) ([]string, error) {
	res := make([]string, 0, 64)

	quoted := false
	buffer := ""
	for _, s := range in {
		if quoted && s != '"' {
			buffer += string(s)
		} else if s == ' ' {
			if buffer != "" {
				res = append(res, buffer)
				buffer = ""
			}
		} else if s == '"' {
			if quoted {
				res = append(res, buffer)
				buffer = ""
				quoted = false
			} else {
				quoted = true
			}
		} else {
			buffer += string(s)
		}
	}

	if quoted {
		return nil, ErrEscaping
	}

	if buffer != "" {
		res = append(res, buffer)
	}

	return res, nil
}
