package data

import "fmt"

type TorinoValue interface {
	String() string
	torinoValue()
}

type TorinoInt struct {
	Value int64
}

func (t *TorinoInt) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t *TorinoInt) torinoValue() {}

type TorinoString struct {
	Value string
}

func (t *TorinoString) String() string {
	// TODO: This won't play nicely with backslash escapes.
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

func (t *TorinoBool) torinoValue() {}

type TorinoNone struct {
}

func (t *TorinoNone) torinoValue() {}

func (t *TorinoNone) String() string {
	return "none"
}
