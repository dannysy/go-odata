package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectExpandStatement(t *testing.T) {
	expand := NewExpandBuilder().
		With("organization",
			NewSelect("id", "name"),
			NewExpandBuilder().With("users", NewSelect("id", "name")).Build()).
		Build()
	assert.Equal(
		t,
		"$expand=organization($select=id,name;$expand=users($select=id,name))",
		expand.CollectToString(),
	)
}
