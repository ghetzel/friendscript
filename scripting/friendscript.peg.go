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
	ruleNL
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
	"NL",
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
	rules  [126]func() bool
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
		/* 1 NL <- <'\n'> */
		nil,
		/* 2 _ <- <(' ' / '\t' / '\r' / '\n')*> */
		func() bool {
			{
				position14 := position
			l15:
				{
					position16, tokenIndex16 := position, tokenIndex
					{
						position17, tokenIndex17 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l18
						}
						position++
						goto l17
					l18:
						position, tokenIndex = position17, tokenIndex17
						if buffer[position] != rune('\t') {
							goto l19
						}
						position++
						goto l17
					l19:
						position, tokenIndex = position17, tokenIndex17
						if buffer[position] != rune('\r') {
							goto l20
						}
						position++
						goto l17
					l20:
						position, tokenIndex = position17, tokenIndex17
						if buffer[position] != rune('\n') {
							goto l16
						}
						position++
					}
				l17:
					goto l15
				l16:
					position, tokenIndex = position16, tokenIndex16
				}
				add(rule_, position14)
			}
			return true
		},
		/* 3 __ <- <(' ' / '\t' / '\r' / '\n')+> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				{
					position25, tokenIndex25 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex = position25, tokenIndex25
					if buffer[position] != rune('\t') {
						goto l27
					}
					position++
					goto l25
				l27:
					position, tokenIndex = position25, tokenIndex25
					if buffer[position] != rune('\r') {
						goto l28
					}
					position++
					goto l25
				l28:
					position, tokenIndex = position25, tokenIndex25
					if buffer[position] != rune('\n') {
						goto l21
					}
					position++
				}
			l25:
			l23:
				{
					position24, tokenIndex24 := position, tokenIndex
					{
						position29, tokenIndex29 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l30
						}
						position++
						goto l29
					l30:
						position, tokenIndex = position29, tokenIndex29
						if buffer[position] != rune('\t') {
							goto l31
						}
						position++
						goto l29
					l31:
						position, tokenIndex = position29, tokenIndex29
						if buffer[position] != rune('\r') {
							goto l32
						}
						position++
						goto l29
					l32:
						position, tokenIndex = position29, tokenIndex29
						if buffer[position] != rune('\n') {
							goto l24
						}
						position++
					}
				l29:
					goto l23
				l24:
					position, tokenIndex = position24, tokenIndex24
				}
				add(rule__, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 4 ASSIGN <- <(_ ('-' '>') _)> */
		nil,
		/* 5 BEGIN <- <(_ ('b' 'e' 'g' 'i' 'n'))> */
		nil,
		/* 6 BREAK <- <(_ ('b' 'r' 'e' 'a' 'k') _)> */
		nil,
		/* 7 CLOSE <- <(_ '}' _)> */
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
		/* 8 COLON <- <(_ ':' _)> */
		nil,
		/* 9 COMMA <- <(_ ',' _)> */
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
		/* 10 Comment <- <(_ '#' (!'\n' .)*)> */
		nil,
		/* 11 CONT <- <(_ ('c' 'o' 'n' 't' 'i' 'n' 'u' 'e') _)> */
		nil,
		/* 12 COUNT <- <(_ ('c' 'o' 'u' 'n' 't') _)> */
		nil,
		/* 13 DECLARE <- <(_ ('d' 'e' 'c' 'l' 'a' 'r' 'e') __)> */
		nil,
		/* 14 DOT <- <'.'> */
		nil,
		/* 15 ELSE <- <(_ ('e' 'l' 's' 'e') _)> */
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
		/* 16 END <- <(_ ('e' 'n' 'd') _)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if !_rules[rule_]() {
					goto l48
				}
				if buffer[position] != rune('e') {
					goto l48
				}
				position++
				if buffer[position] != rune('n') {
					goto l48
				}
				position++
				if buffer[position] != rune('d') {
					goto l48
				}
				position++
				if !_rules[rule_]() {
					goto l48
				}
				add(ruleEND, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 17 IF <- <(_ ('i' 'f') _)> */
		nil,
		/* 18 IN <- <(__ ('i' 'n') __)> */
		nil,
		/* 19 INCLUDE <- <(_ ('i' 'n' 'c' 'l' 'u' 'd' 'e') __)> */
		nil,
		/* 20 LOOP <- <(_ ('l' 'o' 'o' 'p') _)> */
		nil,
		/* 21 NOOP <- <SEMI> */
		nil,
		/* 22 NOT <- <(_ ('n' 'o' 't') __)> */
		nil,
		/* 23 OPEN <- <(_ '{' _)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				if !_rules[rule_]() {
					goto l56
				}
				if buffer[position] != rune('{') {
					goto l56
				}
				position++
				if !_rules[rule_]() {
					goto l56
				}
				add(ruleOPEN, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 24 SCOPE <- <(':' ':')> */
		nil,
		/* 25 SEMI <- <(_ ';' _)> */
		func() bool {
			position59, tokenIndex59 := position, tokenIndex
			{
				position60 := position
				if !_rules[rule_]() {
					goto l59
				}
				if buffer[position] != rune(';') {
					goto l59
				}
				position++
				if !_rules[rule_]() {
					goto l59
				}
				add(ruleSEMI, position60)
			}
			return true
		l59:
			position, tokenIndex = position59, tokenIndex59
			return false
		},
		/* 26 SHEBANG <- <('#' '!' (!'\n' .)+ '\n')> */
		nil,
		/* 27 SKIPVAR <- <(_ '_' _)> */
		nil,
		/* 28 UNSET <- <(_ ('u' 'n' 's' 'e' 't') __)> */
		nil,
		/* 29 Operator <- <(_ (Exponentiate / Multiply / Divide / Modulus / Add / Subtract / BitwiseAnd / BitwiseOr / BitwiseNot / BitwiseXor) _)> */
		nil,
		/* 30 Exponentiate <- <(_ ('*' '*') _)> */
		nil,
		/* 31 Multiply <- <(_ '*' _)> */
		nil,
		/* 32 Divide <- <(_ '/' _)> */
		nil,
		/* 33 Modulus <- <(_ '%' _)> */
		nil,
		/* 34 Add <- <(_ '+' _)> */
		nil,
		/* 35 Subtract <- <(_ '-' _)> */
		nil,
		/* 36 BitwiseAnd <- <(_ '&' _)> */
		nil,
		/* 37 BitwiseOr <- <(_ '|' _)> */
		nil,
		/* 38 BitwiseNot <- <(_ '~' _)> */
		nil,
		/* 39 BitwiseXor <- <(_ '^' _)> */
		nil,
		/* 40 MatchOperator <- <(Match / Unmatch)> */
		nil,
		/* 41 Unmatch <- <(_ ('!' '~') _)> */
		nil,
		/* 42 Match <- <(_ ('=' '~') _)> */
		nil,
		/* 43 AssignmentOperator <- <(_ (AssignEq / StarEq / DivEq / PlusEq / MinusEq / AndEq / OrEq / Append) _)> */
		nil,
		/* 44 AssignEq <- <(_ '=' _)> */
		nil,
		/* 45 StarEq <- <(_ ('*' '=') _)> */
		nil,
		/* 46 DivEq <- <(_ ('/' '=') _)> */
		nil,
		/* 47 PlusEq <- <(_ ('+' '=') _)> */
		nil,
		/* 48 MinusEq <- <(_ ('-' '=') _)> */
		nil,
		/* 49 AndEq <- <(_ ('&' '=') _)> */
		nil,
		/* 50 OrEq <- <(_ ('|' '=') _)> */
		nil,
		/* 51 Append <- <(_ ('<' '<') _)> */
		nil,
		/* 52 ComparisonOperator <- <(_ (Equality / NonEquality / GreaterEqual / LessEqual / GreaterThan / LessThan / Membership / NonMembership) _)> */
		nil,
		/* 53 Equality <- <(_ ('=' '=') _)> */
		nil,
		/* 54 NonEquality <- <(_ ('!' '=') _)> */
		nil,
		/* 55 GreaterThan <- <(_ '>' _)> */
		nil,
		/* 56 GreaterEqual <- <(_ ('>' '=') _)> */
		nil,
		/* 57 LessEqual <- <(_ ('<' '=') _)> */
		nil,
		/* 58 LessThan <- <(_ '<' _)> */
		nil,
		/* 59 Membership <- <(_ ('i' 'n') _)> */
		nil,
		/* 60 NonMembership <- <(_ ('n' 'o' 't') __ ('i' 'n') _)> */
		nil,
		/* 61 Variable <- <(('$' VariableNameSequence) / SKIPVAR)> */
		func() bool {
			position96, tokenIndex96 := position, tokenIndex
			{
				position97 := position
				{
					position98, tokenIndex98 := position, tokenIndex
					if buffer[position] != rune('$') {
						goto l99
					}
					position++
					{
						position100 := position
					l101:
						{
							position102, tokenIndex102 := position, tokenIndex
							if !_rules[ruleVariableName]() {
								goto l102
							}
							{
								position103 := position
								if buffer[position] != rune('.') {
									goto l102
								}
								position++
								add(ruleDOT, position103)
							}
							goto l101
						l102:
							position, tokenIndex = position102, tokenIndex102
						}
						if !_rules[ruleVariableName]() {
							goto l99
						}
						add(ruleVariableNameSequence, position100)
					}
					goto l98
				l99:
					position, tokenIndex = position98, tokenIndex98
					{
						position104 := position
						if !_rules[rule_]() {
							goto l96
						}
						if buffer[position] != rune('_') {
							goto l96
						}
						position++
						if !_rules[rule_]() {
							goto l96
						}
						add(ruleSKIPVAR, position104)
					}
				}
			l98:
				add(ruleVariable, position97)
			}
			return true
		l96:
			position, tokenIndex = position96, tokenIndex96
			return false
		},
		/* 62 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 63 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position106, tokenIndex106 := position, tokenIndex
			{
				position107 := position
				if !_rules[ruleIdentifier]() {
					goto l106
				}
				{
					position108, tokenIndex108 := position, tokenIndex
					if buffer[position] != rune('[') {
						goto l108
					}
					position++
					if !_rules[rule_]() {
						goto l108
					}
					{
						position110 := position
						if !_rules[ruleExpression]() {
							goto l108
						}
						add(ruleVariableIndex, position110)
					}
					if !_rules[rule_]() {
						goto l108
					}
					if buffer[position] != rune(']') {
						goto l108
					}
					position++
					goto l109
				l108:
					position, tokenIndex = position108, tokenIndex108
				}
			l109:
				add(ruleVariableName, position107)
			}
			return true
		l106:
			position, tokenIndex = position106, tokenIndex106
			return false
		},
		/* 64 VariableIndex <- <Expression> */
		nil,
		/* 65 Block <- <(_ (Comment / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position112, tokenIndex112 := position, tokenIndex
			{
				position113 := position
				if !_rules[rule_]() {
					goto l112
				}
				{
					position114, tokenIndex114 := position, tokenIndex
					{
						position116 := position
						if !_rules[rule_]() {
							goto l115
						}
						if buffer[position] != rune('#') {
							goto l115
						}
						position++
					l117:
						{
							position118, tokenIndex118 := position, tokenIndex
							{
								position119, tokenIndex119 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l119
								}
								position++
								goto l118
							l119:
								position, tokenIndex = position119, tokenIndex119
							}
							if !matchDot() {
								goto l118
							}
							goto l117
						l118:
							position, tokenIndex = position118, tokenIndex118
						}
						add(ruleComment, position116)
					}
					goto l114
				l115:
					position, tokenIndex = position114, tokenIndex114
					{
						position121 := position
						{
							position122, tokenIndex122 := position, tokenIndex
							{
								position124 := position
								{
									position125 := position
									if !_rules[rule_]() {
										goto l123
									}
									if buffer[position] != rune('b') {
										goto l123
									}
									position++
									if buffer[position] != rune('r') {
										goto l123
									}
									position++
									if buffer[position] != rune('e') {
										goto l123
									}
									position++
									if buffer[position] != rune('a') {
										goto l123
									}
									position++
									if buffer[position] != rune('k') {
										goto l123
									}
									position++
									if !_rules[rule_]() {
										goto l123
									}
									add(ruleBREAK, position125)
								}
								{
									position126, tokenIndex126 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l126
									}
									goto l127
								l126:
									position, tokenIndex = position126, tokenIndex126
								}
							l127:
								add(ruleFlowControlBreak, position124)
							}
							goto l122
						l123:
							position, tokenIndex = position122, tokenIndex122
							{
								position128 := position
								{
									position129 := position
									if !_rules[rule_]() {
										goto l120
									}
									if buffer[position] != rune('c') {
										goto l120
									}
									position++
									if buffer[position] != rune('o') {
										goto l120
									}
									position++
									if buffer[position] != rune('n') {
										goto l120
									}
									position++
									if buffer[position] != rune('t') {
										goto l120
									}
									position++
									if buffer[position] != rune('i') {
										goto l120
									}
									position++
									if buffer[position] != rune('n') {
										goto l120
									}
									position++
									if buffer[position] != rune('u') {
										goto l120
									}
									position++
									if buffer[position] != rune('e') {
										goto l120
									}
									position++
									if !_rules[rule_]() {
										goto l120
									}
									add(ruleCONT, position129)
								}
								{
									position130, tokenIndex130 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l130
									}
									goto l131
								l130:
									position, tokenIndex = position130, tokenIndex130
								}
							l131:
								add(ruleFlowControlContinue, position128)
							}
						}
					l122:
						add(ruleFlowControlWord, position121)
					}
					goto l114
				l120:
					position, tokenIndex = position114, tokenIndex114
					{
						position132 := position
						{
							position133, tokenIndex133 := position, tokenIndex
							{
								position135 := position
								if !_rules[ruleSEMI]() {
									goto l134
								}
								add(ruleNOOP, position135)
							}
							goto l133
						l134:
							position, tokenIndex = position133, tokenIndex133
							if !_rules[ruleAssignment]() {
								goto l136
							}
							goto l133
						l136:
							position, tokenIndex = position133, tokenIndex133
							{
								position138 := position
								{
									position139, tokenIndex139 := position, tokenIndex
									{
										position141 := position
										{
											position142 := position
											if !_rules[rule_]() {
												goto l140
											}
											if buffer[position] != rune('u') {
												goto l140
											}
											position++
											if buffer[position] != rune('n') {
												goto l140
											}
											position++
											if buffer[position] != rune('s') {
												goto l140
											}
											position++
											if buffer[position] != rune('e') {
												goto l140
											}
											position++
											if buffer[position] != rune('t') {
												goto l140
											}
											position++
											if !_rules[rule__]() {
												goto l140
											}
											add(ruleUNSET, position142)
										}
										if !_rules[ruleVariableSequence]() {
											goto l140
										}
										add(ruleDirectiveUnset, position141)
									}
									goto l139
								l140:
									position, tokenIndex = position139, tokenIndex139
									{
										position144 := position
										{
											position145 := position
											if !_rules[rule_]() {
												goto l143
											}
											if buffer[position] != rune('i') {
												goto l143
											}
											position++
											if buffer[position] != rune('n') {
												goto l143
											}
											position++
											if buffer[position] != rune('c') {
												goto l143
											}
											position++
											if buffer[position] != rune('l') {
												goto l143
											}
											position++
											if buffer[position] != rune('u') {
												goto l143
											}
											position++
											if buffer[position] != rune('d') {
												goto l143
											}
											position++
											if buffer[position] != rune('e') {
												goto l143
											}
											position++
											if !_rules[rule__]() {
												goto l143
											}
											add(ruleINCLUDE, position145)
										}
										if !_rules[ruleString]() {
											goto l143
										}
										add(ruleDirectiveInclude, position144)
									}
									goto l139
								l143:
									position, tokenIndex = position139, tokenIndex139
									{
										position146 := position
										{
											position147 := position
											if !_rules[rule_]() {
												goto l137
											}
											if buffer[position] != rune('d') {
												goto l137
											}
											position++
											if buffer[position] != rune('e') {
												goto l137
											}
											position++
											if buffer[position] != rune('c') {
												goto l137
											}
											position++
											if buffer[position] != rune('l') {
												goto l137
											}
											position++
											if buffer[position] != rune('a') {
												goto l137
											}
											position++
											if buffer[position] != rune('r') {
												goto l137
											}
											position++
											if buffer[position] != rune('e') {
												goto l137
											}
											position++
											if !_rules[rule__]() {
												goto l137
											}
											add(ruleDECLARE, position147)
										}
										if !_rules[ruleVariableSequence]() {
											goto l137
										}
										add(ruleDirectiveDeclare, position146)
									}
								}
							l139:
								add(ruleDirective, position138)
							}
							goto l133
						l137:
							position, tokenIndex = position133, tokenIndex133
							{
								position149 := position
								if !_rules[ruleIfStanza]() {
									goto l148
								}
							l150:
								{
									position151, tokenIndex151 := position, tokenIndex
									{
										position152 := position
										if !_rules[ruleELSE]() {
											goto l151
										}
										if !_rules[ruleIfStanza]() {
											goto l151
										}
										add(ruleElseIfStanza, position152)
									}
									goto l150
								l151:
									position, tokenIndex = position151, tokenIndex151
								}
								{
									position153, tokenIndex153 := position, tokenIndex
									{
										position155 := position
										if !_rules[ruleELSE]() {
											goto l153
										}
										if !_rules[ruleOPEN]() {
											goto l153
										}
									l156:
										{
											position157, tokenIndex157 := position, tokenIndex
											if !_rules[ruleBlock]() {
												goto l157
											}
											goto l156
										l157:
											position, tokenIndex = position157, tokenIndex157
										}
										if !_rules[ruleCLOSE]() {
											goto l153
										}
										add(ruleElseStanza, position155)
									}
									goto l154
								l153:
									position, tokenIndex = position153, tokenIndex153
								}
							l154:
								add(ruleConditional, position149)
							}
							goto l133
						l148:
							position, tokenIndex = position133, tokenIndex133
							{
								position159 := position
								{
									position160 := position
									if !_rules[rule_]() {
										goto l158
									}
									if buffer[position] != rune('l') {
										goto l158
									}
									position++
									if buffer[position] != rune('o') {
										goto l158
									}
									position++
									if buffer[position] != rune('o') {
										goto l158
									}
									position++
									if buffer[position] != rune('p') {
										goto l158
									}
									position++
									if !_rules[rule_]() {
										goto l158
									}
									add(ruleLOOP, position160)
								}
								{
									position161, tokenIndex161 := position, tokenIndex
									if !_rules[ruleOPEN]() {
										goto l162
									}
								l163:
									{
										position164, tokenIndex164 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l164
										}
										goto l163
									l164:
										position, tokenIndex = position164, tokenIndex164
									}
									if !_rules[ruleCLOSE]() {
										goto l162
									}
									goto l161
								l162:
									position, tokenIndex = position161, tokenIndex161
									{
										position166 := position
										{
											position167 := position
											if !_rules[rule_]() {
												goto l165
											}
											if buffer[position] != rune('c') {
												goto l165
											}
											position++
											if buffer[position] != rune('o') {
												goto l165
											}
											position++
											if buffer[position] != rune('u') {
												goto l165
											}
											position++
											if buffer[position] != rune('n') {
												goto l165
											}
											position++
											if buffer[position] != rune('t') {
												goto l165
											}
											position++
											if !_rules[rule_]() {
												goto l165
											}
											add(ruleCOUNT, position167)
										}
										{
											position168, tokenIndex168 := position, tokenIndex
											if !_rules[ruleInteger]() {
												goto l169
											}
											goto l168
										l169:
											position, tokenIndex = position168, tokenIndex168
											if !_rules[ruleVariable]() {
												goto l165
											}
										}
									l168:
										add(ruleLoopConditionFixedLength, position166)
									}
									if !_rules[ruleOPEN]() {
										goto l165
									}
								l170:
									{
										position171, tokenIndex171 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l171
										}
										goto l170
									l171:
										position, tokenIndex = position171, tokenIndex171
									}
									if !_rules[ruleCLOSE]() {
										goto l165
									}
									goto l161
								l165:
									position, tokenIndex = position161, tokenIndex161
									{
										position173 := position
										{
											position174 := position
											if !_rules[ruleVariableSequence]() {
												goto l172
											}
											add(ruleLoopIterableLHS, position174)
										}
										{
											position175 := position
											if !_rules[rule__]() {
												goto l172
											}
											if buffer[position] != rune('i') {
												goto l172
											}
											position++
											if buffer[position] != rune('n') {
												goto l172
											}
											position++
											if !_rules[rule__]() {
												goto l172
											}
											add(ruleIN, position175)
										}
										{
											position176 := position
											{
												position177, tokenIndex177 := position, tokenIndex
												if !_rules[ruleCommand]() {
													goto l178
												}
												goto l177
											l178:
												position, tokenIndex = position177, tokenIndex177
												if !_rules[ruleVariable]() {
													goto l172
												}
											}
										l177:
											add(ruleLoopIterableRHS, position176)
										}
										add(ruleLoopConditionIterable, position173)
									}
									if !_rules[ruleOPEN]() {
										goto l172
									}
								l179:
									{
										position180, tokenIndex180 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l180
										}
										goto l179
									l180:
										position, tokenIndex = position180, tokenIndex180
									}
									if !_rules[ruleCLOSE]() {
										goto l172
									}
									goto l161
								l172:
									position, tokenIndex = position161, tokenIndex161
									{
										position182 := position
										if !_rules[ruleCommand]() {
											goto l181
										}
										if !_rules[ruleSEMI]() {
											goto l181
										}
										if !_rules[ruleConditionalExpression]() {
											goto l181
										}
										if !_rules[ruleSEMI]() {
											goto l181
										}
										if !_rules[ruleCommand]() {
											goto l181
										}
										add(ruleLoopConditionBounded, position182)
									}
									if !_rules[ruleOPEN]() {
										goto l181
									}
								l183:
									{
										position184, tokenIndex184 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l184
										}
										goto l183
									l184:
										position, tokenIndex = position184, tokenIndex184
									}
									if !_rules[ruleCLOSE]() {
										goto l181
									}
									goto l161
								l181:
									position, tokenIndex = position161, tokenIndex161
									{
										position185 := position
										if !_rules[ruleConditionalExpression]() {
											goto l158
										}
										add(ruleLoopConditionTruthy, position185)
									}
									if !_rules[ruleOPEN]() {
										goto l158
									}
								l186:
									{
										position187, tokenIndex187 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l187
										}
										goto l186
									l187:
										position, tokenIndex = position187, tokenIndex187
									}
									if !_rules[ruleCLOSE]() {
										goto l158
									}
								}
							l161:
								add(ruleLoop, position159)
							}
							goto l133
						l158:
							position, tokenIndex = position133, tokenIndex133
							if !_rules[ruleCommand]() {
								goto l112
							}
						}
					l133:
						add(ruleStatementBlock, position132)
					}
				}
			l114:
				{
					position188, tokenIndex188 := position, tokenIndex
					if !_rules[ruleSEMI]() {
						goto l188
					}
					goto l189
				l188:
					position, tokenIndex = position188, tokenIndex188
				}
			l189:
				if !_rules[rule_]() {
					goto l112
				}
				add(ruleBlock, position113)
			}
			return true
		l112:
			position, tokenIndex = position112, tokenIndex112
			return false
		},
		/* 66 FlowControlWord <- <(FlowControlBreak / FlowControlContinue)> */
		nil,
		/* 67 FlowControlBreak <- <(BREAK PositiveInteger?)> */
		nil,
		/* 68 FlowControlContinue <- <(CONT PositiveInteger?)> */
		nil,
		/* 69 StatementBlock <- <(NOOP / Assignment / Directive / Conditional / Loop / Command)> */
		nil,
		/* 70 Assignment <- <(AssignmentLHS AssignmentOperator AssignmentRHS)> */
		func() bool {
			position194, tokenIndex194 := position, tokenIndex
			{
				position195 := position
				{
					position196 := position
					if !_rules[ruleVariableSequence]() {
						goto l194
					}
					add(ruleAssignmentLHS, position196)
				}
				{
					position197 := position
					if !_rules[rule_]() {
						goto l194
					}
					{
						position198, tokenIndex198 := position, tokenIndex
						{
							position200 := position
							if !_rules[rule_]() {
								goto l199
							}
							if buffer[position] != rune('=') {
								goto l199
							}
							position++
							if !_rules[rule_]() {
								goto l199
							}
							add(ruleAssignEq, position200)
						}
						goto l198
					l199:
						position, tokenIndex = position198, tokenIndex198
						{
							position202 := position
							if !_rules[rule_]() {
								goto l201
							}
							if buffer[position] != rune('*') {
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
							add(ruleStarEq, position202)
						}
						goto l198
					l201:
						position, tokenIndex = position198, tokenIndex198
						{
							position204 := position
							if !_rules[rule_]() {
								goto l203
							}
							if buffer[position] != rune('/') {
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
							add(ruleDivEq, position204)
						}
						goto l198
					l203:
						position, tokenIndex = position198, tokenIndex198
						{
							position206 := position
							if !_rules[rule_]() {
								goto l205
							}
							if buffer[position] != rune('+') {
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
							add(rulePlusEq, position206)
						}
						goto l198
					l205:
						position, tokenIndex = position198, tokenIndex198
						{
							position208 := position
							if !_rules[rule_]() {
								goto l207
							}
							if buffer[position] != rune('-') {
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
							add(ruleMinusEq, position208)
						}
						goto l198
					l207:
						position, tokenIndex = position198, tokenIndex198
						{
							position210 := position
							if !_rules[rule_]() {
								goto l209
							}
							if buffer[position] != rune('&') {
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
							add(ruleAndEq, position210)
						}
						goto l198
					l209:
						position, tokenIndex = position198, tokenIndex198
						{
							position212 := position
							if !_rules[rule_]() {
								goto l211
							}
							if buffer[position] != rune('|') {
								goto l211
							}
							position++
							if buffer[position] != rune('=') {
								goto l211
							}
							position++
							if !_rules[rule_]() {
								goto l211
							}
							add(ruleOrEq, position212)
						}
						goto l198
					l211:
						position, tokenIndex = position198, tokenIndex198
						{
							position213 := position
							if !_rules[rule_]() {
								goto l194
							}
							if buffer[position] != rune('<') {
								goto l194
							}
							position++
							if buffer[position] != rune('<') {
								goto l194
							}
							position++
							if !_rules[rule_]() {
								goto l194
							}
							add(ruleAppend, position213)
						}
					}
				l198:
					if !_rules[rule_]() {
						goto l194
					}
					add(ruleAssignmentOperator, position197)
				}
				{
					position214 := position
					if !_rules[ruleExpressionSequence]() {
						goto l194
					}
					add(ruleAssignmentRHS, position214)
				}
				add(ruleAssignment, position195)
			}
			return true
		l194:
			position, tokenIndex = position194, tokenIndex194
			return false
		},
		/* 71 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 72 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 73 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position217, tokenIndex217 := position, tokenIndex
			{
				position218 := position
			l219:
				{
					position220, tokenIndex220 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l220
					}
					if !_rules[ruleCOMMA]() {
						goto l220
					}
					goto l219
				l220:
					position, tokenIndex = position220, tokenIndex220
				}
				if !_rules[ruleVariable]() {
					goto l217
				}
				add(ruleVariableSequence, position218)
			}
			return true
		l217:
			position, tokenIndex = position217, tokenIndex217
			return false
		},
		/* 74 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position221, tokenIndex221 := position, tokenIndex
			{
				position222 := position
			l223:
				{
					position224, tokenIndex224 := position, tokenIndex
					if !_rules[ruleExpression]() {
						goto l224
					}
					if !_rules[ruleCOMMA]() {
						goto l224
					}
					goto l223
				l224:
					position, tokenIndex = position224, tokenIndex224
				}
				if !_rules[ruleExpression]() {
					goto l221
				}
				add(ruleExpressionSequence, position222)
			}
			return true
		l221:
			position, tokenIndex = position221, tokenIndex221
			return false
		},
		/* 75 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position225, tokenIndex225 := position, tokenIndex
			{
				position226 := position
				if !_rules[rule_]() {
					goto l225
				}
				{
					position227 := position
					{
						position228 := position
						{
							position229, tokenIndex229 := position, tokenIndex
							if !_rules[ruleType]() {
								goto l230
							}
							goto l229
						l230:
							position, tokenIndex = position229, tokenIndex229
							if !_rules[ruleVariable]() {
								goto l225
							}
						}
					l229:
						add(ruleValueYielding, position228)
					}
					add(ruleExpressionLHS, position227)
				}
				{
					position231, tokenIndex231 := position, tokenIndex
					{
						position233 := position
						{
							position234 := position
							if !_rules[rule_]() {
								goto l231
							}
							{
								position235, tokenIndex235 := position, tokenIndex
								{
									position237 := position
									if !_rules[rule_]() {
										goto l236
									}
									if buffer[position] != rune('*') {
										goto l236
									}
									position++
									if buffer[position] != rune('*') {
										goto l236
									}
									position++
									if !_rules[rule_]() {
										goto l236
									}
									add(ruleExponentiate, position237)
								}
								goto l235
							l236:
								position, tokenIndex = position235, tokenIndex235
								{
									position239 := position
									if !_rules[rule_]() {
										goto l238
									}
									if buffer[position] != rune('*') {
										goto l238
									}
									position++
									if !_rules[rule_]() {
										goto l238
									}
									add(ruleMultiply, position239)
								}
								goto l235
							l238:
								position, tokenIndex = position235, tokenIndex235
								{
									position241 := position
									if !_rules[rule_]() {
										goto l240
									}
									if buffer[position] != rune('/') {
										goto l240
									}
									position++
									if !_rules[rule_]() {
										goto l240
									}
									add(ruleDivide, position241)
								}
								goto l235
							l240:
								position, tokenIndex = position235, tokenIndex235
								{
									position243 := position
									if !_rules[rule_]() {
										goto l242
									}
									if buffer[position] != rune('%') {
										goto l242
									}
									position++
									if !_rules[rule_]() {
										goto l242
									}
									add(ruleModulus, position243)
								}
								goto l235
							l242:
								position, tokenIndex = position235, tokenIndex235
								{
									position245 := position
									if !_rules[rule_]() {
										goto l244
									}
									if buffer[position] != rune('+') {
										goto l244
									}
									position++
									if !_rules[rule_]() {
										goto l244
									}
									add(ruleAdd, position245)
								}
								goto l235
							l244:
								position, tokenIndex = position235, tokenIndex235
								{
									position247 := position
									if !_rules[rule_]() {
										goto l246
									}
									if buffer[position] != rune('-') {
										goto l246
									}
									position++
									if !_rules[rule_]() {
										goto l246
									}
									add(ruleSubtract, position247)
								}
								goto l235
							l246:
								position, tokenIndex = position235, tokenIndex235
								{
									position249 := position
									if !_rules[rule_]() {
										goto l248
									}
									if buffer[position] != rune('&') {
										goto l248
									}
									position++
									if !_rules[rule_]() {
										goto l248
									}
									add(ruleBitwiseAnd, position249)
								}
								goto l235
							l248:
								position, tokenIndex = position235, tokenIndex235
								{
									position251 := position
									if !_rules[rule_]() {
										goto l250
									}
									if buffer[position] != rune('|') {
										goto l250
									}
									position++
									if !_rules[rule_]() {
										goto l250
									}
									add(ruleBitwiseOr, position251)
								}
								goto l235
							l250:
								position, tokenIndex = position235, tokenIndex235
								{
									position253 := position
									if !_rules[rule_]() {
										goto l252
									}
									if buffer[position] != rune('~') {
										goto l252
									}
									position++
									if !_rules[rule_]() {
										goto l252
									}
									add(ruleBitwiseNot, position253)
								}
								goto l235
							l252:
								position, tokenIndex = position235, tokenIndex235
								{
									position254 := position
									if !_rules[rule_]() {
										goto l231
									}
									if buffer[position] != rune('^') {
										goto l231
									}
									position++
									if !_rules[rule_]() {
										goto l231
									}
									add(ruleBitwiseXor, position254)
								}
							}
						l235:
							if !_rules[rule_]() {
								goto l231
							}
							add(ruleOperator, position234)
						}
						if !_rules[ruleExpression]() {
							goto l231
						}
						add(ruleExpressionRHS, position233)
					}
					goto l232
				l231:
					position, tokenIndex = position231, tokenIndex231
				}
			l232:
				if !_rules[rule_]() {
					goto l225
				}
				add(ruleExpression, position226)
			}
			return true
		l225:
			position, tokenIndex = position225, tokenIndex225
			return false
		},
		/* 76 ExpressionLHS <- <ValueYielding> */
		nil,
		/* 77 ExpressionRHS <- <(Operator Expression)> */
		nil,
		/* 78 ValueYielding <- <(Type / Variable)> */
		nil,
		/* 79 Directive <- <(DirectiveUnset / DirectiveInclude / DirectiveDeclare)> */
		nil,
		/* 80 DirectiveUnset <- <(UNSET VariableSequence)> */
		nil,
		/* 81 DirectiveInclude <- <(INCLUDE String)> */
		nil,
		/* 82 DirectiveDeclare <- <(DECLARE VariableSequence)> */
		nil,
		/* 83 Command <- <(_ CommandName (__ ((CommandFirstArg __ CommandSecondArg) / CommandFirstArg / CommandSecondArg))? (_ CommandResultAssignment)?)> */
		func() bool {
			position262, tokenIndex262 := position, tokenIndex
			{
				position263 := position
				if !_rules[rule_]() {
					goto l262
				}
				{
					position264 := position
					{
						position265, tokenIndex265 := position, tokenIndex
						if !_rules[ruleIdentifier]() {
							goto l265
						}
						{
							position267 := position
							if buffer[position] != rune(':') {
								goto l265
							}
							position++
							if buffer[position] != rune(':') {
								goto l265
							}
							position++
							add(ruleSCOPE, position267)
						}
						goto l266
					l265:
						position, tokenIndex = position265, tokenIndex265
					}
				l266:
					if !_rules[ruleIdentifier]() {
						goto l262
					}
					add(ruleCommandName, position264)
				}
				{
					position268, tokenIndex268 := position, tokenIndex
					if !_rules[rule__]() {
						goto l268
					}
					{
						position270, tokenIndex270 := position, tokenIndex
						if !_rules[ruleCommandFirstArg]() {
							goto l271
						}
						if !_rules[rule__]() {
							goto l271
						}
						if !_rules[ruleCommandSecondArg]() {
							goto l271
						}
						goto l270
					l271:
						position, tokenIndex = position270, tokenIndex270
						if !_rules[ruleCommandFirstArg]() {
							goto l272
						}
						goto l270
					l272:
						position, tokenIndex = position270, tokenIndex270
						if !_rules[ruleCommandSecondArg]() {
							goto l268
						}
					}
				l270:
					goto l269
				l268:
					position, tokenIndex = position268, tokenIndex268
				}
			l269:
				{
					position273, tokenIndex273 := position, tokenIndex
					if !_rules[rule_]() {
						goto l273
					}
					{
						position275 := position
						{
							position276 := position
							if !_rules[rule_]() {
								goto l273
							}
							if buffer[position] != rune('-') {
								goto l273
							}
							position++
							if buffer[position] != rune('>') {
								goto l273
							}
							position++
							if !_rules[rule_]() {
								goto l273
							}
							add(ruleASSIGN, position276)
						}
						if !_rules[ruleVariable]() {
							goto l273
						}
						add(ruleCommandResultAssignment, position275)
					}
					goto l274
				l273:
					position, tokenIndex = position273, tokenIndex273
				}
			l274:
				add(ruleCommand, position263)
			}
			return true
		l262:
			position, tokenIndex = position262, tokenIndex262
			return false
		},
		/* 84 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 85 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position278, tokenIndex278 := position, tokenIndex
			{
				position279 := position
				{
					position280, tokenIndex280 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l281
					}
					goto l280
				l281:
					position, tokenIndex = position280, tokenIndex280
					if !_rules[ruleType]() {
						goto l278
					}
				}
			l280:
				add(ruleCommandFirstArg, position279)
			}
			return true
		l278:
			position, tokenIndex = position278, tokenIndex278
			return false
		},
		/* 86 CommandSecondArg <- <Object> */
		func() bool {
			position282, tokenIndex282 := position, tokenIndex
			{
				position283 := position
				if !_rules[ruleObject]() {
					goto l282
				}
				add(ruleCommandSecondArg, position283)
			}
			return true
		l282:
			position, tokenIndex = position282, tokenIndex282
			return false
		},
		/* 87 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 88 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 89 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position286, tokenIndex286 := position, tokenIndex
			{
				position287 := position
				{
					position288 := position
					if !_rules[rule_]() {
						goto l286
					}
					if buffer[position] != rune('i') {
						goto l286
					}
					position++
					if buffer[position] != rune('f') {
						goto l286
					}
					position++
					if !_rules[rule_]() {
						goto l286
					}
					add(ruleIF, position288)
				}
				if !_rules[ruleConditionalExpression]() {
					goto l286
				}
				if !_rules[ruleOPEN]() {
					goto l286
				}
			l289:
				{
					position290, tokenIndex290 := position, tokenIndex
					if !_rules[ruleBlock]() {
						goto l290
					}
					goto l289
				l290:
					position, tokenIndex = position290, tokenIndex290
				}
				if !_rules[ruleCLOSE]() {
					goto l286
				}
				add(ruleIfStanza, position287)
			}
			return true
		l286:
			position, tokenIndex = position286, tokenIndex286
			return false
		},
		/* 90 ElseIfStanza <- <(ELSE IfStanza)> */
		nil,
		/* 91 ElseStanza <- <(ELSE OPEN Block* CLOSE)> */
		nil,
		/* 92 Loop <- <(LOOP ((OPEN Block* CLOSE) / (LoopConditionFixedLength OPEN Block* CLOSE) / (LoopConditionIterable OPEN Block* CLOSE) / (LoopConditionBounded OPEN Block* CLOSE) / (LoopConditionTruthy OPEN Block* CLOSE)))> */
		nil,
		/* 93 LoopConditionFixedLength <- <(COUNT (Integer / Variable))> */
		nil,
		/* 94 LoopConditionIterable <- <(LoopIterableLHS IN LoopIterableRHS)> */
		nil,
		/* 95 LoopIterableLHS <- <VariableSequence> */
		nil,
		/* 96 LoopIterableRHS <- <(Command / Variable)> */
		nil,
		/* 97 LoopConditionBounded <- <(Command SEMI ConditionalExpression SEMI Command)> */
		nil,
		/* 98 LoopConditionTruthy <- <ConditionalExpression> */
		nil,
		/* 99 ConditionalExpression <- <(NOT? (ConditionWithAssignment / ConditionWithCommand / ConditionWithRegex / ConditionWithComparator))> */
		func() bool {
			position300, tokenIndex300 := position, tokenIndex
			{
				position301 := position
				{
					position302, tokenIndex302 := position, tokenIndex
					{
						position304 := position
						if !_rules[rule_]() {
							goto l302
						}
						if buffer[position] != rune('n') {
							goto l302
						}
						position++
						if buffer[position] != rune('o') {
							goto l302
						}
						position++
						if buffer[position] != rune('t') {
							goto l302
						}
						position++
						if !_rules[rule__]() {
							goto l302
						}
						add(ruleNOT, position304)
					}
					goto l303
				l302:
					position, tokenIndex = position302, tokenIndex302
				}
			l303:
				{
					position305, tokenIndex305 := position, tokenIndex
					{
						position307 := position
						if !_rules[ruleAssignment]() {
							goto l306
						}
						if !_rules[ruleSEMI]() {
							goto l306
						}
						if !_rules[ruleConditionalExpression]() {
							goto l306
						}
						add(ruleConditionWithAssignment, position307)
					}
					goto l305
				l306:
					position, tokenIndex = position305, tokenIndex305
					{
						position309 := position
						if !_rules[ruleCommand]() {
							goto l308
						}
						{
							position310, tokenIndex310 := position, tokenIndex
							if !_rules[ruleSEMI]() {
								goto l310
							}
							if !_rules[ruleConditionalExpression]() {
								goto l310
							}
							goto l311
						l310:
							position, tokenIndex = position310, tokenIndex310
						}
					l311:
						add(ruleConditionWithCommand, position309)
					}
					goto l305
				l308:
					position, tokenIndex = position305, tokenIndex305
					{
						position313 := position
						if !_rules[ruleExpression]() {
							goto l312
						}
						{
							position314 := position
							{
								position315, tokenIndex315 := position, tokenIndex
								{
									position317 := position
									if !_rules[rule_]() {
										goto l316
									}
									if buffer[position] != rune('=') {
										goto l316
									}
									position++
									if buffer[position] != rune('~') {
										goto l316
									}
									position++
									if !_rules[rule_]() {
										goto l316
									}
									add(ruleMatch, position317)
								}
								goto l315
							l316:
								position, tokenIndex = position315, tokenIndex315
								{
									position318 := position
									if !_rules[rule_]() {
										goto l312
									}
									if buffer[position] != rune('!') {
										goto l312
									}
									position++
									if buffer[position] != rune('~') {
										goto l312
									}
									position++
									if !_rules[rule_]() {
										goto l312
									}
									add(ruleUnmatch, position318)
								}
							}
						l315:
							add(ruleMatchOperator, position314)
						}
						if !_rules[ruleRegularExpression]() {
							goto l312
						}
						add(ruleConditionWithRegex, position313)
					}
					goto l305
				l312:
					position, tokenIndex = position305, tokenIndex305
					{
						position319 := position
						{
							position320 := position
							if !_rules[ruleExpression]() {
								goto l300
							}
							add(ruleConditionWithComparatorLHS, position320)
						}
						{
							position321, tokenIndex321 := position, tokenIndex
							{
								position323 := position
								{
									position324 := position
									if !_rules[rule_]() {
										goto l321
									}
									{
										position325, tokenIndex325 := position, tokenIndex
										{
											position327 := position
											if !_rules[rule_]() {
												goto l326
											}
											if buffer[position] != rune('=') {
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
											add(ruleEquality, position327)
										}
										goto l325
									l326:
										position, tokenIndex = position325, tokenIndex325
										{
											position329 := position
											if !_rules[rule_]() {
												goto l328
											}
											if buffer[position] != rune('!') {
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
											add(ruleNonEquality, position329)
										}
										goto l325
									l328:
										position, tokenIndex = position325, tokenIndex325
										{
											position331 := position
											if !_rules[rule_]() {
												goto l330
											}
											if buffer[position] != rune('>') {
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
											add(ruleGreaterEqual, position331)
										}
										goto l325
									l330:
										position, tokenIndex = position325, tokenIndex325
										{
											position333 := position
											if !_rules[rule_]() {
												goto l332
											}
											if buffer[position] != rune('<') {
												goto l332
											}
											position++
											if buffer[position] != rune('=') {
												goto l332
											}
											position++
											if !_rules[rule_]() {
												goto l332
											}
											add(ruleLessEqual, position333)
										}
										goto l325
									l332:
										position, tokenIndex = position325, tokenIndex325
										{
											position335 := position
											if !_rules[rule_]() {
												goto l334
											}
											if buffer[position] != rune('>') {
												goto l334
											}
											position++
											if !_rules[rule_]() {
												goto l334
											}
											add(ruleGreaterThan, position335)
										}
										goto l325
									l334:
										position, tokenIndex = position325, tokenIndex325
										{
											position337 := position
											if !_rules[rule_]() {
												goto l336
											}
											if buffer[position] != rune('<') {
												goto l336
											}
											position++
											if !_rules[rule_]() {
												goto l336
											}
											add(ruleLessThan, position337)
										}
										goto l325
									l336:
										position, tokenIndex = position325, tokenIndex325
										{
											position339 := position
											if !_rules[rule_]() {
												goto l338
											}
											if buffer[position] != rune('i') {
												goto l338
											}
											position++
											if buffer[position] != rune('n') {
												goto l338
											}
											position++
											if !_rules[rule_]() {
												goto l338
											}
											add(ruleMembership, position339)
										}
										goto l325
									l338:
										position, tokenIndex = position325, tokenIndex325
										{
											position340 := position
											if !_rules[rule_]() {
												goto l321
											}
											if buffer[position] != rune('n') {
												goto l321
											}
											position++
											if buffer[position] != rune('o') {
												goto l321
											}
											position++
											if buffer[position] != rune('t') {
												goto l321
											}
											position++
											if !_rules[rule__]() {
												goto l321
											}
											if buffer[position] != rune('i') {
												goto l321
											}
											position++
											if buffer[position] != rune('n') {
												goto l321
											}
											position++
											if !_rules[rule_]() {
												goto l321
											}
											add(ruleNonMembership, position340)
										}
									}
								l325:
									if !_rules[rule_]() {
										goto l321
									}
									add(ruleComparisonOperator, position324)
								}
								if !_rules[ruleExpression]() {
									goto l321
								}
								add(ruleConditionWithComparatorRHS, position323)
							}
							goto l322
						l321:
							position, tokenIndex = position321, tokenIndex321
						}
					l322:
						add(ruleConditionWithComparator, position319)
					}
				}
			l305:
				add(ruleConditionalExpression, position301)
			}
			return true
		l300:
			position, tokenIndex = position300, tokenIndex300
			return false
		},
		/* 100 ConditionWithAssignment <- <(Assignment SEMI ConditionalExpression)> */
		nil,
		/* 101 ConditionWithCommand <- <(Command (SEMI ConditionalExpression)?)> */
		nil,
		/* 102 ConditionWithRegex <- <(Expression MatchOperator RegularExpression)> */
		nil,
		/* 103 ConditionWithComparator <- <(ConditionWithComparatorLHS ConditionWithComparatorRHS?)> */
		nil,
		/* 104 ConditionWithComparatorLHS <- <Expression> */
		nil,
		/* 105 ConditionWithComparatorRHS <- <(ComparisonOperator Expression)> */
		nil,
		/* 106 ScalarType <- <(Boolean / Float / Integer / String / NullValue)> */
		nil,
		/* 107 Identifier <- <(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / ([0-9] / [0-9]) / '_')*)> */
		func() bool {
			position348, tokenIndex348 := position, tokenIndex
			{
				position349 := position
				{
					position350, tokenIndex350 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex = position350, tokenIndex350
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l352
					}
					position++
					goto l350
				l352:
					position, tokenIndex = position350, tokenIndex350
					if buffer[position] != rune('_') {
						goto l348
					}
					position++
				}
			l350:
			l353:
				{
					position354, tokenIndex354 := position, tokenIndex
					{
						position355, tokenIndex355 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l356
						}
						position++
						goto l355
					l356:
						position, tokenIndex = position355, tokenIndex355
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l357
						}
						position++
						goto l355
					l357:
						position, tokenIndex = position355, tokenIndex355
						{
							position359, tokenIndex359 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l360
							}
							position++
							goto l359
						l360:
							position, tokenIndex = position359, tokenIndex359
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l358
							}
							position++
						}
					l359:
						goto l355
					l358:
						position, tokenIndex = position355, tokenIndex355
						if buffer[position] != rune('_') {
							goto l354
						}
						position++
					}
				l355:
					goto l353
				l354:
					position, tokenIndex = position354, tokenIndex354
				}
				add(ruleIdentifier, position349)
			}
			return true
		l348:
			position, tokenIndex = position348, tokenIndex348
			return false
		},
		/* 108 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 109 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 110 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position363, tokenIndex363 := position, tokenIndex
			{
				position364 := position
				{
					position365, tokenIndex365 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l365
					}
					position++
					goto l366
				l365:
					position, tokenIndex = position365, tokenIndex365
				}
			l366:
				if !_rules[rulePositiveInteger]() {
					goto l363
				}
				add(ruleInteger, position364)
			}
			return true
		l363:
			position, tokenIndex = position363, tokenIndex363
			return false
		},
		/* 111 PositiveInteger <- <[0-9]+> */
		func() bool {
			position367, tokenIndex367 := position, tokenIndex
			{
				position368 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l367
				}
				position++
			l369:
				{
					position370, tokenIndex370 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex = position370, tokenIndex370
				}
				add(rulePositiveInteger, position368)
			}
			return true
		l367:
			position, tokenIndex = position367, tokenIndex367
			return false
		},
		/* 112 String <- <(StringLiteral / StringInterpolated / Heredoc)> */
		func() bool {
			position371, tokenIndex371 := position, tokenIndex
			{
				position372 := position
				{
					position373, tokenIndex373 := position, tokenIndex
					{
						position375 := position
						if buffer[position] != rune('\'') {
							goto l374
						}
						position++
					l376:
						{
							position377, tokenIndex377 := position, tokenIndex
							{
								position378, tokenIndex378 := position, tokenIndex
								if buffer[position] != rune('\'') {
									goto l378
								}
								position++
								goto l377
							l378:
								position, tokenIndex = position378, tokenIndex378
							}
							if !matchDot() {
								goto l377
							}
							goto l376
						l377:
							position, tokenIndex = position377, tokenIndex377
						}
						if buffer[position] != rune('\'') {
							goto l374
						}
						position++
						add(ruleStringLiteral, position375)
					}
					goto l373
				l374:
					position, tokenIndex = position373, tokenIndex373
					{
						position380 := position
						if buffer[position] != rune('"') {
							goto l379
						}
						position++
					l381:
						{
							position382, tokenIndex382 := position, tokenIndex
							{
								position383, tokenIndex383 := position, tokenIndex
								if buffer[position] != rune('"') {
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
						if buffer[position] != rune('"') {
							goto l379
						}
						position++
						add(ruleStringInterpolated, position380)
					}
					goto l373
				l379:
					position, tokenIndex = position373, tokenIndex373
					{
						position384 := position
						{
							position385 := position
							if !_rules[rule_]() {
								goto l371
							}
							if buffer[position] != rune('b') {
								goto l371
							}
							position++
							if buffer[position] != rune('e') {
								goto l371
							}
							position++
							if buffer[position] != rune('g') {
								goto l371
							}
							position++
							if buffer[position] != rune('i') {
								goto l371
							}
							position++
							if buffer[position] != rune('n') {
								goto l371
							}
							position++
							add(ruleBEGIN, position385)
						}
						{
							position386 := position
							if buffer[position] != rune('\n') {
								goto l371
							}
							position++
							add(ruleNL, position386)
						}
						{
							position387 := position
						l388:
							{
								position389, tokenIndex389 := position, tokenIndex
								{
									position390, tokenIndex390 := position, tokenIndex
									if !_rules[ruleEND]() {
										goto l390
									}
									goto l389
								l390:
									position, tokenIndex = position390, tokenIndex390
								}
								if !matchDot() {
									goto l389
								}
								goto l388
							l389:
								position, tokenIndex = position389, tokenIndex389
							}
							add(ruleHeredocBody, position387)
						}
						if !_rules[ruleEND]() {
							goto l371
						}
						add(ruleHeredoc, position384)
					}
				}
			l373:
				add(ruleString, position372)
			}
			return true
		l371:
			position, tokenIndex = position371, tokenIndex371
			return false
		},
		/* 113 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		nil,
		/* 114 StringInterpolated <- <('"' (!'"' .)* '"')> */
		nil,
		/* 115 Heredoc <- <(BEGIN NL HeredocBody END)> */
		nil,
		/* 116 HeredocBody <- <(!END .)*> */
		nil,
		/* 117 NullValue <- <('n' 'u' 'l' 'l')> */
		nil,
		/* 118 Object <- <(OPEN (_ KeyValuePair _)* CLOSE)> */
		func() bool {
			position396, tokenIndex396 := position, tokenIndex
			{
				position397 := position
				if !_rules[ruleOPEN]() {
					goto l396
				}
			l398:
				{
					position399, tokenIndex399 := position, tokenIndex
					if !_rules[rule_]() {
						goto l399
					}
					{
						position400 := position
						{
							position401 := position
							if !_rules[ruleIdentifier]() {
								goto l399
							}
							add(ruleKey, position401)
						}
						{
							position402 := position
							if !_rules[rule_]() {
								goto l399
							}
							if buffer[position] != rune(':') {
								goto l399
							}
							position++
							if !_rules[rule_]() {
								goto l399
							}
							add(ruleCOLON, position402)
						}
						{
							position403 := position
							{
								position404, tokenIndex404 := position, tokenIndex
								if !_rules[ruleArray]() {
									goto l405
								}
								goto l404
							l405:
								position, tokenIndex = position404, tokenIndex404
								if !_rules[ruleObject]() {
									goto l406
								}
								goto l404
							l406:
								position, tokenIndex = position404, tokenIndex404
								if !_rules[ruleExpression]() {
									goto l399
								}
							}
						l404:
							add(ruleKValue, position403)
						}
						{
							position407, tokenIndex407 := position, tokenIndex
							if !_rules[ruleCOMMA]() {
								goto l407
							}
							goto l408
						l407:
							position, tokenIndex = position407, tokenIndex407
						}
					l408:
						add(ruleKeyValuePair, position400)
					}
					if !_rules[rule_]() {
						goto l399
					}
					goto l398
				l399:
					position, tokenIndex = position399, tokenIndex399
				}
				if !_rules[ruleCLOSE]() {
					goto l396
				}
				add(ruleObject, position397)
			}
			return true
		l396:
			position, tokenIndex = position396, tokenIndex396
			return false
		},
		/* 119 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position409, tokenIndex409 := position, tokenIndex
			{
				position410 := position
				if buffer[position] != rune('[') {
					goto l409
				}
				position++
				if !_rules[rule_]() {
					goto l409
				}
				if !_rules[ruleExpressionSequence]() {
					goto l409
				}
				{
					position411, tokenIndex411 := position, tokenIndex
					if !_rules[ruleCOMMA]() {
						goto l411
					}
					goto l412
				l411:
					position, tokenIndex = position411, tokenIndex411
				}
			l412:
				if buffer[position] != rune(']') {
					goto l409
				}
				position++
				add(ruleArray, position410)
			}
			return true
		l409:
			position, tokenIndex = position409, tokenIndex409
			return false
		},
		/* 120 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position413, tokenIndex413 := position, tokenIndex
			{
				position414 := position
				if buffer[position] != rune('/') {
					goto l413
				}
				position++
				{
					position417, tokenIndex417 := position, tokenIndex
					if buffer[position] != rune('/') {
						goto l417
					}
					position++
					goto l413
				l417:
					position, tokenIndex = position417, tokenIndex417
				}
				if !matchDot() {
					goto l413
				}
			l415:
				{
					position416, tokenIndex416 := position, tokenIndex
					{
						position418, tokenIndex418 := position, tokenIndex
						if buffer[position] != rune('/') {
							goto l418
						}
						position++
						goto l416
					l418:
						position, tokenIndex = position418, tokenIndex418
					}
					if !matchDot() {
						goto l416
					}
					goto l415
				l416:
					position, tokenIndex = position416, tokenIndex416
				}
				if buffer[position] != rune('/') {
					goto l413
				}
				position++
			l419:
				{
					position420, tokenIndex420 := position, tokenIndex
					{
						position421, tokenIndex421 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l422
						}
						position++
						goto l421
					l422:
						position, tokenIndex = position421, tokenIndex421
						if buffer[position] != rune('l') {
							goto l423
						}
						position++
						goto l421
					l423:
						position, tokenIndex = position421, tokenIndex421
						if buffer[position] != rune('m') {
							goto l424
						}
						position++
						goto l421
					l424:
						position, tokenIndex = position421, tokenIndex421
						if buffer[position] != rune('s') {
							goto l425
						}
						position++
						goto l421
					l425:
						position, tokenIndex = position421, tokenIndex421
						if buffer[position] != rune('u') {
							goto l420
						}
						position++
					}
				l421:
					goto l419
				l420:
					position, tokenIndex = position420, tokenIndex420
				}
				add(ruleRegularExpression, position414)
			}
			return true
		l413:
			position, tokenIndex = position413, tokenIndex413
			return false
		},
		/* 121 KeyValuePair <- <(Key COLON KValue COMMA?)> */
		nil,
		/* 122 Key <- <Identifier> */
		nil,
		/* 123 KValue <- <(Array / Object / Expression)> */
		nil,
		/* 124 Type <- <(Array / Object / RegularExpression / ScalarType)> */
		func() bool {
			position429, tokenIndex429 := position, tokenIndex
			{
				position430 := position
				{
					position431, tokenIndex431 := position, tokenIndex
					if !_rules[ruleArray]() {
						goto l432
					}
					goto l431
				l432:
					position, tokenIndex = position431, tokenIndex431
					if !_rules[ruleObject]() {
						goto l433
					}
					goto l431
				l433:
					position, tokenIndex = position431, tokenIndex431
					if !_rules[ruleRegularExpression]() {
						goto l434
					}
					goto l431
				l434:
					position, tokenIndex = position431, tokenIndex431
					{
						position435 := position
						{
							position436, tokenIndex436 := position, tokenIndex
							{
								position438 := position
								{
									position439, tokenIndex439 := position, tokenIndex
									if buffer[position] != rune('t') {
										goto l440
									}
									position++
									if buffer[position] != rune('r') {
										goto l440
									}
									position++
									if buffer[position] != rune('u') {
										goto l440
									}
									position++
									if buffer[position] != rune('e') {
										goto l440
									}
									position++
									goto l439
								l440:
									position, tokenIndex = position439, tokenIndex439
									if buffer[position] != rune('f') {
										goto l437
									}
									position++
									if buffer[position] != rune('a') {
										goto l437
									}
									position++
									if buffer[position] != rune('l') {
										goto l437
									}
									position++
									if buffer[position] != rune('s') {
										goto l437
									}
									position++
									if buffer[position] != rune('e') {
										goto l437
									}
									position++
								}
							l439:
								add(ruleBoolean, position438)
							}
							goto l436
						l437:
							position, tokenIndex = position436, tokenIndex436
							{
								position442 := position
								if !_rules[ruleInteger]() {
									goto l441
								}
								{
									position443, tokenIndex443 := position, tokenIndex
									if buffer[position] != rune('.') {
										goto l443
									}
									position++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l443
									}
									position++
								l445:
									{
										position446, tokenIndex446 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l446
										}
										position++
										goto l445
									l446:
										position, tokenIndex = position446, tokenIndex446
									}
									goto l444
								l443:
									position, tokenIndex = position443, tokenIndex443
								}
							l444:
								add(ruleFloat, position442)
							}
							goto l436
						l441:
							position, tokenIndex = position436, tokenIndex436
							if !_rules[ruleInteger]() {
								goto l447
							}
							goto l436
						l447:
							position, tokenIndex = position436, tokenIndex436
							if !_rules[ruleString]() {
								goto l448
							}
							goto l436
						l448:
							position, tokenIndex = position436, tokenIndex436
							{
								position449 := position
								if buffer[position] != rune('n') {
									goto l429
								}
								position++
								if buffer[position] != rune('u') {
									goto l429
								}
								position++
								if buffer[position] != rune('l') {
									goto l429
								}
								position++
								if buffer[position] != rune('l') {
									goto l429
								}
								position++
								add(ruleNullValue, position449)
							}
						}
					l436:
						add(ruleScalarType, position435)
					}
				}
			l431:
				add(ruleType, position430)
			}
			return true
		l429:
			position, tokenIndex = position429, tokenIndex429
			return false
		},
	}
	p.rules = _rules
}
