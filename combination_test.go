package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectCombinationStatement(t *testing.T) {
	c := NewCombination(
		NewCombination(
			NewComparison("a", "b", ComparatorEQ),
			NewComparison("c", "d", ComparatorNEQ),
			And,
		),
		NewComparison("e", "f", ComparatorGT),
		Or,
	)
	assert.Equal(t, "((a eq b) and (c ne d)) or (e gt f)", c.CriteriaToString())
}
