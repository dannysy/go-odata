package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectComparisonStatement(t *testing.T) {
	c := NewComparison("a", "b", ComparatorNot)
	assert.Equal(t, "a not b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorEQ)
	assert.Equal(t, "a eq b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorNEQ)
	assert.Equal(t, "a ne b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorGT)
	assert.Equal(t, "a gt b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorGTE)
	assert.Equal(t, "a ge b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorLT)
	assert.Equal(t, "a lt b", c.CriteriaToString())
	c = NewComparison("a", "b", ComparatorLTE)
	assert.Equal(t, "a le b", c.CriteriaToString())
}
