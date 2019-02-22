package data

import (
	"fmt"
	"strconv"
	"strings"
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
	F func(...TorinoValue) (TorinoValue, error)
}

func (t *TorinoBuiltin) Torino() {}

func (t *TorinoBuiltin) String() string {
	return "<built-in function>"
}

func (t *TorinoBuiltin) Repr() string {
	return t.String()
}

type TorinoList struct {
	Values []TorinoValue
}

func (t *TorinoList) Torino() {}

func (t *TorinoList) String() string {
	var str strings.Builder

	str.WriteString("[")
	for i, val := range t.Values {
		str.WriteString(val.Repr())
		if i != len(t.Values)-1 {
			str.WriteString(", ")
		}
	}
	str.WriteString("]")
	return str.String()
}

func (t *TorinoList) Repr() string {
	return t.String()
}

type TorinoMap struct {
	// Keys are stored as their repr value, which is hacky but simple and allows
	// any TorinoValue to be a key.
	Values map[string]TorinoValue
}

func (t *TorinoMap) Torino() {}

func (t *TorinoMap) String() string {
	var str strings.Builder

	str.WriteString("{")
	i := 0
	for key, val := range t.Values {
		str.WriteString(key)
		str.WriteString(": ")
		str.WriteString(val.Repr())

		if i != len(t.Values)-1 {
			str.WriteString(", ")
		}
		i += 1
	}
	str.WriteString("}")
	return str.String()
}

func (t *TorinoMap) Repr() string {
	return t.String()
}

func (t *TorinoMap) Get(key TorinoValue) (TorinoValue, bool) {
	val, ok := t.Values[key.Repr()]
	return val, ok
}

func (t *TorinoMap) Put(key TorinoValue, val TorinoValue) {
	t.Values[key.Repr()] = val
}
