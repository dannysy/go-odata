package odata

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldBuildCorrectOrderByStatement(t *testing.T) {

	orderBy := NewOrderByBuilder().
		WithColumns("a", "b", "c").
		WithDirectionalColumn("d", Asc).
		WithDirectionalColumn("e", Desc).
		Build()
	got, err := url.QueryUnescape(orderBy.CollectToString())
	assert.NoError(t, err)
	assert.Equal(t, "$orderby=a,b,c,d asc,e desc", got)
}
