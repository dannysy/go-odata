package odata

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBuildCorrectEntityStatement(t *testing.T) {
	entity := NewEntityBuilder("Products").WithId("1").
		WithSelect(
			NewSelect("id", "name"),
		).
		WithExpand(
			NewExpandBuilder().With("category", NewSelect("id", "name")).Build(),
		).
		Build()

	got, err := url.QueryUnescape(entity.CollectToString())
	assert.NoError(t, err)
	assert.Equal(
		t,
		`Products(1)?$select=id,name&$expand=category($select=id,name)`,
		got,
	)
}

func TestShouldBuildCorrectListStatement(t *testing.T) {
	entity := NewListBuilder("Products").With(
		NewSelect("id", "name"),
		NewExpandBuilder().With("category", NewSelect("id", "name")).Build(),
		NewFilter(
			NewCombination(
				NewComparison("a", "b", ComparatorEQ),
				NewComparison("c", "d", ComparatorNEQ),
				And,
			),
		),
		NewTop(10),
		NewOrderByBuilder().WithColumns("id").WithDirectionalColumn("name", Asc).Build(),
	).Build()

	got, err := url.QueryUnescape(entity.CollectToString())
	assert.NoError(t, err)
	assert.Equal(
		t,
		`Products?$select=id,name&$expand=category($select=id,name)&$filter=(a eq b) and (c ne d)&$top=10&$orderby=id,name asc`,
		got,
	)
}

func TestShouldBuildCorrectEntitiesCountStatement(t *testing.T) {
	entity := NewCountBuilder("Products").WithFilter(
		NewFilter(
			NewCombination(
				NewComparison("a", "b", ComparatorEQ),
				NewComparison("c", "d", ComparatorNEQ),
				And,
			),
		),
	).Build()

	got, err := url.QueryUnescape(entity.CollectToString())
	assert.NoError(t, err)
	assert.Equal(
		t,
		`Products/$count?$filter=(a eq b) and (c ne d)`,
		got,
	)
}
