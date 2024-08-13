package odata

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectFilterStatement(t *testing.T) {
	f := NewFilter(
		NewCombination(
			NewCombination(
				NewComparison("a", "b", ComparatorEQ),
				NewComparison("c", "d", ComparatorNEQ),
				And,
			),
			NewComparison("e", "f", ComparatorGT),
			Or,
		),
	)
	got, err := url.QueryUnescape(f.CollectToString())
	assert.NoError(t, err)
	assert.Equal(t, "$filter=((a eq b) and (c ne d)) or (e gt f)", got)
}
