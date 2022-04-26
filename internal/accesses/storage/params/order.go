package params

import (
	"github.com/Masterminds/squirrel"
	"pass-keeper/internal/accesses/storage"
)

func NewOrder(ord, by string) storage.Param {
	return func(builder *squirrel.SelectBuilder) {
		builder.OrderByClause(ord, by)
	}
}
