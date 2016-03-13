package enumutil

type enumFactory struct {
	*enumBase
}

func (t *enumFactory) Add(enum EnumString, text string) *enumFactory {
	t.mape2s[enum] = text
	t.maps2e[text] = enum
	return t
}

func (t *enumFactory) Build() *enumBase {
	return t.enumBase
}

// NewEnumFactory start a enum factory
func NewEnumFactory() *enumFactory {
	return &enumFactory{
		&enumBase{
			mape2s: map[EnumString]string{},
			maps2e: map[string]interface{}{},
		},
	}
}
