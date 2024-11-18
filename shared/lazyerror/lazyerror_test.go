package lazyerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unwrap(err error, n int) error {
	for i := 0; i < n; i++ {
		err = errors.Unwrap(err)
	}
	return err
}

func TestErrors(t *testing.T) {
	t.Parallel()

	err := New("err")
	err1 := Errorf("err1: %w", err)
	err2 := Errorf("err2: %w", err1)
	err3 := Errorf("err3: %w", err2)

	expected := "[lazyerror_test.go:20 lazyerror.TestErrors] err"
	assert.Equal(t, expected, err.Error())
	expected = "[lazyerror_test.go:21 lazyerror.TestErrors] err1: " +
		"[lazyerror_test.go:20 lazyerror.TestErrors] err"
	assert.Equal(t, expected, err1.Error())
	expected = "[lazyerror_test.go:22 lazyerror.TestErrors] err2: " +
		"[lazyerror_test.go:21 lazyerror.TestErrors] err1: " +
		"[lazyerror_test.go:20 lazyerror.TestErrors] err"
	assert.Equal(t, expected, err2.Error())
	expected = "[lazyerror_test.go:23 lazyerror.TestErrors] err3: " +
		"[lazyerror_test.go:22 lazyerror.TestErrors] err2: " +
		"[lazyerror_test.go:21 lazyerror.TestErrors] err1: " +
		"[lazyerror_test.go:20 lazyerror.TestErrors] err"
	assert.Equal(t, expected, err3.Error())

	assert.NotEqual(t, err, unwrap(err1, 1))
	assert.Equal(t, err, unwrap(err1, 2))
	assert.NotEqual(t, nil, unwrap(err1, 3))
	assert.Equal(t, nil, unwrap(err1, 4))

	assert.NotEqual(t, err1, unwrap(err2, 1))
	assert.Equal(t, err1, unwrap(err2, 2))
	assert.NotEqual(t, err, unwrap(err2, 3))
	assert.Equal(t, err, unwrap(err2, 4))
	assert.NotEqual(t, nil, unwrap(err2, 5))
	assert.Equal(t, nil, unwrap(err2, 6))

	assert.NotEqual(t, err2, unwrap(err3, 1))
	assert.Equal(t, err2, unwrap(err3, 2))
	assert.NotEqual(t, err1, unwrap(err3, 3))
	assert.Equal(t, err1, unwrap(err3, 4))
	assert.NotEqual(t, err, unwrap(err3, 5))
	assert.Equal(t, err, unwrap(err3, 6))
	assert.NotEqual(t, nil, unwrap(err3, 7))
	assert.Equal(t, nil, unwrap(err3, 8))

	assert.True(t, errors.Is(err3, err3))
	assert.True(t, errors.Is(err3, err2))
	assert.True(t, errors.Is(err3, err1))
	assert.True(t, errors.Is(err3, err))
}
