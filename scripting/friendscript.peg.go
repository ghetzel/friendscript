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
	ruleTRIQUOT
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
	ruleHeredocTriquote
	ruleTriquoteBody
	ruleHeredocBody
	ruleHeredoc
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
	"TRIQUOT",
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
	"HeredocTriquote",
	"TriquoteBody",
	"HeredocBody",
	"Heredoc",
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
	rules  [129]func() bool
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
		func() bool {
			position12, tokenIndex12 := position, tokenIndex
			{
				position13 := position
				if buffer[position] != rune('\n') {
					goto l12
				}
				position++
				add(ruleNL, position13)
			}
			return true
		l12:
			position, tokenIndex = position12, tokenIndex12
			return false
		},
		/* 2 _ <- <(' ' / '\t' / '\r' / '\n')*> */
		func() bool {
			{
				position15 := position
			l16:
				{
					position17, tokenIndex17 := position, tokenIndex
					{
						position18, tokenIndex18 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l19
						}
						position++
						goto l18
					l19:
						position, tokenIndex = position18, tokenIndex18
						if buffer[position] != rune('\t') {
							goto l20
						}
						position++
						goto l18
					l20:
						position, tokenIndex = position18, tokenIndex18
						if buffer[position] != rune('\r') {
							goto l21
						}
						position++
						goto l18
					l21:
						position, tokenIndex = position18, tokenIndex18
						if buffer[position] != rune('\n') {
							goto l17
						}
						position++
					}
				l18:
					goto l16
				l17:
					position, tokenIndex = position17, tokenIndex17
				}
				add(rule_, position15)
			}
			return true
		},
		/* 3 __ <- <(' ' / '\t' / '\r' / '\n')+> */
		func() bool {
			position22, tokenIndex22 := position, tokenIndex
			{
				position23 := position
				{
					position26, tokenIndex26 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex = position26, tokenIndex26
					if buffer[position] != rune('\t') {
						goto l28
					}
					position++
					goto l26
				l28:
					position, tokenIndex = position26, tokenIndex26
					if buffer[position] != rune('\r') {
						goto l29
					}
					position++
					goto l26
				l29:
					position, tokenIndex = position26, tokenIndex26
					if buffer[position] != rune('\n') {
						goto l22
					}
					position++
				}
			l26:
			l24:
				{
					position25, tokenIndex25 := position, tokenIndex
					{
						position30, tokenIndex30 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l31
						}
						position++
						goto l30
					l31:
						position, tokenIndex = position30, tokenIndex30
						if buffer[position] != rune('\t') {
							goto l32
						}
						position++
						goto l30
					l32:
						position, tokenIndex = position30, tokenIndex30
						if buffer[position] != rune('\r') {
							goto l33
						}
						position++
						goto l30
					l33:
						position, tokenIndex = position30, tokenIndex30
						if buffer[position] != rune('\n') {
							goto l25
						}
						position++
					}
				l30:
					goto l24
				l25:
					position, tokenIndex = position25, tokenIndex25
				}
				add(rule__, position23)
			}
			return true
		l22:
			position, tokenIndex = position22, tokenIndex22
			return false
		},
		/* 4 ASSIGN <- <(_ ('-' '>') _)> */
		nil,
		/* 5 TRIQUOT <- <(_ ('"' '"' '"') _)> */
		func() bool {
			position35, tokenIndex35 := position, tokenIndex
			{
				position36 := position
				if !_rules[rule_]() {
					goto l35
				}
				if buffer[position] != rune('"') {
					goto l35
				}
				position++
				if buffer[position] != rune('"') {
					goto l35
				}
				position++
				if buffer[position] != rune('"') {
					goto l35
				}
				position++
				if !_rules[rule_]() {
					goto l35
				}
				add(ruleTRIQUOT, position36)
			}
			return true
		l35:
			position, tokenIndex = position35, tokenIndex35
			return false
		},
		/* 6 BEGIN <- <(_ ('b' 'e' 'g' 'i' 'n'))> */
		nil,
		/* 7 BREAK <- <(_ ('b' 'r' 'e' 'a' 'k') _)> */
		nil,
		/* 8 CLOSE <- <(_ '}' _)> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				if !_rules[rule_]() {
					goto l39
				}
				if buffer[position] != rune('}') {
					goto l39
				}
				position++
				if !_rules[rule_]() {
					goto l39
				}
				add(ruleCLOSE, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 9 COLON <- <(_ ':' _)> */
		nil,
		/* 10 COMMA <- <(_ ',' _)> */
		func() bool {
			position42, tokenIndex42 := position, tokenIndex
			{
				position43 := position
				if !_rules[rule_]() {
					goto l42
				}
				if buffer[position] != rune(',') {
					goto l42
				}
				position++
				if !_rules[rule_]() {
					goto l42
				}
				add(ruleCOMMA, position43)
			}
			return true
		l42:
			position, tokenIndex = position42, tokenIndex42
			return false
		},
		/* 11 Comment <- <(_ '#' (!'\n' .)*)> */
		nil,
		/* 12 CONT <- <(_ ('c' 'o' 'n' 't' 'i' 'n' 'u' 'e') _)> */
		nil,
		/* 13 COUNT <- <(_ ('c' 'o' 'u' 'n' 't') _)> */
		nil,
		/* 14 DECLARE <- <(_ ('d' 'e' 'c' 'l' 'a' 'r' 'e') __)> */
		nil,
		/* 15 DOT <- <'.'> */
		nil,
		/* 16 ELSE <- <(_ ('e' 'l' 's' 'e') _)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				if !_rules[rule_]() {
					goto l49
				}
				if buffer[position] != rune('e') {
					goto l49
				}
				position++
				if buffer[position] != rune('l') {
					goto l49
				}
				position++
				if buffer[position] != rune('s') {
					goto l49
				}
				position++
				if buffer[position] != rune('e') {
					goto l49
				}
				position++
				if !_rules[rule_]() {
					goto l49
				}
				add(ruleELSE, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 17 END <- <(NL _ ('e' 'n' 'd') _)> */
		func() bool {
			position51, tokenIndex51 := position, tokenIndex
			{
				position52 := position
				if !_rules[ruleNL]() {
					goto l51
				}
				if !_rules[rule_]() {
					goto l51
				}
				if buffer[position] != rune('e') {
					goto l51
				}
				position++
				if buffer[position] != rune('n') {
					goto l51
				}
				position++
				if buffer[position] != rune('d') {
					goto l51
				}
				position++
				if !_rules[rule_]() {
					goto l51
				}
				add(ruleEND, position52)
			}
			return true
		l51:
			position, tokenIndex = position51, tokenIndex51
			return false
		},
		/* 18 IF <- <(_ ('i' 'f') _)> */
		nil,
		/* 19 IN <- <(__ ('i' 'n') __)> */
		nil,
		/* 20 INCLUDE <- <(_ ('i' 'n' 'c' 'l' 'u' 'd' 'e') __)> */
		nil,
		/* 21 LOOP <- <(_ ('l' 'o' 'o' 'p') _)> */
		nil,
		/* 22 NOOP <- <SEMI> */
		nil,
		/* 23 NOT <- <(_ ('n' 'o' 't') __)> */
		nil,
		/* 24 OPEN <- <(_ '{' _)> */
		func() bool {
			position59, tokenIndex59 := position, tokenIndex
			{
				position60 := position
				if !_rules[rule_]() {
					goto l59
				}
				if buffer[position] != rune('{') {
					goto l59
				}
				position++
				if !_rules[rule_]() {
					goto l59
				}
				add(ruleOPEN, position60)
			}
			return true
		l59:
			position, tokenIndex = position59, tokenIndex59
			return false
		},
		/* 25 SCOPE <- <(':' ':')> */
		nil,
		/* 26 SEMI <- <(_ ';' _)> */
		func() bool {
			position62, tokenIndex62 := position, tokenIndex
			{
				position63 := position
				if !_rules[rule_]() {
					goto l62
				}
				if buffer[position] != rune(';') {
					goto l62
				}
				position++
				if !_rules[rule_]() {
					goto l62
				}
				add(ruleSEMI, position63)
			}
			return true
		l62:
			position, tokenIndex = position62, tokenIndex62
			return false
		},
		/* 27 SHEBANG <- <('#' '!' (!'\n' .)+ '\n')> */
		nil,
		/* 28 SKIPVAR <- <(_ '_' _)> */
		nil,
		/* 29 UNSET <- <(_ ('u' 'n' 's' 'e' 't') __)> */
		nil,
		/* 30 Operator <- <(_ (Exponentiate / Multiply / Divide / Modulus / Add / Subtract / BitwiseAnd / BitwiseOr / BitwiseNot / BitwiseXor) _)> */
		nil,
		/* 31 Exponentiate <- <(_ ('*' '*') _)> */
		nil,
		/* 32 Multiply <- <(_ '*' _)> */
		nil,
		/* 33 Divide <- <(_ '/' _)> */
		nil,
		/* 34 Modulus <- <(_ '%' _)> */
		nil,
		/* 35 Add <- <(_ '+' _)> */
		nil,
		/* 36 Subtract <- <(_ '-' _)> */
		nil,
		/* 37 BitwiseAnd <- <(_ '&' _)> */
		nil,
		/* 38 BitwiseOr <- <(_ '|' _)> */
		nil,
		/* 39 BitwiseNot <- <(_ '~' _)> */
		nil,
		/* 40 BitwiseXor <- <(_ '^' _)> */
		nil,
		/* 41 MatchOperator <- <(Match / Unmatch)> */
		nil,
		/* 42 Unmatch <- <(_ ('!' '~') _)> */
		nil,
		/* 43 Match <- <(_ ('=' '~') _)> */
		nil,
		/* 44 AssignmentOperator <- <(_ (AssignEq / StarEq / DivEq / PlusEq / MinusEq / AndEq / OrEq / Append) _)> */
		nil,
		/* 45 AssignEq <- <(_ '=' _)> */
		nil,
		/* 46 StarEq <- <(_ ('*' '=') _)> */
		nil,
		/* 47 DivEq <- <(_ ('/' '=') _)> */
		nil,
		/* 48 PlusEq <- <(_ ('+' '=') _)> */
		nil,
		/* 49 MinusEq <- <(_ ('-' '=') _)> */
		nil,
		/* 50 AndEq <- <(_ ('&' '=') _)> */
		nil,
		/* 51 OrEq <- <(_ ('|' '=') _)> */
		nil,
		/* 52 Append <- <(_ ('<' '<') _)> */
		nil,
		/* 53 ComparisonOperator <- <(_ (Equality / NonEquality / GreaterEqual / LessEqual / GreaterThan / LessThan / Membership / NonMembership) _)> */
		nil,
		/* 54 Equality <- <(_ ('=' '=') _)> */
		nil,
		/* 55 NonEquality <- <(_ ('!' '=') _)> */
		nil,
		/* 56 GreaterThan <- <(_ '>' _)> */
		nil,
		/* 57 GreaterEqual <- <(_ ('>' '=') _)> */
		nil,
		/* 58 LessEqual <- <(_ ('<' '=') _)> */
		nil,
		/* 59 LessThan <- <(_ '<' _)> */
		nil,
		/* 60 Membership <- <(_ ('i' 'n') _)> */
		nil,
		/* 61 NonMembership <- <(_ ('n' 'o' 't') __ ('i' 'n') _)> */
		nil,
		/* 62 Variable <- <(('$' VariableNameSequence) / SKIPVAR)> */
		func() bool {
			position99, tokenIndex99 := position, tokenIndex
			{
				position100 := position
				{
					position101, tokenIndex101 := position, tokenIndex
					if buffer[position] != rune('$') {
						goto l102
					}
					position++
					{
						position103 := position
					l104:
						{
							position105, tokenIndex105 := position, tokenIndex
							if !_rules[ruleVariableName]() {
								goto l105
							}
							{
								position106 := position
								if buffer[position] != rune('.') {
									goto l105
								}
								position++
								add(ruleDOT, position106)
							}
							goto l104
						l105:
							position, tokenIndex = position105, tokenIndex105
						}
						if !_rules[ruleVariableName]() {
							goto l102
						}
						add(ruleVariableNameSequence, position103)
					}
					goto l101
				l102:
					position, tokenIndex = position101, tokenIndex101
					{
						position107 := position
						if !_rules[rule_]() {
							goto l99
						}
						if buffer[position] != rune('_') {
							goto l99
						}
						position++
						if !_rules[rule_]() {
							goto l99
						}
						add(ruleSKIPVAR, position107)
					}
				}
			l101:
				add(ruleVariable, position100)
			}
			return true
		l99:
			position, tokenIndex = position99, tokenIndex99
			return false
		},
		/* 63 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 64 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position109, tokenIndex109 := position, tokenIndex
			{
				position110 := position
				if !_rules[ruleIdentifier]() {
					goto l109
				}
				{
					position111, tokenIndex111 := position, tokenIndex
					if buffer[position] != rune('[') {
						goto l111
					}
					position++
					if !_rules[rule_]() {
						goto l111
					}
					{
						position113 := position
						if !_rules[ruleExpression]() {
							goto l111
						}
						add(ruleVariableIndex, position113)
					}
					if !_rules[rule_]() {
						goto l111
					}
					if buffer[position] != rune(']') {
						goto l111
					}
					position++
					goto l112
				l111:
					position, tokenIndex = position111, tokenIndex111
				}
			l112:
				add(ruleVariableName, position110)
			}
			return true
		l109:
			position, tokenIndex = position109, tokenIndex109
			return false
		},
		/* 65 VariableIndex <- <Expression> */
		nil,
		/* 66 Block <- <(_ (Comment / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position115, tokenIndex115 := position, tokenIndex
			{
				position116 := position
				if !_rules[rule_]() {
					goto l115
				}
				{
					position117, tokenIndex117 := position, tokenIndex
					{
						position119 := position
						if !_rules[rule_]() {
							goto l118
						}
						if buffer[position] != rune('#') {
							goto l118
						}
						position++
					l120:
						{
							position121, tokenIndex121 := position, tokenIndex
							{
								position122, tokenIndex122 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l122
								}
								position++
								goto l121
							l122:
								position, tokenIndex = position122, tokenIndex122
							}
							if !matchDot() {
								goto l121
							}
							goto l120
						l121:
							position, tokenIndex = position121, tokenIndex121
						}
						add(ruleComment, position119)
					}
					goto l117
				l118:
					position, tokenIndex = position117, tokenIndex117
					{
						position124 := position
						{
							position125, tokenIndex125 := position, tokenIndex
							{
								position127 := position
								{
									position128 := position
									if !_rules[rule_]() {
										goto l126
									}
									if buffer[position] != rune('b') {
										goto l126
									}
									position++
									if buffer[position] != rune('r') {
										goto l126
									}
									position++
									if buffer[position] != rune('e') {
										goto l126
									}
									position++
									if buffer[position] != rune('a') {
										goto l126
									}
									position++
									if buffer[position] != rune('k') {
										goto l126
									}
									position++
									if !_rules[rule_]() {
										goto l126
									}
									add(ruleBREAK, position128)
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
								add(ruleFlowControlBreak, position127)
							}
							goto l125
						l126:
							position, tokenIndex = position125, tokenIndex125
							{
								position131 := position
								{
									position132 := position
									if !_rules[rule_]() {
										goto l123
									}
									if buffer[position] != rune('c') {
										goto l123
									}
									position++
									if buffer[position] != rune('o') {
										goto l123
									}
									position++
									if buffer[position] != rune('n') {
										goto l123
									}
									position++
									if buffer[position] != rune('t') {
										goto l123
									}
									position++
									if buffer[position] != rune('i') {
										goto l123
									}
									position++
									if buffer[position] != rune('n') {
										goto l123
									}
									position++
									if buffer[position] != rune('u') {
										goto l123
									}
									position++
									if buffer[position] != rune('e') {
										goto l123
									}
									position++
									if !_rules[rule_]() {
										goto l123
									}
									add(ruleCONT, position132)
								}
								{
									position133, tokenIndex133 := position, tokenIndex
									if !_rules[rulePositiveInteger]() {
										goto l133
									}
									goto l134
								l133:
									position, tokenIndex = position133, tokenIndex133
								}
							l134:
								add(ruleFlowControlContinue, position131)
							}
						}
					l125:
						add(ruleFlowControlWord, position124)
					}
					goto l117
				l123:
					position, tokenIndex = position117, tokenIndex117
					{
						position135 := position
						{
							position136, tokenIndex136 := position, tokenIndex
							{
								position138 := position
								if !_rules[ruleSEMI]() {
									goto l137
								}
								add(ruleNOOP, position138)
							}
							goto l136
						l137:
							position, tokenIndex = position136, tokenIndex136
							if !_rules[ruleAssignment]() {
								goto l139
							}
							goto l136
						l139:
							position, tokenIndex = position136, tokenIndex136
							{
								position141 := position
								{
									position142, tokenIndex142 := position, tokenIndex
									{
										position144 := position
										{
											position145 := position
											if !_rules[rule_]() {
												goto l143
											}
											if buffer[position] != rune('u') {
												goto l143
											}
											position++
											if buffer[position] != rune('n') {
												goto l143
											}
											position++
											if buffer[position] != rune('s') {
												goto l143
											}
											position++
											if buffer[position] != rune('e') {
												goto l143
											}
											position++
											if buffer[position] != rune('t') {
												goto l143
											}
											position++
											if !_rules[rule__]() {
												goto l143
											}
											add(ruleUNSET, position145)
										}
										if !_rules[ruleVariableSequence]() {
											goto l143
										}
										add(ruleDirectiveUnset, position144)
									}
									goto l142
								l143:
									position, tokenIndex = position142, tokenIndex142
									{
										position147 := position
										{
											position148 := position
											if !_rules[rule_]() {
												goto l146
											}
											if buffer[position] != rune('i') {
												goto l146
											}
											position++
											if buffer[position] != rune('n') {
												goto l146
											}
											position++
											if buffer[position] != rune('c') {
												goto l146
											}
											position++
											if buffer[position] != rune('l') {
												goto l146
											}
											position++
											if buffer[position] != rune('u') {
												goto l146
											}
											position++
											if buffer[position] != rune('d') {
												goto l146
											}
											position++
											if buffer[position] != rune('e') {
												goto l146
											}
											position++
											if !_rules[rule__]() {
												goto l146
											}
											add(ruleINCLUDE, position148)
										}
										if !_rules[ruleString]() {
											goto l146
										}
										add(ruleDirectiveInclude, position147)
									}
									goto l142
								l146:
									position, tokenIndex = position142, tokenIndex142
									{
										position149 := position
										{
											position150 := position
											if !_rules[rule_]() {
												goto l140
											}
											if buffer[position] != rune('d') {
												goto l140
											}
											position++
											if buffer[position] != rune('e') {
												goto l140
											}
											position++
											if buffer[position] != rune('c') {
												goto l140
											}
											position++
											if buffer[position] != rune('l') {
												goto l140
											}
											position++
											if buffer[position] != rune('a') {
												goto l140
											}
											position++
											if buffer[position] != rune('r') {
												goto l140
											}
											position++
											if buffer[position] != rune('e') {
												goto l140
											}
											position++
											if !_rules[rule__]() {
												goto l140
											}
											add(ruleDECLARE, position150)
										}
										if !_rules[ruleVariableSequence]() {
											goto l140
										}
										add(ruleDirectiveDeclare, position149)
									}
								}
							l142:
								add(ruleDirective, position141)
							}
							goto l136
						l140:
							position, tokenIndex = position136, tokenIndex136
							{
								position152 := position
								if !_rules[ruleIfStanza]() {
									goto l151
								}
							l153:
								{
									position154, tokenIndex154 := position, tokenIndex
									{
										position155 := position
										if !_rules[ruleELSE]() {
											goto l154
										}
										if !_rules[ruleIfStanza]() {
											goto l154
										}
										add(ruleElseIfStanza, position155)
									}
									goto l153
								l154:
									position, tokenIndex = position154, tokenIndex154
								}
								{
									position156, tokenIndex156 := position, tokenIndex
									{
										position158 := position
										if !_rules[ruleELSE]() {
											goto l156
										}
										if !_rules[ruleOPEN]() {
											goto l156
										}
									l159:
										{
											position160, tokenIndex160 := position, tokenIndex
											if !_rules[ruleBlock]() {
												goto l160
											}
											goto l159
										l160:
											position, tokenIndex = position160, tokenIndex160
										}
										if !_rules[ruleCLOSE]() {
											goto l156
										}
										add(ruleElseStanza, position158)
									}
									goto l157
								l156:
									position, tokenIndex = position156, tokenIndex156
								}
							l157:
								add(ruleConditional, position152)
							}
							goto l136
						l151:
							position, tokenIndex = position136, tokenIndex136
							{
								position162 := position
								{
									position163 := position
									if !_rules[rule_]() {
										goto l161
									}
									if buffer[position] != rune('l') {
										goto l161
									}
									position++
									if buffer[position] != rune('o') {
										goto l161
									}
									position++
									if buffer[position] != rune('o') {
										goto l161
									}
									position++
									if buffer[position] != rune('p') {
										goto l161
									}
									position++
									if !_rules[rule_]() {
										goto l161
									}
									add(ruleLOOP, position163)
								}
								{
									position164, tokenIndex164 := position, tokenIndex
									if !_rules[ruleOPEN]() {
										goto l165
									}
								l166:
									{
										position167, tokenIndex167 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l167
										}
										goto l166
									l167:
										position, tokenIndex = position167, tokenIndex167
									}
									if !_rules[ruleCLOSE]() {
										goto l165
									}
									goto l164
								l165:
									position, tokenIndex = position164, tokenIndex164
									{
										position169 := position
										{
											position170 := position
											if !_rules[rule_]() {
												goto l168
											}
											if buffer[position] != rune('c') {
												goto l168
											}
											position++
											if buffer[position] != rune('o') {
												goto l168
											}
											position++
											if buffer[position] != rune('u') {
												goto l168
											}
											position++
											if buffer[position] != rune('n') {
												goto l168
											}
											position++
											if buffer[position] != rune('t') {
												goto l168
											}
											position++
											if !_rules[rule_]() {
												goto l168
											}
											add(ruleCOUNT, position170)
										}
										{
											position171, tokenIndex171 := position, tokenIndex
											if !_rules[ruleInteger]() {
												goto l172
											}
											goto l171
										l172:
											position, tokenIndex = position171, tokenIndex171
											if !_rules[ruleVariable]() {
												goto l168
											}
										}
									l171:
										add(ruleLoopConditionFixedLength, position169)
									}
									if !_rules[ruleOPEN]() {
										goto l168
									}
								l173:
									{
										position174, tokenIndex174 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l174
										}
										goto l173
									l174:
										position, tokenIndex = position174, tokenIndex174
									}
									if !_rules[ruleCLOSE]() {
										goto l168
									}
									goto l164
								l168:
									position, tokenIndex = position164, tokenIndex164
									{
										position176 := position
										{
											position177 := position
											if !_rules[ruleVariableSequence]() {
												goto l175
											}
											add(ruleLoopIterableLHS, position177)
										}
										{
											position178 := position
											if !_rules[rule__]() {
												goto l175
											}
											if buffer[position] != rune('i') {
												goto l175
											}
											position++
											if buffer[position] != rune('n') {
												goto l175
											}
											position++
											if !_rules[rule__]() {
												goto l175
											}
											add(ruleIN, position178)
										}
										{
											position179 := position
											{
												position180, tokenIndex180 := position, tokenIndex
												if !_rules[ruleCommand]() {
													goto l181
												}
												goto l180
											l181:
												position, tokenIndex = position180, tokenIndex180
												if !_rules[ruleVariable]() {
													goto l175
												}
											}
										l180:
											add(ruleLoopIterableRHS, position179)
										}
										add(ruleLoopConditionIterable, position176)
									}
									if !_rules[ruleOPEN]() {
										goto l175
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
										goto l175
									}
									goto l164
								l175:
									position, tokenIndex = position164, tokenIndex164
									{
										position185 := position
										if !_rules[ruleCommand]() {
											goto l184
										}
										if !_rules[ruleSEMI]() {
											goto l184
										}
										if !_rules[ruleConditionalExpression]() {
											goto l184
										}
										if !_rules[ruleSEMI]() {
											goto l184
										}
										if !_rules[ruleCommand]() {
											goto l184
										}
										add(ruleLoopConditionBounded, position185)
									}
									if !_rules[ruleOPEN]() {
										goto l184
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
										goto l184
									}
									goto l164
								l184:
									position, tokenIndex = position164, tokenIndex164
									{
										position188 := position
										if !_rules[ruleConditionalExpression]() {
											goto l161
										}
										add(ruleLoopConditionTruthy, position188)
									}
									if !_rules[ruleOPEN]() {
										goto l161
									}
								l189:
									{
										position190, tokenIndex190 := position, tokenIndex
										if !_rules[ruleBlock]() {
											goto l190
										}
										goto l189
									l190:
										position, tokenIndex = position190, tokenIndex190
									}
									if !_rules[ruleCLOSE]() {
										goto l161
									}
								}
							l164:
								add(ruleLoop, position162)
							}
							goto l136
						l161:
							position, tokenIndex = position136, tokenIndex136
							if !_rules[ruleCommand]() {
								goto l115
							}
						}
					l136:
						add(ruleStatementBlock, position135)
					}
				}
			l117:
				{
					position191, tokenIndex191 := position, tokenIndex
					if !_rules[ruleSEMI]() {
						goto l191
					}
					goto l192
				l191:
					position, tokenIndex = position191, tokenIndex191
				}
			l192:
				if !_rules[rule_]() {
					goto l115
				}
				add(ruleBlock, position116)
			}
			return true
		l115:
			position, tokenIndex = position115, tokenIndex115
			return false
		},
		/* 67 FlowControlWord <- <(FlowControlBreak / FlowControlContinue)> */
		nil,
		/* 68 FlowControlBreak <- <(BREAK PositiveInteger?)> */
		nil,
		/* 69 FlowControlContinue <- <(CONT PositiveInteger?)> */
		nil,
		/* 70 StatementBlock <- <(NOOP / Assignment / Directive / Conditional / Loop / Command)> */
		nil,
		/* 71 Assignment <- <(AssignmentLHS AssignmentOperator AssignmentRHS)> */
		func() bool {
			position197, tokenIndex197 := position, tokenIndex
			{
				position198 := position
				{
					position199 := position
					if !_rules[ruleVariableSequence]() {
						goto l197
					}
					add(ruleAssignmentLHS, position199)
				}
				{
					position200 := position
					if !_rules[rule_]() {
						goto l197
					}
					{
						position201, tokenIndex201 := position, tokenIndex
						{
							position203 := position
							if !_rules[rule_]() {
								goto l202
							}
							if buffer[position] != rune('=') {
								goto l202
							}
							position++
							if !_rules[rule_]() {
								goto l202
							}
							add(ruleAssignEq, position203)
						}
						goto l201
					l202:
						position, tokenIndex = position201, tokenIndex201
						{
							position205 := position
							if !_rules[rule_]() {
								goto l204
							}
							if buffer[position] != rune('*') {
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
							add(ruleStarEq, position205)
						}
						goto l201
					l204:
						position, tokenIndex = position201, tokenIndex201
						{
							position207 := position
							if !_rules[rule_]() {
								goto l206
							}
							if buffer[position] != rune('/') {
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
							add(ruleDivEq, position207)
						}
						goto l201
					l206:
						position, tokenIndex = position201, tokenIndex201
						{
							position209 := position
							if !_rules[rule_]() {
								goto l208
							}
							if buffer[position] != rune('+') {
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
							add(rulePlusEq, position209)
						}
						goto l201
					l208:
						position, tokenIndex = position201, tokenIndex201
						{
							position211 := position
							if !_rules[rule_]() {
								goto l210
							}
							if buffer[position] != rune('-') {
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
							add(ruleMinusEq, position211)
						}
						goto l201
					l210:
						position, tokenIndex = position201, tokenIndex201
						{
							position213 := position
							if !_rules[rule_]() {
								goto l212
							}
							if buffer[position] != rune('&') {
								goto l212
							}
							position++
							if buffer[position] != rune('=') {
								goto l212
							}
							position++
							if !_rules[rule_]() {
								goto l212
							}
							add(ruleAndEq, position213)
						}
						goto l201
					l212:
						position, tokenIndex = position201, tokenIndex201
						{
							position215 := position
							if !_rules[rule_]() {
								goto l214
							}
							if buffer[position] != rune('|') {
								goto l214
							}
							position++
							if buffer[position] != rune('=') {
								goto l214
							}
							position++
							if !_rules[rule_]() {
								goto l214
							}
							add(ruleOrEq, position215)
						}
						goto l201
					l214:
						position, tokenIndex = position201, tokenIndex201
						{
							position216 := position
							if !_rules[rule_]() {
								goto l197
							}
							if buffer[position] != rune('<') {
								goto l197
							}
							position++
							if buffer[position] != rune('<') {
								goto l197
							}
							position++
							if !_rules[rule_]() {
								goto l197
							}
							add(ruleAppend, position216)
						}
					}
				l201:
					if !_rules[rule_]() {
						goto l197
					}
					add(ruleAssignmentOperator, position200)
				}
				{
					position217 := position
					if !_rules[ruleExpressionSequence]() {
						goto l197
					}
					add(ruleAssignmentRHS, position217)
				}
				add(ruleAssignment, position198)
			}
			return true
		l197:
			position, tokenIndex = position197, tokenIndex197
			return false
		},
		/* 72 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 73 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 74 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position220, tokenIndex220 := position, tokenIndex
			{
				position221 := position
			l222:
				{
					position223, tokenIndex223 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l223
					}
					if !_rules[ruleCOMMA]() {
						goto l223
					}
					goto l222
				l223:
					position, tokenIndex = position223, tokenIndex223
				}
				if !_rules[ruleVariable]() {
					goto l220
				}
				add(ruleVariableSequence, position221)
			}
			return true
		l220:
			position, tokenIndex = position220, tokenIndex220
			return false
		},
		/* 75 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
			l226:
				{
					position227, tokenIndex227 := position, tokenIndex
					if !_rules[ruleExpression]() {
						goto l227
					}
					if !_rules[ruleCOMMA]() {
						goto l227
					}
					goto l226
				l227:
					position, tokenIndex = position227, tokenIndex227
				}
				if !_rules[ruleExpression]() {
					goto l224
				}
				add(ruleExpressionSequence, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
			return false
		},
		/* 76 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position228, tokenIndex228 := position, tokenIndex
			{
				position229 := position
				if !_rules[rule_]() {
					goto l228
				}
				{
					position230 := position
					{
						position231 := position
						{
							position232, tokenIndex232 := position, tokenIndex
							if !_rules[ruleType]() {
								goto l233
							}
							goto l232
						l233:
							position, tokenIndex = position232, tokenIndex232
							if !_rules[ruleVariable]() {
								goto l228
							}
						}
					l232:
						add(ruleValueYielding, position231)
					}
					add(ruleExpressionLHS, position230)
				}
				{
					position234, tokenIndex234 := position, tokenIndex
					{
						position236 := position
						{
							position237 := position
							if !_rules[rule_]() {
								goto l234
							}
							{
								position238, tokenIndex238 := position, tokenIndex
								{
									position240 := position
									if !_rules[rule_]() {
										goto l239
									}
									if buffer[position] != rune('*') {
										goto l239
									}
									position++
									if buffer[position] != rune('*') {
										goto l239
									}
									position++
									if !_rules[rule_]() {
										goto l239
									}
									add(ruleExponentiate, position240)
								}
								goto l238
							l239:
								position, tokenIndex = position238, tokenIndex238
								{
									position242 := position
									if !_rules[rule_]() {
										goto l241
									}
									if buffer[position] != rune('*') {
										goto l241
									}
									position++
									if !_rules[rule_]() {
										goto l241
									}
									add(ruleMultiply, position242)
								}
								goto l238
							l241:
								position, tokenIndex = position238, tokenIndex238
								{
									position244 := position
									if !_rules[rule_]() {
										goto l243
									}
									if buffer[position] != rune('/') {
										goto l243
									}
									position++
									if !_rules[rule_]() {
										goto l243
									}
									add(ruleDivide, position244)
								}
								goto l238
							l243:
								position, tokenIndex = position238, tokenIndex238
								{
									position246 := position
									if !_rules[rule_]() {
										goto l245
									}
									if buffer[position] != rune('%') {
										goto l245
									}
									position++
									if !_rules[rule_]() {
										goto l245
									}
									add(ruleModulus, position246)
								}
								goto l238
							l245:
								position, tokenIndex = position238, tokenIndex238
								{
									position248 := position
									if !_rules[rule_]() {
										goto l247
									}
									if buffer[position] != rune('+') {
										goto l247
									}
									position++
									if !_rules[rule_]() {
										goto l247
									}
									add(ruleAdd, position248)
								}
								goto l238
							l247:
								position, tokenIndex = position238, tokenIndex238
								{
									position250 := position
									if !_rules[rule_]() {
										goto l249
									}
									if buffer[position] != rune('-') {
										goto l249
									}
									position++
									if !_rules[rule_]() {
										goto l249
									}
									add(ruleSubtract, position250)
								}
								goto l238
							l249:
								position, tokenIndex = position238, tokenIndex238
								{
									position252 := position
									if !_rules[rule_]() {
										goto l251
									}
									if buffer[position] != rune('&') {
										goto l251
									}
									position++
									if !_rules[rule_]() {
										goto l251
									}
									add(ruleBitwiseAnd, position252)
								}
								goto l238
							l251:
								position, tokenIndex = position238, tokenIndex238
								{
									position254 := position
									if !_rules[rule_]() {
										goto l253
									}
									if buffer[position] != rune('|') {
										goto l253
									}
									position++
									if !_rules[rule_]() {
										goto l253
									}
									add(ruleBitwiseOr, position254)
								}
								goto l238
							l253:
								position, tokenIndex = position238, tokenIndex238
								{
									position256 := position
									if !_rules[rule_]() {
										goto l255
									}
									if buffer[position] != rune('~') {
										goto l255
									}
									position++
									if !_rules[rule_]() {
										goto l255
									}
									add(ruleBitwiseNot, position256)
								}
								goto l238
							l255:
								position, tokenIndex = position238, tokenIndex238
								{
									position257 := position
									if !_rules[rule_]() {
										goto l234
									}
									if buffer[position] != rune('^') {
										goto l234
									}
									position++
									if !_rules[rule_]() {
										goto l234
									}
									add(ruleBitwiseXor, position257)
								}
							}
						l238:
							if !_rules[rule_]() {
								goto l234
							}
							add(ruleOperator, position237)
						}
						if !_rules[ruleExpression]() {
							goto l234
						}
						add(ruleExpressionRHS, position236)
					}
					goto l235
				l234:
					position, tokenIndex = position234, tokenIndex234
				}
			l235:
				if !_rules[rule_]() {
					goto l228
				}
				add(ruleExpression, position229)
			}
			return true
		l228:
			position, tokenIndex = position228, tokenIndex228
			return false
		},
		/* 77 ExpressionLHS <- <ValueYielding> */
		nil,
		/* 78 ExpressionRHS <- <(Operator Expression)> */
		nil,
		/* 79 ValueYielding <- <(Type / Variable)> */
		nil,
		/* 80 Directive <- <(DirectiveUnset / DirectiveInclude / DirectiveDeclare)> */
		nil,
		/* 81 DirectiveUnset <- <(UNSET VariableSequence)> */
		nil,
		/* 82 DirectiveInclude <- <(INCLUDE String)> */
		nil,
		/* 83 DirectiveDeclare <- <(DECLARE VariableSequence)> */
		nil,
		/* 84 Command <- <(_ CommandName (__ ((CommandFirstArg __ CommandSecondArg) / CommandFirstArg / CommandSecondArg))? (_ CommandResultAssignment)?)> */
		func() bool {
			position265, tokenIndex265 := position, tokenIndex
			{
				position266 := position
				if !_rules[rule_]() {
					goto l265
				}
				{
					position267 := position
					{
						position268, tokenIndex268 := position, tokenIndex
						if !_rules[ruleIdentifier]() {
							goto l268
						}
						{
							position270 := position
							if buffer[position] != rune(':') {
								goto l268
							}
							position++
							if buffer[position] != rune(':') {
								goto l268
							}
							position++
							add(ruleSCOPE, position270)
						}
						goto l269
					l268:
						position, tokenIndex = position268, tokenIndex268
					}
				l269:
					if !_rules[ruleIdentifier]() {
						goto l265
					}
					add(ruleCommandName, position267)
				}
				{
					position271, tokenIndex271 := position, tokenIndex
					if !_rules[rule__]() {
						goto l271
					}
					{
						position273, tokenIndex273 := position, tokenIndex
						if !_rules[ruleCommandFirstArg]() {
							goto l274
						}
						if !_rules[rule__]() {
							goto l274
						}
						if !_rules[ruleCommandSecondArg]() {
							goto l274
						}
						goto l273
					l274:
						position, tokenIndex = position273, tokenIndex273
						if !_rules[ruleCommandFirstArg]() {
							goto l275
						}
						goto l273
					l275:
						position, tokenIndex = position273, tokenIndex273
						if !_rules[ruleCommandSecondArg]() {
							goto l271
						}
					}
				l273:
					goto l272
				l271:
					position, tokenIndex = position271, tokenIndex271
				}
			l272:
				{
					position276, tokenIndex276 := position, tokenIndex
					if !_rules[rule_]() {
						goto l276
					}
					{
						position278 := position
						{
							position279 := position
							if !_rules[rule_]() {
								goto l276
							}
							if buffer[position] != rune('-') {
								goto l276
							}
							position++
							if buffer[position] != rune('>') {
								goto l276
							}
							position++
							if !_rules[rule_]() {
								goto l276
							}
							add(ruleASSIGN, position279)
						}
						if !_rules[ruleVariable]() {
							goto l276
						}
						add(ruleCommandResultAssignment, position278)
					}
					goto l277
				l276:
					position, tokenIndex = position276, tokenIndex276
				}
			l277:
				add(ruleCommand, position266)
			}
			return true
		l265:
			position, tokenIndex = position265, tokenIndex265
			return false
		},
		/* 85 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 86 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position281, tokenIndex281 := position, tokenIndex
			{
				position282 := position
				{
					position283, tokenIndex283 := position, tokenIndex
					if !_rules[ruleVariable]() {
						goto l284
					}
					goto l283
				l284:
					position, tokenIndex = position283, tokenIndex283
					if !_rules[ruleType]() {
						goto l281
					}
				}
			l283:
				add(ruleCommandFirstArg, position282)
			}
			return true
		l281:
			position, tokenIndex = position281, tokenIndex281
			return false
		},
		/* 87 CommandSecondArg <- <Object> */
		func() bool {
			position285, tokenIndex285 := position, tokenIndex
			{
				position286 := position
				if !_rules[ruleObject]() {
					goto l285
				}
				add(ruleCommandSecondArg, position286)
			}
			return true
		l285:
			position, tokenIndex = position285, tokenIndex285
			return false
		},
		/* 88 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 89 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 90 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position289, tokenIndex289 := position, tokenIndex
			{
				position290 := position
				{
					position291 := position
					if !_rules[rule_]() {
						goto l289
					}
					if buffer[position] != rune('i') {
						goto l289
					}
					position++
					if buffer[position] != rune('f') {
						goto l289
					}
					position++
					if !_rules[rule_]() {
						goto l289
					}
					add(ruleIF, position291)
				}
				if !_rules[ruleConditionalExpression]() {
					goto l289
				}
				if !_rules[ruleOPEN]() {
					goto l289
				}
			l292:
				{
					position293, tokenIndex293 := position, tokenIndex
					if !_rules[ruleBlock]() {
						goto l293
					}
					goto l292
				l293:
					position, tokenIndex = position293, tokenIndex293
				}
				if !_rules[ruleCLOSE]() {
					goto l289
				}
				add(ruleIfStanza, position290)
			}
			return true
		l289:
			position, tokenIndex = position289, tokenIndex289
			return false
		},
		/* 91 ElseIfStanza <- <(ELSE IfStanza)> */
		nil,
		/* 92 ElseStanza <- <(ELSE OPEN Block* CLOSE)> */
		nil,
		/* 93 Loop <- <(LOOP ((OPEN Block* CLOSE) / (LoopConditionFixedLength OPEN Block* CLOSE) / (LoopConditionIterable OPEN Block* CLOSE) / (LoopConditionBounded OPEN Block* CLOSE) / (LoopConditionTruthy OPEN Block* CLOSE)))> */
		nil,
		/* 94 LoopConditionFixedLength <- <(COUNT (Integer / Variable))> */
		nil,
		/* 95 LoopConditionIterable <- <(LoopIterableLHS IN LoopIterableRHS)> */
		nil,
		/* 96 LoopIterableLHS <- <VariableSequence> */
		nil,
		/* 97 LoopIterableRHS <- <(Command / Variable)> */
		nil,
		/* 98 LoopConditionBounded <- <(Command SEMI ConditionalExpression SEMI Command)> */
		nil,
		/* 99 LoopConditionTruthy <- <ConditionalExpression> */
		nil,
		/* 100 ConditionalExpression <- <(NOT? (ConditionWithAssignment / ConditionWithCommand / ConditionWithRegex / ConditionWithComparator))> */
		func() bool {
			position303, tokenIndex303 := position, tokenIndex
			{
				position304 := position
				{
					position305, tokenIndex305 := position, tokenIndex
					{
						position307 := position
						if !_rules[rule_]() {
							goto l305
						}
						if buffer[position] != rune('n') {
							goto l305
						}
						position++
						if buffer[position] != rune('o') {
							goto l305
						}
						position++
						if buffer[position] != rune('t') {
							goto l305
						}
						position++
						if !_rules[rule__]() {
							goto l305
						}
						add(ruleNOT, position307)
					}
					goto l306
				l305:
					position, tokenIndex = position305, tokenIndex305
				}
			l306:
				{
					position308, tokenIndex308 := position, tokenIndex
					{
						position310 := position
						if !_rules[ruleAssignment]() {
							goto l309
						}
						if !_rules[ruleSEMI]() {
							goto l309
						}
						if !_rules[ruleConditionalExpression]() {
							goto l309
						}
						add(ruleConditionWithAssignment, position310)
					}
					goto l308
				l309:
					position, tokenIndex = position308, tokenIndex308
					{
						position312 := position
						if !_rules[ruleCommand]() {
							goto l311
						}
						{
							position313, tokenIndex313 := position, tokenIndex
							if !_rules[ruleSEMI]() {
								goto l313
							}
							if !_rules[ruleConditionalExpression]() {
								goto l313
							}
							goto l314
						l313:
							position, tokenIndex = position313, tokenIndex313
						}
					l314:
						add(ruleConditionWithCommand, position312)
					}
					goto l308
				l311:
					position, tokenIndex = position308, tokenIndex308
					{
						position316 := position
						if !_rules[ruleExpression]() {
							goto l315
						}
						{
							position317 := position
							{
								position318, tokenIndex318 := position, tokenIndex
								{
									position320 := position
									if !_rules[rule_]() {
										goto l319
									}
									if buffer[position] != rune('=') {
										goto l319
									}
									position++
									if buffer[position] != rune('~') {
										goto l319
									}
									position++
									if !_rules[rule_]() {
										goto l319
									}
									add(ruleMatch, position320)
								}
								goto l318
							l319:
								position, tokenIndex = position318, tokenIndex318
								{
									position321 := position
									if !_rules[rule_]() {
										goto l315
									}
									if buffer[position] != rune('!') {
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
									add(ruleUnmatch, position321)
								}
							}
						l318:
							add(ruleMatchOperator, position317)
						}
						if !_rules[ruleRegularExpression]() {
							goto l315
						}
						add(ruleConditionWithRegex, position316)
					}
					goto l308
				l315:
					position, tokenIndex = position308, tokenIndex308
					{
						position322 := position
						{
							position323 := position
							if !_rules[ruleExpression]() {
								goto l303
							}
							add(ruleConditionWithComparatorLHS, position323)
						}
						{
							position324, tokenIndex324 := position, tokenIndex
							{
								position326 := position
								{
									position327 := position
									if !_rules[rule_]() {
										goto l324
									}
									{
										position328, tokenIndex328 := position, tokenIndex
										{
											position330 := position
											if !_rules[rule_]() {
												goto l329
											}
											if buffer[position] != rune('=') {
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
											add(ruleEquality, position330)
										}
										goto l328
									l329:
										position, tokenIndex = position328, tokenIndex328
										{
											position332 := position
											if !_rules[rule_]() {
												goto l331
											}
											if buffer[position] != rune('!') {
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
											add(ruleNonEquality, position332)
										}
										goto l328
									l331:
										position, tokenIndex = position328, tokenIndex328
										{
											position334 := position
											if !_rules[rule_]() {
												goto l333
											}
											if buffer[position] != rune('>') {
												goto l333
											}
											position++
											if buffer[position] != rune('=') {
												goto l333
											}
											position++
											if !_rules[rule_]() {
												goto l333
											}
											add(ruleGreaterEqual, position334)
										}
										goto l328
									l333:
										position, tokenIndex = position328, tokenIndex328
										{
											position336 := position
											if !_rules[rule_]() {
												goto l335
											}
											if buffer[position] != rune('<') {
												goto l335
											}
											position++
											if buffer[position] != rune('=') {
												goto l335
											}
											position++
											if !_rules[rule_]() {
												goto l335
											}
											add(ruleLessEqual, position336)
										}
										goto l328
									l335:
										position, tokenIndex = position328, tokenIndex328
										{
											position338 := position
											if !_rules[rule_]() {
												goto l337
											}
											if buffer[position] != rune('>') {
												goto l337
											}
											position++
											if !_rules[rule_]() {
												goto l337
											}
											add(ruleGreaterThan, position338)
										}
										goto l328
									l337:
										position, tokenIndex = position328, tokenIndex328
										{
											position340 := position
											if !_rules[rule_]() {
												goto l339
											}
											if buffer[position] != rune('<') {
												goto l339
											}
											position++
											if !_rules[rule_]() {
												goto l339
											}
											add(ruleLessThan, position340)
										}
										goto l328
									l339:
										position, tokenIndex = position328, tokenIndex328
										{
											position342 := position
											if !_rules[rule_]() {
												goto l341
											}
											if buffer[position] != rune('i') {
												goto l341
											}
											position++
											if buffer[position] != rune('n') {
												goto l341
											}
											position++
											if !_rules[rule_]() {
												goto l341
											}
											add(ruleMembership, position342)
										}
										goto l328
									l341:
										position, tokenIndex = position328, tokenIndex328
										{
											position343 := position
											if !_rules[rule_]() {
												goto l324
											}
											if buffer[position] != rune('n') {
												goto l324
											}
											position++
											if buffer[position] != rune('o') {
												goto l324
											}
											position++
											if buffer[position] != rune('t') {
												goto l324
											}
											position++
											if !_rules[rule__]() {
												goto l324
											}
											if buffer[position] != rune('i') {
												goto l324
											}
											position++
											if buffer[position] != rune('n') {
												goto l324
											}
											position++
											if !_rules[rule_]() {
												goto l324
											}
											add(ruleNonMembership, position343)
										}
									}
								l328:
									if !_rules[rule_]() {
										goto l324
									}
									add(ruleComparisonOperator, position327)
								}
								if !_rules[ruleExpression]() {
									goto l324
								}
								add(ruleConditionWithComparatorRHS, position326)
							}
							goto l325
						l324:
							position, tokenIndex = position324, tokenIndex324
						}
					l325:
						add(ruleConditionWithComparator, position322)
					}
				}
			l308:
				add(ruleConditionalExpression, position304)
			}
			return true
		l303:
			position, tokenIndex = position303, tokenIndex303
			return false
		},
		/* 101 ConditionWithAssignment <- <(Assignment SEMI ConditionalExpression)> */
		nil,
		/* 102 ConditionWithCommand <- <(Command (SEMI ConditionalExpression)?)> */
		nil,
		/* 103 ConditionWithRegex <- <(Expression MatchOperator RegularExpression)> */
		nil,
		/* 104 ConditionWithComparator <- <(ConditionWithComparatorLHS ConditionWithComparatorRHS?)> */
		nil,
		/* 105 ConditionWithComparatorLHS <- <Expression> */
		nil,
		/* 106 ConditionWithComparatorRHS <- <(ComparisonOperator Expression)> */
		nil,
		/* 107 ScalarType <- <(Boolean / Float / Integer / String / NullValue)> */
		nil,
		/* 108 Identifier <- <(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / ([0-9] / [0-9]) / '_')*)> */
		func() bool {
			position351, tokenIndex351 := position, tokenIndex
			{
				position352 := position
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
					if buffer[position] != rune('_') {
						goto l351
					}
					position++
				}
			l353:
			l356:
				{
					position357, tokenIndex357 := position, tokenIndex
					{
						position358, tokenIndex358 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l359
						}
						position++
						goto l358
					l359:
						position, tokenIndex = position358, tokenIndex358
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l360
						}
						position++
						goto l358
					l360:
						position, tokenIndex = position358, tokenIndex358
						{
							position362, tokenIndex362 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex = position362, tokenIndex362
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l361
							}
							position++
						}
					l362:
						goto l358
					l361:
						position, tokenIndex = position358, tokenIndex358
						if buffer[position] != rune('_') {
							goto l357
						}
						position++
					}
				l358:
					goto l356
				l357:
					position, tokenIndex = position357, tokenIndex357
				}
				add(ruleIdentifier, position352)
			}
			return true
		l351:
			position, tokenIndex = position351, tokenIndex351
			return false
		},
		/* 109 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 110 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 111 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position366, tokenIndex366 := position, tokenIndex
			{
				position367 := position
				{
					position368, tokenIndex368 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l368
					}
					position++
					goto l369
				l368:
					position, tokenIndex = position368, tokenIndex368
				}
			l369:
				if !_rules[rulePositiveInteger]() {
					goto l366
				}
				add(ruleInteger, position367)
			}
			return true
		l366:
			position, tokenIndex = position366, tokenIndex366
			return false
		},
		/* 112 PositiveInteger <- <[0-9]+> */
		func() bool {
			position370, tokenIndex370 := position, tokenIndex
			{
				position371 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l370
				}
				position++
			l372:
				{
					position373, tokenIndex373 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex = position373, tokenIndex373
				}
				add(rulePositiveInteger, position371)
			}
			return true
		l370:
			position, tokenIndex = position370, tokenIndex370
			return false
		},
		/* 113 String <- <(HeredocTriquote / StringLiteral / StringInterpolated / Heredoc)> */
		func() bool {
			position374, tokenIndex374 := position, tokenIndex
			{
				position375 := position
				{
					position376, tokenIndex376 := position, tokenIndex
					{
						position378 := position
						if !_rules[ruleTRIQUOT]() {
							goto l377
						}
						{
							position379 := position
						l380:
							{
								position381, tokenIndex381 := position, tokenIndex
								{
									position382, tokenIndex382 := position, tokenIndex
									if !_rules[ruleTRIQUOT]() {
										goto l382
									}
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
							add(ruleTriquoteBody, position379)
						}
						if !_rules[ruleTRIQUOT]() {
							goto l377
						}
						add(ruleHeredocTriquote, position378)
					}
					goto l376
				l377:
					position, tokenIndex = position376, tokenIndex376
					{
						position384 := position
						if buffer[position] != rune('\'') {
							goto l383
						}
						position++
					l385:
						{
							position386, tokenIndex386 := position, tokenIndex
							{
								position387, tokenIndex387 := position, tokenIndex
								if buffer[position] != rune('\'') {
									goto l387
								}
								position++
								goto l386
							l387:
								position, tokenIndex = position387, tokenIndex387
							}
							if !matchDot() {
								goto l386
							}
							goto l385
						l386:
							position, tokenIndex = position386, tokenIndex386
						}
						if buffer[position] != rune('\'') {
							goto l383
						}
						position++
						add(ruleStringLiteral, position384)
					}
					goto l376
				l383:
					position, tokenIndex = position376, tokenIndex376
					{
						position389 := position
						if buffer[position] != rune('"') {
							goto l388
						}
						position++
					l390:
						{
							position391, tokenIndex391 := position, tokenIndex
							{
								position392, tokenIndex392 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l392
								}
								position++
								goto l391
							l392:
								position, tokenIndex = position392, tokenIndex392
							}
							if !matchDot() {
								goto l391
							}
							goto l390
						l391:
							position, tokenIndex = position391, tokenIndex391
						}
						if buffer[position] != rune('"') {
							goto l388
						}
						position++
						add(ruleStringInterpolated, position389)
					}
					goto l376
				l388:
					position, tokenIndex = position376, tokenIndex376
					{
						position393 := position
						{
							position394 := position
							if !_rules[rule_]() {
								goto l374
							}
							if buffer[position] != rune('b') {
								goto l374
							}
							position++
							if buffer[position] != rune('e') {
								goto l374
							}
							position++
							if buffer[position] != rune('g') {
								goto l374
							}
							position++
							if buffer[position] != rune('i') {
								goto l374
							}
							position++
							if buffer[position] != rune('n') {
								goto l374
							}
							position++
							add(ruleBEGIN, position394)
						}
						if !_rules[ruleNL]() {
							goto l374
						}
						{
							position395 := position
						l396:
							{
								position397, tokenIndex397 := position, tokenIndex
								{
									position398, tokenIndex398 := position, tokenIndex
									if !_rules[ruleEND]() {
										goto l398
									}
									goto l397
								l398:
									position, tokenIndex = position398, tokenIndex398
								}
								if !matchDot() {
									goto l397
								}
								goto l396
							l397:
								position, tokenIndex = position397, tokenIndex397
							}
							add(ruleHeredocBody, position395)
						}
						if !_rules[ruleEND]() {
							goto l374
						}
						add(ruleHeredoc, position393)
					}
				}
			l376:
				add(ruleString, position375)
			}
			return true
		l374:
			position, tokenIndex = position374, tokenIndex374
			return false
		},
		/* 114 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		nil,
		/* 115 StringInterpolated <- <('"' (!'"' .)* '"')> */
		nil,
		/* 116 HeredocTriquote <- <(TRIQUOT TriquoteBody TRIQUOT)> */
		nil,
		/* 117 TriquoteBody <- <(!TRIQUOT .)*> */
		nil,
		/* 118 HeredocBody <- <(!END .)*> */
		nil,
		/* 119 Heredoc <- <(BEGIN NL HeredocBody END)> */
		nil,
		/* 120 NullValue <- <('n' 'u' 'l' 'l')> */
		nil,
		/* 121 Object <- <(OPEN (_ KeyValuePair _)* CLOSE)> */
		func() bool {
			position406, tokenIndex406 := position, tokenIndex
			{
				position407 := position
				if !_rules[ruleOPEN]() {
					goto l406
				}
			l408:
				{
					position409, tokenIndex409 := position, tokenIndex
					if !_rules[rule_]() {
						goto l409
					}
					{
						position410 := position
						{
							position411 := position
							if !_rules[ruleIdentifier]() {
								goto l409
							}
							add(ruleKey, position411)
						}
						{
							position412 := position
							if !_rules[rule_]() {
								goto l409
							}
							if buffer[position] != rune(':') {
								goto l409
							}
							position++
							if !_rules[rule_]() {
								goto l409
							}
							add(ruleCOLON, position412)
						}
						{
							position413 := position
							{
								position414, tokenIndex414 := position, tokenIndex
								if !_rules[ruleArray]() {
									goto l415
								}
								goto l414
							l415:
								position, tokenIndex = position414, tokenIndex414
								if !_rules[ruleObject]() {
									goto l416
								}
								goto l414
							l416:
								position, tokenIndex = position414, tokenIndex414
								if !_rules[ruleExpression]() {
									goto l409
								}
							}
						l414:
							add(ruleKValue, position413)
						}
						{
							position417, tokenIndex417 := position, tokenIndex
							if !_rules[ruleCOMMA]() {
								goto l417
							}
							goto l418
						l417:
							position, tokenIndex = position417, tokenIndex417
						}
					l418:
						add(ruleKeyValuePair, position410)
					}
					if !_rules[rule_]() {
						goto l409
					}
					goto l408
				l409:
					position, tokenIndex = position409, tokenIndex409
				}
				if !_rules[ruleCLOSE]() {
					goto l406
				}
				add(ruleObject, position407)
			}
			return true
		l406:
			position, tokenIndex = position406, tokenIndex406
			return false
		},
		/* 122 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position419, tokenIndex419 := position, tokenIndex
			{
				position420 := position
				if buffer[position] != rune('[') {
					goto l419
				}
				position++
				if !_rules[rule_]() {
					goto l419
				}
				if !_rules[ruleExpressionSequence]() {
					goto l419
				}
				{
					position421, tokenIndex421 := position, tokenIndex
					if !_rules[ruleCOMMA]() {
						goto l421
					}
					goto l422
				l421:
					position, tokenIndex = position421, tokenIndex421
				}
			l422:
				if buffer[position] != rune(']') {
					goto l419
				}
				position++
				add(ruleArray, position420)
			}
			return true
		l419:
			position, tokenIndex = position419, tokenIndex419
			return false
		},
		/* 123 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position423, tokenIndex423 := position, tokenIndex
			{
				position424 := position
				if buffer[position] != rune('/') {
					goto l423
				}
				position++
				{
					position427, tokenIndex427 := position, tokenIndex
					if buffer[position] != rune('/') {
						goto l427
					}
					position++
					goto l423
				l427:
					position, tokenIndex = position427, tokenIndex427
				}
				if !matchDot() {
					goto l423
				}
			l425:
				{
					position426, tokenIndex426 := position, tokenIndex
					{
						position428, tokenIndex428 := position, tokenIndex
						if buffer[position] != rune('/') {
							goto l428
						}
						position++
						goto l426
					l428:
						position, tokenIndex = position428, tokenIndex428
					}
					if !matchDot() {
						goto l426
					}
					goto l425
				l426:
					position, tokenIndex = position426, tokenIndex426
				}
				if buffer[position] != rune('/') {
					goto l423
				}
				position++
			l429:
				{
					position430, tokenIndex430 := position, tokenIndex
					{
						position431, tokenIndex431 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l432
						}
						position++
						goto l431
					l432:
						position, tokenIndex = position431, tokenIndex431
						if buffer[position] != rune('l') {
							goto l433
						}
						position++
						goto l431
					l433:
						position, tokenIndex = position431, tokenIndex431
						if buffer[position] != rune('m') {
							goto l434
						}
						position++
						goto l431
					l434:
						position, tokenIndex = position431, tokenIndex431
						if buffer[position] != rune('s') {
							goto l435
						}
						position++
						goto l431
					l435:
						position, tokenIndex = position431, tokenIndex431
						if buffer[position] != rune('u') {
							goto l430
						}
						position++
					}
				l431:
					goto l429
				l430:
					position, tokenIndex = position430, tokenIndex430
				}
				add(ruleRegularExpression, position424)
			}
			return true
		l423:
			position, tokenIndex = position423, tokenIndex423
			return false
		},
		/* 124 KeyValuePair <- <(Key COLON KValue COMMA?)> */
		nil,
		/* 125 Key <- <Identifier> */
		nil,
		/* 126 KValue <- <(Array / Object / Expression)> */
		nil,
		/* 127 Type <- <(Array / Object / RegularExpression / ScalarType)> */
		func() bool {
			position439, tokenIndex439 := position, tokenIndex
			{
				position440 := position
				{
					position441, tokenIndex441 := position, tokenIndex
					if !_rules[ruleArray]() {
						goto l442
					}
					goto l441
				l442:
					position, tokenIndex = position441, tokenIndex441
					if !_rules[ruleObject]() {
						goto l443
					}
					goto l441
				l443:
					position, tokenIndex = position441, tokenIndex441
					if !_rules[ruleRegularExpression]() {
						goto l444
					}
					goto l441
				l444:
					position, tokenIndex = position441, tokenIndex441
					{
						position445 := position
						{
							position446, tokenIndex446 := position, tokenIndex
							{
								position448 := position
								{
									position449, tokenIndex449 := position, tokenIndex
									if buffer[position] != rune('t') {
										goto l450
									}
									position++
									if buffer[position] != rune('r') {
										goto l450
									}
									position++
									if buffer[position] != rune('u') {
										goto l450
									}
									position++
									if buffer[position] != rune('e') {
										goto l450
									}
									position++
									goto l449
								l450:
									position, tokenIndex = position449, tokenIndex449
									if buffer[position] != rune('f') {
										goto l447
									}
									position++
									if buffer[position] != rune('a') {
										goto l447
									}
									position++
									if buffer[position] != rune('l') {
										goto l447
									}
									position++
									if buffer[position] != rune('s') {
										goto l447
									}
									position++
									if buffer[position] != rune('e') {
										goto l447
									}
									position++
								}
							l449:
								add(ruleBoolean, position448)
							}
							goto l446
						l447:
							position, tokenIndex = position446, tokenIndex446
							{
								position452 := position
								if !_rules[ruleInteger]() {
									goto l451
								}
								{
									position453, tokenIndex453 := position, tokenIndex
									if buffer[position] != rune('.') {
										goto l453
									}
									position++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l453
									}
									position++
								l455:
									{
										position456, tokenIndex456 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l456
										}
										position++
										goto l455
									l456:
										position, tokenIndex = position456, tokenIndex456
									}
									goto l454
								l453:
									position, tokenIndex = position453, tokenIndex453
								}
							l454:
								add(ruleFloat, position452)
							}
							goto l446
						l451:
							position, tokenIndex = position446, tokenIndex446
							if !_rules[ruleInteger]() {
								goto l457
							}
							goto l446
						l457:
							position, tokenIndex = position446, tokenIndex446
							if !_rules[ruleString]() {
								goto l458
							}
							goto l446
						l458:
							position, tokenIndex = position446, tokenIndex446
							{
								position459 := position
								if buffer[position] != rune('n') {
									goto l439
								}
								position++
								if buffer[position] != rune('u') {
									goto l439
								}
								position++
								if buffer[position] != rune('l') {
									goto l439
								}
								position++
								if buffer[position] != rune('l') {
									goto l439
								}
								position++
								add(ruleNullValue, position459)
							}
						}
					l446:
						add(ruleScalarType, position445)
					}
				}
			l441:
				add(ruleType, position440)
			}
			return true
		l439:
			position, tokenIndex = position439, tokenIndex439
			return false
		},
	}
	p.rules = _rules
}
