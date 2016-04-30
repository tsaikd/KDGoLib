package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorFactory(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factory := NewFactory("factory error")
	require.Equal("factory error", factory.Error())

	err := factory.New(nil)
	require.Error(err)
	require.Equal("factory error", err.Error())
	require.True(factory.Match(err))
	require.Equal("ErrorFactory_test.go", err.FileName())

	switch FactoryOf(err) {
	case factory:
	default:
		t.Fatal("Invalid factory switch case")
	}
}

func Test_ErrorFactory_with_param(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factory := NewFactory("factory error with param %d")
	require.Equal("factory error with param %d", factory.Error())

	err := factory.New(nil, 123)
	require.Error(err)
	require.Equal("factory error with param 123", err.Error())
	require.True(factory.Match(err))
}

func Test_ErrorFactory_chain(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factory1 := NewFactory("factory error chain 1")
	factory2 := NewFactory("factory error chain 2")

	err1 := factory1.New(nil)
	err2 := factory2.New(err1)
	require.Error(err2)
	require.Equal("factory error chain 2; factory error chain 1", err2.Error())
	require.True(factory2.Match(err2))
	require.False(factory1.Match(err2))
	require.True(factory1.In(err2))
}

func Test_ErrorFactory_chain_with_origin_error(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factory := NewFactory("factory error chain")

	err := factory.New(errors.New("origin error"))
	require.Error(err)
	require.Equal("factory error chain; origin error", err.Error())
	require.True(factory.Match(err))
	require.NotNil(err.Parent())
	require.Equal("origin error", err.Parent().Error())
	require.Equal("ErrorFactory_test.go", err.Parent().FileName())
	require.Equal(65, err.Parent().Line())
	require.Nil(err.Parent().Parent())
	require.Nil(err.Parent().Factory())
}

func Test_ErrorFactory_sorted(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factories := AllSortedNamedFactories()
	require.Len(factories, 7)
	require.Equal(ErrorWalkLoop, factories[0])
	require.Equal(testFactory, factories[6])
}
