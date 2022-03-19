package config

import "encoding/json"

type Part struct {
	values          map[string]string
	installFields   []string
	temporaryFields []string
	defaultValues   map[string]string
	fieldNames      map[string]string
	validate        func() error
}

func NewPart() *Part {
	return &Part{
		values:        map[string]string{},
		installFields: []string{},
		defaultValues: map[string]string{},
		fieldNames:    map[string]string{},
	}
}

func (h *Part) MarshalJSON() ([]byte, error) {
	v := h.values

	for _, k := range h.TemporaryFields() {
		delete(v, k)
	}

	return json.Marshal(v)
}

func (h *Part) Load(m map[string]string) {
	h.values = m
}

func (h *Part) InstallFields() []string {
	return h.installFields
}

func (h *Part) SetInstallFields(f []string) {
	h.installFields = f
}

func (h *Part) TemporaryFields() []string {
	return h.temporaryFields
}

func (h *Part) SetTemporaryFields(f []string) {
	h.temporaryFields = f
}

func (h *Part) FieldNames() map[string]string {
	return h.fieldNames
}

func (h *Part) SetFieldNames(m map[string]string) {
	h.fieldNames = m
}

func (h *Part) DefaultValues() map[string]string {
	return h.defaultValues
}

func (h *Part) SetDefaultValues(v map[string]string) {
	h.defaultValues = v
}

func (h *Part) Get(key string) string {
	if val, exists := h.values[key]; exists {
		return val
	}

	val, _ := h.Default(key)

	return val
}

func (h *Part) Default(key string) (val string, exists bool) {
	val, exists = h.DefaultValues()[key]

	return val, exists
}

func (h *Part) Set(key string, value string) {
	if h.values == nil {
		h.values = map[string]string{}
	}

	h.values[key] = value
}

func (h *Part) Validate() error {
	if h.validate != nil {
		return h.validate()
	}

	return nil
}

func (h *Part) SetValidate(f func() error) {
	h.validate = f
}
