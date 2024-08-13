package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectCountStatement(t *testing.T) {
	count := NewCount(true)
	assert.Equal(t, "$count=true", count.CollectToString())
}
