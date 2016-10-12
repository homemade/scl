package main

import "strings"

type pathSlice []string

func (s pathSlice) String() string {
	return strings.Join(s, ":")
}

func (s *pathSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}
