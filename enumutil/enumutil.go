package enumutil

type enumFactory struct {
	*enumBase
}

func (t *enumFactory) Add(enum interface{}, text string) *enumFactory {
	t.mape2s[enum] = text
	t.maps2e[text] = enum
	return t
}

func (t *enumFactory) Build() *enumBase {
	return t.enumBase
}

func NewEnumFactory() *enumFactory {
	return &enumFactory{
		&enumBase{
			mape2s: map[interface{}]string{},
			maps2e: map[string]interface{}{},
		},
	}
}
