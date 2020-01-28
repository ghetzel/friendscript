package scripting

//go:generate peg -inline friendscript.peg

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleFriendscript
	rule_
	rule__
	ruleASSIGN
	ruleTRIQUOT
	ruleBREAK
	ruleCLOSE
	ruleCOLON
	ruleCOMMA
	ruleCOMMENT
	ruleCONT
	ruleCOUNT
	ruleDECLARE
	ruleDOT
	ruleELSE
	ruleIF
	ruleIN
	ruleINCLUDE
	ruleLOOP
	ruleNOOP
	ruleNOT
	ruleOPEN
	ruleSCOPE
	ruleSEMI
	ruleSHEBANG
	ruleSKIPVAR
	ruleUNSET
	ruleScalarType
	ruleIdentifier
	ruleFloat
	ruleBoolean
	ruleInteger
	rulePositiveInteger
	ruleString
	ruleStringLiteral
	ruleStringInterpolated
	ruleTriquote
	ruleTriquoteBody
	ruleNullValue
	ruleObject
	ruleArray
	ruleRegularExpression
	ruleKeyValuePair
	ruleKey
	ruleKValue
	ruleType
	ruleExponentiate
	ruleMultiply
	ruleDivide
	ruleModulus
	ruleAdd
	ruleSubtract
	ruleBitwiseAnd
	ruleBitwiseOr
	ruleBitwiseNot
	ruleBitwiseXor
	ruleMatchOperator
	ruleUnmatch
	ruleMatch
	ruleOperator
	ruleAssignmentOperator
	ruleAssignEq
	ruleStarEq
	ruleDivEq
	rulePlusEq
	ruleMinusEq
	ruleAndEq
	ruleOrEq
	ruleAppend
	ruleComparisonOperator
	ruleEquality
	ruleNonEquality
	ruleGreaterThan
	ruleGreaterEqual
	ruleLessEqual
	ruleLessThan
	ruleMembership
	ruleNonMembership
	ruleVariable
	ruleVariableNameSequence
	ruleVariableName
	ruleVariableIndex
	ruleBlock
	ruleFlowControlWord
	ruleFlowControlBreak
	ruleFlowControlContinue
	ruleStatementBlock
	ruleAssignment
	ruleAssignmentLHS
	ruleAssignmentRHS
	ruleVariableSequence
	ruleExpressionSequence
	ruleExpression
	ruleExpressionLHS
	ruleExpressionRHS
	ruleValueYielding
	ruleDirective
	ruleDirectiveUnset
	ruleDirectiveInclude
	ruleDirectiveDeclare
	ruleCommand
	ruleCommandName
	ruleCommandFirstArg
	ruleCommandSecondArg
	ruleCommandResultAssignment
	ruleConditional
	ruleIfStanza
	ruleElseIfStanza
	ruleElseStanza
	ruleLoop
	ruleLoopConditionFixedLength
	ruleLoopConditionIterable
	ruleLoopIterableLHS
	ruleLoopIterableRHS
	ruleLoopConditionBounded
	ruleLoopConditionTruthy
	ruleConditionalExpression
	ruleConditionWithAssignment
	ruleConditionWithCommand
	ruleConditionWithRegex
	ruleConditionWithComparator
	ruleConditionWithComparatorLHS
	ruleConditionWithComparatorRHS
)

