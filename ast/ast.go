package ast

import (
	"github.com/user/courseWork/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) stmtNode() {

}

func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) exprNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
