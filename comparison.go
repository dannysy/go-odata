package odata

import (
	"fmt"
)

type Comparator string

const (
	ComparatorEQ  Comparator = "eq"
	ComparatorNEQ Comparator = "ne"
	ComparatorGT  Comparator = "gt"
	ComparatorGTE Comparator = "ge"
	ComparatorLT  Comparator = "lt"
	ComparatorLTE Comparator = "le"
	ComparatorNot Comparator = "not"
)

type Comparison struct {
	left     string
	right    string
	comparer Comparator
}

func NewComparison(left, right string, comparer Comparator) *Comparison {
	return &Comparison{
		left:     left,
		right:    right,
		comparer: comparer,
	}
}

func (c *Comparison) CriteriaToString() string {
	return fmt.Sprintf("%s %s %s", c.left, c.comparer, c.right)
}
