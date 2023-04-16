package response

type Error struct {
	msg string
	src error
	val Validations
}

func NewError(msg string) *Error { return &Error{msg: msg} }

type Validations map[string][]string

func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return e.src }

func (e *Error) WithMessage(msg string) *Error        { e.msg = msg; return e }
func (e *Error) WithSource(err error) *Error          { e.src = err; return e }
func (e *Error) WithValidations(v Validations) *Error { e.val = v; return e }
func (e *Error) Payload() any {
	return struct {
		Message     string      `json:"message"`
		Source      string      `json:"source,omitempty"`
		Validations Validations `json:"validations,omitempty"`
	}{Message: e.msg, Source: e.src.Error(), Validations: e.val}
}

func (v Validations) Add(field, msg string) { v[field] = append(v[field], msg) }
func (v Validations) Set(field, msg string) { v[field] = []string{msg} }
func (v Validations) Has(field string) bool { _, ok := v[field]; return ok }
func (v Validations) Get(field string) []string {
	if v.Has(field) {
		return v[field]
	}
	return nil
}
