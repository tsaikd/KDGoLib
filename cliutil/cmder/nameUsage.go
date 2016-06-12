package cmder

type nameUsage map[string]bool

func (t nameUsage) IsUsed(name string) bool {
	_, used := t[name]
	return used
}

func (t nameUsage) Use(name string) {
	t[name] = true
}
