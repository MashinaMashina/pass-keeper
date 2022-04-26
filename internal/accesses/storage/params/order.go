package params

import (
	"github.com/Masterminds/squirrel"
	"pass-keeper/internal/accesses/storage"
)

func NewOrder(ord, by string) storage.Param {
	return func(builder *squirrel.SelectBuilder) {
		*builder = builder.OrderByClause(ord, by)
	}
}
