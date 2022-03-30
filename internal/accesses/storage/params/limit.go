package params

import (
	"pass-keeper/internal/accesses/storage"
	"strconv"
)

type limit struct {
	limit int
}

func NewLimit(l int) storage.Param {
	return &limit{
		limit: l,
	}
}

func (_this *limit) ParamType() string {
	return "limit"
}

func (_this *limit) Value() []string {
	return []string{strconv.Itoa(_this.limit)}
}
