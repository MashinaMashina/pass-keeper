package params

import "pass-keeper/internal/accesses/storage"

type order struct {
	order string
	by    string
}

func NewOrder(ord, by string) storage.Param {
	return &order{
		order: ord,
		by:    by,
	}
}

func (_this *order) ParamType() string {
	return "order"
}

func (_this *order) Value() []string {
	return []string{_this.order, _this.by}
}
