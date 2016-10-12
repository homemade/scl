package scl

import (
	"fmt"
	"regexp"
	"strings"
)

var inlineVariableMatcher = regexp.MustCompile(`([^$\\]|^)\$([a-zA-Z0-9_]+)`)

type variable struct {
	name  string
	value string
}

type mixin struct {
	declaration *scannerLine
	arguments   []variable
	defaults    []string
}

type scope struct {
	parent      *scope
	branch      *scannerLine
	branchScope *scope
	variables   map[string]*variable
	mixins      map[string]*mixin
}

func newScope() *scope {
	return &scope{
		variables: make(map[string]*variable),
		mixins:    make(map[string]*mixin),
	}
}

func (s *scope) setArgumentVariable(name, value string) {
	s.variables[name] = &variable{name, value}
}

func (s *scope) setVariable(name, value string) {

	v, ok := s.variables[name]

	if !ok || v == nil {
		s.variables[name] = &variable{name, value}
	} else {
		s.variables[name].value = value
	}
}

func (s *scope) variable(name string) string {

	value, ok := s.variables[name]

	if !ok || value == nil {
		return ""
	}

	return s.variables[name].value
}

func (s *scope) setMixin(name string, declaration *scannerLine, argumentTokens []token, defaults []string) {

	mixin := &mixin{
		declaration: declaration,
		defaults:    defaults,
	}

	for _, t := range argumentTokens {
		mixin.arguments = append(mixin.arguments, variable{name: t.content})
	}

	s.mixins[name] = mixin
}

func (s *scope) removeMixin(name string) {
	delete(s.mixins, name)
}

func (s *scope) mixin(name string) (*mixin, error) {

	m, ok := s.mixins[name]

	if !ok {
		return nil, fmt.Errorf("Mixin %s not declared in this scope", name)
	}

	return m, nil
}

func (s *scope) interpolateLiteral(literal string) (outp string, err error) {

	outp = inlineVariableMatcher.ReplaceAllStringFunc(literal, func(match string) string {

		if err != nil {
			return match
		}

		replacement := ""
		prefix := ""

		if match[0] == '$' {
			replacement = string(match[1:])
		} else {
			replacement = string(match[2:])
			prefix = string(match[0])
		}

		if v := s.variable(replacement); v != "" {
			return prefix + v
		}

		err = fmt.Errorf("Unknown variable '$%s'", replacement)

		return match
	})

	// Variables with multiple leading dollars are not matched as variables.
	// To fix this, dollar characters preceded by backslashes are escaped.
	outp = strings.Replace(outp, `\$`, `$`, -1)

	return
}

func (s *scope) clone() *scope {

	s2 := newScope()
	s2.parent = s
	s2.branch = s.branch
	s2.branchScope = s.branchScope

	for k, v := range s.variables {
		s2.variables[k] = v
	}

	for k, v := range s.mixins {
		s2.mixins[k] = v
	}

	return s2
}
