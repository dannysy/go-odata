package odata

import (
	"fmt"
	"net/url"
)

const filterKey = "$filter"

type Filter struct {
	criteria CriteriableToString
}

func NewFilter(criteria CriteriableToString) *Filter {
	return &Filter{
		criteria: criteria,
	}
}

func (f *Filter) CollectToString() string {
	return fmt.Sprintf("%s=%s", filterKey, url.QueryEscape(f.criteria.CriteriaToString()))
}
