package object

import "fmt"

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBG = "BOOLEAN"
	NIL_OBJ    = "NIL"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBG }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type NIL struct{}

func (n *NIL) Type() ObjectType { return NIL_OBJ }
func (n *NIL) Inspect() string  { return "nil" }
