package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectSelectStatement(t *testing.T) {
	s := NewSelect("a", "b", "c")

	assert.Equal(t, "$select=a,b,c", s.CollectToString())
}
