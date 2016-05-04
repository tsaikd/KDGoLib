package testcase

func addCase(workingCases *ContextType, testcase Case) {
	for _, depend := range testcase.Depends() {
		addCase(workingCases, depend)
	}

	name := testcase.Name()
	if !workingCases.IsExists(name) {
		workingCases.Set(name, testcase)
	}
}

func prepare(workingCases *ContextType, registedCases Queue) {
	registedCases.Walk(func(name string, element interface{}) bool {
		testcase := element.(Case)
		addCase(workingCases, testcase)
		return false
	})
}

func setup(doneCases *ContextType, workingCases *ContextType, context *ContextType, testcase Case) {
	workingCases.debug("setup", testcase.Name())
	value := testcase.Setup(context, testcase)
	context.Set(testcase.Name(), value)
	doneCases.Set(testcase.Name(), testcase)
}

func teardown(doneCases *ContextType, workingCases *ContextType, context *ContextType) {
	for element := doneCases.Last(); element != nil; element = doneCases.Last() {
		testcase := element.(Case)
		if isNeed(workingCases, testcase.Name()) {
			return
		}

		workingCases.debug("teardown", testcase.Name())
		testcase.TearDown(context, testcase, context.Get(testcase.Name()))
		context.Remove(testcase.Name())
		doneCases.Remove(testcase.Name())
	}
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

// isNeed search all test cases and dependency that name is used
func isNeed(t Queue, name string) bool {
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
