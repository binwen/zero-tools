package stringx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	cases := []struct {
		slice  []string
		value  string
		expect bool
	}{
		{[]string{"1"}, "1", true},
		{[]string{"1"}, "2", false},
		{[]string{"1", "2"}, "1", true},
		{[]string{"1", "2"}, "3", false},
		{nil, "3", false},
		{nil, "", false},
	}

	for _, each := range cases {
		t.Run(path.Join(each.slice...), func(t *testing.T) {
			actual := Contains(each.slice, each.value)
			assert.Equal(t, each.expect, actual)
		})
	}
}
