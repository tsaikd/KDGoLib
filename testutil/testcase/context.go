package testcase

import (
	"testing"

	"github.com/tsaikd/KDGoLib/logutil"
)

// NewContext create Context
func NewContext(logger logutil.LevelLogger) Context {
	return &ContextType{
		logger: logger,
	}
}

// Context is context for running test case
type Context interface {
	Regist(testcase Case) Case
	IsNeed(name string) bool
	Get(name string) interface{}
	StartTest(t *testing.T, beforeEach func(Context, Case), afterEach func(Context, Case))
	T() *testing.T
}

// ContextType implement Context
type ContextType struct {
	SortedMap

	logger       logutil.LevelLogger
	t            *testing.T
	injectingMap map[string]bool
}

// Regist a test case
func (t *ContextType) Regist(testcase Case) Case {
	t.Set(testcase.Name(), testcase)
	return testcase
}

func (t *ContextType) injecting(testcase Case, callback func()) {
	name := testcase.Name()
	if t.injectingMap == nil {
		t.injectingMap = map[string]bool{}
	}
	if _, exist := t.injectingMap[name]; exist {
		return
	}
	t.injectingMap[name] = true
	callback()
	delete(t.injectingMap, name)
}

func (t *ContextType) setup(testcase Case) {
	// maybe setup by other dependency
	if t.IsExists(testcase.Name()) {
		return
	}

	t.debug("prepare", testcase.Name())

	// prepare dependency
	for _, depend := range testcase.Depends() {
		if !t.IsExists(depend.Name()) {
			t.injecting(depend, func() {
				t.setup(depend)
			})
		}
	}

	t.debug("setup", testcase.Name())
	value := testcase.Setup(t, testcase)
	t.Set(testcase.Name(), value)
}

func (t *ContextType) teardown(testcase Case, restTestCase Context) {
	if !restTestCase.IsNeed(testcase.Name()) {
		t.debug("teardown", testcase.Name())
		testcase.TearDown(t, testcase, t.Get(testcase.Name()))
		t.Remove(testcase.Name())
	}

	depends := testcase.Depends()
	for i := len(depends) - 1; i >= 0; i-- {
		depend := depends[i]
		t.injecting(depend, func() {
			t.teardown(depend, restTestCase)
		})
	}
}

func (t *ContextType) debug(msgs ...interface{}) {
	if t.logger == nil {
		return
	}
	t.logger.Debugln(msgs...)
}

func isTestCaseNeed(testcase Case, name string, checkedMap map[string]bool) bool {
	if need, exist := checkedMap[testcase.Name()]; exist {
		return need
	}

	if name == testcase.Name() {
		checkedMap[testcase.Name()] = true
		return true
	}
	for _, depend := range testcase.Depends() {
		if isTestCaseNeed(depend, name, checkedMap) {
			checkedMap[testcase.Name()] = true
			return true
		}
	}

	checkedMap[testcase.Name()] = false
	return false
}

// IsNeed search all test cases and dependency that name is used
func (t *ContextType) IsNeed(name string) bool {
	need := false
	checkedMap := map[string]bool{}
	t.Walk(func(key string, value interface{}) bool {
		testcase := value.(Case)
		if isTestCaseNeed(testcase, name, checkedMap) {
			need = true
			return true
		}
		return false
	})
	return need
}

// StartTest start all registed testcase
func (t *ContextType) StartTest(test *testing.T, beforeEach func(Context, Case), afterEach func(Context, Case)) {
	context := ContextType{
		logger: t.logger,
		t:      test,
	}
	t.Sort()

	testedCase := map[string]Case{}

	for testcase := t.Shift(); testcase != nil; testcase = t.Shift() {
		tcase := testcase.(Case)
		testedCase[tcase.Name()] = tcase
		if beforeEach != nil {
			beforeEach(t, tcase)
		}
		context.setup(tcase)
		context.teardown(tcase, t)
		if afterEach != nil {
			afterEach(t, tcase)
		}
	}

	// teardown rest testcase which still exist because dependency
	maxrun := len(testedCase)
	for len(testedCase) > 0 {
		maxrun--
		if maxrun < 0 {
			panic("teardown rest testcase reach max run")
		}
		for name, testcase := range testedCase {
			if context.IsExists(name) {
				context.teardown(testcase, t)
			}
			if !context.IsExists(name) {
				delete(testedCase, name)
			}
		}
	}
}

// T return *testing.T injected by StartTest
func (t *ContextType) T() *testing.T {
	return t.t
}
