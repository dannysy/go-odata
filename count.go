package odata

import (
	"fmt"
	"strconv"
)

const countKey = "$count"

type Count struct {
	withResult bool
}

func NewCount(withResult bool) *Count {
	return &Count{withResult: withResult}
}

func (c *Count) CollectToString() string {
	return fmt.Sprintf("%s=%s", countKey, strconv.FormatBool(c.withResult))
}

func TotalCount(entityName string) string {
	return fmt.Sprintf("%s/$count", entityName)
}
