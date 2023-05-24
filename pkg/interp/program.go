package interp

import "procinspect/pkg/semantic"

type (
	Program struct {
		Script     *semantic.Script
		Procedures []*Procedure
	}

	Procedure struct {
		Name string
	}
)
