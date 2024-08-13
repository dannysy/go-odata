package odata

import (
	"fmt"
)

type Operator string

const (
	And Operator = "and"
	Or  Operator = "or"
)

type Combination struct {
	left     CriteriableToString
	right    CriteriableToString
	operator Operator
}

func NewCombination(left, right CriteriableToString, operator Operator) *Combination {
	return &Combination{
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (c *Combination) CriteriaToString() string {
	return fmt.Sprintf("(%s) %s (%s)", c.left.CriteriaToString(), c.operator, c.right.CriteriaToString())
}
