package params

import (
	"github.com/Masterminds/squirrel"
	"pass-keeper/internal/accesses/storage"
)

func NewEq(field, cond string) storage.Param {
	return func(builder *squirrel.SelectBuilder) {
		*builder = builder.Where(field+" = ?", cond)
	}
}
