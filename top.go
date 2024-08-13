package odata

import "fmt"

const topKey = "$top"

type Top struct {
	n int
}

func NewTop(n int) *Top {
	return &Top{n: n}
}

func (t *Top) CollectToString() string {
	return fmt.Sprintf("%s=%d", topKey, t.n)

}
