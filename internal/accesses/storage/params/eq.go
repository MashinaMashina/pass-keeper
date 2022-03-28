package params

import "pass-keeper/internal/accesses/storage"

type eq struct {
	field string
	cond  string
}

func NewEq(field, cond string) storage.Param {
	return &eq{
		field: field,
		cond:  cond,
	}
}

func (_this *eq) ParamType() string {
	return "eq"
}

func (_this *eq) Value() []string {
	return []string{_this.field, _this.cond}
}
