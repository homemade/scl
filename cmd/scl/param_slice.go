package main

import "strings"

type paramSlice []*param

func (ps paramSlice) String() string {
	outputSlice := make([]string, len(ps))
	for i, p := range ps {
		outputSlice[i] = p.String()
	}

	return strings.Join(outputSlice, ", ")
}

func (ps *paramSlice) Set(value string) error {
	p := new(param)

	if err := p.SetFromString(value); err != nil {
		return err
	}

	*ps = append(*ps, p)

	return nil
}
