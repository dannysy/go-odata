package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectTopStatement(t *testing.T) {
	top := NewTop(10)
	assert.Equal(t, "$top=10", top.CollectToString())
}
