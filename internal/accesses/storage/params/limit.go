package params

import (
	"github.com/Masterminds/squirrel"
	"pass-keeper/internal/accesses/storage"
)

func NewLimit(l uint64) storage.Param {
	return func(builder *squirrel.SelectBuilder) {
		*builder = builder.Limit(l)
	}
}
