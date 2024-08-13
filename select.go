package odata

import (
	"fmt"
	"strings"
)

const selectKey = "$select"

type Select struct {
	columns []string
}

func NewSelect(columns ...string) *Select {
	return &Select{
		columns: columns,
	}
}

func (s *Select) CollectToString() string {
	return fmt.Sprintf("%s=%s", selectKey, strings.Join(s.columns, ","))
}
