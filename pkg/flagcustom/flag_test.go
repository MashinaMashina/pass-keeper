package flagcustom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlag(t *testing.T) {
	flags, err := ParseFlags(`-ssh login@host --load "sess name"  -pw qwerty`)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, []string{"-ssh", "login@host", "--load", "sess name", "-pw", "qwerty"}, flags)

	_, err = ParseFlags(`-ssh login@host --load "sess name" -pw qwerty -p "any`)
	if err == nil {
		t.Error("there must be an error, input arguments is invalid")
		return
	}

	if err != ErrEscaping {
		t.Error("error must be ErrEscaping")
		return
	}
}
