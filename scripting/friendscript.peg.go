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
	ruleBEGIN
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
	ruleEND
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
	ruleHeredoc
	ruleHeredocBody
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
	"BEGIN",
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
	"END",
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
	"Heredoc",
	"HeredocBody",
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
	rules  [125]func() bool
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
		/* 4 BEGIN <- <(_ ('b' 'e' 'g' 'i' 'n'))> */
		nil,
		/* 5 BREAK <- <(_ ('b' 'r' 'e' 'a' 'k') _)> */
		nil,
		/* 6 CLOSE <- <(_ '}' _)> */
		func() bool {
			position35, tokenIndex35 := position, tokenIndex
			{
				position36 := position
				if !_rules[rule_]() {
					goto l35
				}
				if buffer[position] != rune('}') {
					goto l35
				}
				position++
				if !_rules[rule_]() {
					goto l35
				}
				add(ruleCLOSE, position36)
			}
			return true
		l35:
			position, tokenIndex = position35, tokenIndex35
			return false
		},
		/* 7 COLON <- <(_ ':' _)> */
		nil,
		/* 8 COMMA <- <(_ ',' _)> */
		func() bool {
			position38, tokenIndex38 := position, tokenIndex
			{
				position39 := position
				if !_rules[rule_]() {
					goto l38
				}
				if buffer[position] != rune(',') {
					goto l38
				}
				position++
				if !_rules[rule_]() {
					goto l38
				}
				add(ruleCOMMA, position39)
			}
			return true
		l38:
			position, tokenIndex = position38, tokenIndex38
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
			position45, tokenIndex45 := position, tokenIndex
			{
				position46 := position
				if !_rules[rule_]() {
					goto l45
				}
				if buffer[position] != rune('e') {
					goto l45
				}
				position++
				if buffer[position] != rune('l') {
					goto l45
				}
				position++
				if buffer[position] != rune('s') {
					goto l45
				}
				position++
				if buffer[position] != rune('e') {
					goto l45
				}
				position++
				if !_rules[rule_]() {
					goto l45
				}
				add(ruleELSE, position46)
			}
			return true
		l45:
			position, tokenIndex = position45, tokenIndex45
			return false
		},
		/* 15 END <- <(_ ('e' 'n' 'd') _)> */
		func() bool {
			position47, tokenIndex47 := position, tokenIndex
			{
				position48 := position
				if !_rules[rule_]() {
					goto l47
				}
				if buffer[position] != rune('e') {
					goto l47
				}
				position++
				if buffer[position] != rune('n') {
					goto l47
				}
				position++
				if buffer[position] != rune('d') {
					goto l47
				}
				position++
				if !_rules[rule_]() {
					goto l47
				}
				add(ruleEND, position48)
			}
			return true
		l47:
			position, tokenIndex = position47, tokenIndex47
			return false
		},
		/* 16 IF <- <(_ ('i' 'f') _)> */
		nil,
		/* 17 IN <- <(__ ('i' 'n') __)> */
		nil,
		/* 18 INCLUDE <- <(_ ('i' 'n' 'c' 'l' 'u' 'd' 'e') __)> */
		nil,
		/* 19 LOOP <- <(_ ('l' 'o' 'o' 'p') _)> */
		nil,
		/* 20 NOOP <- <SEMI> */
		nil,
		/* 21 NOT <- <(_ ('n' 'o' 't') __)> */
		nil,
		/* 22 OPEN <- <(_ '{' _)> */
		func() bool {
			position55, tokenIndex55 := position, tokenIndex
			{
				position56 := position
				if !_rules[rule_]() {
					goto l55
				}
				if buffer[position] != rune('{') {
					goto l55
				}
				position++
				if !_rules[rule_]() {
					goto l55
				}
				add(ruleOPEN, position56)
			}
			return true
		l55:
			position, tokenIndex = position55, tokenIndex55
			return false
		},
		/* 23 SCOPE <- <(':' ':')> */
		nil,
		/* 24 SEMI <- <(_ ';' _)> */
		func() bool {
			position58, tokenIndex58 := position, tokenIndex
			{
				position59 := position
				if !_rules[rule_]() {
					goto l58
				}
				if buffer[position] != rune(';') {
					goto l58
				}
				position++
				if !_rules[rule_]() {
					goto l58
				}
				add(ruleSEMI, position59)
			}
			return true
		l58:
			position, tokenIndex = position58, tokenIndex58
			return false
		},
		/* 25 SHEBANG <- <('#' '!' (!'\n' .)+ '\n')> */
		nil,
		/* 26 SKIPVAR <- <(_ '_' _)> */
		nil,
		/* 27 UNSET <- <(_ ('u' 'n' 's' 'e' 't') __)> */
		nil,
		/* 28 Operator <- <(_ (Exponentiate / Multiply / Divide / Modulus / Add / Subtract / BitwiseAnd / BitwiseOr / BitwiseNot / BitwiseXor) _)> */
		nil,
		/* 29 Exponentiate <- <(_ ('*' '*') _)> */
		nil,
		/* 30 Multiply <- <(_ '*' _)> */
		nil,
		/* 31 Divide <- <(_ '/' _)> */
		nil,
		/* 32 Modulus <- <(_ '%' _)> */
		nil,
		/* 33 Add <- <(_ '+' _)> */
		nil,
		/* 34 Subtract <- <(_ '-' _)> */
		nil,
		/* 35 BitwiseAnd <- <(_ '&' _)> */
		nil,
		/* 36 BitwiseOr <- <(_ '|' _)> */
		nil,
		/* 37 BitwiseNot <- <(_ '~' _)> */
		nil,
		/* 38 BitwiseXor <- <(_ '^' _)> */
		nil,
		/* 39 MatchOperator <- <(Match / Unmatch)> */
		nil,
		/* 40 Unmatch <- <(_ ('!' '~') _)> */
		nil,
		/* 41 Match <- <(_ ('=' '~') _)> */
		nil,
		/* 42 AssignmentOperator <- <(_ (AssignEq / StarEq / DivEq / PlusEq / MinusEq / AndEq / OrEq / Append) _)> */
		nil,
		/* 43 AssignEq <- <(_ '=' _)> */
		nil,
		/* 44 StarEq <- <(_ ('*' '=') _)> */
		nil,
		/* 45 DivEq <- <(_ ('/' '=') _)> */
		nil,
		/* 46 PlusEq <- <(_ ('+' '=') _)> */
		nil,
		/* 47 MinusEq <- <(_ ('-' '=') _)> */
		nil,
		/* 48 AndEq <- <(_ ('&' '=') _)> */
		nil,
		/* 49 OrEq <- <(_ ('|' '=') _)> */
		nil,
		/* 50 Append <- <(_ ('<' '<') _)> */
		nil,
		/* 51 ComparisonOperator <- <(_ (Equality / NonEquality / GreaterEqual / LessEqual / GreaterThan / LessThan / Membership / NonMembership) _)> */
		nil,
		/* 52 Equality <- <(_ ('=' '=') _)> */
		nil,
		/* 53 NonEquality <- <(_ ('!' '=') _)> */
		nil,
		/* 54 GreaterThan <- <(_ '>' _)> */
		nil,
		/* 55 GreaterEqual <- <(_ ('>' '=') _)> */
		nil,
		/* 56 LessEqual <- <(_ ('<' '=') _)> */
		nil,
		/* 57 LessThan <- <(_ '<' _)> */
		nil,
		/* 58 Membership <- <(_ ('i' 'n') _)> */
		nil,
		/* 59 NonMembership <- <(_ ('n' 'o' 't') __ ('i' 'n') _)> */
		nil,
		/* 60 Variable <- <(('$' VariableNameSequence) / SKIPVAR)> */
		func() bool {
			position95, tokenIndex95 := position, tokenIndex
			{
				position96 := position
				{
					position97, tokenIndex97 := position, tokenIndex
					if buffer[position] != rune('$') {
						goto l98
					}
					position++
					{
						position99 := position
					l100:
						{
							position101, tokenIndex101 := position, tokenIndex
							if !_rules[ruleVariableName]() {
								goto l101
							}
							{
								position102 := position
								if buffer[position] != rune('.') {
									goto l101
								}
								position++
								add(ruleDOT, position102)
							}
							goto l100
						l101:
							position, tokenIndex = position101, tokenIndex101
						}
						if !_rules[ruleVariableName]() {
							goto l98
						}
						add(ruleVariableNameSequence, position99)
					}
					goto l97
				l98:
					position, tokenIndex = position97, tokenIndex97
					{
						position103 := position
						if !_rules[rule_]() {
							goto l95
						}
						if buffer[position] != rune('_') {
							goto l95
						}
						position++
						if !_rules[rule_]() {
							goto l95
						}
						add(ruleSKIPVAR, position103)
					}
				}
			l97:
				add(ruleVariable, position96)
			}
			return true
		l95:
			position, tokenIndex = position95, tokenIndex95
			return false
		},
		/* 61 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 62 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position105, tokenIndex105 := position, tokenIndex
			{
				position106 := position
				if !_rules[ruleIdentifier]() {
					goto l105
				}
				{
					position107, tokenIndex107 := position, tokenIndex
					if buffer[position] != rune('[') {
						goto l107
					}
					position++
					if !_rules[rule_]() {
						goto l107
					}
					{
						position109 := position
						if !_rules[ruleExpression]() {
							goto l107
						}
						add(ruleVariableIndex, position109)
					}
					if !_rules[rule_]() {
						goto l107
					}
					if buffer[position] != rune(']') {
						goto l107
					}
					position++
					goto l108
				l107:
					position, tokenIndex = position107, tokenIndex107
				}
			l108:
				add(ruleVariableName, position106)
			}
			return true
		l105:
			position, tokenIndex = position105, tokenIndex105
			return false
		},
		/* 63 VariableIndex <- <Expression> */
		nil,
		/* 64 Block <- <(_ (Comment / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position111, tokenIndex111 := position, tokenIndex
			{
				position112 := position
				if !_rules[rule_]() {
					goto l111
				}
				{
					position113, tokenIndex113 := position, tokenIndex
					{
						position115 := position
						if !_rules[rule_]() {
							goto l114
						}
						if buffer[position] != rune('#') {
							goto l114
						}
						position++
					l116:
						{
							position117, tokenIndex117 := position, tokenIndex
							{
								position118, tokenIndex118 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l118
								}
								position++
								goto l117
							l118:
								position, tokenIndex = position118, tokenIndex118
							}
							if !matchDot() {
								goto l117
							}
							goto l116
						l117:
							position, tokenIndex = position117, tokenIndex117
						}
						add(ruleComment, position115)
					}
					goto l113
				l114:
					position, tokenIndex = position113, tokenIndex113
					{
						position120 := position
						{
							position121, tokenIndex121 := position, tokenIndex
							{
								position123 := position
								{
									position124 := position
									if !_rules[rule_]() {
										goto l122
									}
									if buffer[position] != rune('b') {
										goto l122
									}
									position++
									if buffer[position] != rune('r') {
										goto l122
									}
									position++
									if buffer[position] != rune('e') {
										goto l122
									}
									position++
									if buffer[position] != rune('a') {
										goto l122
									}
									position++
									if buffer[position] != rune('k') {
										goto l122
									}
									position++
									if !_rules[rule_]() {
										goto l122
									}
									add(ruleBREAK, position124)
								}
								{
									position125, tokenIndex125 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l125
									}
									goto l126
								l125:
									position, tokenIndex = position125, tokenIndex125
								}
							l126:
								add(ruleFlowControlBreak, position123)
							}
							goto l121
						l122:
							position, tokenIndex = position121, tokenIndex121
							{
								position127 := position
								{
									position128 := position
									if !_rules[rule_]() {
										goto l119
									}
									if buffer[position] != rune('c') {
										goto l119
									}
									position++
									if buffer[position] != rune('o') {
										goto l119
									}
									position++
									if buffer[position] != rune('n') {
										goto l119
									}
									position++
									if buffer[position] != rune('t') {
										goto l119
									}
									position++
									if buffer[position] != rune('i') {
										goto l119
									}
									position++
									if buffer[position] != rune('n') {
										goto l119
									}
									position++
									if buffer[position] != rune('u') {
										goto l119
									}
									position++
									if buffer[position] != rune('e') {
										goto l119
									}
									position++
									if !_rules[rule_]() {
										goto l119
									}
									add(ruleCONT, position128)
								}
								{
									position129, tokenIndex129 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l129
									}
									goto l130
								l129:
									position, tokenIndex = position129, tokenIndex129
								}
							l130:
								add(ruleFlowControlContinue, position127)
							}
						}
					l121:
						add(ruleFlowControlWord, position120)
					}
					goto l113
				l119:
					position, tokenIndex = position113, tokenIndex113
					{
						position131 := position
						{
							position132, tokenIndex132 := position, tokenIndex
							{
								position134 := position
								if !_rules[ruleSEMI]() {
									goto l133
								}
								add(ruleNOOP, position134)
							}
							goto l132
						l133:
							position, tokenIndex = position132, tokenIndex132
							if !_rules[ruleAssignment]() {
								goto l135
							}
							goto l132
						l135:
							position, tokenIndex = position132, tokenIndex132
							{
								position137 := position
								{
									position138, tokenIndex138 := position, tokenIndex
									{
										position140 := position
										{
											position141 := position
											if !_rules[rule_]() {
												goto l139
											}
											if buffer[position] != rune('u') {
												goto l139
											}
											position++
											if buffer[position] != rune('n') {
												goto l139
											}
											position++
											if buffer[position] != rune('s') {
												goto l139
											}
											position++
											if buffer[position] != rune('e') {
												goto l139
											}
											position++
											if buffer[position] != rune('t') {
												goto l139
											}
											position++
											if !_rules[rule__]() {
												goto l139
											}
											add(ruleUNSET, position141)
										}
										if !_rules[ruleVariableSequence]() {
											goto l139
										}
										add(ruleDirectiveUnset, position140)
									}
									goto l138
								l139:
									position, tokenIndex = position138, tokenIndex138
									{
										position143 := position
										{
											position144 := position
											if !_rules[rule_]() {
												goto l142
											}
											if buffer[position] != rune('i') {
												goto l142
											}
											position++
											if buffer[position] != rune('n') {
												goto l142
											}
											position++
											if buffer[position] != rune('c') {
												goto l142
											}
											position++
											if buffer[position] != rune('l') {
												goto l142
											}
											position++
											if buffer[position] != rune('u') {
												goto l142
											}
											position++
											if buffer[position] != rune('d') {
												goto l142
											}
											position++
											if buffer[position] != rune('e') {
												goto l142
											}
											position++
											if !_rules[rule__]() {
												goto l142
											}
											add(ruleINCLUDE, position144)
										}
										if !_rules[ruleString]() {
											goto l142
										}
										add(ruleDirectiveInclude, position143)
									}
									goto l138
								l142:
									position, tokenIndex = position138, tokenIndex138
									{
										position145 := position
										{
											position146 := position
											if !_rules[rule_]() {
												goto l136
											}
											if buffer[position] != rune('d') {
												goto l136
											}
											position++
											if buffer[position] != rune('e') {
												goto l136
											}
											position++
											if buffer[position] != rune('c') {
												goto l136
											}
											position++
											if buffer[position] != rune('l') {
												goto l136
											}
											position++
											if buffer[position] != rune('a') {
												goto l136
											}
											position++
											if buffer[position] != rune('r') {
												goto l136
											}
											position++
											if buffer[position] != rune('e') {
												goto l136
											}
											position++
											if !_rules[rule__]() {
												goto l136
											}
											add(ruleDECLARE, position146)
										}
										if !_rules[ruleVariableSequence]() {
											goto l136
										}
										add(ruleDirectiveDeclare, position145)
									}
								}
							l138:
								add(ruleDirective, position137)
							}
							goto l132
						l136:
							position, tokenIndex = position132, tokenIndex132
							{
								position148 := position
								if !_rules[ruleIfStanza]() {
									goto l147
								}
							l149:
								{
									position150, tokenIndex150 := position, tokenIndex
									{
										position151 := position
										if !_rules[ruleELSE]() {
											goto l150
										}
										if !_rules[ruleIfStanza]() {
											goto l150
										}
										add(ruleElseIfStanza, position151)
									}
									goto l149
								l150:
									position, tokenIndex = position150, tokenIndex150
								}
								{
									position152, tokenIndex152 := position, tokenIndex
									{
										position154 := position
										if !_rules[ruleELSE]() {
											goto l152
										}
										if !_rules[ruleOPEN]() {
											goto l152
										}
									l155:
										{
											position156, tokenIndex156 := position, tokenIndex
											if !_rules[ruleBlock]() {
												goto l156
											}
											goto l155
										l156:
											position, tokenIndex = position156, tokenIndex156
										}
										if !_rules[ruleCLOSE]() {
											goto l152
										}
										add(ruleElseStanza, position154)
									}
									goto l153
								l152:
									position, tokenIndex = position152, tokenIndex152
								}
							l153:
								add(ruleConditional, position148)
							}
							goto l132
						l147:
							position, tokenIndex = position132, tokenIndex132
							{
								position158 := position
								{
									position159 := position
									if !_rules[rule_]() {
										goto l157
									}
									if buffer[position] != rune('l') {
										goto l157
									}
									position++
									if buffer[position] != rune('o') {
										goto l157
									}
									position++
									if buffer[position] != rune('o') {
										goto l157
									}
									position++
									if buffer[position] != rune('p') {
										goto l157
									}
									position++
									if !_rules[rule_]() {
										goto l157
									}
									add(ruleLOOP, position159)
								}
								{
									position160, tokenIndex160 := position, tokenIndex
									if !_rules[ruleOPEN]() {
										goto l161
									}
								l162:
									{
										position163, tokenIndex163 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l163
										}
										goto l162
									l163:
										position, tokenIndex = position163, tokenIndex163
									}
									if !_rules[ruleCLOSE]() {
										goto l161
									}
									goto l160
								l161:
									position, tokenIndex = position160, tokenIndex160
									{
										position165 := position
										{
											position166 := position
											if !_rules[rule_]() {
												goto l164
											}
											if buffer[position] != rune('c') {
												goto l164
											}
											position++
											if buffer[position] != rune('o') {
												goto l164
											}
											position++
											if buffer[position] != rune('u') {
												goto l164
											}
											position++
											if buffer[position] != rune('n') {
												goto l164
											}
											position++
											if buffer[position] != rune('t') {
												goto l164
											}
											position++
											if !_rules[rule_]() {
												goto l164
											}
											add(ruleCOUNT, position166)
										}
										{
											position167, tokenIndex167 := position, tokenIndex
											if !_rules[ruleInteger]() {
												goto l168
											}
											goto l167
										l168:
											position, tokenIndex = position167, tokenIndex167
											if !_rules[ruleVariable]() {
												goto l164
											}
										}
									l167:
										add(ruleLoopConditionFixedLength, position165)
									}
									if !_rules[ruleOPEN]() {
										goto l164
									}
								l169:
									{
										position170, tokenIndex170 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l170
										}
										goto l169
									l170:
										position, tokenIndex = position170, tokenIndex170
									}
									if !_rules[ruleCLOSE]() {
										goto l164
									}
									goto l160
								l164:
									position, tokenIndex = position160, tokenIndex160
									{
										position172 := position
										{
											position173 := position
											if !_rules[ruleVariableSequence]() {
												goto l171
											}
											add(ruleLoopIterableLHS, position173)
										}
										{
											position174 := position
											if !_rules[rule__]() {
												goto l171
											}
											if buffer[position] != rune('i') {
												goto l171
											}
											position++
											if buffer[position] != rune('n') {
												goto l171
											}
											position++
											if !_rules[rule__]() {
												goto l171
											}
											add(ruleIN, position174)
										}
										{
											position175 := position
											{
												position176, tokenIndex176 := position, tokenIndex
												if !_rules[ruleCommand]() {
													goto l177
												}
												goto l176
											l177:
												position, tokenIndex = position176, tokenIndex176
												if !_rules[ruleVariable]() {
													goto l171
												}
											}
										l176:
											add(ruleLoopIterableRHS, position175)
										}
										add(ruleLoopConditionIterable, position172)
									}
									if !_rules[ruleOPEN]() {
										goto l171
									}
								l178:
									{
										position179, tokenIndex179 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l179
										}
										goto l178
									l179:
										position, tokenIndex = position179, tokenIndex179
									}
									if !_rules[ruleCLOSE]() {
										goto l171
									}
									goto l160
								l171:
									position, tokenIndex = position160, tokenIndex160
									{
										position181 := position
										if !_rules[ruleCommand]() {
											goto l180
										}
										if !_rules[ruleSEMI]() {
											goto l180
										}
										if !_rules[ruleConditionalExpression]() {
											goto l180
										}
										if !_rules[ruleSEMI]() {
											goto l180
										}
										if !_rules[ruleCommand]() {
											goto l180
										}
										add(ruleLoopConditionBounded, position181)
									}
									if !_rules[ruleOPEN]() {
										goto l180
									}
								l182:
									{
										position183, tokenIndex183 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l183
										}
										goto l182
									l183:
										position, tokenIndex = position183, tokenIndex183
									}
									if !_rules[ruleCLOSE]() {
										goto l180
									}
									goto l160
								l180:
									position, tokenIndex = position160, tokenIndex160
									{
										position184 := position
										if !_rules[ruleConditionalExpression]() {
											goto l157
										}
										add(ruleLoopConditionTruthy, position184)
									}
									if !_rules[ruleOPEN]() {
										goto l157
									}
								l185:
									{
										position186, tokenIndex186 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l186
										}
										goto l185
									l186:
										position, tokenIndex = position186, tokenIndex186
									}
									if !_rules[ruleCLOSE]() {
										goto l157
									}
								}
							l160:
								add(ruleLoop, position158)
							}
							goto l132
						l157:
							position, tokenIndex = position132, tokenIndex132
							if !_rules[ruleCommand]() {
								goto l111
							}
						}
					l132:
						add(ruleStatementBlock, position131)
					}
				}
			l113:
				{
					position187, tokenIndex187 := position, tokenIndex
					if !_rules[ruleSEMI]() {
						goto l187
					}
					goto l188
				l187:
					position, tokenIndex = position187, tokenIndex187
				}
			l188:
				if !_rules[rule_]() {
					goto l111
				}
				add(ruleBlock, position112)
			}
			return true
		l111:
			position, tokenIndex = position111, tokenIndex111
			return false
		},
		/* 65 FlowControlWord <- <(FlowControlBreak / FlowControlContinue)> */
		nil,
		/* 66 FlowControlBreak <- <(BREAK PositiveInteger?)> */
		nil,
		/* 67 FlowControlContinue <- <(CONT PositiveInteger?)> */
		nil,
		/* 68 StatementBlock <- <(NOOP / Assignment / Directive / Conditional / Loop / Command)> */
		nil,
		/* 69 Assignment <- <(AssignmentLHS AssignmentOperator AssignmentRHS)> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				{
					position195 := position
					if !_rules[ruleVariableSequence]() {
						goto l193
					}
					add(ruleAssignmentLHS, position195)
				}
				{
					position196 := position
					if !_rules[rule_]() {
						goto l193
					}
					{
						position197, tokenIndex197 := position, tokenIndex
						{
							position199 := position
							if !_rules[rule_]() {
								goto l198
							}
							if buffer[position] != rune('=') {
								goto l198
							}
							position++
							if !_rules[rule_]() {
								goto l198
							}
							add(ruleAssignEq, position199)
						}
						goto l197
					l198:
						position, tokenIndex = position197, tokenIndex197
						{
							position201 := position
							if !_rules[rule_]() {
								goto l200
							}
							if buffer[position] != rune('*') {
								goto l200
							}
							position++
							if buffer[position] != rune('=') {
								goto l200
							}
							position++
							if !_rules[rule_]() {
								goto l200
							}
							add(ruleStarEq, position201)
						}
						goto l197
					l200:
						position, tokenIndex = position197, tokenIndex197
						{
							position203 := position
							if !_rules[rule_]() {
								goto l202
							}
							if buffer[position] != rune('/') {
								goto l202
							}
							position++
							if buffer[position] != rune('=') {
								goto l202
							}
							position++
							if !_rules[rule_]() {
								goto l202
							}
							add(ruleDivEq, position203)
						}
						goto l197
					l202:
						position, tokenIndex = position197, tokenIndex197
						{
							position205 := position
							if !_rules[rule_]() {
								goto l204
							}
							if buffer[position] != rune('+') {
								goto l204
							}
							position++
							if buffer[position] != rune('=') {
								goto l204
							}
							position++
							if !_rules[rule_]() {
								goto l204
							}
							add(rulePlusEq, position205)
						}
						goto l197
					l204:
						position, tokenIndex = position197, tokenIndex197
						{
							position207 := position
							if !_rules[rule_]() {
								goto l206
							}
							if buffer[position] != rune('-') {
								goto l206
							}
							position++
							if buffer[position] != rune('=') {
								goto l206
							}
							position++
							if !_rules[rule_]() {
								goto l206
							}
							add(ruleMinusEq, position207)
						}
						goto l197
					l206:
						position, tokenIndex = position197, tokenIndex197
						{
							position209 := position
							if !_rules[rule_]() {
								goto l208
							}
							if buffer[position] != rune('&') {
								goto l208
							}
							position++
							if buffer[position] != rune('=') {
								goto l208
							}
							position++
							if !_rules[rule_]() {
								goto l208
							}
							add(ruleAndEq, position209)
						}
						goto l197
					l208:
						position, tokenIndex = position197, tokenIndex197
						{
							position211 := position
							if !_rules[rule_]() {
								goto l210
							}
							if buffer[position] != rune('|') {
								goto l210
							}
							position++
							if buffer[position] != rune('=') {
								goto l210
							}
							position++
							if !_rules[rule_]() {
								goto l210
							}
							add(ruleOrEq, position211)
						}
						goto l197
					l210:
						position, tokenIndex = position197, tokenIndex197
						{
							position212 := position
							if !_rules[rule_]() {
								goto l193
							}
							if buffer[position] != rune('<') {
								goto l193
							}
							position++
							if buffer[position] != rune('<') {
								goto l193
							}
							position++
							if !_rules[rule_]() {
								goto l193
							}
							add(ruleAppend, position212)
						}
					}
				l197:
					if !_rules[rule_]() {
						goto l193
					}
					add(ruleAssignmentOperator, position196)
				}
				{
					position213 := position
					if !_rules[ruleExpressionSequence]() {
						goto l193
					}
					add(ruleAssignmentRHS, position213)
				}
				add(ruleAssignment, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 70 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 71 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 72 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position216, tokenIndex216 := position, tokenIndex
			{
				position217 := position
			l218:
				{
					position219, tokenIndex219 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l219
					}
					if !_rules[ruleCOMMA]() {
						goto l219
					}
					goto l218
				l219:
					position, tokenIndex = position219, tokenIndex219
				}
				if !_rules[ruleVariable]() {
					goto l216
				}
				add(ruleVariableSequence, position217)
			}
			return true
		l216:
			position, tokenIndex = position216, tokenIndex216
			return false
		},
		/* 73 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position220, tokenIndex220 := position, tokenIndex
			{
				position221 := position
			l222:
				{
					position223, tokenIndex223 := position, tokenIndex
					if !_rules[ruleExpression]() {
						goto l223
					}
					if !_rules[ruleCOMMA]() {
						goto l223
					}
					goto l222
				l223:
					position, tokenIndex = position223, tokenIndex223
				}
				if !_rules[ruleExpression]() {
					goto l220
				}
				add(ruleExpressionSequence, position221)
			}
			return true
		l220:
			position, tokenIndex = position220, tokenIndex220
			return false
		},
		/* 74 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
				if !_rules[rule_]() {
					goto l224
				}
				{
					position226 := position
					{
						position227 := position
						{
							position228, tokenIndex228 := position, tokenIndex
							if !_rules[ruleType]() {
								goto l229
							}
							goto l228
						l229:
							position, tokenIndex = position228, tokenIndex228
							if !_rules[ruleVariable]() {
								goto l224
							}
						}
					l228:
						add(ruleValueYielding, position227)
					}
					add(ruleExpressionLHS, position226)
				}
				{
					position230, tokenIndex230 := position, tokenIndex
					{
						position232 := position
						{
							position233 := position
							if !_rules[rule_]() {
								goto l230
							}
							{
								position234, tokenIndex234 := position, tokenIndex
								{
									position236 := position
									if !_rules[rule_]() {
										goto l235
									}
									if buffer[position] != rune('*') {
										goto l235
									}
									position++
									if buffer[position] != rune('*') {
										goto l235
									}
									position++
									if !_rules[rule_]() {
										goto l235
									}
									add(ruleExponentiate, position236)
								}
								goto l234
							l235:
								position, tokenIndex = position234, tokenIndex234
								{
									position238 := position
									if !_rules[rule_]() {
										goto l237
									}
									if buffer[position] != rune('*') {
										goto l237
									}
									position++
									if !_rules[rule_]() {
										goto l237
									}
									add(ruleMultiply, position238)
								}
								goto l234
							l237:
								position, tokenIndex = position234, tokenIndex234
								{
									position240 := position
									if !_rules[rule_]() {
										goto l239
									}
									if buffer[position] != rune('/') {
										goto l239
									}
									position++
									if !_rules[rule_]() {
										goto l239
									}
									add(ruleDivide, position240)
								}
								goto l234
							l239:
								position, tokenIndex = position234, tokenIndex234
								{
									position242 := position
									if !_rules[rule_]() {
										goto l241
									}
									if buffer[position] != rune('%') {
										goto l241
									}
									position++
									if !_rules[rule_]() {
										goto l241
									}
									add(ruleModulus, position242)
								}
								goto l234
							l241:
								position, tokenIndex = position234, tokenIndex234
								{
									position244 := position
									if !_rules[rule_]() {
										goto l243
									}
									if buffer[position] != rune('+') {
										goto l243
									}
									position++
									if !_rules[rule_]() {
										goto l243
									}
									add(ruleAdd, position244)
								}
								goto l234
							l243:
								position, tokenIndex = position234, tokenIndex234
								{
									position246 := position
									if !_rules[rule_]() {
										goto l245
									}
									if buffer[position] != rune('-') {
										goto l245
									}
									position++
									if !_rules[rule_]() {
										goto l245
									}
									add(ruleSubtract, position246)
								}
								goto l234
							l245:
								position, tokenIndex = position234, tokenIndex234
								{
									position248 := position
									if !_rules[rule_]() {
										goto l247
									}
									if buffer[position] != rune('&') {
										goto l247
									}
									position++
									if !_rules[rule_]() {
										goto l247
									}
									add(ruleBitwiseAnd, position248)
								}
								goto l234
							l247:
								position, tokenIndex = position234, tokenIndex234
								{
									position250 := position
									if !_rules[rule_]() {
										goto l249
									}
									if buffer[position] != rune('|') {
										goto l249
									}
									position++
									if !_rules[rule_]() {
										goto l249
									}
									add(ruleBitwiseOr, position250)
								}
								goto l234
							l249:
								position, tokenIndex = position234, tokenIndex234
								{
									position252 := position
									if !_rules[rule_]() {
										goto l251
									}
									if buffer[position] != rune('~') {
										goto l251
									}
									position++
									if !_rules[rule_]() {
										goto l251
									}
									add(ruleBitwiseNot, position252)
								}
								goto l234
							l251:
								position, tokenIndex = position234, tokenIndex234
								{
									position253 := position
									if !_rules[rule_]() {
										goto l230
									}
									if buffer[position] != rune('^') {
										goto l230
									}
									position++
									if !_rules[rule_]() {
										goto l230
									}
									add(ruleBitwiseXor, position253)
								}
							}
						l234:
							if !_rules[rule_]() {
								goto l230
							}
							add(ruleOperator, position233)
						}
						if !_rules[ruleExpression]() {
							goto l230
						}
						add(ruleExpressionRHS, position232)
					}
					goto l231
				l230:
					position, tokenIndex = position230, tokenIndex230
				}
			l231:
				if !_rules[rule_]() {
					goto l224
				}
				add(ruleExpression, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
			return false
		},
		/* 75 ExpressionLHS <- <ValueYielding> */
		nil,
		/* 76 ExpressionRHS <- <(Operator Expression)> */
		nil,
		/* 77 ValueYielding <- <(Type / Variable)> */
		nil,
		/* 78 Directive <- <(DirectiveUnset / DirectiveInclude / DirectiveDeclare)> */
		nil,
		/* 79 DirectiveUnset <- <(UNSET VariableSequence)> */
		nil,
		/* 80 DirectiveInclude <- <(INCLUDE String)> */
		nil,
		/* 81 DirectiveDeclare <- <(DECLARE VariableSequence)> */
		nil,
		/* 82 Command <- <(_ CommandName __ ((CommandFirstArg __ CommandSecondArg) / CommandFirstArg / CommandSecondArg)? (_ CommandResultAssignment)?)> */
		func() bool {
			position261, tokenIndex261 := position, tokenIndex
			{
				position262 := position
				if !_rules[rule_]() {
					goto l261
				}
				{
					position263 := position
					{
						position264, tokenIndex264 := position, tokenIndex
						if !_rules[ruleIdentifier]() {
							goto l264
						}
						{
							position266 := position
							if buffer[position] != rune(':') {
								goto l264
							}
							position++
							if buffer[position] != rune(':') {
								goto l264
							}
							position++
							add(ruleSCOPE, position266)
						}
						goto l265
					l264:
						position, tokenIndex = position264, tokenIndex264
					}
				l265:
					if !_rules[ruleIdentifier]() {
						goto l261
					}
					add(ruleCommandName, position263)
				}
				if !_rules[rule__]() {
					goto l261
				}
				{
					position267, tokenIndex267 := position, tokenIndex
					{
						position269, tokenIndex269 := position, tokenIndex
						if !_rules[ruleCommandFirstArg]() {
							goto l270
						}
						if !_rules[rule__]() {
							goto l270
						}
						if !_rules[ruleCommandSecondArg]() {
							goto l270
						}
						goto l269
					l270:
						position, tokenIndex = position269, tokenIndex269
						if !_rules[ruleCommandFirstArg]() {
							goto l271
						}
						goto l269
					l271:
						position, tokenIndex = position269, tokenIndex269
						if !_rules[ruleCommandSecondArg]() {
							goto l267
						}
					}
				l269:
					goto l268
				l267:
					position, tokenIndex = position267, tokenIndex267
				}
			l268:
				{
					position272, tokenIndex272 := position, tokenIndex
					if !_rules[rule_]() {
						goto l272
					}
					{
						position274 := position
						{
							position275 := position
							if !_rules[rule_]() {
								goto l272
							}
							if buffer[position] != rune('-') {
								goto l272
							}
							position++
							if buffer[position] != rune('>') {
								goto l272
							}
							position++
							if !_rules[rule_]() {
								goto l272
							}
							add(ruleASSIGN, position275)
						}
						if !_rules[ruleVariable]() {
							goto l272
						}
						add(ruleCommandResultAssignment, position274)
					}
					goto l273
				l272:
					position, tokenIndex = position272, tokenIndex272
				}
			l273:
				add(ruleCommand, position262)
			}
			return true
		l261:
			position, tokenIndex = position261, tokenIndex261
			return false
		},
		/* 83 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 84 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position277, tokenIndex277 := position, tokenIndex
			{
				position278 := position
				{
					position279, tokenIndex279 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l280
					}
					goto l279
				l280:
					position, tokenIndex = position279, tokenIndex279
					if !_rules[ruleType]() {
						goto l277
					}
				}
			l279:
				add(ruleCommandFirstArg, position278)
			}
			return true
		l277:
			position, tokenIndex = position277, tokenIndex277
			return false
		},
		/* 85 CommandSecondArg <- <Object> */
		func() bool {
			position281, tokenIndex281 := position, tokenIndex
			{
				position282 := position
				if !_rules[ruleObject]() {
					goto l281
				}
				add(ruleCommandSecondArg, position282)
			}
			return true
		l281:
			position, tokenIndex = position281, tokenIndex281
			return false
		},
		/* 86 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 87 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 88 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position285, tokenIndex285 := position, tokenIndex
			{
				position286 := position
				{
					position287 := position
					if !_rules[rule_]() {
						goto l285
					}
					if buffer[position] != rune('i') {
						goto l285
					}
					position++
					if buffer[position] != rune('f') {
						goto l285
					}
					position++
					if !_rules[rule_]() {
						goto l285
					}
					add(ruleIF, position287)
				}
				if !_rules[ruleConditionalExpression]() {
					goto l285
				}
				if !_rules[ruleOPEN]() {
					goto l285
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
					goto l285
				}
				add(ruleIfStanza, position286)
			}
			return true
		l285:
			position, tokenIndex = position285, tokenIndex285
			return false
		},
		/* 89 ElseIfStanza <- <(ELSE IfStanza)> */
		nil,
		/* 90 ElseStanza <- <(ELSE OPEN Block* CLOSE)> */
		nil,
		/* 91 Loop <- <(LOOP ((OPEN Block* CLOSE) / (LoopConditionFixedLength OPEN Block* CLOSE) / (LoopConditionIterable OPEN Block* CLOSE) / (LoopConditionBounded OPEN Block* CLOSE) / (LoopConditionTruthy OPEN Block* CLOSE)))> */
		nil,
		/* 92 LoopConditionFixedLength <- <(COUNT (Integer / Variable))> */
		nil,
		/* 93 LoopConditionIterable <- <(LoopIterableLHS IN LoopIterableRHS)> */
		nil,
		/* 94 LoopIterableLHS <- <VariableSequence> */
		nil,
		/* 95 LoopIterableRHS <- <(Command / Variable)> */
		nil,
		/* 96 LoopConditionBounded <- <(Command SEMI ConditionalExpression SEMI Command)> */
		nil,
		/* 97 LoopConditionTruthy <- <ConditionalExpression> */
		nil,
		/* 98 ConditionalExpression <- <(NOT? (ConditionWithAssignment / ConditionWithCommand / ConditionWithRegex / ConditionWithComparator))> */
		func() bool {
			position299, tokenIndex299 := position, tokenIndex
			{
				position300 := position
				{
					position301, tokenIndex301 := position, tokenIndex
					{
						position303 := position
						if !_rules[rule_]() {
							goto l301
						}
						if buffer[position] != rune('n') {
							goto l301
						}
						position++
						if buffer[position] != rune('o') {
							goto l301
						}
						position++
						if buffer[position] != rune('t') {
							goto l301
						}
						position++
						if !_rules[rule__]() {
							goto l301
						}
						add(ruleNOT, position303)
					}
					goto l302
				l301:
					position, tokenIndex = position301, tokenIndex301
				}
			l302:
				{
					position304, tokenIndex304 := position, tokenIndex
					{
						position306 := position
						if !_rules[ruleAssignment]() {
							goto l305
						}
						if !_rules[ruleSEMI]() {
							goto l305
						}
						if !_rules[ruleConditionalExpression]() {
							goto l305
						}
						add(ruleConditionWithAssignment, position306)
					}
					goto l304
				l305:
					position, tokenIndex = position304, tokenIndex304
					{
						position308 := position
						if !_rules[ruleCommand]() {
							goto l307
						}
						{
							position309, tokenIndex309 := position, tokenIndex
							if !_rules[ruleSEMI]() {
								goto l309
							}
							if !_rules[ruleConditionalExpression]() {
								goto l309
							}
							goto l310
						l309:
							position, tokenIndex = position309, tokenIndex309
						}
					l310:
						add(ruleConditionWithCommand, position308)
					}
					goto l304
				l307:
					position, tokenIndex = position304, tokenIndex304
					{
						position312 := position
						if !_rules[ruleExpression]() {
							goto l311
						}
						{
							position313 := position
							{
								position314, tokenIndex314 := position, tokenIndex
								{
									position316 := position
									if !_rules[rule_]() {
										goto l315
									}
									if buffer[position] != rune('=') {
										goto l315
									}
									position++
									if buffer[position] != rune('~') {
										goto l315
									}
									position++
									if !_rules[rule_]() {
										goto l315
									}
									add(ruleMatch, position316)
								}
								goto l314
							l315:
								position, tokenIndex = position314, tokenIndex314
								{
									position317 := position
									if !_rules[rule_]() {
										goto l311
									}
									if buffer[position] != rune('!') {
										goto l311
									}
									position++
									if buffer[position] != rune('~') {
										goto l311
									}
									position++
									if !_rules[rule_]() {
										goto l311
									}
									add(ruleUnmatch, position317)
								}
							}
						l314:
							add(ruleMatchOperator, position313)
						}
						if !_rules[ruleRegularExpression]() {
							goto l311
						}
						add(ruleConditionWithRegex, position312)
					}
					goto l304
				l311:
					position, tokenIndex = position304, tokenIndex304
					{
						position318 := position
						{
							position319 := position
							if !_rules[ruleExpression]() {
								goto l299
							}
							add(ruleConditionWithComparatorLHS, position319)
						}
						{
							position320, tokenIndex320 := position, tokenIndex
							{
								position322 := position
								{
									position323 := position
									if !_rules[rule_]() {
										goto l320
									}
									{
										position324, tokenIndex324 := position, tokenIndex
										{
											position326 := position
											if !_rules[rule_]() {
												goto l325
											}
											if buffer[position] != rune('=') {
												goto l325
											}
											position++
											if buffer[position] != rune('=') {
												goto l325
											}
											position++
											if !_rules[rule_]() {
												goto l325
											}
											add(ruleEquality, position326)
										}
										goto l324
									l325:
										position, tokenIndex = position324, tokenIndex324
										{
											position328 := position
											if !_rules[rule_]() {
												goto l327
											}
											if buffer[position] != rune('!') {
												goto l327
											}
											position++
											if buffer[position] != rune('=') {
												goto l327
											}
											position++
											if !_rules[rule_]() {
												goto l327
											}
											add(ruleNonEquality, position328)
										}
										goto l324
									l327:
										position, tokenIndex = position324, tokenIndex324
										{
											position330 := position
											if !_rules[rule_]() {
												goto l329
											}
											if buffer[position] != rune('>') {
												goto l329
											}
											position++
											if buffer[position] != rune('=') {
												goto l329
											}
											position++
											if !_rules[rule_]() {
												goto l329
											}
											add(ruleGreaterEqual, position330)
										}
										goto l324
									l329:
										position, tokenIndex = position324, tokenIndex324
										{
											position332 := position
											if !_rules[rule_]() {
												goto l331
											}
											if buffer[position] != rune('<') {
												goto l331
											}
											position++
											if buffer[position] != rune('=') {
												goto l331
											}
											position++
											if !_rules[rule_]() {
												goto l331
											}
											add(ruleLessEqual, position332)
										}
										goto l324
									l331:
										position, tokenIndex = position324, tokenIndex324
										{
											position334 := position
											if !_rules[rule_]() {
												goto l333
											}
											if buffer[position] != rune('>') {
												goto l333
											}
											position++
											if !_rules[rule_]() {
												goto l333
											}
											add(ruleGreaterThan, position334)
										}
										goto l324
									l333:
										position, tokenIndex = position324, tokenIndex324
										{
											position336 := position
											if !_rules[rule_]() {
												goto l335
											}
											if buffer[position] != rune('<') {
												goto l335
											}
											position++
											if !_rules[rule_]() {
												goto l335
											}
											add(ruleLessThan, position336)
										}
										goto l324
									l335:
										position, tokenIndex = position324, tokenIndex324
										{
											position338 := position
											if !_rules[rule_]() {
												goto l337
											}
											if buffer[position] != rune('i') {
												goto l337
											}
											position++
											if buffer[position] != rune('n') {
												goto l337
											}
											position++
											if !_rules[rule_]() {
												goto l337
											}
											add(ruleMembership, position338)
										}
										goto l324
									l337:
										position, tokenIndex = position324, tokenIndex324
										{
											position339 := position
											if !_rules[rule_]() {
												goto l320
											}
											if buffer[position] != rune('n') {
												goto l320
											}
											position++
											if buffer[position] != rune('o') {
												goto l320
											}
											position++
											if buffer[position] != rune('t') {
												goto l320
											}
											position++
											if !_rules[rule__]() {
												goto l320
											}
											if buffer[position] != rune('i') {
												goto l320
											}
											position++
											if buffer[position] != rune('n') {
												goto l320
											}
											position++
											if !_rules[rule_]() {
												goto l320
											}
											add(ruleNonMembership, position339)
										}
									}
								l324:
									if !_rules[rule_]() {
										goto l320
									}
									add(ruleComparisonOperator, position323)
								}
								if !_rules[ruleExpression]() {
									goto l320
								}
								add(ruleConditionWithComparatorRHS, position322)
							}
							goto l321
						l320:
							position, tokenIndex = position320, tokenIndex320
						}
					l321:
						add(ruleConditionWithComparator, position318)
					}
				}
			l304:
				add(ruleConditionalExpression, position300)
			}
			return true
		l299:
			position, tokenIndex = position299, tokenIndex299
			return false
		},
		/* 99 ConditionWithAssignment <- <(Assignment SEMI ConditionalExpression)> */
		nil,
		/* 100 ConditionWithCommand <- <(Command (SEMI ConditionalExpression)?)> */
		nil,
		/* 101 ConditionWithRegex <- <(Expression MatchOperator RegularExpression)> */
		nil,
		/* 102 ConditionWithComparator <- <(ConditionWithComparatorLHS ConditionWithComparatorRHS?)> */
		nil,
		/* 103 ConditionWithComparatorLHS <- <Expression> */
		nil,
		/* 104 ConditionWithComparatorRHS <- <(ComparisonOperator Expression)> */
		nil,
		/* 105 ScalarType <- <(Boolean / Float / Integer / String / NullValue)> */
		nil,
		/* 106 Identifier <- <(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / ([0-9] / [0-9]) / '_')*)> */
		func() bool {
			position347, tokenIndex347 := position, tokenIndex
			{
				position348 := position
				{
					position349, tokenIndex349 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex = position349, tokenIndex349
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l351
					}
					position++
					goto l349
				l351:
					position, tokenIndex = position349, tokenIndex349
					if buffer[position] != rune('_') {
						goto l347
					}
					position++
				}
			l349:
			l352:
				{
					position353, tokenIndex353 := position, tokenIndex
					{
						position354, tokenIndex354 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l355
						}
						position++
						goto l354
					l355:
						position, tokenIndex = position354, tokenIndex354
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l356
						}
						position++
						goto l354
					l356:
						position, tokenIndex = position354, tokenIndex354
						{
							position358, tokenIndex358 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex = position358, tokenIndex358
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l357
							}
							position++
						}
					l358:
						goto l354
					l357:
						position, tokenIndex = position354, tokenIndex354
						if buffer[position] != rune('_') {
							goto l353
						}
						position++
					}
				l354:
					goto l352
				l353:
					position, tokenIndex = position353, tokenIndex353
				}
				add(ruleIdentifier, position348)
			}
			return true
		l347:
			position, tokenIndex = position347, tokenIndex347
			return false
		},
		/* 107 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 108 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 109 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position362, tokenIndex362 := position, tokenIndex
			{
				position363 := position
				{
					position364, tokenIndex364 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l364
					}
					position++
					goto l365
				l364:
					position, tokenIndex = position364, tokenIndex364
				}
			l365:
				if !_rules[rulePositiveInteger]() {
					goto l362
				}
				add(ruleInteger, position363)
			}
			return true
		l362:
			position, tokenIndex = position362, tokenIndex362
			return false
		},
		/* 110 PositiveInteger <- <[0-9]+> */
		func() bool {
			position366, tokenIndex366 := position, tokenIndex
			{
				position367 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l366
				}
				position++
			l368:
				{
					position369, tokenIndex369 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex = position369, tokenIndex369
				}
				add(rulePositiveInteger, position367)
			}
			return true
		l366:
			position, tokenIndex = position366, tokenIndex366
			return false
		},
		/* 111 String <- <(StringLiteral / StringInterpolated / Heredoc)> */
		func() bool {
			position370, tokenIndex370 := position, tokenIndex
			{
				position371 := position
				{
					position372, tokenIndex372 := position, tokenIndex
					{
						position374 := position
						if buffer[position] != rune('\'') {
							goto l373
						}
						position++
					l375:
						{
							position376, tokenIndex376 := position, tokenIndex
							{
								position377, tokenIndex377 := position, tokenIndex
								if buffer[position] != rune('\'') {
									goto l377
								}
								position++
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
						if buffer[position] != rune('\'') {
							goto l373
						}
						position++
						add(ruleStringLiteral, position374)
					}
					goto l372
				l373:
					position, tokenIndex = position372, tokenIndex372
					{
						position379 := position
						if buffer[position] != rune('"') {
							goto l378
						}
						position++
					l380:
						{
							position381, tokenIndex381 := position, tokenIndex
							{
								position382, tokenIndex382 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l382
								}
								position++
								goto l381
							l382:
								position, tokenIndex = position382, tokenIndex382
							}
							if !matchDot() {
								goto l381
							}
							goto l380
						l381:
							position, tokenIndex = position381, tokenIndex381
						}
						if buffer[position] != rune('"') {
							goto l378
						}
						position++
						add(ruleStringInterpolated, position379)
					}
					goto l372
				l378:
					position, tokenIndex = position372, tokenIndex372
					{
						position383 := position
						{
							position384 := position
							if !_rules[rule_]() {
								goto l370
							}
							if buffer[position] != rune('b') {
								goto l370
							}
							position++
							if buffer[position] != rune('e') {
								goto l370
							}
							position++
							if buffer[position] != rune('g') {
								goto l370
							}
							position++
							if buffer[position] != rune('i') {
								goto l370
							}
							position++
							if buffer[position] != rune('n') {
								goto l370
							}
							position++
							add(ruleBEGIN, position384)
						}
						if buffer[position] != rune('\n') {
							goto l370
						}
						position++
						{
							position385 := position
						l386:
							{
								position387, tokenIndex387 := position, tokenIndex
								{
									position388, tokenIndex388 := position, tokenIndex
									if !_rules[ruleEND]() {
										goto l388
									}
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
							add(ruleHeredocBody, position385)
						}
						if !_rules[ruleEND]() {
							goto l370
						}
						add(ruleHeredoc, position383)
					}
				}
			l372:
				add(ruleString, position371)
			}
			return true
		l370:
			position, tokenIndex = position370, tokenIndex370
			return false
		},
		/* 112 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		nil,
		/* 113 StringInterpolated <- <('"' (!'"' .)* '"')> */
		nil,
		/* 114 Heredoc <- <(BEGIN '\n' HeredocBody END)> */
		nil,
		/* 115 HeredocBody <- <(!END .)*> */
		nil,
		/* 116 NullValue <- <('n' 'u' 'l' 'l')> */
		nil,
		/* 117 Object <- <(OPEN (_ KeyValuePair _)* CLOSE)> */
		func() bool {
			position394, tokenIndex394 := position, tokenIndex
			{
				position395 := position
				if !_rules[ruleOPEN]() {
					goto l394
				}
			l396:
				{
					position397, tokenIndex397 := position, tokenIndex
					if !_rules[rule_]() {
						goto l397
					}
					{
						position398 := position
						{
							position399 := position
							if !_rules[ruleIdentifier]() {
								goto l397
							}
							add(ruleKey, position399)
						}
						{
							position400 := position
							if !_rules[rule_]() {
								goto l397
							}
							if buffer[position] != rune(':') {
								goto l397
							}
							position++
							if !_rules[rule_]() {
								goto l397
							}
							add(ruleCOLON, position400)
						}
						{
							position401 := position
							{
								position402, tokenIndex402 := position, tokenIndex
								if !_rules[ruleArray]() {
									goto l403
								}
								goto l402
							l403:
								position, tokenIndex = position402, tokenIndex402
								if !_rules[ruleObject]() {
									goto l404
								}
								goto l402
							l404:
								position, tokenIndex = position402, tokenIndex402
								if !_rules[ruleExpression]() {
									goto l397
								}
							}
						l402:
							add(ruleKValue, position401)
						}
						{
							position405, tokenIndex405 := position, tokenIndex
							if !_rules[ruleCOMMA]() {
								goto l405
							}
							goto l406
						l405:
							position, tokenIndex = position405, tokenIndex405
						}
					l406:
						add(ruleKeyValuePair, position398)
					}
					if !_rules[rule_]() {
						goto l397
					}
					goto l396
				l397:
					position, tokenIndex = position397, tokenIndex397
				}
				if !_rules[ruleCLOSE]() {
					goto l394
				}
				add(ruleObject, position395)
			}
			return true
		l394:
			position, tokenIndex = position394, tokenIndex394
			return false
		},
		/* 118 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position407, tokenIndex407 := position, tokenIndex
			{
				position408 := position
				if buffer[position] != rune('[') {
					goto l407
				}
				position++
				if !_rules[rule_]() {
					goto l407
				}
				if !_rules[ruleExpressionSequence]() {
					goto l407
				}
				{
					position409, tokenIndex409 := position, tokenIndex
					if !_rules[ruleCOMMA]() {
						goto l409
					}
					goto l410
				l409:
					position, tokenIndex = position409, tokenIndex409
				}
			l410:
				if buffer[position] != rune(']') {
					goto l407
				}
				position++
				add(ruleArray, position408)
			}
			return true
		l407:
			position, tokenIndex = position407, tokenIndex407
			return false
		},
		/* 119 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position411, tokenIndex411 := position, tokenIndex
			{
				position412 := position
				if buffer[position] != rune('/') {
					goto l411
				}
				position++
				{
					position415, tokenIndex415 := position, tokenIndex
					if buffer[position] != rune('/') {
						goto l415
					}
					position++
					goto l411
				l415:
					position, tokenIndex = position415, tokenIndex415
				}
				if !matchDot() {
					goto l411
				}
			l413:
				{
					position414, tokenIndex414 := position, tokenIndex
					{
						position416, tokenIndex416 := position, tokenIndex
						if buffer[position] != rune('/') {
							goto l416
						}
						position++
						goto l414
					l416:
						position, tokenIndex = position416, tokenIndex416
					}
					if !matchDot() {
						goto l414
					}
					goto l413
				l414:
					position, tokenIndex = position414, tokenIndex414
				}
				if buffer[position] != rune('/') {
					goto l411
				}
				position++
			l417:
				{
					position418, tokenIndex418 := position, tokenIndex
					{
						position419, tokenIndex419 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l420
						}
						position++
						goto l419
					l420:
						position, tokenIndex = position419, tokenIndex419
						if buffer[position] != rune('l') {
							goto l421
						}
						position++
						goto l419
					l421:
						position, tokenIndex = position419, tokenIndex419
						if buffer[position] != rune('m') {
							goto l422
						}
						position++
						goto l419
					l422:
						position, tokenIndex = position419, tokenIndex419
						if buffer[position] != rune('s') {
							goto l423
						}
						position++
						goto l419
					l423:
						position, tokenIndex = position419, tokenIndex419
						if buffer[position] != rune('u') {
							goto l418
						}
						position++
					}
				l419:
					goto l417
				l418:
					position, tokenIndex = position418, tokenIndex418
				}
				add(ruleRegularExpression, position412)
			}
			return true
		l411:
			position, tokenIndex = position411, tokenIndex411
			return false
		},
		/* 120 KeyValuePair <- <(Key COLON KValue COMMA?)> */
		nil,
		/* 121 Key <- <Identifier> */
		nil,
		/* 122 KValue <- <(Array / Object / Expression)> */
		nil,
		/* 123 Type <- <(Array / Object / RegularExpression / ScalarType)> */
		func() bool {
			position427, tokenIndex427 := position, tokenIndex
			{
				position428 := position
				{
					position429, tokenIndex429 := position, tokenIndex
					if !_rules[ruleArray]() {
						goto l430
					}
					goto l429
				l430:
					position, tokenIndex = position429, tokenIndex429
					if !_rules[ruleObject]() {
						goto l431
					}
					goto l429
				l431:
					position, tokenIndex = position429, tokenIndex429
					if !_rules[ruleRegularExpression]() {
						goto l432
					}
					goto l429
				l432:
					position, tokenIndex = position429, tokenIndex429
					{
						position433 := position
						{
							position434, tokenIndex434 := position, tokenIndex
							{
								position436 := position
								{
									position437, tokenIndex437 := position, tokenIndex
									if buffer[position] != rune('t') {
										goto l438
									}
									position++
									if buffer[position] != rune('r') {
										goto l438
									}
									position++
									if buffer[position] != rune('u') {
										goto l438
									}
									position++
									if buffer[position] != rune('e') {
										goto l438
									}
									position++
									goto l437
								l438:
									position, tokenIndex = position437, tokenIndex437
									if buffer[position] != rune('f') {
										goto l435
									}
									position++
									if buffer[position] != rune('a') {
										goto l435
									}
									position++
									if buffer[position] != rune('l') {
										goto l435
									}
									position++
									if buffer[position] != rune('s') {
										goto l435
									}
									position++
									if buffer[position] != rune('e') {
										goto l435
									}
									position++
								}
							l437:
								add(ruleBoolean, position436)
							}
							goto l434
						l435:
							position, tokenIndex = position434, tokenIndex434
							{
								position440 := position
								if !_rules[ruleInteger]() {
									goto l439
								}
								{
									position441, tokenIndex441 := position, tokenIndex
									if buffer[position] != rune('.') {
										goto l441
									}
									position++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l441
									}
									position++
								l443:
									{
										position444, tokenIndex444 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l444
										}
										position++
										goto l443
									l444:
										position, tokenIndex = position444, tokenIndex444
									}
									goto l442
								l441:
									position, tokenIndex = position441, tokenIndex441
								}
							l442:
								add(ruleFloat, position440)
							}
							goto l434
						l439:
							position, tokenIndex = position434, tokenIndex434
							if !_rules[ruleInteger]() {
								goto l445
							}
							goto l434
						l445:
							position, tokenIndex = position434, tokenIndex434
							if !_rules[ruleString]() {
								goto l446
							}
							goto l434
						l446:
							position, tokenIndex = position434, tokenIndex434
							{
								position447 := position
								if buffer[position] != rune('n') {
									goto l427
								}
								position++
								if buffer[position] != rune('u') {
									goto l427
								}
								position++
								if buffer[position] != rune('l') {
									goto l427
								}
								position++
								if buffer[position] != rune('l') {
									goto l427
								}
								position++
								add(ruleNullValue, position447)
							}
						}
					l434:
						add(ruleScalarType, position433)
					}
				}
			l429:
				add(ruleType, position428)
			}
			return true
		l427:
			position, tokenIndex = position427, tokenIndex427
			return false
		},
	}
	p.rules = _rules
}
