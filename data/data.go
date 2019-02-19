package data

import "fmt"

type TorinoValue interface {
	String() string
	Repr() string
	torinoValue()
}

type TorinoInt struct {
	Value int
}

func (t *TorinoInt) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t *TorinoInt) Repr() string {
	return t.String()
}

func (t *TorinoInt) torinoValue() {}

type TorinoString struct {
	Value string
}

func (t *TorinoString) String() string {
	return fmt.Sprintf("%s", t.Value)
}

func (t *TorinoString) Repr() string {
	// TODO: This won't work well with backslash escapes.
	return fmt.Sprintf("\"%s\"", t.Value)
}

func (t *TorinoString) torinoValue() {}

type TorinoBool struct {
	Value bool
}

func (t *TorinoBool) String() string {
	if t.Value {
		return "true"
	} else {
		return "false"
	}
}

func (t *TorinoBool) Repr() string {
	return t.String()
}

func (t *TorinoBool) torinoValue() {}

type TorinoNone struct {
}

func (t *TorinoNone) torinoValue() {}

func (t *TorinoNone) String() string {
	return "none"
}

func (t *TorinoNone) Repr() string {
	return t.String()
}

type TorinoBuiltin struct {
	F func(...TorinoValue) TorinoValue
}

func (t *TorinoBuiltin) torinoValue() {}

func (t *TorinoBuiltin) String() string {
	return "<built-in function>"
}

func (t *TorinoBuiltin) Repr() string {
	return t.String()
}
