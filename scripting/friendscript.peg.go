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
	ruleComment
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
	ruleOperator
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
	"Comment",
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
	"Operator",
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
		/* 9 Comment <- <(_ '#' (!'\n' .)*)> */
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
		/* 27 Operator <- <(_ (Exponentiate / Multiply / Divide / Modulus / Add / Subtract / BitwiseAnd / BitwiseOr / BitwiseNot / BitwiseXor) _)> */
		nil,
		/* 28 Exponentiate <- <(_ ('*' '*') _)> */
		nil,
		/* 29 Multiply <- <(_ '*' _)> */
		nil,
		/* 30 Divide <- <(_ '/' _)> */
		nil,
		/* 31 Modulus <- <(_ '%' _)> */
		nil,
		/* 32 Add <- <(_ '+' _)> */
		nil,
		/* 33 Subtract <- <(_ '-' _)> */
		nil,
		/* 34 BitwiseAnd <- <(_ '&' _)> */
		nil,
		/* 35 BitwiseOr <- <(_ '|' _)> */
		nil,
		/* 36 BitwiseNot <- <(_ '~' _)> */
		nil,
		/* 37 BitwiseXor <- <(_ '^' _)> */
		nil,
		/* 38 MatchOperator <- <(Match / Unmatch)> */
		nil,
		/* 39 Unmatch <- <(_ ('!' '~') _)> */
		nil,
		/* 40 Match <- <(_ ('=' '~') _)> */
		nil,
		/* 41 AssignmentOperator <- <(_ (AssignEq / StarEq / DivEq / PlusEq / MinusEq / AndEq / OrEq / Append) _)> */
		nil,
		/* 42 AssignEq <- <(_ '=' _)> */
		nil,
		/* 43 StarEq <- <(_ ('*' '=') _)> */
		nil,
		/* 44 DivEq <- <(_ ('/' '=') _)> */
		nil,
		/* 45 PlusEq <- <(_ ('+' '=') _)> */
		nil,
		/* 46 MinusEq <- <(_ ('-' '=') _)> */
		nil,
		/* 47 AndEq <- <(_ ('&' '=') _)> */
		nil,
		/* 48 OrEq <- <(_ ('|' '=') _)> */
		nil,
		/* 49 Append <- <(_ ('<' '<') _)> */
		nil,
		/* 50 ComparisonOperator <- <(_ (Equality / NonEquality / GreaterEqual / LessEqual / GreaterThan / LessThan / Membership / NonMembership) _)> */
		nil,
		/* 51 Equality <- <(_ ('=' '=') _)> */
		nil,
		/* 52 NonEquality <- <(_ ('!' '=') _)> */
		nil,
		/* 53 GreaterThan <- <(_ '>' _)> */
		nil,
		/* 54 GreaterEqual <- <(_ ('>' '=') _)> */
		nil,
		/* 55 LessEqual <- <(_ ('<' '=') _)> */
		nil,
		/* 56 LessThan <- <(_ '<' _)> */
		nil,
		/* 57 Membership <- <(_ ('i' 'n') _)> */
		nil,
		/* 58 NonMembership <- <(_ ('n' 'o' 't') __ ('i' 'n') _)> */
		nil,
		/* 59 Variable <- <(('$' VariableNameSequence) / SKIPVAR)> */
		func() bool {
			position94, tokenIndex94 := position, tokenIndex
			{
				position95 := position
				{
					position96, tokenIndex96 := position, tokenIndex
					if buffer[position] != rune('$') {
						goto l97
					}
					position++
					{
						position98 := position
					l99:
						{
							position100, tokenIndex100 := position, tokenIndex
							if !_rules[ruleVariableName]() {
								goto l100
							}
							{
								position101 := position
								if buffer[position] != rune('.') {
									goto l100
								}
								position++
								add(ruleDOT, position101)
							}
							goto l99
						l100:
							position, tokenIndex = position100, tokenIndex100
						}
						if !_rules[ruleVariableName]() {
							goto l97
						}
						add(ruleVariableNameSequence, position98)
					}
					goto l96
				l97:
					position, tokenIndex = position96, tokenIndex96
					{
						position102 := position
						if !_rules[rule_]() {
							goto l94
						}
						if buffer[position] != rune('_') {
							goto l94
						}
						position++
						if !_rules[rule_]() {
							goto l94
						}
						add(ruleSKIPVAR, position102)
					}
				}
			l96:
				add(ruleVariable, position95)
			}
			return true
		l94:
			position, tokenIndex = position94, tokenIndex94
			return false
		},
		/* 60 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 61 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position104, tokenIndex104 := position, tokenIndex
			{
				position105 := position
				if !_rules[ruleIdentifier]() {
					goto l104
				}
				{
					position106, tokenIndex106 := position, tokenIndex
					if buffer[position] != rune('[') {
						goto l106
					}
					position++
					if !_rules[rule_]() {
						goto l106
					}
					{
						position108 := position
						if !_rules[ruleExpression]() {
							goto l106
						}
						add(ruleVariableIndex, position108)
					}
					if !_rules[rule_]() {
						goto l106
					}
					if buffer[position] != rune(']') {
						goto l106
					}
					position++
					goto l107
				l106:
					position, tokenIndex = position106, tokenIndex106
				}
			l107:
				add(ruleVariableName, position105)
			}
			return true
		l104:
			position, tokenIndex = position104, tokenIndex104
			return false
		},
		/* 62 VariableIndex <- <Expression> */
		nil,
		/* 63 Block <- <(_ (Comment / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position110, tokenIndex110 := position, tokenIndex
			{
				position111 := position
				if !_rules[rule_]() {
					goto l110
				}
				{
					position112, tokenIndex112 := position, tokenIndex
					{
						position114 := position
						if !_rules[rule_]() {
							goto l113
						}
						if buffer[position] != rune('#') {
							goto l113
						}
						position++
					l115:
						{
							position116, tokenIndex116 := position, tokenIndex
							{
								position117, tokenIndex117 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l117
								}
								position++
								goto l116
							l117:
								position, tokenIndex = position117, tokenIndex117
							}
							if !matchDot() {
								goto l116
							}
							goto l115
						l116:
							position, tokenIndex = position116, tokenIndex116
						}
						add(ruleComment, position114)
					}
					goto l112
				l113:
					position, tokenIndex = position112, tokenIndex112
					{
						position119 := position
						{
							position120, tokenIndex120 := position, tokenIndex
							{
								position122 := position
								{
									position123 := position
									if !_rules[rule_]() {
										goto l121
									}
									if buffer[position] != rune('b') {
										goto l121
									}
									position++
									if buffer[position] != rune('r') {
										goto l121
									}
									position++
									if buffer[position] != rune('e') {
										goto l121
									}
									position++
									if buffer[position] != rune('a') {
										goto l121
									}
									position++
									if buffer[position] != rune('k') {
										goto l121
									}
									position++
									if !_rules[rule_]() {
										goto l121
									}
									add(ruleBREAK, position123)
								}
								{
									position124, tokenIndex124 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l124
									}
									goto l125
								l124:
									position, tokenIndex = position124, tokenIndex124
								}
							l125:
								add(ruleFlowControlBreak, position122)
							}
							goto l120
						l121:
							position, tokenIndex = position120, tokenIndex120
							{
								position126 := position
								{
									position127 := position
									if !_rules[rule_]() {
										goto l118
									}
									if buffer[position] != rune('c') {
										goto l118
									}
									position++
									if buffer[position] != rune('o') {
										goto l118
									}
									position++
									if buffer[position] != rune('n') {
										goto l118
									}
									position++
									if buffer[position] != rune('t') {
										goto l118
									}
									position++
									if buffer[position] != rune('i') {
										goto l118
									}
									position++
									if buffer[position] != rune('n') {
										goto l118
									}
									position++
									if buffer[position] != rune('u') {
										goto l118
									}
									position++
									if buffer[position] != rune('e') {
										goto l118
									}
									position++
									if !_rules[rule_]() {
										goto l118
									}
									add(ruleCONT, position127)
								}
								{
									position128, tokenIndex128 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l128
									}
									goto l129
								l128:
									position, tokenIndex = position128, tokenIndex128
								}
							l129:
								add(ruleFlowControlContinue, position126)
							}
						}
					l120:
						add(ruleFlowControlWord, position119)
					}
					goto l112
				l118:
					position, tokenIndex = position112, tokenIndex112
					{
						position130 := position
						{
							position131, tokenIndex131 := position, tokenIndex
							{
								position133 := position
								if !_rules[ruleSEMI]() {
									goto l132
								}
								add(ruleNOOP, position133)
							}
							goto l131
						l132:
							position, tokenIndex = position131, tokenIndex131
							if !_rules[ruleAssignment]() {
								goto l134
							}
							goto l131
						l134:
							position, tokenIndex = position131, tokenIndex131
							{
								position136 := position
								{
									position137, tokenIndex137 := position, tokenIndex
									{
										position139 := position
										{
											position140 := position
											if !_rules[rule_]() {
												goto l138
											}
											if buffer[position] != rune('u') {
												goto l138
											}
											position++
											if buffer[position] != rune('n') {
												goto l138
											}
											position++
											if buffer[position] != rune('s') {
												goto l138
											}
											position++
											if buffer[position] != rune('e') {
												goto l138
											}
											position++
											if buffer[position] != rune('t') {
												goto l138
											}
											position++
											if !_rules[rule__]() {
												goto l138
											}
											add(ruleUNSET, position140)
										}
										if !_rules[ruleVariableSequence]() {
											goto l138
										}
										add(ruleDirectiveUnset, position139)
									}
									goto l137
								l138:
									position, tokenIndex = position137, tokenIndex137
									{
										position142 := position
										{
											position143 := position
											if !_rules[rule_]() {
												goto l141
											}
											if buffer[position] != rune('i') {
												goto l141
											}
											position++
											if buffer[position] != rune('n') {
												goto l141
											}
											position++
											if buffer[position] != rune('c') {
												goto l141
											}
											position++
											if buffer[position] != rune('l') {
												goto l141
											}
											position++
											if buffer[position] != rune('u') {
												goto l141
											}
											position++
											if buffer[position] != rune('d') {
												goto l141
											}
											position++
											if buffer[position] != rune('e') {
												goto l141
											}
											position++
											if !_rules[rule__]() {
												goto l141
											}
											add(ruleINCLUDE, position143)
										}
										if !_rules[ruleString]() {
											goto l141
										}
										add(ruleDirectiveInclude, position142)
									}
									goto l137
								l141:
									position, tokenIndex = position137, tokenIndex137
									{
										position144 := position
										{
											position145 := position
											if !_rules[rule_]() {
												goto l135
											}
											if buffer[position] != rune('d') {
												goto l135
											}
											position++
											if buffer[position] != rune('e') {
												goto l135
											}
											position++
											if buffer[position] != rune('c') {
												goto l135
											}
											position++
											if buffer[position] != rune('l') {
												goto l135
											}
											position++
											if buffer[position] != rune('a') {
												goto l135
											}
											position++
											if buffer[position] != rune('r') {
												goto l135
											}
											position++
											if buffer[position] != rune('e') {
												goto l135
											}
											position++
											if !_rules[rule__]() {
												goto l135
											}
											add(ruleDECLARE, position145)
										}
										if !_rules[ruleVariableSequence]() {
											goto l135
										}
										add(ruleDirectiveDeclare, position144)
									}
								}
							l137:
								add(ruleDirective, position136)
							}
							goto l131
						l135:
							position, tokenIndex = position131, tokenIndex131
							{
								position147 := position
								if !_rules[ruleIfStanza]() {
									goto l146
								}
							l148:
								{
									position149, tokenIndex149 := position, tokenIndex
									{
										position150 := position
										if !_rules[ruleELSE]() {
											goto l149
										}
										if !_rules[ruleIfStanza]() {
											goto l149
										}
										add(ruleElseIfStanza, position150)
									}
									goto l148
								l149:
									position, tokenIndex = position149, tokenIndex149
								}
								{
									position151, tokenIndex151 := position, tokenIndex
									{
										position153 := position
										if !_rules[ruleELSE]() {
											goto l151
										}
										if !_rules[ruleOPEN]() {
											goto l151
										}
									l154:
										{
											position155, tokenIndex155 := position, tokenIndex
											if !_rules[ruleBlock]() {
												goto l155
											}
											goto l154
										l155:
											position, tokenIndex = position155, tokenIndex155
										}
										if !_rules[ruleCLOSE]() {
											goto l151
										}
										add(ruleElseStanza, position153)
									}
									goto l152
								l151:
									position, tokenIndex = position151, tokenIndex151
								}
							l152:
								add(ruleConditional, position147)
							}
							goto l131
						l146:
							position, tokenIndex = position131, tokenIndex131
							{
								position157 := position
								{
									position158 := position
									if !_rules[rule_]() {
										goto l156
									}
									if buffer[position] != rune('l') {
										goto l156
									}
									position++
									if buffer[position] != rune('o') {
										goto l156
									}
									position++
									if buffer[position] != rune('o') {
										goto l156
									}
									position++
									if buffer[position] != rune('p') {
										goto l156
									}
									position++
									if !_rules[rule_]() {
										goto l156
									}
									add(ruleLOOP, position158)
								}
								{
									position159, tokenIndex159 := position, tokenIndex
									if !_rules[ruleOPEN]() {
										goto l160
									}
								l161:
									{
										position162, tokenIndex162 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l162
										}
										goto l161
									l162:
										position, tokenIndex = position162, tokenIndex162
									}
									if !_rules[ruleCLOSE]() {
										goto l160
									}
									goto l159
								l160:
									position, tokenIndex = position159, tokenIndex159
									{
										position164 := position
										{
											position165 := position
											if !_rules[rule_]() {
												goto l163
											}
											if buffer[position] != rune('c') {
												goto l163
											}
											position++
											if buffer[position] != rune('o') {
												goto l163
											}
											position++
											if buffer[position] != rune('u') {
												goto l163
											}
											position++
											if buffer[position] != rune('n') {
												goto l163
											}
											position++
											if buffer[position] != rune('t') {
												goto l163
											}
											position++
											if !_rules[rule_]() {
												goto l163
											}
											add(ruleCOUNT, position165)
										}
										{
											position166, tokenIndex166 := position, tokenIndex
											if !_rules[ruleInteger]() {
												goto l167
											}
											goto l166
										l167:
											position, tokenIndex = position166, tokenIndex166
											if !_rules[ruleVariable]() {
												goto l163
											}
										}
									l166:
										add(ruleLoopConditionFixedLength, position164)
									}
									if !_rules[ruleOPEN]() {
										goto l163
									}
								l168:
									{
										position169, tokenIndex169 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l169
										}
										goto l168
									l169:
										position, tokenIndex = position169, tokenIndex169
									}
									if !_rules[ruleCLOSE]() {
										goto l163
									}
									goto l159
								l163:
									position, tokenIndex = position159, tokenIndex159
									{
										position171 := position
										{
											position172 := position
											if !_rules[ruleVariableSequence]() {
												goto l170
											}
											add(ruleLoopIterableLHS, position172)
										}
										{
											position173 := position
											if !_rules[rule__]() {
												goto l170
											}
											if buffer[position] != rune('i') {
												goto l170
											}
											position++
											if buffer[position] != rune('n') {
												goto l170
											}
											position++
											if !_rules[rule__]() {
												goto l170
											}
											add(ruleIN, position173)
										}
										{
											position174 := position
											{
												position175, tokenIndex175 := position, tokenIndex
												if !_rules[ruleCommand]() {
													goto l176
												}
												goto l175
											l176:
												position, tokenIndex = position175, tokenIndex175
												if !_rules[ruleVariable]() {
													goto l170
												}
											}
										l175:
											add(ruleLoopIterableRHS, position174)
										}
										add(ruleLoopConditionIterable, position171)
									}
									if !_rules[ruleOPEN]() {
										goto l170
									}
								l177:
									{
										position178, tokenIndex178 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l178
										}
										goto l177
									l178:
										position, tokenIndex = position178, tokenIndex178
									}
									if !_rules[ruleCLOSE]() {
										goto l170
									}
									goto l159
								l170:
									position, tokenIndex = position159, tokenIndex159
									{
										position180 := position
										if !_rules[ruleCommand]() {
											goto l179
										}
										if !_rules[ruleSEMI]() {
											goto l179
										}
										if !_rules[ruleConditionalExpression]() {
											goto l179
										}
										if !_rules[ruleSEMI]() {
											goto l179
										}
										if !_rules[ruleCommand]() {
											goto l179
										}
										add(ruleLoopConditionBounded, position180)
									}
									if !_rules[ruleOPEN]() {
										goto l179
									}
								l181:
									{
										position182, tokenIndex182 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l182
										}
										goto l181
									l182:
										position, tokenIndex = position182, tokenIndex182
									}
									if !_rules[ruleCLOSE]() {
										goto l179
									}
									goto l159
								l179:
									position, tokenIndex = position159, tokenIndex159
									{
										position183 := position
										if !_rules[ruleConditionalExpression]() {
											goto l156
										}
										add(ruleLoopConditionTruthy, position183)
									}
									if !_rules[ruleOPEN]() {
										goto l156
									}
								l184:
									{
										position185, tokenIndex185 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l185
										}
										goto l184
									l185:
										position, tokenIndex = position185, tokenIndex185
									}
									if !_rules[ruleCLOSE]() {
										goto l156
									}
								}
							l159:
								add(ruleLoop, position157)
							}
							goto l131
						l156:
							position, tokenIndex = position131, tokenIndex131
							if !_rules[ruleCommand]() {
								goto l110
							}
						}
					l131:
						add(ruleStatementBlock, position130)
					}
				}
			l112:
				{
					position186, tokenIndex186 := position, tokenIndex
					if !_rules[ruleSEMI]() {
						goto l186
					}
					goto l187
				l186:
					position, tokenIndex = position186, tokenIndex186
				}
			l187:
				if !_rules[rule_]() {
					goto l110
				}
				add(ruleBlock, position111)
			}
			return true
		l110:
			position, tokenIndex = position110, tokenIndex110
			return false
		},
		/* 64 FlowControlWord <- <(FlowControlBreak / FlowControlContinue)> */
		nil,
		/* 65 FlowControlBreak <- <(BREAK PositiveInteger?)> */
		nil,
		/* 66 FlowControlContinue <- <(CONT PositiveInteger?)> */
		nil,
		/* 67 StatementBlock <- <(NOOP / Assignment / Directive / Conditional / Loop / Command)> */
		nil,
		/* 68 Assignment <- <(AssignmentLHS AssignmentOperator AssignmentRHS)> */
		func() bool {
			position192, tokenIndex192 := position, tokenIndex
			{
				position193 := position
				{
					position194 := position
					if !_rules[ruleVariableSequence]() {
						goto l192
					}
					add(ruleAssignmentLHS, position194)
				}
				{
					position195 := position
					if !_rules[rule_]() {
						goto l192
					}
					{
						position196, tokenIndex196 := position, tokenIndex
						{
							position198 := position
							if !_rules[rule_]() {
								goto l197
							}
							if buffer[position] != rune('=') {
								goto l197
							}
							position++
							if !_rules[rule_]() {
								goto l197
							}
							add(ruleAssignEq, position198)
						}
						goto l196
					l197:
						position, tokenIndex = position196, tokenIndex196
						{
							position200 := position
							if !_rules[rule_]() {
								goto l199
							}
							if buffer[position] != rune('*') {
								goto l199
							}
							position++
							if buffer[position] != rune('=') {
								goto l199
							}
							position++
							if !_rules[rule_]() {
								goto l199
							}
							add(ruleStarEq, position200)
						}
						goto l196
					l199:
						position, tokenIndex = position196, tokenIndex196
						{
							position202 := position
							if !_rules[rule_]() {
								goto l201
							}
							if buffer[position] != rune('/') {
								goto l201
							}
							position++
							if buffer[position] != rune('=') {
								goto l201
							}
							position++
							if !_rules[rule_]() {
								goto l201
							}
							add(ruleDivEq, position202)
						}
						goto l196
					l201:
						position, tokenIndex = position196, tokenIndex196
						{
							position204 := position
							if !_rules[rule_]() {
								goto l203
							}
							if buffer[position] != rune('+') {
								goto l203
							}
							position++
							if buffer[position] != rune('=') {
								goto l203
							}
							position++
							if !_rules[rule_]() {
								goto l203
							}
							add(rulePlusEq, position204)
						}
						goto l196
					l203:
						position, tokenIndex = position196, tokenIndex196
						{
							position206 := position
							if !_rules[rule_]() {
								goto l205
							}
							if buffer[position] != rune('-') {
								goto l205
							}
							position++
							if buffer[position] != rune('=') {
								goto l205
							}
							position++
							if !_rules[rule_]() {
								goto l205
							}
							add(ruleMinusEq, position206)
						}
						goto l196
					l205:
						position, tokenIndex = position196, tokenIndex196
						{
							position208 := position
							if !_rules[rule_]() {
								goto l207
							}
							if buffer[position] != rune('&') {
								goto l207
							}
							position++
							if buffer[position] != rune('=') {
								goto l207
							}
							position++
							if !_rules[rule_]() {
								goto l207
							}
							add(ruleAndEq, position208)
						}
						goto l196
					l207:
						position, tokenIndex = position196, tokenIndex196
						{
							position210 := position
							if !_rules[rule_]() {
								goto l209
							}
							if buffer[position] != rune('|') {
								goto l209
							}
							position++
							if buffer[position] != rune('=') {
								goto l209
							}
							position++
							if !_rules[rule_]() {
								goto l209
							}
							add(ruleOrEq, position210)
						}
						goto l196
					l209:
						position, tokenIndex = position196, tokenIndex196
						{
							position211 := position
							if !_rules[rule_]() {
								goto l192
							}
							if buffer[position] != rune('<') {
								goto l192
							}
							position++
							if buffer[position] != rune('<') {
								goto l192
							}
							position++
							if !_rules[rule_]() {
								goto l192
							}
							add(ruleAppend, position211)
						}
					}
				l196:
					if !_rules[rule_]() {
						goto l192
					}
					add(ruleAssignmentOperator, position195)
				}
				{
					position212 := position
					if !_rules[ruleExpressionSequence]() {
						goto l192
					}
					add(ruleAssignmentRHS, position212)
				}
				add(ruleAssignment, position193)
			}
			return true
		l192:
			position, tokenIndex = position192, tokenIndex192
			return false
		},
		/* 69 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 70 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 71 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position215, tokenIndex215 := position, tokenIndex
			{
				position216 := position
			l217:
				{
					position218, tokenIndex218 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l218
					}
					if !_rules[ruleCOMMA]() {
						goto l218
					}
					goto l217
				l218:
					position, tokenIndex = position218, tokenIndex218
				}
				if !_rules[ruleVariable]() {
					goto l215
				}
				add(ruleVariableSequence, position216)
			}
			return true
		l215:
			position, tokenIndex = position215, tokenIndex215
			return false
		},
		/* 72 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position219, tokenIndex219 := position, tokenIndex
			{
				position220 := position
			l221:
				{
					position222, tokenIndex222 := position, tokenIndex
					if !_rules[ruleExpression]() {
						goto l222
					}
					if !_rules[ruleCOMMA]() {
						goto l222
					}
					goto l221
				l222:
					position, tokenIndex = position222, tokenIndex222
				}
				if !_rules[ruleExpression]() {
					goto l219
				}
				add(ruleExpressionSequence, position220)
			}
			return true
		l219:
			position, tokenIndex = position219, tokenIndex219
			return false
		},
		/* 73 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				if !_rules[rule_]() {
					goto l223
				}
				{
					position225 := position
					{
						position226 := position
						{
							position227, tokenIndex227 := position, tokenIndex
							if !_rules[ruleType]() {
								goto l228
							}
							goto l227
						l228:
							position, tokenIndex = position227, tokenIndex227
							if !_rules[ruleVariable]() {
								goto l223
							}
						}
					l227:
						add(ruleValueYielding, position226)
					}
					add(ruleExpressionLHS, position225)
				}
				{
					position229, tokenIndex229 := position, tokenIndex
					{
						position231 := position
						{
							position232 := position
							if !_rules[rule_]() {
								goto l229
							}
							{
								position233, tokenIndex233 := position, tokenIndex
								{
									position235 := position
									if !_rules[rule_]() {
										goto l234
									}
									if buffer[position] != rune('*') {
										goto l234
									}
									position++
									if buffer[position] != rune('*') {
										goto l234
									}
									position++
									if !_rules[rule_]() {
										goto l234
									}
									add(ruleExponentiate, position235)
								}
								goto l233
							l234:
								position, tokenIndex = position233, tokenIndex233
								{
									position237 := position
									if !_rules[rule_]() {
										goto l236
									}
									if buffer[position] != rune('*') {
										goto l236
									}
									position++
									if !_rules[rule_]() {
										goto l236
									}
									add(ruleMultiply, position237)
								}
								goto l233
							l236:
								position, tokenIndex = position233, tokenIndex233
								{
									position239 := position
									if !_rules[rule_]() {
										goto l238
									}
									if buffer[position] != rune('/') {
										goto l238
									}
									position++
									if !_rules[rule_]() {
										goto l238
									}
									add(ruleDivide, position239)
								}
								goto l233
							l238:
								position, tokenIndex = position233, tokenIndex233
								{
									position241 := position
									if !_rules[rule_]() {
										goto l240
									}
									if buffer[position] != rune('%') {
										goto l240
									}
									position++
									if !_rules[rule_]() {
										goto l240
									}
									add(ruleModulus, position241)
								}
								goto l233
							l240:
								position, tokenIndex = position233, tokenIndex233
								{
									position243 := position
									if !_rules[rule_]() {
										goto l242
									}
									if buffer[position] != rune('+') {
										goto l242
									}
									position++
									if !_rules[rule_]() {
										goto l242
									}
									add(ruleAdd, position243)
								}
								goto l233
							l242:
								position, tokenIndex = position233, tokenIndex233
								{
									position245 := position
									if !_rules[rule_]() {
										goto l244
									}
									if buffer[position] != rune('-') {
										goto l244
									}
									position++
									if !_rules[rule_]() {
										goto l244
									}
									add(ruleSubtract, position245)
								}
								goto l233
							l244:
								position, tokenIndex = position233, tokenIndex233
								{
									position247 := position
									if !_rules[rule_]() {
										goto l246
									}
									if buffer[position] != rune('&') {
										goto l246
									}
									position++
									if !_rules[rule_]() {
										goto l246
									}
									add(ruleBitwiseAnd, position247)
								}
								goto l233
							l246:
								position, tokenIndex = position233, tokenIndex233
								{
									position249 := position
									if !_rules[rule_]() {
										goto l248
									}
									if buffer[position] != rune('|') {
										goto l248
									}
									position++
									if !_rules[rule_]() {
										goto l248
									}
									add(ruleBitwiseOr, position249)
								}
								goto l233
							l248:
								position, tokenIndex = position233, tokenIndex233
								{
									position251 := position
									if !_rules[rule_]() {
										goto l250
									}
									if buffer[position] != rune('~') {
										goto l250
									}
									position++
									if !_rules[rule_]() {
										goto l250
									}
									add(ruleBitwiseNot, position251)
								}
								goto l233
							l250:
								position, tokenIndex = position233, tokenIndex233
								{
									position252 := position
									if !_rules[rule_]() {
										goto l229
									}
									if buffer[position] != rune('^') {
										goto l229
									}
									position++
									if !_rules[rule_]() {
										goto l229
									}
									add(ruleBitwiseXor, position252)
								}
							}
						l233:
							if !_rules[rule_]() {
								goto l229
							}
							add(ruleOperator, position232)
						}
						if !_rules[ruleExpression]() {
							goto l229
						}
						add(ruleExpressionRHS, position231)
					}
					goto l230
				l229:
					position, tokenIndex = position229, tokenIndex229
				}
			l230:
				if !_rules[rule_]() {
					goto l223
				}
				add(ruleExpression, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 74 ExpressionLHS <- <ValueYielding> */
		nil,
		/* 75 ExpressionRHS <- <(Operator Expression)> */
		nil,
		/* 76 ValueYielding <- <(Type / Variable)> */
		nil,
		/* 77 Directive <- <(DirectiveUnset / DirectiveInclude / DirectiveDeclare)> */
		nil,
		/* 78 DirectiveUnset <- <(UNSET VariableSequence)> */
		nil,
		/* 79 DirectiveInclude <- <(INCLUDE String)> */
		nil,
		/* 80 DirectiveDeclare <- <(DECLARE VariableSequence)> */
		nil,
		/* 81 Command <- <(_ CommandName (__ ((CommandFirstArg __ CommandSecondArg) / CommandFirstArg / CommandSecondArg))? (_ CommandResultAssignment)?)> */
		func() bool {
			position260, tokenIndex260 := position, tokenIndex
			{
				position261 := position
				if !_rules[rule_]() {
					goto l260
				}
				{
					position262 := position
					{
						position263, tokenIndex263 := position, tokenIndex
						if !_rules[ruleIdentifier]() {
							goto l263
						}
						{
							position265 := position
							if buffer[position] != rune(':') {
								goto l263
							}
							position++
							if buffer[position] != rune(':') {
								goto l263
							}
							position++
							add(ruleSCOPE, position265)
						}
						goto l264
					l263:
						position, tokenIndex = position263, tokenIndex263
					}
				l264:
					if !_rules[ruleIdentifier]() {
						goto l260
					}
					add(ruleCommandName, position262)
				}
				{
					position266, tokenIndex266 := position, tokenIndex
					if !_rules[rule__]() {
						goto l266
					}
					{
						position268, tokenIndex268 := position, tokenIndex
						if !_rules[ruleCommandFirstArg]() {
							goto l269
						}
						if !_rules[rule__]() {
							goto l269
						}
						if !_rules[ruleCommandSecondArg]() {
							goto l269
						}
						goto l268
					l269:
						position, tokenIndex = position268, tokenIndex268
						if !_rules[ruleCommandFirstArg]() {
							goto l270
						}
						goto l268
					l270:
						position, tokenIndex = position268, tokenIndex268
						if !_rules[ruleCommandSecondArg]() {
							goto l266
						}
					}
				l268:
					goto l267
				l266:
					position, tokenIndex = position266, tokenIndex266
				}
			l267:
				{
					position271, tokenIndex271 := position, tokenIndex
					if !_rules[rule_]() {
						goto l271
					}
					{
						position273 := position
						{
							position274 := position
							if !_rules[rule_]() {
								goto l271
							}
							if buffer[position] != rune('-') {
								goto l271
							}
							position++
							if buffer[position] != rune('>') {
								goto l271
							}
							position++
							if !_rules[rule_]() {
								goto l271
							}
							add(ruleASSIGN, position274)
						}
						if !_rules[ruleVariable]() {
							goto l271
						}
						add(ruleCommandResultAssignment, position273)
					}
					goto l272
				l271:
					position, tokenIndex = position271, tokenIndex271
				}
			l272:
				add(ruleCommand, position261)
			}
			return true
		l260:
			position, tokenIndex = position260, tokenIndex260
			return false
		},
		/* 82 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 83 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position276, tokenIndex276 := position, tokenIndex
			{
				position277 := position
				{
					position278, tokenIndex278 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l279
					}
					goto l278
				l279:
					position, tokenIndex = position278, tokenIndex278
					if !_rules[ruleType]() {
						goto l276
					}
				}
			l278:
				add(ruleCommandFirstArg, position277)
			}
			return true
		l276:
			position, tokenIndex = position276, tokenIndex276
			return false
		},
		/* 84 CommandSecondArg <- <Object> */
		func() bool {
			position280, tokenIndex280 := position, tokenIndex
			{
				position281 := position
				if !_rules[ruleObject]() {
					goto l280
				}
				add(ruleCommandSecondArg, position281)
			}
			return true
		l280:
			position, tokenIndex = position280, tokenIndex280
			return false
		},
		/* 85 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 86 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 87 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position284, tokenIndex284 := position, tokenIndex
			{
				position285 := position
				{
					position286 := position
					if !_rules[rule_]() {
						goto l284
					}
					if buffer[position] != rune('i') {
						goto l284
					}
					position++
					if buffer[position] != rune('f') {
						goto l284
					}
					position++
					if !_rules[rule_]() {
						goto l284
					}
					add(ruleIF, position286)
				}
				if !_rules[ruleConditionalExpression]() {
					goto l284
				}
				if !_rules[ruleOPEN]() {
					goto l284
				}
			l287:
				{
					position288, tokenIndex288 := position, tokenIndex
					if !_rules[ruleBlock]() {
						goto l288
					}
					goto l287
				l288:
					position, tokenIndex = position288, tokenIndex288
				}
				if !_rules[ruleCLOSE]() {
					goto l284
				}
				add(ruleIfStanza, position285)
			}
			return true
		l284:
			position, tokenIndex = position284, tokenIndex284
			return false
		},
		/* 88 ElseIfStanza <- <(ELSE IfStanza)> */
		nil,
		/* 89 ElseStanza <- <(ELSE OPEN Block* CLOSE)> */
		nil,
		/* 90 Loop <- <(LOOP ((OPEN Block* CLOSE) / (LoopConditionFixedLength OPEN Block* CLOSE) / (LoopConditionIterable OPEN Block* CLOSE) / (LoopConditionBounded OPEN Block* CLOSE) / (LoopConditionTruthy OPEN Block* CLOSE)))> */
		nil,
		/* 91 LoopConditionFixedLength <- <(COUNT (Integer / Variable))> */
		nil,
		/* 92 LoopConditionIterable <- <(LoopIterableLHS IN LoopIterableRHS)> */
		nil,
		/* 93 LoopIterableLHS <- <VariableSequence> */
		nil,
		/* 94 LoopIterableRHS <- <(Command / Variable)> */
		nil,
		/* 95 LoopConditionBounded <- <(Command SEMI ConditionalExpression SEMI Command)> */
		nil,
		/* 96 LoopConditionTruthy <- <ConditionalExpression> */
		nil,
		/* 97 ConditionalExpression <- <(NOT? (ConditionWithAssignment / ConditionWithCommand / ConditionWithRegex / ConditionWithComparator))> */
		func() bool {
			position298, tokenIndex298 := position, tokenIndex
			{
				position299 := position
				{
					position300, tokenIndex300 := position, tokenIndex
					{
						position302 := position
						if !_rules[rule_]() {
							goto l300
						}
						if buffer[position] != rune('n') {
							goto l300
						}
						position++
						if buffer[position] != rune('o') {
							goto l300
						}
						position++
						if buffer[position] != rune('t') {
							goto l300
						}
						position++
						if !_rules[rule__]() {
							goto l300
						}
						add(ruleNOT, position302)
					}
					goto l301
				l300:
					position, tokenIndex = position300, tokenIndex300
				}
			l301:
				{
					position303, tokenIndex303 := position, tokenIndex
					{
						position305 := position
						if !_rules[ruleAssignment]() {
							goto l304
						}
						if !_rules[ruleSEMI]() {
							goto l304
						}
						if !_rules[ruleConditionalExpression]() {
							goto l304
						}
						add(ruleConditionWithAssignment, position305)
					}
					goto l303
				l304:
					position, tokenIndex = position303, tokenIndex303
					{
						position307 := position
						if !_rules[ruleCommand]() {
							goto l306
						}
						{
							position308, tokenIndex308 := position, tokenIndex
							if !_rules[ruleSEMI]() {
								goto l308
							}
							if !_rules[ruleConditionalExpression]() {
								goto l308
							}
							goto l309
						l308:
							position, tokenIndex = position308, tokenIndex308
						}
					l309:
						add(ruleConditionWithCommand, position307)
					}
					goto l303
				l306:
					position, tokenIndex = position303, tokenIndex303
					{
						position311 := position
						if !_rules[ruleExpression]() {
							goto l310
						}
						{
							position312 := position
							{
								position313, tokenIndex313 := position, tokenIndex
								{
									position315 := position
									if !_rules[rule_]() {
										goto l314
									}
									if buffer[position] != rune('=') {
										goto l314
									}
									position++
									if buffer[position] != rune('~') {
										goto l314
									}
									position++
									if !_rules[rule_]() {
										goto l314
									}
									add(ruleMatch, position315)
								}
								goto l313
							l314:
								position, tokenIndex = position313, tokenIndex313
								{
									position316 := position
									if !_rules[rule_]() {
										goto l310
									}
									if buffer[position] != rune('!') {
										goto l310
									}
									position++
									if buffer[position] != rune('~') {
										goto l310
									}
									position++
									if !_rules[rule_]() {
										goto l310
									}
									add(ruleUnmatch, position316)
								}
							}
						l313:
							add(ruleMatchOperator, position312)
						}
						if !_rules[ruleRegularExpression]() {
							goto l310
						}
						add(ruleConditionWithRegex, position311)
					}
					goto l303
				l310:
					position, tokenIndex = position303, tokenIndex303
					{
						position317 := position
						{
							position318 := position
							if !_rules[ruleExpression]() {
								goto l298
							}
							add(ruleConditionWithComparatorLHS, position318)
						}
						{
							position319, tokenIndex319 := position, tokenIndex
							{
								position321 := position
								{
									position322 := position
									if !_rules[rule_]() {
										goto l319
									}
									{
										position323, tokenIndex323 := position, tokenIndex
										{
											position325 := position
											if !_rules[rule_]() {
												goto l324
											}
											if buffer[position] != rune('=') {
												goto l324
											}
											position++
											if buffer[position] != rune('=') {
												goto l324
											}
											position++
											if !_rules[rule_]() {
												goto l324
											}
											add(ruleEquality, position325)
										}
										goto l323
									l324:
										position, tokenIndex = position323, tokenIndex323
										{
											position327 := position
											if !_rules[rule_]() {
												goto l326
											}
											if buffer[position] != rune('!') {
												goto l326
											}
											position++
											if buffer[position] != rune('=') {
												goto l326
											}
											position++
											if !_rules[rule_]() {
												goto l326
											}
											add(ruleNonEquality, position327)
										}
										goto l323
									l326:
										position, tokenIndex = position323, tokenIndex323
										{
											position329 := position
											if !_rules[rule_]() {
												goto l328
											}
											if buffer[position] != rune('>') {
												goto l328
											}
											position++
											if buffer[position] != rune('=') {
												goto l328
											}
											position++
											if !_rules[rule_]() {
												goto l328
											}
											add(ruleGreaterEqual, position329)
										}
										goto l323
									l328:
										position, tokenIndex = position323, tokenIndex323
										{
											position331 := position
											if !_rules[rule_]() {
												goto l330
											}
											if buffer[position] != rune('<') {
												goto l330
											}
											position++
											if buffer[position] != rune('=') {
												goto l330
											}
											position++
											if !_rules[rule_]() {
												goto l330
											}
											add(ruleLessEqual, position331)
										}
										goto l323
									l330:
										position, tokenIndex = position323, tokenIndex323
										{
											position333 := position
											if !_rules[rule_]() {
												goto l332
											}
											if buffer[position] != rune('>') {
												goto l332
											}
											position++
											if !_rules[rule_]() {
												goto l332
											}
											add(ruleGreaterThan, position333)
										}
										goto l323
									l332:
										position, tokenIndex = position323, tokenIndex323
										{
											position335 := position
											if !_rules[rule_]() {
												goto l334
											}
											if buffer[position] != rune('<') {
												goto l334
											}
											position++
											if !_rules[rule_]() {
												goto l334
											}
											add(ruleLessThan, position335)
										}
										goto l323
									l334:
										position, tokenIndex = position323, tokenIndex323
										{
											position337 := position
											if !_rules[rule_]() {
												goto l336
											}
											if buffer[position] != rune('i') {
												goto l336
											}
											position++
											if buffer[position] != rune('n') {
												goto l336
											}
											position++
											if !_rules[rule_]() {
												goto l336
											}
											add(ruleMembership, position337)
										}
										goto l323
									l336:
										position, tokenIndex = position323, tokenIndex323
										{
											position338 := position
											if !_rules[rule_]() {
												goto l319
											}
											if buffer[position] != rune('n') {
												goto l319
											}
											position++
											if buffer[position] != rune('o') {
												goto l319
											}
											position++
											if buffer[position] != rune('t') {
												goto l319
											}
											position++
											if !_rules[rule__]() {
												goto l319
											}
											if buffer[position] != rune('i') {
												goto l319
											}
											position++
											if buffer[position] != rune('n') {
												goto l319
											}
											position++
											if !_rules[rule_]() {
												goto l319
											}
											add(ruleNonMembership, position338)
										}
									}
								l323:
									if !_rules[rule_]() {
										goto l319
									}
									add(ruleComparisonOperator, position322)
								}
								if !_rules[ruleExpression]() {
									goto l319
								}
								add(ruleConditionWithComparatorRHS, position321)
							}
							goto l320
						l319:
							position, tokenIndex = position319, tokenIndex319
						}
					l320:
						add(ruleConditionWithComparator, position317)
					}
				}
			l303:
				add(ruleConditionalExpression, position299)
			}
			return true
		l298:
			position, tokenIndex = position298, tokenIndex298
			return false
		},
		/* 98 ConditionWithAssignment <- <(Assignment SEMI ConditionalExpression)> */
		nil,
		/* 99 ConditionWithCommand <- <(Command (SEMI ConditionalExpression)?)> */
		nil,
		/* 100 ConditionWithRegex <- <(Expression MatchOperator RegularExpression)> */
		nil,
		/* 101 ConditionWithComparator <- <(ConditionWithComparatorLHS ConditionWithComparatorRHS?)> */
		nil,
		/* 102 ConditionWithComparatorLHS <- <Expression> */
		nil,
		/* 103 ConditionWithComparatorRHS <- <(ComparisonOperator Expression)> */
		nil,
		/* 104 ScalarType <- <(Boolean / Float / Integer / String / NullValue)> */
		nil,
		/* 105 Identifier <- <(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / ([0-9] / [0-9]) / '_')*)> */
		func() bool {
			position346, tokenIndex346 := position, tokenIndex
			{
				position347 := position
				{
					position348, tokenIndex348 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex = position348, tokenIndex348
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l350
					}
					position++
					goto l348
				l350:
					position, tokenIndex = position348, tokenIndex348
					if buffer[position] != rune('_') {
						goto l346
					}
					position++
				}
			l348:
			l351:
				{
					position352, tokenIndex352 := position, tokenIndex
					{
						position353, tokenIndex353 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l354
						}
						position++
						goto l353
					l354:
						position, tokenIndex = position353, tokenIndex353
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l355
						}
						position++
						goto l353
					l355:
						position, tokenIndex = position353, tokenIndex353
						{
							position357, tokenIndex357 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l358
							}
							position++
							goto l357
						l358:
							position, tokenIndex = position357, tokenIndex357
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l356
							}
							position++
						}
					l357:
						goto l353
					l356:
						position, tokenIndex = position353, tokenIndex353
						if buffer[position] != rune('_') {
							goto l352
						}
						position++
					}
				l353:
					goto l351
				l352:
					position, tokenIndex = position352, tokenIndex352
				}
				add(ruleIdentifier, position347)
			}
			return true
		l346:
			position, tokenIndex = position346, tokenIndex346
			return false
		},
		/* 106 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 107 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 108 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position361, tokenIndex361 := position, tokenIndex
			{
				position362 := position
				{
					position363, tokenIndex363 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l363
					}
					position++
					goto l364
				l363:
					position, tokenIndex = position363, tokenIndex363
				}
			l364:
				if !_rules[rulePositiveInteger]() {
					goto l361
				}
				add(ruleInteger, position362)
			}
			return true
		l361:
			position, tokenIndex = position361, tokenIndex361
			return false
		},
		/* 109 PositiveInteger <- <[0-9]+> */
		func() bool {
			position365, tokenIndex365 := position, tokenIndex
			{
				position366 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l365
				}
				position++
			l367:
				{
					position368, tokenIndex368 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex = position368, tokenIndex368
				}
				add(rulePositiveInteger, position366)
			}
			return true
		l365:
			position, tokenIndex = position365, tokenIndex365
			return false
		},
		/* 110 String <- <(Triquote / StringLiteral / StringInterpolated)> */
		func() bool {
			position369, tokenIndex369 := position, tokenIndex
			{
				position370 := position
				{
					position371, tokenIndex371 := position, tokenIndex
					{
						position373 := position
						if !_rules[ruleTRIQUOT]() {
							goto l372
						}
						{
							position374 := position
						l375:
							{
								position376, tokenIndex376 := position, tokenIndex
								{
									position377, tokenIndex377 := position, tokenIndex
									if !_rules[ruleTRIQUOT]() {
										goto l377
									}
									goto l376
								l377:
									position, tokenIndex = position377, tokenIndex377
								}
								if !matchDot() {
									goto l376
								}
								goto l375
							l376:
								position, tokenIndex = position376, tokenIndex376
							}
							add(ruleTriquoteBody, position374)
						}
						if !_rules[ruleTRIQUOT]() {
							goto l372
						}
						add(ruleTriquote, position373)
					}
					goto l371
				l372:
					position, tokenIndex = position371, tokenIndex371
					if !_rules[ruleStringLiteral]() {
						goto l378
					}
					goto l371
				l378:
					position, tokenIndex = position371, tokenIndex371
					if !_rules[ruleStringInterpolated]() {
						goto l369
					}
				}
			l371:
				add(ruleString, position370)
			}
			return true
		l369:
			position, tokenIndex = position369, tokenIndex369
			return false
		},
		/* 111 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		func() bool {
			position379, tokenIndex379 := position, tokenIndex
			{
				position380 := position
				if buffer[position] != rune('\'') {
					goto l379
				}
				position++
			l381:
				{
					position382, tokenIndex382 := position, tokenIndex
					{
						position383, tokenIndex383 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l383
						}
						position++
						goto l382
					l383:
						position, tokenIndex = position383, tokenIndex383
					}
					if !matchDot() {
						goto l382
					}
					goto l381
				l382:
					position, tokenIndex = position382, tokenIndex382
				}
				if buffer[position] != rune('\'') {
					goto l379
				}
				position++
				add(ruleStringLiteral, position380)
			}
			return true
		l379:
			position, tokenIndex = position379, tokenIndex379
			return false
		},
		/* 112 StringInterpolated <- <('"' (!'"' .)* '"')> */
		func() bool {
			position384, tokenIndex384 := position, tokenIndex
			{
				position385 := position
				if buffer[position] != rune('"') {
					goto l384
				}
				position++
			l386:
				{
					position387, tokenIndex387 := position, tokenIndex
					{
						position388, tokenIndex388 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l388
						}
						position++
						goto l387
					l388:
						position, tokenIndex = position388, tokenIndex388
					}
					if !matchDot() {
						goto l387
					}
					goto l386
				l387:
					position, tokenIndex = position387, tokenIndex387
				}
				if buffer[position] != rune('"') {
					goto l384
				}
				position++
				add(ruleStringInterpolated, position385)
			}
			return true
		l384:
			position, tokenIndex = position384, tokenIndex384
			return false
		},
		/* 113 Triquote <- <(TRIQUOT TriquoteBody TRIQUOT)> */
		nil,
		/* 114 TriquoteBody <- <(!TRIQUOT .)*> */
		nil,
		/* 115 NullValue <- <('n' 'u' 'l' 'l')> */
		nil,
		/* 116 Object <- <(OPEN (_ KeyValuePair _)* CLOSE)> */
		func() bool {
			position392, tokenIndex392 := position, tokenIndex
			{
				position393 := position
				if !_rules[ruleOPEN]() {
					goto l392
				}
			l394:
				{
					position395, tokenIndex395 := position, tokenIndex
					if !_rules[rule_]() {
						goto l395
					}
					{
						position396 := position
						{
							position397 := position
							{
								position398, tokenIndex398 := position, tokenIndex
								if !_rules[ruleIdentifier]() {
									goto l399
								}
								goto l398
							l399:
								position, tokenIndex = position398, tokenIndex398
								if !_rules[ruleStringLiteral]() {
									goto l400
								}
								goto l398
							l400:
								position, tokenIndex = position398, tokenIndex398
								if !_rules[ruleStringInterpolated]() {
									goto l395
								}
							}
						l398:
							add(ruleKey, position397)
						}
						{
							position401 := position
							if !_rules[rule_]() {
								goto l395
							}
							if buffer[position] != rune(':') {
								goto l395
							}
							position++
							if !_rules[rule_]() {
								goto l395
							}
							add(ruleCOLON, position401)
						}
						{
							position402 := position
							{
								position403, tokenIndex403 := position, tokenIndex
								if !_rules[ruleArray]() {
									goto l404
								}
								goto l403
							l404:
								position, tokenIndex = position403, tokenIndex403
								if !_rules[ruleObject]() {
									goto l405
								}
								goto l403
							l405:
								position, tokenIndex = position403, tokenIndex403
								if !_rules[ruleExpression]() {
									goto l395
								}
							}
						l403:
							add(ruleKValue, position402)
						}
						{
							position406, tokenIndex406 := position, tokenIndex
							if !_rules[ruleCOMMA]() {
								goto l406
							}
							goto l407
						l406:
							position, tokenIndex = position406, tokenIndex406
						}
					l407:
						add(ruleKeyValuePair, position396)
					}
					if !_rules[rule_]() {
						goto l395
					}
					goto l394
				l395:
					position, tokenIndex = position395, tokenIndex395
				}
				if !_rules[ruleCLOSE]() {
					goto l392
				}
				add(ruleObject, position393)
			}
			return true
		l392:
			position, tokenIndex = position392, tokenIndex392
			return false
		},
		/* 117 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position408, tokenIndex408 := position, tokenIndex
			{
				position409 := position
				if buffer[position] != rune('[') {
					goto l408
				}
				position++
				if !_rules[rule_]() {
					goto l408
				}
				if !_rules[ruleExpressionSequence]() {
					goto l408
				}
				{
					position410, tokenIndex410 := position, tokenIndex
					if !_rules[ruleCOMMA]() {
						goto l410
					}
					goto l411
				l410:
					position, tokenIndex = position410, tokenIndex410
				}
			l411:
				if buffer[position] != rune(']') {
					goto l408
				}
				position++
				add(ruleArray, position409)
			}
			return true
		l408:
			position, tokenIndex = position408, tokenIndex408
			return false
		},
		/* 118 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position412, tokenIndex412 := position, tokenIndex
			{
				position413 := position
				if buffer[position] != rune('/') {
					goto l412
				}
				position++
				{
					position416, tokenIndex416 := position, tokenIndex
					if buffer[position] != rune('/') {
						goto l416
					}
					position++
					goto l412
				l416:
					position, tokenIndex = position416, tokenIndex416
				}
				if !matchDot() {
					goto l412
				}
			l414:
				{
					position415, tokenIndex415 := position, tokenIndex
					{
						position417, tokenIndex417 := position, tokenIndex
						if buffer[position] != rune('/') {
							goto l417
						}
						position++
						goto l415
					l417:
						position, tokenIndex = position417, tokenIndex417
					}
					if !matchDot() {
						goto l415
					}
					goto l414
				l415:
					position, tokenIndex = position415, tokenIndex415
				}
				if buffer[position] != rune('/') {
					goto l412
				}
				position++
			l418:
				{
					position419, tokenIndex419 := position, tokenIndex
					{
						position420, tokenIndex420 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l421
						}
						position++
						goto l420
					l421:
						position, tokenIndex = position420, tokenIndex420
						if buffer[position] != rune('l') {
							goto l422
						}
						position++
						goto l420
					l422:
						position, tokenIndex = position420, tokenIndex420
						if buffer[position] != rune('m') {
							goto l423
						}
						position++
						goto l420
					l423:
						position, tokenIndex = position420, tokenIndex420
						if buffer[position] != rune('s') {
							goto l424
						}
						position++
						goto l420
					l424:
						position, tokenIndex = position420, tokenIndex420
						if buffer[position] != rune('u') {
							goto l419
						}
						position++
					}
				l420:
					goto l418
				l419:
					position, tokenIndex = position419, tokenIndex419
				}
				add(ruleRegularExpression, position413)
			}
			return true
		l412:
			position, tokenIndex = position412, tokenIndex412
			return false
		},
		/* 119 KeyValuePair <- <(Key COLON KValue COMMA?)> */
		nil,
		/* 120 Key <- <(Identifier / StringLiteral / StringInterpolated)> */
		nil,
		/* 121 KValue <- <(Array / Object / Expression)> */
		nil,
		/* 122 Type <- <(Array / Object / RegularExpression / ScalarType)> */
		func() bool {
			position428, tokenIndex428 := position, tokenIndex
			{
				position429 := position
				{
					position430, tokenIndex430 := position, tokenIndex
					if !_rules[ruleArray]() {
						goto l431
					}
					goto l430
				l431:
					position, tokenIndex = position430, tokenIndex430
					if !_rules[ruleObject]() {
						goto l432
					}
					goto l430
				l432:
					position, tokenIndex = position430, tokenIndex430
					if !_rules[ruleRegularExpression]() {
						goto l433
					}
					goto l430
				l433:
					position, tokenIndex = position430, tokenIndex430
					{
						position434 := position
						{
							position435, tokenIndex435 := position, tokenIndex
							{
								position437 := position
								{
									position438, tokenIndex438 := position, tokenIndex
									if buffer[position] != rune('t') {
										goto l439
									}
									position++
									if buffer[position] != rune('r') {
										goto l439
									}
									position++
									if buffer[position] != rune('u') {
										goto l439
									}
									position++
									if buffer[position] != rune('e') {
										goto l439
									}
									position++
									goto l438
								l439:
									position, tokenIndex = position438, tokenIndex438
									if buffer[position] != rune('f') {
										goto l436
									}
									position++
									if buffer[position] != rune('a') {
										goto l436
									}
									position++
									if buffer[position] != rune('l') {
										goto l436
									}
									position++
									if buffer[position] != rune('s') {
										goto l436
									}
									position++
									if buffer[position] != rune('e') {
										goto l436
									}
									position++
								}
							l438:
								add(ruleBoolean, position437)
							}
							goto l435
						l436:
							position, tokenIndex = position435, tokenIndex435
							{
								position441 := position
								if !_rules[ruleInteger]() {
									goto l440
								}
								{
									position442, tokenIndex442 := position, tokenIndex
									if buffer[position] != rune('.') {
										goto l442
									}
									position++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l442
									}
									position++
								l444:
									{
										position445, tokenIndex445 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l445
										}
										position++
										goto l444
									l445:
										position, tokenIndex = position445, tokenIndex445
									}
									goto l443
								l442:
									position, tokenIndex = position442, tokenIndex442
								}
							l443:
								add(ruleFloat, position441)
							}
							goto l435
						l440:
							position, tokenIndex = position435, tokenIndex435
							if !_rules[ruleInteger]() {
								goto l446
							}
							goto l435
						l446:
							position, tokenIndex = position435, tokenIndex435
							if !_rules[ruleString]() {
								goto l447
							}
							goto l435
						l447:
							position, tokenIndex = position435, tokenIndex435
							{
								position448 := position
								if buffer[position] != rune('n') {
									goto l428
								}
								position++
								if buffer[position] != rune('u') {
									goto l428
								}
								position++
								if buffer[position] != rune('l') {
									goto l428
								}
								position++
								if buffer[position] != rune('l') {
									goto l428
								}
								position++
								add(ruleNullValue, position448)
							}
						}
					l435:
						add(ruleScalarType, position434)
					}
				}
			l430:
				add(ruleType, position429)
			}
			return true
		l428:
			position, tokenIndex = position428, tokenIndex428
			return false
		},
	}
	p.rules = _rules
}
