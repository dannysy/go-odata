package odata

import (
	"strings"
)

type Entity struct {
	name        string
	id          string
	shouldCount bool
	operations  []CollectableToString
}

type builder struct {
	entity *Entity
}

func NewEntityBuilder(name string) EntityBuilder {
	return newBuilder(name)
}

func NewListBuilder(name string) ListBuilder {
	return newBuilder(name)
}

func NewCountBuilder(name string) CountBuilder {
	b := newBuilder(name)
	b.entity.shouldCount = true
	return b
}

func (b *builder) WithId(id string) EntityBuilder {
	b.entity.id = id
	return b
}

func (b *builder) With(ops ...CollectableToString) ListBuilder {
	b.entity.operations = append(b.entity.operations, ops...)
	return b
}

func (b *builder) WithSelect(s *Select) EntityBuilder {
	b.With(s)
	return b
}

func (b *builder) WithExpand(expand *Expand) EntityBuilder {
	b.With(expand)
	return b
}

func (b *builder) WithFilter(filter *Filter) CountBuilder {
	b.With(filter)
	return b
}

func (b *builder) Build() *Entity {
	return b.entity
}

func (e *Entity) CollectToString() string {
	sb := strings.Builder{}
	sb.WriteString(e.name)

	// Id can be requested only with $select & $expand
	if e.id != "" {
		sb.WriteString("(" + e.id + ")")
	}
	// Number of items can be requested only with $filter
	if e.shouldCount {
		sb.WriteString("/$count")
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

func newBuilder(name string) *builder {
	return &builder{
		entity: &Entity{
			name:       name,
			operations: make([]CollectableToString, 0, 4),
		},
	}
}
