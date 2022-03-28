package params

import "pass-keeper/internal/accesses/storage"

type like struct {
	field   string
	pattern string
}

func NewLike(field, pattern string) storage.Param {
	return &like{
		field:   field,
		pattern: pattern,
	}
}

func (_this *like) ParamType() string {
	return "like"
}

func (_this *like) Value() []string {
	return []string{_this.field, _this.pattern}
}