var rul3s = [...]string{
	"Unknown",
	"Friendscript",
	"_",
	"__",
	"ASSIGN",
	"TRIQUOT",
	"BREAK",
	"CLOSE",
	"COLON",
	"COMMA",
	"COMMENT",
	"CONT",
	"COUNT",
	"DECLARE",
	"DOT",
	"ELSE",
	"IF",
	"IN",
	"INCLUDE",
	"LOOP",
	"NOOP",
	"NOT",
	"OPEN",
	"SCOPE",
	"SEMI",
	"SHEBANG",
	"SKIPVAR",
	"UNSET",
	"ScalarType",
	"Identifier",
	"Float",
	"Boolean",
	"Integer",
	"PositiveInteger",
	"String",
	"StringLiteral",
	"StringInterpolated",
	"Triquote",
	"TriquoteBody",
	"NullValue",
	"Object",
	"Array",
	"RegularExpression",
	"KeyValuePair",
	"Key",
	"KValue",
	"Type",
	"Exponentiate",
	"Multiply",
	"Divide",
	"Modulus",
	"Add",
	"Subtract",
	"BitwiseAnd",
	"BitwiseOr",
	"BitwiseNot",
	"BitwiseXor",
	"MatchOperator",
	"Unmatch",
	"Match",
	"Operator",
	"AssignmentOperator",
	"AssignEq",
	"StarEq",
	"DivEq",
	"PlusEq",
	"MinusEq",
	"AndEq",
	"OrEq",
	"Append",
	"ComparisonOperator",
	"Equality",
	"NonEquality",
	"GreaterThan",
	"GreaterEqual",
	"LessEqual",
	"LessThan",
	"Membership",
	"NonMembership",
	"Variable",
	"VariableNameSequence",
	"VariableName",
	"VariableIndex",
	"Block",
	"FlowControlWord",
	"FlowControlBreak",
	"FlowControlContinue",
	"StatementBlock",
	"Assignment",
	"AssignmentLHS",
	"AssignmentRHS",
	"VariableSequence",
	"ExpressionSequence",
	"Expression",
	"ExpressionLHS",
	"ExpressionRHS",
	"ValueYielding",
	"Directive",
	"DirectiveUnset",
	"DirectiveInclude",
	"DirectiveDeclare",
	"Command",
	"CommandName",
	"CommandFirstArg",
	"CommandSecondArg",
	"CommandResultAssignment",
	"Conditional",
	"IfStanza",
	"ElseIfStanza",
	"ElseStanza",
	"Loop",
	"LoopConditionFixedLength",
	"LoopConditionIterable",
	"LoopIterableLHS",
	"LoopIterableRHS",
	"LoopConditionBounded",
	"LoopConditionTruthy",
	"ConditionalExpression",
	"ConditionWithAssignment",
	"ConditionWithCommand",
	"ConditionWithRegex",
	"ConditionWithComparator",
	"ConditionWithComparatorLHS",
	"ConditionWithComparatorRHS",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Friendscript struct {
	runtime

	Buffer string
	buffer []rune
	rules  [124]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Friendscript) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Friendscript) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Friendscript
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Friendscript) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Friendscript) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Friendscript <- <(_ SHEBANG? _ Block* !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rule_]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					{
						position4 := position
						if buffer[position] != rune('#') {
							goto l2
						}
						position++
						if buffer[position] != rune('!') {
							goto l2
						}
						position++
						{
							position7, tokenIndex7 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l7
							}
							position++
							goto l2
						l7:
							position, tokenIndex = position7, tokenIndex7
						}
						if !matchDot() {
							goto l2
						}
					l5:
						{
							position6, tokenIndex6 := position, tokenIndex
							{
								position8, tokenIndex8 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l8
								}
								position++
								goto l6
							l8:
								position, tokenIndex = position8, tokenIndex8
							}
							if !matchDot() {
								goto l6
							}
							goto l5
						l6:
							position, tokenIndex = position6, tokenIndex6
						}
						if buffer[position] != rune('\n') {
							goto l2
						}
						position++
						add(ruleSHEBANG, position4)
					}
					goto l3
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
			l3:
				if !_rules[rule_]() {
					goto l0
				}
			l9:
				{
					position10, tokenIndex10 := position, tokenIndex
					if !_rules[ruleBlock]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex = position10, tokenIndex10
				}
				{
					position11, tokenIndex11 := position, tokenIndex
					if !matchDot() {
						goto l11
					}
					goto l0
				l11:
					position, tokenIndex = position11, tokenIndex11
				}
				add(ruleFriendscript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 _ <- <(' ' / '\t' / '\r' / '\n')*> */
		func() bool {
			{
				position13 := position
			l14:
				{
					position15, tokenIndex15 := position, tokenIndex
					{
						position16, tokenIndex16 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l17
						}
						position++
						goto l16
					l17:
						position, tokenIndex = position16, tokenIndex16
						if buffer[position] != rune('\t') {
							goto l18
						}
						position++
						goto l16
					l18:
						position, tokenIndex = position16, tokenIndex16
						if buffer[position] != rune('\r') {
							goto l19
						}
						position++
						goto l16
					l19:
						position, tokenIndex = position16, tokenIndex16
						if buffer[position] != rune('\n') {
							goto l15
						}
						position++
					}
				l16:
					goto l14
				l15:
					position, tokenIndex = position15, tokenIndex15
				}
				add(rule_, position13)
			}
			return true
		},
		/* 2 __ <- <(' ' / '\t' / '\r' / '\n')+> */
		func() bool {
			position20, tokenIndex20 := position, tokenIndex
			{
				position21 := position
				{
					position24, tokenIndex24 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l25
					}
					position++
					goto l24
				l25:
					position, tokenIndex = position24, tokenIndex24
					if buffer[position] != rune('\t') {
						goto l26
					}
					position++
					goto l24
				l26:
					position, tokenIndex = position24, tokenIndex24
					if buffer[position] != rune('\r') {
						goto l27
					}
					position++
					goto l24
				l27:
					position, tokenIndex = position24, tokenIndex24
					if buffer[position] != rune('\n') {
						goto l20
					}
					position++
				}
			l24:
			l22:
				{
					position23, tokenIndex23 := position, tokenIndex
					{
						position28, tokenIndex28 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l29
						}
						position++
						goto l28
					l29:
						position, tokenIndex = position28, tokenIndex28
						if buffer[position] != rune('\t') {
							goto l30
						}
						position++
						goto l28
					l30:
						position, tokenIndex = position28, tokenIndex28
						if buffer[position] != rune('\r') {
							goto l31
						}
						position++
						goto l28
					l31:
						position, tokenIndex = position28, tokenIndex28
						if buffer[position] != rune('\n') {
							goto l23
						}
						position++
					}
				l28:
					goto l22
				l23:
					position, tokenIndex = position23, tokenIndex23
				}
				add(rule__, position21)
			}
			return true
		l20:
			position, tokenIndex = position20, tokenIndex20
			return false
		},
		/* 3 ASSIGN <- <(_ ('-' '>') _)> */
		nil,
		/* 4 TRIQUOT <- <(_ ('"' '"' '"') _)> */
		func() bool {
			position33, tokenIndex33 := position, tokenIndex
			{
				position34 := position
				if !_rules[rule_]() {
					goto l33
				}
				if buffer[position] != rune('"') {
					goto l33
				}
				position++
				if buffer[position] != rune('"') {
					goto l33
				}
				position++
				if buffer[position] != rune('"') {
					goto l33
				}
				position++
				if !_rules[rule_]() {
					goto l33
				}
				add(ruleTRIQUOT, position34)
			}
			return true
		l33:
			position, tokenIndex = position33, tokenIndex33
			return false
		},
		/* 5 BREAK <- <(_ ('b' 'r' 'e' 'a' 'k') _)> */
		nil,
		/* 6 CLOSE <- <(_ '}' _)> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				if !_rules[rule_]() {
					goto l36
				}
				if buffer[position] != rune('}') {
					goto l36
				}
				position++
				if !_rules[rule_]() {
					goto l36
				}
				add(ruleCLOSE, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 7 COLON <- <(_ ':' _)> */
		nil,
		/* 8 COMMA <- <(_ ',' _)> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				if !_rules[rule_]() {
					goto l39
				}
				if buffer[position] != rune(',') {
					goto l39
				}
				position++
				if !_rules[rule_]() {
					goto l39
				}
				add(ruleCOMMA, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 9 COMMENT <- <(_ '#' (!'\n' .)*)> */
		nil,
		/* 10 CONT <- <(_ ('c' 'o' 'n' 't' 'i' 'n' 'u' 'e') _)> */
		nil,
		/* 11 COUNT <- <(_ ('c' 'o' 'u' 'n' 't') _)> */
		nil,
		/* 12 DECLARE <- <(_ ('d' 'e' 'c' 'l' 'a' 'r' 'e') __)> */
		nil,
		/* 13 DOT <- <'.'> */
		nil,
		/* 14 ELSE <- <(_ ('e' 'l' 's' 'e') _)> */
		func() bool {
			position46, tokenIndex46 := position, tokenIndex
			{
				position47 := position
				if !_rules[rule_]() {
					goto l46
				}
				if buffer[position] != rune('e') {
					goto l46
				}
				position++
				if buffer[position] != rune('l') {
					goto l46
				}
				position++
				if buffer[position] != rune('s') {
					goto l46
				}
				position++
				if buffer[position] != rune('e') {
					goto l46
				}
				position++
				if !_rules[rule_]() {
					goto l46
				}
				add(ruleELSE, position47)
			}
			return true
		l46:
			position, tokenIndex = position46, tokenIndex46
			return false
		},
		/* 15 IF <- <(_ ('i' 'f') _)> */
		nil,
		/* 16 IN <- <(__ ('i' 'n') __)> */
		nil,
		/* 17 INCLUDE <- <(_ ('i' 'n' 'c' 'l' 'u' 'd' 'e') __)> */
		nil,
		/* 18 LOOP <- <(_ ('l' 'o' 'o' 'p') _)> */
		nil,
		/* 19 NOOP <- <SEMI> */
		nil,
		/* 20 NOT <- <(_ ('n' 'o' 't') __)> */
		nil,
		/* 21 OPEN <- <(_ '{' _)> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				if !_rules[rule_]() {
					goto l54
				}
				if buffer[position] != rune('{') {
					goto l54
				}
				position++
				if !_rules[rule_]() {
					goto l54
				}
				add(ruleOPEN, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 22 SCOPE <- <(':' ':')> */
		nil,
		/* 23 SEMI <- <(_ ';' _)> */
		func() bool {
			position57, tokenIndex57 := position, tokenIndex
			{
				position58 := position
				if !_rules[rule_]() {
					goto l57
				}
				if buffer[position] != rune(';') {
					goto l57
				}
				position++
				if !_rules[rule_]() {
					goto l57
				}
				add(ruleSEMI, position58)
			}
			return true
		l57:
			position, tokenIndex = position57, tokenIndex57
			return false
		},
		/* 24 SHEBANG <- <('#' '!' (!'\n' .)+ '\n')> */
		nil,
		/* 25 SKIPVAR <- <(_ '_' _)> */
		nil,
		/* 26 UNSET <- <(_ ('u' 'n' 's' 'e' 't') __)> */
		nil,
		/* 27 ScalarType <- <(Boolean / Float / Integer / String / NullValue)> */
		nil,
		/* 28 Identifier <- <(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / ([0-9] / [0-9]) / '_')*)> */
		func() bool {
			position63, tokenIndex63 := position, tokenIndex
			{
				position64 := position
				{
					position65, tokenIndex65 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex = position65, tokenIndex65
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l67
					}
					position++
					goto l65
				l67:
					position, tokenIndex = position65, tokenIndex65
					if buffer[position] != rune('_') {
						goto l63
					}
					position++
				}
			l65:
			l68:
				{
					position69, tokenIndex69 := position, tokenIndex
					{
						position70, tokenIndex70 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex = position70, tokenIndex70
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l72
						}
						position++
						goto l70
					l72:
						position, tokenIndex = position70, tokenIndex70
						{
							position74, tokenIndex74 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l75
							}
							position++
							goto l74
						l75:
							position, tokenIndex = position74, tokenIndex74
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l73
							}
							position++
						}
					l74:
						goto l70
					l73:
						position, tokenIndex = position70, tokenIndex70
						if buffer[position] != rune('_') {
							goto l69
						}
						position++
					}
				l70:
					goto l68
				l69:
					position, tokenIndex = position69, tokenIndex69
				}
				add(ruleIdentifier, position64)
			}
			return true
		l63:
			position, tokenIndex = position63, tokenIndex63
			return false
		},
		/* 29 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 30 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 31 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position78, tokenIndex78 := position, tokenIndex
			{
				position79 := position
				{
					position80, tokenIndex80 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l80
					}
					position++
					goto l81
				l80:
					position, tokenIndex = position80, tokenIndex80
				}
			l81:
				if !_rules[rulePositiveInteger]() {
					goto l78
				}
				add(ruleInteger, position79)
			}
			return true
		l78:
			position, tokenIndex = position78, tokenIndex78
			return false
		},
		/* 32 PositiveInteger <- <[0-9]+> */
		func() bool {
			position82, tokenIndex82 := position, tokenIndex
			{
				position83 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l82
				}
				position++
			l84:
				{
					position85, tokenIndex85 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex = position85, tokenIndex85
				}
				add(rulePositiveInteger, position83)
			}
			return true
		l82:
			position, tokenIndex = position82, tokenIndex82
			return false
		},
		/* 33 String <- <(Triquote / StringLiteral / StringInterpolated)> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				{
					position88, tokenIndex88 := position, tokenIndex
					{
						position90 := position
						if !_rules[ruleTRIQUOT]() {
							goto l89
						}
						{
							position91 := position
						l92:
							{
								position93, tokenIndex93 := position, tokenIndex
								{
									position94, tokenIndex94 := position, tokenIndex
									if !_rules[ruleTRIQUOT]() {
										goto l94
									}
									goto l93
								l94:
									position, tokenIndex = position94, tokenIndex94
								}
								if !matchDot() {
									goto l93
								}
								goto l92
							l93:
								position, tokenIndex = position93, tokenIndex93
							}
							add(ruleTriquoteBody, position91)
						}
						if !_rules[ruleTRIQUOT]() {
							goto l89
						}
						add(ruleTriquote, position90)
					}
					goto l88
				l89:
					position, tokenIndex = position88, tokenIndex88
					if !_rules[ruleStringLiteral]() {
						goto l95
					}
					goto l88
				l95:
					position, tokenIndex = position88, tokenIndex88
					if !_rules[ruleStringInterpolated]() {
						goto l86
					}
				}
			l88:
				add(ruleString, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 34 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		func() bool {
			position96, tokenIndex96 := position, tokenIndex
			{
				position97 := position
				if buffer[position] != rune('\'') {
					goto l96
				}
				position++
			l98:
				{
					position99, tokenIndex99 := position, tokenIndex
					{
						position100, tokenIndex100 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l100
						}
						position++
						goto l99
					l100:
						position, tokenIndex = position100, tokenIndex100
					}
					if !matchDot() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex = position99, tokenIndex99
				}
				if buffer[position] != rune('\'') {
					goto l96
				}
				position++
				add(ruleStringLiteral, position97)
			}
			return true
		l96:
			position, tokenIndex = position96, tokenIndex96
			return false
		},
		/* 35 StringInterpolated <- <('"' (!'"' .)* '"')> */
		func() bool {
			position101, tokenIndex101 := position, tokenIndex
			{
				position102 := position
				if buffer[position] != rune('"') {
					goto l101
				}
				position++
			l103:
				{
					position104, tokenIndex104 := position, tokenIndex
					{
						position105, tokenIndex105 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l105
						}
						position++
						goto l104
					l105:
						position, tokenIndex = position105, tokenIndex105
					}
					if !matchDot() {
						goto l104
					}
					goto l103
				l104:
					position, tokenIndex = position104, tokenIndex104
				}
				if buffer[position] != rune('"') {
					goto l101
				}
				position++
				add(ruleStringInterpolated, position102)
			}
			return true
		l101:
			position, tokenIndex = position101, tokenIndex101
			return false
		},
		/* 36 Triquote <- <(TRIQUOT TriquoteBody TRIQUOT)> */
		nil,
		/* 37 TriquoteBody <- <(!TRIQUOT .)*> */
		nil,
		/* 38 NullValue <- <('n' 'u' 'l' 'l')> */
		nil,
		/* 39 Object <- <(OPEN (_ KeyValuePair _)* CLOSE)> */
		func() bool {
			position109, tokenIndex109 := position, tokenIndex
			{
				position110 := position
				if !_rules[ruleOPEN]() {
					goto l109
				}
			l111:
				{
					position112, tokenIndex112 := position, tokenIndex
					if !_rules[rule_]() {
						goto l112
					}
					{
						position113 := position
						{
							position114 := position
							{
								position115, tokenIndex115 := position, tokenIndex
								if !_rules[ruleIdentifier]() {
									goto l116
								}
								goto l115
							l116:
								position, tokenIndex = position115, tokenIndex115
								if !_rules[ruleStringLiteral]() {
									goto l117
								}
								goto l115
							l117:
								position, tokenIndex = position115, tokenIndex115
								if !_rules[ruleStringInterpolated]() {
									goto l112
								}
							}
						l115:
							add(ruleKey, position114)
						}
						{
							position118 := position
							if !_rules[rule_]() {
								goto l112
							}
							if buffer[position] != rune(':') {
								goto l112
							}
							position++
							if !_rules[rule_]() {
								goto l112
							}
							add(ruleCOLON, position118)
						}
						{
							position119 := position
							{
								position120, tokenIndex120 := position, tokenIndex
								if !_rules[ruleArray]() {
									goto l121
								}
								goto l120
							l121:
								position, tokenIndex = position120, tokenIndex120
								if !_rules[ruleObject]() {
									goto l122
								}
								goto l120
							l122:
								position, tokenIndex = position120, tokenIndex120
								if !_rules[ruleExpression]() {
									goto l112
								}
							}
						l120:
							add(ruleKValue, position119)
						}
						{
							position123, tokenIndex123 := position, tokenIndex
							if !_rules[ruleCOMMA]() {
								goto l123
							}
							goto l124
						l123:
							position, tokenIndex = position123, tokenIndex123
						}
					l124:
						add(ruleKeyValuePair, position113)
					}
					if !_rules[rule_]() {
						goto l112
					}
					goto l111
				l112:
					position, tokenIndex = position112, tokenIndex112
				}
				if !_rules[ruleCLOSE]() {
					goto l109
				}
				add(ruleObject, position110)
			}
			return true
		l109:
			position, tokenIndex = position109, tokenIndex109
			return false
		},
		/* 40 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position125, tokenIndex125 := position, tokenIndex
			{
				position126 := position
				if buffer[position] != rune('[') {
					goto l125
				}
				position++
				if !_rules[rule_]() {
					goto l125
				}
				if !_rules[ruleExpressionSequence]() {
					goto l125
				}
				{
					position127, tokenIndex127 := position, tokenIndex
					if !_rules[ruleCOMMA]() {
						goto l127
					}
					goto l128
				l127:
					position, tokenIndex = position127, tokenIndex127
				}
			l128:
				if buffer[position] != rune(']') {
					goto l125
				}
				position++
				add(ruleArray, position126)
			}
			return true
		l125:
			position, tokenIndex = position125, tokenIndex125
			return false
		},
		/* 41 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position129, tokenIndex129 := position, tokenIndex
			{
				position130 := position
				if buffer[position] != rune('/') {
					goto l129
				}
				position++
				{
					position133, tokenIndex133 := position, tokenIndex
					if buffer[position] != rune('/') {
						goto l133
					}
					position++
					goto l129
				l133:
					position, tokenIndex = position133, tokenIndex133
				}
				if !matchDot() {
					goto l129
				}
			l131:
				{
					position132, tokenIndex132 := position, tokenIndex
					{
						position134, tokenIndex134 := position, tokenIndex
						if buffer[position] != rune('/') {
							goto l134
						}
						position++
						goto l132
					l134:
						position, tokenIndex = position134, tokenIndex134
					}
					if !matchDot() {
						goto l132
					}
					goto l131
				l132:
					position, tokenIndex = position132, tokenIndex132
				}
				if buffer[position] != rune('/') {
					goto l129
				}
				position++
			l135:
				{
					position136, tokenIndex136 := position, tokenIndex
					{
						position137, tokenIndex137 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l138
						}
						position++
						goto l137
					l138:
						position, tokenIndex = position137, tokenIndex137
						if buffer[position] != rune('l') {
							goto l139
						}
						position++
						goto l137
					l139:
						position, tokenIndex = position137, tokenIndex137
						if buffer[position] != rune('m') {
							goto l140
						}
						position++
						goto l137
					l140:
						position, tokenIndex = position137, tokenIndex137
						if buffer[position] != rune('s') {
							goto l141
						}
						position++
						goto l137
					l141:
						position, tokenIndex = position137, tokenIndex137
						if buffer[position] != rune('u') {
							goto l136
						}
						position++
					}
				l137:
					goto l135
				l136:
					position, tokenIndex = position136, tokenIndex136
				}
				add(ruleRegularExpression, position130)
			}
			return true
		l129:
			position, tokenIndex = position129, tokenIndex129
			return false
		},
		/* 42 KeyValuePair <- <(Key COLON KValue COMMA?)> */
		nil,
		/* 43 Key <- <(Identifier / StringLiteral / StringInterpolated)> */
		nil,
		/* 44 KValue <- <(Array / Object / Expression)> */
		nil,
		/* 45 Type <- <(Array / Object / RegularExpression / ScalarType)> */
		func() bool {
			position145, tokenIndex145 := position, tokenIndex
			{
				position146 := position
				{
					position147, tokenIndex147 := position, tokenIndex
					if !_rules[ruleArray]() {
						goto l148
					}
					goto l147
				l148:
					position, tokenIndex = position147, tokenIndex147
					if !_rules[ruleObject]() {
						goto l149
					}
					goto l147
				l149:
					position, tokenIndex = position147, tokenIndex147
					if !_rules[ruleRegularExpression]() {
						goto l150
					}
					goto l147
				l150:
					position, tokenIndex = position147, tokenIndex147
					{
						position151 := position
						{
							position152, tokenIndex152 := position, tokenIndex
							{
								position154 := position
								{
									position155, tokenIndex155 := position, tokenIndex
									if buffer[position] != rune('t') {
										goto l156
									}
									position++
									if buffer[position] != rune('r') {
										goto l156
									}
									position++
									if buffer[position] != rune('u') {
										goto l156
									}
									position++
									if buffer[position] != rune('e') {
										goto l156
									}
									position++
									goto l155
								l156:
									position, tokenIndex = position155, tokenIndex155
									if buffer[position] != rune('f') {
										goto l153
									}
									position++
									if buffer[position] != rune('a') {
										goto l153
									}
									position++
									if buffer[position] != rune('l') {
										goto l153
									}
									position++
									if buffer[position] != rune('s') {
										goto l153
									}
									position++
									if buffer[position] != rune('e') {
										goto l153
									}
									position++
								}
							l155:
								add(ruleBoolean, position154)
							}
							goto l152
						l153:
							position, tokenIndex = position152, tokenIndex152
							{
								position158 := position
								if !_rules[ruleInteger]() {
									goto l157
								}
								{
									position159, tokenIndex159 := position, tokenIndex
									if buffer[position] != rune('.') {
										goto l159
									}
									position++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l159
									}
									position++
								l161:
									{
										position162, tokenIndex162 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l162
										}
										position++
										goto l161
									l162:
										position, tokenIndex = position162, tokenIndex162
									}
									goto l160
								l159:
									position, tokenIndex = position159, tokenIndex159
								}
							l160:
								add(ruleFloat, position158)
							}
							goto l152
						l157:
							position, tokenIndex = position152, tokenIndex152
							if !_rules[ruleInteger]() {
								goto l163
							}
							goto l152
						l163:
							position, tokenIndex = position152, tokenIndex152
							if !_rules[ruleString]() {
								goto l164
							}
							goto l152
						l164:
							position, tokenIndex = position152, tokenIndex152
							{
								position165 := position
								if buffer[position] != rune('n') {
									goto l145
								}
								position++
								if buffer[position] != rune('u') {
									goto l145
								}
								position++
								if buffer[position] != rune('l') {
									goto l145
								}
								position++
								if buffer[position] != rune('l') {
									goto l145
								}
								position++
								add(ruleNullValue, position165)
							}
						}
					l152:
						add(ruleScalarType, position151)
					}
				}
			l147:
				add(ruleType, position146)
			}
			return true
		l145:
			position, tokenIndex = position145, tokenIndex145
			return false
		},
		/* 46 Exponentiate <- <(_ ('*' '*') _)> */
		nil,
		/* 47 Multiply <- <(_ '*' _)> */
		nil,
		/* 48 Divide <- <(_ '/' _)> */
		nil,
		/* 49 Modulus <- <(_ '%' _)> */
		nil,
		/* 50 Add <- <(_ '+' _)> */
		nil,
		/* 51 Subtract <- <(_ '-' _)> */
		nil,
		/* 52 BitwiseAnd <- <(_ '&' _)> */
		nil,
		/* 53 BitwiseOr <- <(_ '|' _)> */
		nil,
		/* 54 BitwiseNot <- <(_ '~' _)> */
		nil,
		/* 55 BitwiseXor <- <(_ '^' _)> */
		nil,
		/* 56 MatchOperator <- <(Match / Unmatch)> */
		nil,
		/* 57 Unmatch <- <(_ ('!' '~') _)> */
		nil,
		/* 58 Match <- <(_ ('=' '~') _)> */
		nil,
		/* 59 Operator <- <(_ (Exponentiate / Multiply / Divide / Modulus / Add / Subtract / BitwiseAnd / BitwiseOr / BitwiseNot / BitwiseXor) _)> */
		nil,
		/* 60 AssignmentOperator <- <(_ (AssignEq / StarEq / DivEq / PlusEq / MinusEq / AndEq / OrEq / Append) _)> */
		nil,
		/* 61 AssignEq <- <(_ '=' _)> */
		nil,
		/* 62 StarEq <- <(_ ('*' '=') _)> */
		nil,
		/* 63 DivEq <- <(_ ('/' '=') _)> */
		nil,
		/* 64 PlusEq <- <(_ ('+' '=') _)> */
		nil,
		/* 65 MinusEq <- <(_ ('-' '=') _)> */
		nil,
		/* 66 AndEq <- <(_ ('&' '=') _)> */
		nil,
		/* 67 OrEq <- <(_ ('|' '=') _)> */
		nil,
		/* 68 Append <- <(_ ('<' '<') _)> */
		nil,
		/* 69 ComparisonOperator <- <(_ (Equality / NonEquality / GreaterEqual / LessEqual / GreaterThan / LessThan / Membership / NonMembership) _)> */
		nil,
		/* 70 Equality <- <(_ ('=' '=') _)> */
		nil,
		/* 71 NonEquality <- <(_ ('!' '=') _)> */
		nil,
		/* 72 GreaterThan <- <(_ '>' _)> */
		nil,
		/* 73 GreaterEqual <- <(_ ('>' '=') _)> */
		nil,
		/* 74 LessEqual <- <(_ ('<' '=') _)> */
		nil,
		/* 75 LessThan <- <(_ '<' _)> */
		nil,
		/* 76 Membership <- <(_ ('i' 'n') _)> */
		nil,
		/* 77 NonMembership <- <(_ ('n' 'o' 't') __ ('i' 'n') _)> */
		nil,
		/* 78 Variable <- <(('$' VariableNameSequence) / SKIPVAR)> */
		func() bool {
			position198, tokenIndex198 := position, tokenIndex
			{
				position199 := position
				{
					position200, tokenIndex200 := position, tokenIndex
					if buffer[position] != rune('$') {
						goto l201
					}
					position++
					{
						position202 := position
					l203:
						{
							position204, tokenIndex204 := position, tokenIndex
							if !_rules[ruleVariableName]() {
								goto l204
							}
							{
								position205 := position
								if buffer[position] != rune('.') {
									goto l204
								}
								position++
								add(ruleDOT, position205)
							}
							goto l203
						l204:
							position, tokenIndex = position204, tokenIndex204
						}
						if !_rules[ruleVariableName]() {
							goto l201
						}
						add(ruleVariableNameSequence, position202)
					}
					goto l200
				l201:
					position, tokenIndex = position200, tokenIndex200
					{
						position206 := position
						if !_rules[rule_]() {
							goto l198
						}
						if buffer[position] != rune('_') {
							goto l198
						}
						position++
						if !_rules[rule_]() {
							goto l198
						}
						add(ruleSKIPVAR, position206)
					}
				}
			l200:
				add(ruleVariable, position199)
			}
			return true
		l198:
			position, tokenIndex = position198, tokenIndex198
			return false
		},
		/* 79 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 80 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position208, tokenIndex208 := position, tokenIndex
			{
				position209 := position
				if !_rules[ruleIdentifier]() {
					goto l208
				}
				{
					position210, tokenIndex210 := position, tokenIndex
					if buffer[position] != rune('[') {
						goto l210
					}
					position++
					if !_rules[rule_]() {
						goto l210
					}
					{
						position212 := position
						if !_rules[ruleExpression]() {
							goto l210
						}
						add(ruleVariableIndex, position212)
					}
					if !_rules[rule_]() {
						goto l210
					}
					if buffer[position] != rune(']') {
						goto l210
					}
					position++
					goto l211
				l210:
					position, tokenIndex = position210, tokenIndex210
				}
			l211:
				add(ruleVariableName, position209)
			}
			return true
		l208:
			position, tokenIndex = position208, tokenIndex208
			return false
		},
		/* 81 VariableIndex <- <Expression> */
		nil,
		/* 82 Block <- <(_ (COMMENT / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position214, tokenIndex214 := position, tokenIndex
			{
				position215 := position
				if !_rules[rule_]() {
					goto l214
				}
				{
					position216, tokenIndex216 := position, tokenIndex
					{
						position218 := position
						if !_rules[rule_]() {
							goto l217
						}
						if buffer[position] != rune('#') {
							goto l217
						}
						position++
					l219:
						{
							position220, tokenIndex220 := position, tokenIndex
							{
								position221, tokenIndex221 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l221
								}
								position++
								goto l220
							l221:
								position, tokenIndex = position221, tokenIndex221
							}
							if !matchDot() {
								goto l220
							}
							goto l219
						l220:
							position, tokenIndex = position220, tokenIndex220
						}
						add(ruleCOMMENT, position218)
					}
					goto l216
				l217:
					position, tokenIndex = position216, tokenIndex216
					{
						position223 := position
						{
							position224, tokenIndex224 := position, tokenIndex
							{
								position226 := position
								{
									position227 := position
									if !_rules[rule_]() {
										goto l225
									}
									if buffer[position] != rune('b') {
										goto l225
									}
									position++
									if buffer[position] != rune('r') {
										goto l225
									}
									position++
									if buffer[position] != rune('e') {
										goto l225
									}
									position++
									if buffer[position] != rune('a') {
										goto l225
									}
									position++
									if buffer[position] != rune('k') {
										goto l225
									}
									position++
									if !_rules[rule_]() {
										goto l225
									}
									add(ruleBREAK, position227)
								}
								{
									position228, tokenIndex228 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l228
									}
									goto l229
								l228:
									position, tokenIndex = position228, tokenIndex228
								}
							l229:
								add(ruleFlowControlBreak, position226)
							}
							goto l224
						l225:
							position, tokenIndex = position224, tokenIndex224
							{
								position230 := position
								{
									position231 := position
									if !_rules[rule_]() {
										goto l222
									}
									if buffer[position] != rune('c') {
										goto l222
									}
									position++
									if buffer[position] != rune('o') {
										goto l222
									}
									position++
									if buffer[position] != rune('n') {
										goto l222
									}
									position++
									if buffer[position] != rune('t') {
										goto l222
									}
									position++
									if buffer[position] != rune('i') {
										goto l222
									}
									position++
									if buffer[position] != rune('n') {
										goto l222
									}
									position++
									if buffer[position] != rune('u') {
										goto l222
									}
									position++
									if buffer[position] != rune('e') {
										goto l222
									}
									position++
									if !_rules[rule_]() {
										goto l222
									}
									add(ruleCONT, position231)
								}
								{
									position232, tokenIndex232 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l232
									}
									goto l233
								l232:
									position, tokenIndex = position232, tokenIndex232
								}
							l233:
								add(ruleFlowControlContinue, position230)
							}
						}
					l224:
						add(ruleFlowControlWord, position223)
					}
					goto l216
				l222:
					position, tokenIndex = position216, tokenIndex216
					{
						position234 := position
						{
							position235, tokenIndex235 := position, tokenIndex
							{
								position237 := position
								if !_rules[ruleSEMI]() {
									goto l236
								}
								add(ruleNOOP, position237)
							}
							goto l235
						l236:
							position, tokenIndex = position235, tokenIndex235
							if !_rules[ruleAssignment]() {
								goto l238
							}
							goto l235
						l238:
							position, tokenIndex = position235, tokenIndex235
							{
								position240 := position
								{
									position241, tokenIndex241 := position, tokenIndex
									{
										position243 := position
										{
											position244 := position
											if !_rules[rule_]() {
												goto l242
											}
											if buffer[position] != rune('u') {
												goto l242
											}
											position++
											if buffer[position] != rune('n') {
												goto l242
											}
											position++
											if buffer[position] != rune('s') {
												goto l242
											}
											position++
											if buffer[position] != rune('e') {
												goto l242
											}
											position++
											if buffer[position] != rune('t') {
												goto l242
											}
											position++
											if !_rules[rule__]() {
												goto l242
											}
											add(ruleUNSET, position244)
										}
										if !_rules[ruleVariableSequence]() {
											goto l242
										}
										add(ruleDirectiveUnset, position243)
									}
									goto l241
								l242:
									position, tokenIndex = position241, tokenIndex241
									{
										position246 := position
										{
											position247 := position
											if !_rules[rule_]() {
												goto l245
											}
											if buffer[position] != rune('i') {
												goto l245
											}
											position++
											if buffer[position] != rune('n') {
												goto l245
											}
											position++
											if buffer[position] != rune('c') {
												goto l245
											}
											position++
											if buffer[position] != rune('l') {
												goto l245
											}
											position++
											if buffer[position] != rune('u') {
												goto l245
											}
											position++
											if buffer[position] != rune('d') {
												goto l245
											}
											position++
											if buffer[position] != rune('e') {
												goto l245
											}
											position++
											if !_rules[rule__]() {
												goto l245
											}
											add(ruleINCLUDE, position247)
										}
										if !_rules[ruleString]() {
											goto l245
										}
										add(ruleDirectiveInclude, position246)
									}
									goto l241
								l245:
									position, tokenIndex = position241, tokenIndex241
									{
										position248 := position
										{
											position249 := position
											if !_rules[rule_]() {
												goto l239
											}
											if buffer[position] != rune('d') {
												goto l239
											}
											position++
											if buffer[position] != rune('e') {
												goto l239
											}
											position++
											if buffer[position] != rune('c') {
												goto l239
											}
											position++
											if buffer[position] != rune('l') {
												goto l239
											}
											position++
											if buffer[position] != rune('a') {
												goto l239
											}
											position++
											if buffer[position] != rune('r') {
												goto l239
											}
											position++
											if buffer[position] != rune('e') {
												goto l239
											}
											position++
											if !_rules[rule__]() {
												goto l239
											}
											add(ruleDECLARE, position249)
										}
										if !_rules[ruleVariableSequence]() {
											goto l239
										}
										add(ruleDirectiveDeclare, position248)
									}
								}
							l241:
								add(ruleDirective, position240)
							}
							goto l235
						l239:
							position, tokenIndex = position235, tokenIndex235
							{
								position251 := position
								if !_rules[ruleIfStanza]() {
									goto l250
								}
							l252:
								{
									position253, tokenIndex253 := position, tokenIndex
									{
										position254 := position
										if !_rules[ruleELSE]() {
											goto l253
										}
										if !_rules[ruleIfStanza]() {
											goto l253
										}
										add(ruleElseIfStanza, position254)
									}
									goto l252
								l253:
									position, tokenIndex = position253, tokenIndex253
								}
								{
									position255, tokenIndex255 := position, tokenIndex
									{
										position257 := position
										if !_rules[ruleELSE]() {
											goto l255
										}
										if !_rules[ruleOPEN]() {
											goto l255
										}
									l258:
										{
											position259, tokenIndex259 := position, tokenIndex
											if !_rules[ruleBlock]() {
												goto l259
											}
											goto l258
										l259:
											position, tokenIndex = position259, tokenIndex259
										}
										if !_rules[ruleCLOSE]() {
											goto l255
										}
										add(ruleElseStanza, position257)
									}
									goto l256
								l255:
									position, tokenIndex = position255, tokenIndex255
								}
							l256:
								add(ruleConditional, position251)
							}
							goto l235
						l250:
							position, tokenIndex = position235, tokenIndex235
							{
								position261 := position
								{
									position262 := position
									if !_rules[rule_]() {
										goto l260
									}
									if buffer[position] != rune('l') {
										goto l260
									}
									position++
									if buffer[position] != rune('o') {
										goto l260
									}
									position++
									if buffer[position] != rune('o') {
										goto l260
									}
									position++
									if buffer[position] != rune('p') {
										goto l260
									}
									position++
									if !_rules[rule_]() {
										goto l260
									}
									add(ruleLOOP, position262)
								}
								{
									position263, tokenIndex263 := position, tokenIndex
									if !_rules[ruleOPEN]() {
										goto l264
									}
								l265:
									{
										position266, tokenIndex266 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l266
										}
										goto l265
									l266:
										position, tokenIndex = position266, tokenIndex266
									}
									if !_rules[ruleCLOSE]() {
										goto l264
									}
									goto l263
								l264:
									position, tokenIndex = position263, tokenIndex263
									{
										position268 := position
										{
											position269 := position
											if !_rules[rule_]() {
												goto l267
											}
											if buffer[position] != rune('c') {
												goto l267
											}
											position++
											if buffer[position] != rune('o') {
												goto l267
											}
											position++
											if buffer[position] != rune('u') {
												goto l267
											}
											position++
											if buffer[position] != rune('n') {
												goto l267
											}
											position++
											if buffer[position] != rune('t') {
												goto l267
											}
											position++
											if !_rules[rule_]() {
												goto l267
											}
											add(ruleCOUNT, position269)
										}
										{
											position270, tokenIndex270 := position, tokenIndex
											if !_rules[ruleInteger]() {
												goto l271
											}
											goto l270
										l271:
											position, tokenIndex = position270, tokenIndex270
											if !_rules[ruleVariable]() {
												goto l267
											}
										}
									l270:
										add(ruleLoopConditionFixedLength, position268)
									}
									if !_rules[ruleOPEN]() {
										goto l267
									}
								l272:
									{
										position273, tokenIndex273 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l273
										}
										goto l272
									l273:
										position, tokenIndex = position273, tokenIndex273
									}
									if !_rules[ruleCLOSE]() {
										goto l267
									}
									goto l263
								l267:
									position, tokenIndex = position263, tokenIndex263
									{
										position275 := position
										{
											position276 := position
											if !_rules[ruleVariableSequence]() {
												goto l274
											}
											add(ruleLoopIterableLHS, position276)
										}
										{
											position277 := position
											if !_rules[rule__]() {
												goto l274
											}
											if buffer[position] != rune('i') {
												goto l274
											}
											position++
											if buffer[position] != rune('n') {
												goto l274
											}
											position++
											if !_rules[rule__]() {
												goto l274
											}
											add(ruleIN, position277)
										}
										{
											position278 := position
											{
												position279, tokenIndex279 := position, tokenIndex
												if !_rules[ruleCommand]() {
													goto l280
												}
												goto l279
											l280:
												position, tokenIndex = position279, tokenIndex279
												if !_rules[ruleVariable]() {
													goto l274
												}
											}
										l279:
											add(ruleLoopIterableRHS, position278)
										}
										add(ruleLoopConditionIterable, position275)
									}
									if !_rules[ruleOPEN]() {
										goto l274
									}
								l281:
									{
										position282, tokenIndex282 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l282
										}
										goto l281
									l282:
										position, tokenIndex = position282, tokenIndex282
									}
									if !_rules[ruleCLOSE]() {
										goto l274
									}
									goto l263
								l274:
									position, tokenIndex = position263, tokenIndex263
									{
										position284 := position
										if !_rules[ruleCommand]() {
											goto l283
										}
										if !_rules[ruleSEMI]() {
											goto l283
										}
										if !_rules[ruleConditionalExpression]() {
											goto l283
										}
										if !_rules[ruleSEMI]() {
											goto l283
										}
										if !_rules[ruleCommand]() {
											goto l283
										}
										add(ruleLoopConditionBounded, position284)
									}
									if !_rules[ruleOPEN]() {
										goto l283
									}
								l285:
									{
										position286, tokenIndex286 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l286
										}
										goto l285
									l286:
										position, tokenIndex = position286, tokenIndex286
									}
									if !_rules[ruleCLOSE]() {
										goto l283
									}
									goto l263
								l283:
									position, tokenIndex = position263, tokenIndex263
									{
										position287 := position
										if !_rules[ruleConditionalExpression]() {
											goto l260
										}
										add(ruleLoopConditionTruthy, position287)
									}
									if !_rules[ruleOPEN]() {
										goto l260
									}
								l288:
									{
										position289, tokenIndex289 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l289
										}
										goto l288
									l289:
										position, tokenIndex = position289, tokenIndex289
									}
									if !_rules[ruleCLOSE]() {
										goto l260
									}
								}
							l263:
								add(ruleLoop, position261)
							}
							goto l235
						l260:
							position, tokenIndex = position235, tokenIndex235
							if !_rules[ruleCommand]() {
								goto l214
							}
						}
					l235:
						add(ruleStatementBlock, position234)
					}
				}
			l216:
				{
					position290, tokenIndex290 := position, tokenIndex
					if !_rules[ruleSEMI]() {
						goto l290
					}
					goto l291
				l290:
					position, tokenIndex = position290, tokenIndex290
				}
			l291:
				if !_rules[rule_]() {
					goto l214
				}
				add(ruleBlock, position215)
			}
			return true
		l214:
			position, tokenIndex = position214, tokenIndex214
			return false
		},
		/* 83 FlowControlWord <- <(FlowControlBreak / FlowControlContinue)> */
		nil,
		/* 84 FlowControlBreak <- <(BREAK PositiveInteger?)> */
		nil,
		/* 85 FlowControlContinue <- <(CONT PositiveInteger?)> */
		nil,
		/* 86 StatementBlock <- <(NOOP / Assignment / Directive / Conditional / Loop / Command)> */
		nil,
		/* 87 Assignment <- <(AssignmentLHS AssignmentOperator AssignmentRHS)> */
		func() bool {
			position296, tokenIndex296 := position, tokenIndex
			{
				position297 := position
				{
					position298 := position
					if !_rules[ruleVariableSequence]() {
						goto l296
					}
					add(ruleAssignmentLHS, position298)
				}
				{
					position299 := position
					if !_rules[rule_]() {
						goto l296
					}
					{
						position300, tokenIndex300 := position, tokenIndex
						{
							position302 := position
							if !_rules[rule_]() {
								goto l301
							}
							if buffer[position] != rune('=') {
								goto l301
							}
							position++
							if !_rules[rule_]() {
								goto l301
							}
							add(ruleAssignEq, position302)
						}
						goto l300
					l301:
						position, tokenIndex = position300, tokenIndex300
						{
							position304 := position
							if !_rules[rule_]() {
								goto l303
							}
							if buffer[position] != rune('*') {
								goto l303
							}
							position++
							if buffer[position] != rune('=') {
								goto l303
							}
							position++
							if !_rules[rule_]() {
								goto l303
							}
							add(ruleStarEq, position304)
						}
						goto l300
					l303:
						position, tokenIndex = position300, tokenIndex300
						{
							position306 := position
							if !_rules[rule_]() {
								goto l305
							}
							if buffer[position] != rune('/') {
								goto l305
							}
							position++
							if buffer[position] != rune('=') {
								goto l305
							}
							position++
							if !_rules[rule_]() {
								goto l305
							}
							add(ruleDivEq, position306)
						}
						goto l300
					l305:
						position, tokenIndex = position300, tokenIndex300
						{
							position308 := position
							if !_rules[rule_]() {
								goto l307
							}
							if buffer[position] != rune('+') {
								goto l307
							}
							position++
							if buffer[position] != rune('=') {
								goto l307
							}
							position++
							if !_rules[rule_]() {
								goto l307
							}
							add(rulePlusEq, position308)
						}
						goto l300
					l307:
						position, tokenIndex = position300, tokenIndex300
						{
							position310 := position
							if !_rules[rule_]() {
								goto l309
							}
							if buffer[position] != rune('-') {
								goto l309
							}
							position++
							if buffer[position] != rune('=') {
								goto l309
							}
							position++
							if !_rules[rule_]() {
								goto l309
							}
							add(ruleMinusEq, position310)
						}
						goto l300
					l309:
						position, tokenIndex = position300, tokenIndex300
						{
							position312 := position
							if !_rules[rule_]() {
								goto l311
							}
							if buffer[position] != rune('&') {
								goto l311
							}
							position++
							if buffer[position] != rune('=') {
								goto l311
							}
							position++
							if !_rules[rule_]() {
								goto l311
							}
							add(ruleAndEq, position312)
						}
						goto l300
					l311:
						position, tokenIndex = position300, tokenIndex300
						{
							position314 := position
							if !_rules[rule_]() {
								goto l313
							}
							if buffer[position] != rune('|') {
								goto l313
							}
							position++
							if buffer[position] != rune('=') {
								goto l313
							}
							position++
							if !_rules[rule_]() {
								goto l313
							}
							add(ruleOrEq, position314)
						}
						goto l300
					l313:
						position, tokenIndex = position300, tokenIndex300
						{
							position315 := position
							if !_rules[rule_]() {
								goto l296
							}
							if buffer[position] != rune('<') {
								goto l296
							}
							position++
							if buffer[position] != rune('<') {
								goto l296
							}
							position++
							if !_rules[rule_]() {
								goto l296
							}
							add(ruleAppend, position315)
						}
					}
				l300:
					if !_rules[rule_]() {
						goto l296
					}
					add(ruleAssignmentOperator, position299)
				}
				{
					position316 := position
					if !_rules[ruleExpressionSequence]() {
						goto l296
					}
					add(ruleAssignmentRHS, position316)
				}
				add(ruleAssignment, position297)
			}
			return true
		l296:
			position, tokenIndex = position296, tokenIndex296
			return false
		},
		/* 88 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 89 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 90 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position319, tokenIndex319 := position, tokenIndex
			{
				position320 := position
			l321:
				{
					position322, tokenIndex322 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l322
					}
					if !_rules[ruleCOMMA]() {
						goto l322
					}
					goto l321
				l322:
					position, tokenIndex = position322, tokenIndex322
				}
				if !_rules[ruleVariable]() {
					goto l319
				}
				add(ruleVariableSequence, position320)
			}
			return true
		l319:
			position, tokenIndex = position319, tokenIndex319
			return false
		},
		/* 91 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position323, tokenIndex323 := position, tokenIndex
			{
				position324 := position
			l325:
				{
					position326, tokenIndex326 := position, tokenIndex
					if !_rules[ruleExpression]() {
						goto l326
					}
					if !_rules[ruleCOMMA]() {
						goto l326
					}
					goto l325
				l326:
					position, tokenIndex = position326, tokenIndex326
				}
				if !_rules[ruleExpression]() {
					goto l323
				}
				add(ruleExpressionSequence, position324)
			}
			return true
		l323:
			position, tokenIndex = position323, tokenIndex323
			return false
		},
		/* 92 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position327, tokenIndex327 := position, tokenIndex
			{
				position328 := position
				if !_rules[rule_]() {
					goto l327
				}
				{
					position329 := position
					{
						position330 := position
						{
							position331, tokenIndex331 := position, tokenIndex
							if !_rules[ruleType]() {
								goto l332
							}
							goto l331
						l332:
							position, tokenIndex = position331, tokenIndex331
							if !_rules[ruleVariable]() {
								goto l327
							}
						}
					l331:
						add(ruleValueYielding, position330)
					}
					add(ruleExpressionLHS, position329)
				}
				{
					position333, tokenIndex333 := position, tokenIndex
					{
						position335 := position
						{
							position336 := position
							if !_rules[rule_]() {
								goto l333
							}
							{
								position337, tokenIndex337 := position, tokenIndex
								{
									position339 := position
									if !_rules[rule_]() {
										goto l338
									}
									if buffer[position] != rune('*') {
										goto l338
									}
									position++
									if buffer[position] != rune('*') {
										goto l338
									}
									position++
									if !_rules[rule_]() {
										goto l338
									}
									add(ruleExponentiate, position339)
								}
								goto l337
							l338:
								position, tokenIndex = position337, tokenIndex337
								{
									position341 := position
									if !_rules[rule_]() {
										goto l340
									}
									if buffer[position] != rune('*') {
										goto l340
									}
									position++
									if !_rules[rule_]() {
										goto l340
									}
									add(ruleMultiply, position341)
								}
								goto l337
							l340:
								position, tokenIndex = position337, tokenIndex337
								{
									position343 := position
									if !_rules[rule_]() {
										goto l342
									}
									if buffer[position] != rune('/') {
										goto l342
									}
									position++
									if !_rules[rule_]() {
										goto l342
									}
									add(ruleDivide, position343)
								}
								goto l337
							l342:
								position, tokenIndex = position337, tokenIndex337
								{
									position345 := position
									if !_rules[rule_]() {
										goto l344
									}
									if buffer[position] != rune('%') {
										goto l344
									}
									position++
									if !_rules[rule_]() {
										goto l344
									}
									add(ruleModulus, position345)
								}
								goto l337
							l344:
								position, tokenIndex = position337, tokenIndex337
								{
									position347 := position
									if !_rules[rule_]() {
										goto l346
									}
									if buffer[position] != rune('+') {
										goto l346
									}
									position++
									if !_rules[rule_]() {
										goto l346
									}
									add(ruleAdd, position347)
								}
								goto l337
							l346:
								position, tokenIndex = position337, tokenIndex337
								{
									position349 := position
									if !_rules[rule_]() {
										goto l348
									}
									if buffer[position] != rune('-') {
										goto l348
									}
									position++
									if !_rules[rule_]() {
										goto l348
									}
									add(ruleSubtract, position349)
								}
								goto l337
							l348:
								position, tokenIndex = position337, tokenIndex337
								{
									position351 := position
									if !_rules[rule_]() {
										goto l350
									}
									if buffer[position] != rune('&') {
										goto l350
									}
									position++
									if !_rules[rule_]() {
										goto l350
									}
									add(ruleBitwiseAnd, position351)
								}
								goto l337
							l350:
								position, tokenIndex = position337, tokenIndex337
								{
									position353 := position
									if !_rules[rule_]() {
										goto l352
									}
									if buffer[position] != rune('|') {
										goto l352
									}
									position++
									if !_rules[rule_]() {
										goto l352
									}
									add(ruleBitwiseOr, position353)
								}
								goto l337
							l352:
								position, tokenIndex = position337, tokenIndex337
								{
									position355 := position
									if !_rules[rule_]() {
										goto l354
									}
									if buffer[position] != rune('~') {
										goto l354
									}
									position++
									if !_rules[rule_]() {
										goto l354
									}
									add(ruleBitwiseNot, position355)
								}
								goto l337
							l354:
								position, tokenIndex = position337, tokenIndex337
								{
									position356 := position
									if !_rules[rule_]() {
										goto l333
									}
									if buffer[position] != rune('^') {
										goto l333
									}
									position++
									if !_rules[rule_]() {
										goto l333
									}
									add(ruleBitwiseXor, position356)
								}
							}
						l337:
							if !_rules[rule_]() {
								goto l333
							}
							add(ruleOperator, position336)
						}
						if !_rules[ruleExpression]() {
							goto l333
						}
						add(ruleExpressionRHS, position335)
					}
					goto l334
				l333:
					position, tokenIndex = position333, tokenIndex333
				}
			l334:
				if !_rules[rule_]() {
					goto l327
				}
				add(ruleExpression, position328)
			}
			return true
		l327:
			position, tokenIndex = position327, tokenIndex327
			return false
		},
		/* 93 ExpressionLHS <- <ValueYielding> */
		nil,
		/* 94 ExpressionRHS <- <(Operator Expression)> */
		nil,
		/* 95 ValueYielding <- <(Type / Variable)> */
		nil,
		/* 96 Directive <- <(DirectiveUnset / DirectiveInclude / DirectiveDeclare)> */
		nil,
		/* 97 DirectiveUnset <- <(UNSET VariableSequence)> */
		nil,
		/* 98 DirectiveInclude <- <(INCLUDE String)> */
		nil,
		/* 99 DirectiveDeclare <- <(DECLARE VariableSequence)> */
		nil,
		/* 100 Command <- <(_ CommandName (__ ((CommandFirstArg __ CommandSecondArg) / CommandFirstArg / CommandSecondArg))? (_ CommandResultAssignment)?)> */
		func() bool {
			position364, tokenIndex364 := position, tokenIndex
			{
				position365 := position
				if !_rules[rule_]() {
					goto l364
				}
				{
					position366 := position
					{
						position367, tokenIndex367 := position, tokenIndex
						if !_rules[ruleIdentifier]() {
							goto l367
						}
						{
							position369 := position
							if buffer[position] != rune(':') {
								goto l367
							}
							position++
							if buffer[position] != rune(':') {
								goto l367
							}
							position++
							add(ruleSCOPE, position369)
						}
						goto l368
					l367:
						position, tokenIndex = position367, tokenIndex367
					}
				l368:
					if !_rules[ruleIdentifier]() {
						goto l364
					}
					add(ruleCommandName, position366)
				}
				{
					position370, tokenIndex370 := position, tokenIndex
					if !_rules[rule__]() {
						goto l370
					}
					{
						position372, tokenIndex372 := position, tokenIndex
						if !_rules[ruleCommandFirstArg]() {
							goto l373
						}
						if !_rules[rule__]() {
							goto l373
						}
						if !_rules[ruleCommandSecondArg]() {
							goto l373
						}
						goto l372
					l373:
						position, tokenIndex = position372, tokenIndex372
						if !_rules[ruleCommandFirstArg]() {
							goto l374
						}
						goto l372
					l374:
						position, tokenIndex = position372, tokenIndex372
						if !_rules[ruleCommandSecondArg]() {
							goto l370
						}
					}
				l372:
					goto l371
				l370:
					position, tokenIndex = position370, tokenIndex370
				}
			l371:
				{
					position375, tokenIndex375 := position, tokenIndex
					if !_rules[rule_]() {
						goto l375
					}
					{
						position377 := position
						{
							position378 := position
							if !_rules[rule_]() {
								goto l375
							}
							if buffer[position] != rune('-') {
								goto l375
							}
							position++
							if buffer[position] != rune('>') {
								goto l375
							}
							position++
							if !_rules[rule_]() {
								goto l375
							}
							add(ruleASSIGN, position378)
						}
						if !_rules[ruleVariable]() {
							goto l375
						}
						add(ruleCommandResultAssignment, position377)
					}
					goto l376
				l375:
					position, tokenIndex = position375, tokenIndex375
				}
			l376:
				add(ruleCommand, position365)
			}
			return true
		l364:
			position, tokenIndex = position364, tokenIndex364
			return false
		},
		/* 101 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 102 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position380, tokenIndex380 := position, tokenIndex
			{
				position381 := position
				{
					position382, tokenIndex382 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l383
					}
					goto l382
				l383:
					position, tokenIndex = position382, tokenIndex382
					if !_rules[ruleType]() {
						goto l380
					}
				}
			l382:
				add(ruleCommandFirstArg, position381)
			}
			return true
		l380:
			position, tokenIndex = position380, tokenIndex380
			return false
		},
		/* 103 CommandSecondArg <- <Object> */
		func() bool {
			position384, tokenIndex384 := position, tokenIndex
			{
				position385 := position
				if !_rules[ruleObject]() {
					goto l384
				}
				add(ruleCommandSecondArg, position385)
			}
			return true
		l384:
			position, tokenIndex = position384, tokenIndex384
			return false
		},
		/* 104 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 105 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 106 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position388, tokenIndex388 := position, tokenIndex
			{
				position389 := position
				{
					position390 := position
					if !_rules[rule_]() {
						goto l388
					}
					if buffer[position] != rune('i') {
						goto l388
					}
					position++
					if buffer[position] != rune('f') {
						goto l388
					}
					position++
					if !_rules[rule_]() {
						goto l388
					}
					add(ruleIF, position390)
				}
				if !_rules[ruleConditionalExpression]() {
					goto l388
				}
				if !_rules[ruleOPEN]() {
					goto l388
				}
			l391:
				{
					position392, tokenIndex392 := position, tokenIndex
					if !_rules[ruleBlock]() {
						goto l392
					}
					goto l391
				l392:
					position, tokenIndex = position392, tokenIndex392
				}
				if !_rules[ruleCLOSE]() {
					goto l388
				}
				add(ruleIfStanza, position389)
			}
			return true
		l388:
			position, tokenIndex = position388, tokenIndex388
			return false
		},
		/* 107 ElseIfStanza <- <(ELSE IfStanza)> */
		nil,
		/* 108 ElseStanza <- <(ELSE OPEN Block* CLOSE)> */
		nil,
		/* 109 Loop <- <(LOOP ((OPEN Block* CLOSE) / (LoopConditionFixedLength OPEN Block* CLOSE) / (LoopConditionIterable OPEN Block* CLOSE) / (LoopConditionBounded OPEN Block* CLOSE) / (LoopConditionTruthy OPEN Block* CLOSE)))> */
		nil,
		/* 110 LoopConditionFixedLength <- <(COUNT (Integer / Variable))> */
		nil,
		/* 111 LoopConditionIterable <- <(LoopIterableLHS IN LoopIterableRHS)> */
		nil,
		/* 112 LoopIterableLHS <- <VariableSequence> */
		nil,
		/* 113 LoopIterableRHS <- <(Command / Variable)> */
		nil,
		/* 114 LoopConditionBounded <- <(Command SEMI ConditionalExpression SEMI Command)> */
		nil,
		/* 115 LoopConditionTruthy <- <ConditionalExpression> */
		nil,
		/* 116 ConditionalExpression <- <(NOT? (ConditionWithAssignment / ConditionWithCommand / ConditionWithRegex / ConditionWithComparator))> */
		func() bool {
			position402, tokenIndex402 := position, tokenIndex
			{
				position403 := position
				{
					position404, tokenIndex404 := position, tokenIndex
					{
						position406 := position
						if !_rules[rule_]() {
							goto l404
						}
						if buffer[position] != rune('n') {
							goto l404
						}
						position++
						if buffer[position] != rune('o') {
							goto l404
						}
						position++
						if buffer[position] != rune('t') {
							goto l404
						}
						position++
						if !_rules[rule__]() {
							goto l404
						}
						add(ruleNOT, position406)
					}
					goto l405
				l404:
					position, tokenIndex = position404, tokenIndex404
				}
			l405:
				{
					position407, tokenIndex407 := position, tokenIndex
					{
						position409 := position
						if !_rules[ruleAssignment]() {
							goto l408
						}
						if !_rules[ruleSEMI]() {
							goto l408
						}
						if !_rules[ruleConditionalExpression]() {
							goto l408
						}
						add(ruleConditionWithAssignment, position409)
					}
					goto l407
				l408:
					position, tokenIndex = position407, tokenIndex407
					{
						position411 := position
						if !_rules[ruleCommand]() {
							goto l410
						}
						{
							position412, tokenIndex412 := position, tokenIndex
							if !_rules[ruleSEMI]() {
								goto l412
							}
							if !_rules[ruleConditionalExpression]() {
								goto l412
							}
							goto l413
						l412:
							position, tokenIndex = position412, tokenIndex412
						}
					l413:
						add(ruleConditionWithCommand, position411)
					}
					goto l407
				l410:
					position, tokenIndex = position407, tokenIndex407
					{
						position415 := position
						if !_rules[ruleExpression]() {
							goto l414
						}
						{
							position416 := position
							{
								position417, tokenIndex417 := position, tokenIndex
								{
									position419 := position
									if !_rules[rule_]() {
										goto l418
									}
									if buffer[position] != rune('=') {
										goto l418
									}
									position++
									if buffer[position] != rune('~') {
										goto l418
									}
									position++
									if !_rules[rule_]() {
										goto l418
									}
									add(ruleMatch, position419)
								}
								goto l417
							l418:
								position, tokenIndex = position417, tokenIndex417
								{
									position420 := position
									if !_rules[rule_]() {
										goto l414
									}
									if buffer[position] != rune('!') {
										goto l414
									}
									position++
									if buffer[position] != rune('~') {
										goto l414
									}
									position++
									if !_rules[rule_]() {
										goto l414
									}
									add(ruleUnmatch, position420)
								}
							}
						l417:
							add(ruleMatchOperator, position416)
						}
						if !_rules[ruleRegularExpression]() {
							goto l414
						}
						add(ruleConditionWithRegex, position415)
					}
					goto l407
				l414:
					position, tokenIndex = position407, tokenIndex407
					{
						position421 := position
						{
							position422 := position
							if !_rules[ruleExpression]() {
								goto l402
							}
							add(ruleConditionWithComparatorLHS, position422)
						}
						{
							position423, tokenIndex423 := position, tokenIndex
							{
								position425 := position
								{
									position426 := position
									if !_rules[rule_]() {
										goto l423
									}
									{
										position427, tokenIndex427 := position, tokenIndex
										{
											position429 := position
											if !_rules[rule_]() {
												goto l428
											}
											if buffer[position] != rune('=') {
												goto l428
											}
											position++
											if buffer[position] != rune('=') {
												goto l428
											}
											position++
											if !_rules[rule_]() {
												goto l428
											}
											add(ruleEquality, position429)
										}
										goto l427
									l428:
										position, tokenIndex = position427, tokenIndex427
										{
											position431 := position
											if !_rules[rule_]() {
												goto l430
											}
											if buffer[position] != rune('!') {
												goto l430
											}
											position++
											if buffer[position] != rune('=') {
												goto l430
											}
											position++
											if !_rules[rule_]() {
												goto l430
											}
											add(ruleNonEquality, position431)
										}
										goto l427
									l430:
										position, tokenIndex = position427, tokenIndex427
										{
											position433 := position
											if !_rules[rule_]() {
												goto l432
											}
											if buffer[position] != rune('>') {
												goto l432
											}
											position++
											if buffer[position] != rune('=') {
												goto l432
											}
											position++
											if !_rules[rule_]() {
												goto l432
											}
											add(ruleGreaterEqual, position433)
										}
										goto l427
									l432:
										position, tokenIndex = position427, tokenIndex427
										{
											position435 := position
											if !_rules[rule_]() {
												goto l434
											}
											if buffer[position] != rune('<') {
												goto l434
											}
											position++
											if buffer[position] != rune('=') {
												goto l434
											}
											position++
											if !_rules[rule_]() {
												goto l434
											}
											add(ruleLessEqual, position435)
										}
										goto l427
									l434:
										position, tokenIndex = position427, tokenIndex427
										{
											position437 := position
											if !_rules[rule_]() {
												goto l436
											}
											if buffer[position] != rune('>') {
												goto l436
											}
											position++
											if !_rules[rule_]() {
												goto l436
											}
											add(ruleGreaterThan, position437)
										}
										goto l427
									l436:
										position, tokenIndex = position427, tokenIndex427
										{
											position439 := position
											if !_rules[rule_]() {
												goto l438
											}
											if buffer[position] != rune('<') {
												goto l438
											}
											position++
											if !_rules[rule_]() {
												goto l438
											}
											add(ruleLessThan, position439)
										}
										goto l427
									l438:
										position, tokenIndex = position427, tokenIndex427
										{
											position441 := position
											if !_rules[rule_]() {
												goto l440
											}
											if buffer[position] != rune('i') {
												goto l440
											}
											position++
											if buffer[position] != rune('n') {
												goto l440
											}
											position++
											if !_rules[rule_]() {
												goto l440
											}
											add(ruleMembership, position441)
										}
										goto l427
									l440:
										position, tokenIndex = position427, tokenIndex427
										{
											position442 := position
											if !_rules[rule_]() {
												goto l423
											}
											if buffer[position] != rune('n') {
												goto l423
											}
											position++
											if buffer[position] != rune('o') {
												goto l423
											}
											position++
											if buffer[position] != rune('t') {
												goto l423
											}
											position++
											if !_rules[rule__]() {
												goto l423
											}
											if buffer[position] != rune('i') {
												goto l423
											}
											position++
											if buffer[position] != rune('n') {
												goto l423
											}
											position++
											if !_rules[rule_]() {
												goto l423
											}
											add(ruleNonMembership, position442)
										}
									}
								l427:
									if !_rules[rule_]() {
										goto l423
									}
									add(ruleComparisonOperator, position426)
								}
								if !_rules[ruleExpression]() {
									goto l423
								}
								add(ruleConditionWithComparatorRHS, position425)
							}
							goto l424
						l423:
							position, tokenIndex = position423, tokenIndex423
						}
					l424:
						add(ruleConditionWithComparator, position421)
					}
				}
			l407:
				add(ruleConditionalExpression, position403)
			}
			return true
		l402:
			position, tokenIndex = position402, tokenIndex402
			return false
		},
		/* 117 ConditionWithAssignment <- <(Assignment SEMI ConditionalExpression)> */
		nil,
		/* 118 ConditionWithCommand <- <(Command (SEMI ConditionalExpression)?)> */
		nil,
		/* 119 ConditionWithRegex <- <(Expression MatchOperator RegularExpression)> */
		nil,
		/* 120 ConditionWithComparator <- <(ConditionWithComparatorLHS ConditionWithComparatorRHS?)> */
		nil,
		/* 121 ConditionWithComparatorLHS <- <Expression> */
		nil,
		/* 122 ConditionWithComparatorRHS <- <(ComparisonOperator Expression)> */
		nil,
	}
	p.rules = _rules
}
