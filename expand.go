package odata

import (
	"fmt"
	"strings"
)

const expandKey = "$expand"

type Expand struct {
	colsAndQueries map[string][]CollectableToString
}

type ExpandBuilder struct {
	e *Expand
}

func NewExpandBuilder() *ExpandBuilder {
	return &ExpandBuilder{
		e: &Expand{colsAndQueries: make(map[string][]CollectableToString)},
	}
}

func (b *ExpandBuilder) With(column string, queries ...CollectableToString) *ExpandBuilder {
	b.e.colsAndQueries[column] = queries
	return b
}

func (b *ExpandBuilder) Build() *Expand {
	return b.e
}

func (e *Expand) CollectToString() string {
	collectExpandsFn := func(inners []CollectableToString) string {
		sb := strings.Builder{}
		for i, inner := range inners {
			sb.WriteString(inner.CollectToString())
			if i != len(inners)-1 {
				sb.WriteString(";")
			}

		}
		return sb.String()
	}
	n := len(e.colsAndQueries)
	if n == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s=", expandKey))
	for col, queries := range e.colsAndQueries {
		if len(queries) == 0 {
			sb.WriteString(col)
		} else {
			sb.WriteString(fmt.Sprintf("%s(%s)", col, collectExpandsFn(queries)))
		}
		n--
		if n > 0 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
