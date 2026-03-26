package scripting

import (
	"fmt"
	"strings"

	"github.com/ghetzel/go-stockutil/log"
)

type Expression struct {
	statement *Statement
	node      *node32
}

func NewExpression(statement *Statement, node *node32) *Expression {
	if node == nil || node.first(ruleValueYielding) == nil {
		log.Fatal("expression node must have a ValueYielding child")
		return nil
	}

	return &Expression{
		statement: statement,
		node:      node,
	}
}

func (self *Expression) Script() *Friendscript {
	return self.statement.Script()
}

func (self *Expression) String() string {
	str := strings.TrimSpace(self.statement.raw(self.node))
	return str
}

func (self *Expression) Value() (any, error) {
	// example: if "x" == "y"

	if lhs := self.node.first(ruleExpressionLHS); lhs != nil { // if "x"
		if value, err := self.resolveValue(
			lhs.firstChild(ruleValueYielding),
		); err == nil { // "x"
			if rhs := self.node.first(ruleExpressionRHS); rhs != nil { // == "y"
				if op, err := parseOperator(rhs.firstChild(ruleOperator)); err == nil { // ==
					if exprNode := rhs.firstChild(ruleExpression); exprNode != nil { // "y"
						return op.evaluate(value, NewExpression(self.statement, exprNode))
					}
				} else if op != opNull {
					return new(emptyValue), err
				}
			}
			return value, nil
		} else {
			return nil, fmt.Errorf("invalid value: %v", err)
		}
	} else {
		return nil, fmt.Errorf("left-hand side of expression did not yield a value")
	}
}

func (self *Expression) resolveValue(node *node32) (any, error) {
	if cmdNode := node.firstN(1, ruleInlineCommand, ruleCommand); cmdNode != nil { // evaluate commands to retrieve value
		return self.statement.evaluateCommand(cmdNode)
	} else if varNode := node.firstN(1, ruleVariable); varNode != nil { // expand variable to reach final value
		return self.statement.resolveVariable(varNode)
	} else if typeNode := node.firstN(1, ruleType); typeNode != nil { // use the literal value
		return self.statement.parseValue(typeNode)
	} else {
		return nil, fmt.Errorf("invalid value argument '%v'", self.statement.raw(node))
	}
}

func exprToValue(in any) (any, error) {
	if IsEmpty(in) {
		return nil, nil
	} else if expr, ok := in.(*Expression); ok && expr != nil {
		if v, err := expr.Value(); err == nil {
			in = v
		} else {
			return nil, err
		}
	} else if expr, ok := in.(Expression); ok {
		if v, err := expr.Value(); err == nil {
			in = v
		} else {
			return nil, err
		}
	}

	return in, nil
}
