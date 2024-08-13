package odata

import (
	"strings"
)

const orderByKey = "$orderby"

type Direction string

const (
	Default Direction = ""
	Asc     Direction = "asc"
	Desc    Direction = "desc"
)

type order struct {
	column    string
	direction Direction
}

type OrderBy struct {
	orders []order
}

type OrderByBuilder struct {
	o *OrderBy
}

func NewOrderByBuilder() *OrderByBuilder {
	return &OrderByBuilder{
		o: &OrderBy{
			orders: make([]order, 0, 4),
		},
	}
}
func (b *OrderByBuilder) WithColumns(columns ...string) *OrderByBuilder {
	for _, column := range columns {
		b.o.orders = append(b.o.orders, order{column: column, direction: Default})
	}
	return b
}
func (b *OrderByBuilder) WithDirectionalColumn(column string, direction Direction) *OrderByBuilder {
	b.o.orders = append(b.o.orders, order{column: column, direction: direction})
	return b
}
func (b *OrderByBuilder) Build() *OrderBy {
	return b.o
}

func (t *OrderBy) CollectToString() string {
	sb := strings.Builder{}
	sb.WriteString(orderByKey + "=")
	for i, o := range t.orders {
		sb.WriteString(o.column)
		if o.direction != Default {
			// need to escape space  - %20 is ' '
			sb.WriteString("%20" + string(o.direction))
		}
		if i < len(t.orders)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
