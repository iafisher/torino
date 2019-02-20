package data

import (
	"fmt"
	"strconv"
)

type TorinoValue interface {
	String() string
	Repr() string
	Torino()
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

func (t *TorinoInt) Torino() {}

type TorinoString struct {
	Value string
}

func (t *TorinoString) String() string {
	return fmt.Sprintf("%s", t.Value)
}

func (t *TorinoString) Repr() string {
	return strconv.Quote(t.Value)
}

func (t *TorinoString) Torino() {}

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

func (t *TorinoBool) Torino() {}

type TorinoNone struct {
}

func (t *TorinoNone) Torino() {}

func (t *TorinoNone) String() string {
	return "none"
}

func (t *TorinoNone) Repr() string {
	return t.String()
}

type TorinoBuiltin struct {
	F func(...TorinoValue) TorinoValue
}

func (t *TorinoBuiltin) Torino() {}

func (t *TorinoBuiltin) String() string {
	return "<built-in function>"
}

func (t *TorinoBuiltin) Repr() string {
	return t.String()
}
