package odata

import (
	"strings"
)

type Entity struct {
	name       string
	id         string
	operations []CollectableToString
}

type EntityBuilder struct {
	e *Entity
}

func NewEntityBuilder(name string) *EntityBuilder {
	return &EntityBuilder{
		e: &Entity{
			name:       name,
			operations: make([]CollectableToString, 0, 4),
		},
	}
}

func (b *EntityBuilder) With(ops ...CollectableToString) *EntityBuilder {
	b.e.operations = append(b.e.operations, ops...)
	return b
}

func (b *EntityBuilder) WithId(id string) *EntityBuilder {
	b.e.id = id
	return b

}

func (b *EntityBuilder) Build() *Entity {
	return b.e
}

func (e *Entity) CollectToString() string {
	sb := strings.Builder{}
	sb.WriteString(e.name)
	if e.id != "" {
		sb.WriteString("(" + e.id + ")")
	}
	sb.WriteString("?")
	for i, op := range e.operations {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(op.CollectToString())
	}
	return sb.String()
}
