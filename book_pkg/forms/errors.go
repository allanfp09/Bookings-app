package forms

type errors map[string][]string

// Add adds error message to a field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// GetErr returns error if exist
func (e errors) GetErr(field string) string {
	val := e[field]
	if len(val) == 0 {
		return ""
	}

	return val[0]
}
