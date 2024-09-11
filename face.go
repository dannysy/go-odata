package odata

type CollectableToString interface {
	CollectToString() string
}

type CriteriableToString interface {
	CriteriaToString() string
}

type Builder interface {
	Build() *Entity
}
type ListBuilder interface {
	With(...CollectableToString) ListBuilder
	Builder
}

type EntityBuilder interface {
	WithId(id string) EntityBuilder
	WithSelect(s *Select) EntityBuilder
	WithExpand(e *Expand) EntityBuilder
	Builder
}

type CountBuilder interface {
	WithFilter(f *Filter) CountBuilder
	Builder
}
