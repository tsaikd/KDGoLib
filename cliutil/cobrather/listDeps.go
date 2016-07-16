package cobrather

// ListDepsOption options for ListDeps
type ListDepsOption int8

const (
	// OIncludeCommand list result contains commands in Module
	OIncludeCommand ListDepsOption = 1 << iota
	// OIncludeDepInCommand list result contains dependencies in commands, but not commands self
	OIncludeDepInCommand
)

// Has return true if ListDepsOption contains opt
func (t ListDepsOption) Has(opt ListDepsOption) bool {
	return t&opt == opt
}

// ListDeps list all dependent modules recursively
func ListDeps(opt ListDepsOption, modules ...*Module) []*Module {
	walked := map[*Module]bool{}
	result := []*Module{}
	for _, module := range modules {
		listDeps(opt, module, &result, walked)
	}
	return result
}

func listDeps(opt ListDepsOption, from *Module, result *[]*Module, walked map[*Module]bool) {
	for _, module := range from.Dependencies {
		listDeps(opt, module, result, walked)

		if !walked[module] {
			walked[module] = true
			*result = append(*result, module)
		}
	}

	if opt.Has(OIncludeCommand) || opt.Has(OIncludeDepInCommand) {
		for _, module := range from.Commands {
			listDeps(opt, module, result, walked)

			if opt.Has(OIncludeCommand) && !walked[module] {
				walked[module] = true
				*result = append(*result, module)
			}
		}
	}
}
