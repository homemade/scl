package main

import (
	"errors"
	"fmt"
	"strings"
)

type param struct {
	name  string
	value string
}

func (p param) String() string {
	return p.name + "=" + p.value
}

// expects 'name=value'
func (p *param) SetFromString(s string) error {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) <= 1 {
		return errors.New("Unable to convert to param: " + s)
	}

	p.name = strings.TrimSpace(parts[0])
	p.value = fmt.Sprintf(`"%s"`, strings.TrimSpace(parts[1]))

	return nil
}
