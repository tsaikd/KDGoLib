package testcase

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/testutil/requireutil"
)

func Test_TestCase(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	testBuffer := bytes.NewBuffer(nil)
	generalSetup := func(c Context, t Case) (value interface{}) {
		testBuffer.WriteString("setup " + t.Name() + "\n")
		return nil
	}
	generalTearDown := func(c Context, t Case, value interface{}) {
		testBuffer.WriteString("teardown " + t.Name() + "\n")
	}
	register := NewContext(nil)

	caseRoot1 := register.Regist(
		NewCase("caseRoot1", generalSetup).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot1)

	caseRoot1Child1 := register.Regist(
		NewCase("caseRoot1Child1", generalSetup).
			SetDepends(caseRoot1).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot1Child1)

	caseRoot1Child2 := register.Regist(
		NewCase("caseRoot1Child2", generalSetup).
			SetDepends(caseRoot1).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot1Child2)

	caseRoot2 := register.Regist(
		NewCase("caseRoot2", generalSetup).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot2)

	caseRoot3 := register.Regist(
		NewCase("caseRoot3", generalSetup).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot3)

	caseRoot3Child1 := register.Regist(
		NewCase("caseRoot3Child1", generalSetup).
			SetDepends(caseRoot3).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot3Child1)

	caseRoot3Child1Child1 := register.Regist(
		NewCase("caseRoot3Child1Child1", generalSetup).
			SetDepends(caseRoot3Child1).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseRoot3Child1Child1)

	caseMix1 := register.Regist(
		NewCase("caseMix1", generalSetup).
			SetDepends(caseRoot1, caseRoot1Child2).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseMix1)

	caseMix2 := register.Regist(
		NewCase("caseMix2", generalSetup).
			SetDepends(caseRoot2, caseRoot1Child1).
			SetTearDown(generalTearDown),
	)
	require.NotNil(caseMix2)

	register.StartTest(t, func(c Context, testcase Case) {
		testBuffer.WriteString("testing " + testcase.Name() + "\n")
	}, nil)

	expectedBuffer := strings.TrimSpace(`
testing caseMix1
setup caseRoot1
setup caseRoot1Child2
setup caseMix1
teardown caseMix1
testing caseMix2
setup caseRoot2
setup caseRoot1Child1
setup caseMix2
teardown caseMix2
testing caseRoot1
testing caseRoot1Child1
teardown caseRoot1Child1
testing caseRoot1Child2
teardown caseRoot1Child2
teardown caseRoot1
testing caseRoot2
teardown caseRoot2
testing caseRoot3
setup caseRoot3
testing caseRoot3Child1
setup caseRoot3Child1
testing caseRoot3Child1Child1
setup caseRoot3Child1Child1
teardown caseRoot3Child1Child1
teardown caseRoot3Child1
teardown caseRoot3
	`) + "\n"
	requireutil.RequireText(t, expectedBuffer, testBuffer.String())
}
