package injectutil

import (
	"errors"
	"testing"

	"github.com/codegangsta/inject"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	assert := assert.New(t)
	inj := inject.New()

	invokefunc := func() error {
		return errors.New("error")
	}

	_, err := inj.Invoke(invokefunc)
	assert.NoError(err)

	_, err = Invoke(inj, invokefunc)
	assert.Error(err)
}
