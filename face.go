package odata

type CollectableToString interface {
	CollectToString() string
}

type CriteriableToString interface {
	CriteriaToString() string
}
