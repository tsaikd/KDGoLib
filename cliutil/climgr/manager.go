package climgr

// Manager manage all registed cli module
type Manager struct {
	moduleNameMap map[string]Module
}

// NewManager create a new manager instance
func NewManager() *Manager {
	return &Manager{
		moduleNameMap: map[string]Module{},
	}
}
