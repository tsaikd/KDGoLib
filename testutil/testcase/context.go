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
	T() *testing.T
	Regist(testcase Case) Case
	Get(name string) interface{}
	StartTest(test *testing.T)
}

// ContextType implement Context
type ContextType struct {
	QueueType

	logger logutil.LevelLogger
	t      *testing.T
}

// T return *testing.T injected by StartTest
func (t ContextType) T() *testing.T {
	return t.t
}

// Regist a test case
func (t *ContextType) Regist(testcase Case) Case {
	t.Set(testcase.Name(), testcase)
	return testcase
}

func (t *ContextType) debug(msgs ...interface{}) {
	if t.logger == nil {
		return
	}
	t.logger.Debugln(msgs...)
}

// StartTest start all registed testcase
func (t *ContextType) StartTest(test *testing.T) {
	registedCases := t
	logger := t.logger

	context := ContextType{
		logger: logger,
		t:      test,
	}
	workingCases := ContextType{
		logger: logger,
		t:      test,
	}
	doneCases := ContextType{
		logger: logger,
		t:      test,
	}

	// prepare workingCases order
	prepare(&workingCases, registedCases)

	for element := workingCases.Shift(); element != nil; element = workingCases.Shift() {
		testcase := element.(Case)
		setup(&doneCases, &workingCases, &context, testcase)
		teardown(&doneCases, &workingCases, &context)
	}

	// teardown rest testcase which still exist because dependency
	teardown(&doneCases, &workingCases, &context)
}
