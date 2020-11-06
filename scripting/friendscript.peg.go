package scripting

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

	rulePre
	ruleIn
	ruleSuf
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

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
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
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type Friendscript struct {
	runtime

	Buffer string
	buffer []rune
	rules  [124]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
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
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *Friendscript) Highlighter() {
	p.PrintSyntax()
}

func (p *Friendscript) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
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
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rule_]() {
					goto l0
				}
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					{
						position4 := position
						depth++
						if buffer[position] != rune('#') {
							goto l2
						}
						position++
						if buffer[position] != rune('!') {
							goto l2
						}
						position++
						{
							position7, tokenIndex7, depth7 := position, tokenIndex, depth
							if buffer[position] != rune('\n') {
								goto l7
							}
							position++
							goto l2
						l7:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
						}
						if !matchDot() {
							goto l2
						}
					l5:
						{
							position6, tokenIndex6, depth6 := position, tokenIndex, depth
							{
								position8, tokenIndex8, depth8 := position, tokenIndex, depth
								if buffer[position] != rune('\n') {
									goto l8
								}
								position++
								goto l6
							l8:
								position, tokenIndex, depth = position8, tokenIndex8, depth8
							}
							if !matchDot() {
								goto l6
							}
							goto l5
						l6:
							position, tokenIndex, depth = position6, tokenIndex6, depth6
						}
						if buffer[position] != rune('\n') {
							goto l2
						}
						position++
						depth--
						add(ruleSHEBANG, position4)
					}
					goto l3
				l2:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
				}
			l3:
				if !_rules[rule_]() {
					goto l0
				}
			l9:
				{
					position10, tokenIndex10, depth10 := position, tokenIndex, depth
					if !_rules[ruleBlock]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex, depth = position10, tokenIndex10, depth10
				}
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !matchDot() {
						goto l11
					}
					goto l0
				l11:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
				}
				depth--
				add(ruleFriendscript, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 _ <- <(' ' / '\t' / '\r' / '\n')*> */
		func() bool {
			{
				position13 := position
				depth++
			l14:
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					{
						position16, tokenIndex16, depth16 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l17
						}
						position++
						goto l16
					l17:
						position, tokenIndex, depth = position16, tokenIndex16, depth16
						if buffer[position] != rune('\t') {
							goto l18
						}
						position++
						goto l16
					l18:
						position, tokenIndex, depth = position16, tokenIndex16, depth16
						if buffer[position] != rune('\r') {
							goto l19
						}
						position++
						goto l16
					l19:
						position, tokenIndex, depth = position16, tokenIndex16, depth16
						if buffer[position] != rune('\n') {
							goto l15
						}
						position++
					}
				l16:
					goto l14
				l15:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
				}
				depth--
				add(rule_, position13)
			}
			return true
		},
		/* 2 __ <- <(' ' / '\t' / '\r' / '\n')+> */
		func() bool {
			position20, tokenIndex20, depth20 := position, tokenIndex, depth
			{
				position21 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l25
					}
					position++
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('\t') {
						goto l26
					}
					position++
					goto l24
				l26:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('\r') {
						goto l27
					}
					position++
					goto l24
				l27:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('\n') {
						goto l20
					}
					position++
				}
			l24:
			l22:
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					{
						position28, tokenIndex28, depth28 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l29
						}
						position++
						goto l28
					l29:
						position, tokenIndex, depth = position28, tokenIndex28, depth28
						if buffer[position] != rune('\t') {
							goto l30
						}
						position++
						goto l28
					l30:
						position, tokenIndex, depth = position28, tokenIndex28, depth28
						if buffer[position] != rune('\r') {
							goto l31
						}
						position++
						goto l28
					l31:
						position, tokenIndex, depth = position28, tokenIndex28, depth28
						if buffer[position] != rune('\n') {
							goto l23
						}
						position++
					}
				l28:
					goto l22
				l23:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
				}
				depth--
				add(rule__, position21)
			}
			return true
		l20:
			position, tokenIndex, depth = position20, tokenIndex20, depth20
			return false
		},
		/* 3 ASSIGN <- <(_ ('-' '>') _)> */
		nil,
		/* 4 TRIQUOT <- <(_ ('"' '"' '"') _)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
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
				depth--
				add(ruleTRIQUOT, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 5 BREAK <- <(_ ('b' 'r' 'e' 'a' 'k') _)> */
		nil,
		/* 6 CLOSE <- <(_ '}' _)> */
		func() bool {
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{
				position37 := position
				depth++
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
				depth--
				add(ruleCLOSE, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 7 COLON <- <(_ ':' _)> */
		nil,
		/* 8 COMMA <- <(_ ',' _)> */
		func() bool {
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
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
				depth--
				add(ruleCOMMA, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
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
			position46, tokenIndex46, depth46 := position, tokenIndex, depth
			{
				position47 := position
				depth++
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
				depth--
				add(ruleELSE, position47)
			}
			return true
		l46:
			position, tokenIndex, depth = position46, tokenIndex46, depth46
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
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
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
				depth--
				add(ruleOPEN, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 22 SCOPE <- <(':' ':')> */
		nil,
		/* 23 SEMI <- <(_ ';' _)> */
		func() bool {
			position57, tokenIndex57, depth57 := position, tokenIndex, depth
			{
				position58 := position
				depth++
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
				depth--
				add(ruleSEMI, position58)
			}
			return true
		l57:
			position, tokenIndex, depth = position57, tokenIndex57, depth57
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
			position63, tokenIndex63, depth63 := position, tokenIndex, depth
			{
				position64 := position
				depth++
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l67
					}
					position++
					goto l65
				l67:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if buffer[position] != rune('_') {
						goto l63
					}
					position++
				}
			l65:
			l68:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l72
						}
						position++
						goto l70
					l72:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						{
							position74, tokenIndex74, depth74 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l75
							}
							position++
							goto l74
						l75:
							position, tokenIndex, depth = position74, tokenIndex74, depth74
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l73
							}
							position++
						}
					l74:
						goto l70
					l73:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						if buffer[position] != rune('_') {
							goto l69
						}
						position++
					}
				l70:
					goto l68
				l69:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
				}
				depth--
				add(ruleIdentifier, position64)
			}
			return true
		l63:
			position, tokenIndex, depth = position63, tokenIndex63, depth63
			return false
		},
		/* 29 Float <- <(Integer ('.' [0-9]+)?)> */
		nil,
		/* 30 Boolean <- <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		nil,
		/* 31 Integer <- <('-'? PositiveInteger)> */
		func() bool {
			position78, tokenIndex78, depth78 := position, tokenIndex, depth
			{
				position79 := position
				depth++
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l80
					}
					position++
					goto l81
				l80:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
				}
			l81:
				if !_rules[rulePositiveInteger]() {
					goto l78
				}
				depth--
				add(ruleInteger, position79)
			}
			return true
		l78:
			position, tokenIndex, depth = position78, tokenIndex78, depth78
			return false
		},
		/* 32 PositiveInteger <- <[0-9]+> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l82
				}
				position++
			l84:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
				}
				depth--
				add(rulePositiveInteger, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
			return false
		},
		/* 33 String <- <(Triquote / StringLiteral / StringInterpolated)> */
		func() bool {
			position86, tokenIndex86, depth86 := position, tokenIndex, depth
			{
				position87 := position
				depth++
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					{
						position90 := position
						depth++
						if !_rules[ruleTRIQUOT]() {
							goto l89
						}
						{
							position91 := position
							depth++
						l92:
							{
								position93, tokenIndex93, depth93 := position, tokenIndex, depth
								{
									position94, tokenIndex94, depth94 := position, tokenIndex, depth
									if !_rules[ruleTRIQUOT]() {
										goto l94
									}
									goto l93
								l94:
									position, tokenIndex, depth = position94, tokenIndex94, depth94
								}
								if !matchDot() {
									goto l93
								}
								goto l92
							l93:
								position, tokenIndex, depth = position93, tokenIndex93, depth93
							}
							depth--
							add(ruleTriquoteBody, position91)
						}
						if !_rules[ruleTRIQUOT]() {
							goto l89
						}
						depth--
						add(ruleTriquote, position90)
					}
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if !_rules[ruleStringLiteral]() {
						goto l95
					}
					goto l88
				l95:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if !_rules[ruleStringInterpolated]() {
						goto l86
					}
				}
			l88:
				depth--
				add(ruleString, position87)
			}
			return true
		l86:
			position, tokenIndex, depth = position86, tokenIndex86, depth86
			return false
		},
		/* 34 StringLiteral <- <('\'' (!'\'' .)* '\'')> */
		func() bool {
			position96, tokenIndex96, depth96 := position, tokenIndex, depth
			{
				position97 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l96
				}
				position++
			l98:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					{
						position100, tokenIndex100, depth100 := position, tokenIndex, depth
						if buffer[position] != rune('\'') {
							goto l100
						}
						position++
						goto l99
					l100:
						position, tokenIndex, depth = position100, tokenIndex100, depth100
					}
					if !matchDot() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
				}
				if buffer[position] != rune('\'') {
					goto l96
				}
				position++
				depth--
				add(ruleStringLiteral, position97)
			}
			return true
		l96:
			position, tokenIndex, depth = position96, tokenIndex96, depth96
			return false
		},
		/* 35 StringInterpolated <- <('"' (!'"' .)* '"')> */
		func() bool {
			position101, tokenIndex101, depth101 := position, tokenIndex, depth
			{
				position102 := position
				depth++
				if buffer[position] != rune('"') {
					goto l101
				}
				position++
			l103:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					{
						position105, tokenIndex105, depth105 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l105
						}
						position++
						goto l104
					l105:
						position, tokenIndex, depth = position105, tokenIndex105, depth105
					}
					if !matchDot() {
						goto l104
					}
					goto l103
				l104:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
				}
				if buffer[position] != rune('"') {
					goto l101
				}
				position++
				depth--
				add(ruleStringInterpolated, position102)
			}
			return true
		l101:
			position, tokenIndex, depth = position101, tokenIndex101, depth101
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
			position109, tokenIndex109, depth109 := position, tokenIndex, depth
			{
				position110 := position
				depth++
				if !_rules[ruleOPEN]() {
					goto l109
				}
			l111:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if !_rules[rule_]() {
						goto l112
					}
					{
						position113 := position
						depth++
						{
							position114 := position
							depth++
							{
								position115, tokenIndex115, depth115 := position, tokenIndex, depth
								if !_rules[ruleIdentifier]() {
									goto l116
								}
								goto l115
							l116:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
								if !_rules[ruleStringLiteral]() {
									goto l117
								}
								goto l115
							l117:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
								if !_rules[ruleStringInterpolated]() {
									goto l112
								}
							}
						l115:
							depth--
							add(ruleKey, position114)
						}
						{
							position118 := position
							depth++
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
							depth--
							add(ruleCOLON, position118)
						}
						{
							position119 := position
							depth++
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if !_rules[ruleArray]() {
									goto l121
								}
								goto l120
							l121:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if !_rules[ruleObject]() {
									goto l122
								}
								goto l120
							l122:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if !_rules[ruleExpression]() {
									goto l112
								}
							}
						l120:
							depth--
							add(ruleKValue, position119)
						}
						{
							position123, tokenIndex123, depth123 := position, tokenIndex, depth
							if !_rules[ruleCOMMA]() {
								goto l123
							}
							goto l124
						l123:
							position, tokenIndex, depth = position123, tokenIndex123, depth123
						}
					l124:
						depth--
						add(ruleKeyValuePair, position113)
					}
					if !_rules[rule_]() {
						goto l112
					}
					goto l111
				l112:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
				}
				if !_rules[ruleCLOSE]() {
					goto l109
				}
				depth--
				add(ruleObject, position110)
			}
			return true
		l109:
			position, tokenIndex, depth = position109, tokenIndex109, depth109
			return false
		},
		/* 40 Array <- <('[' _ ExpressionSequence COMMA? ']')> */
		func() bool {
			position125, tokenIndex125, depth125 := position, tokenIndex, depth
			{
				position126 := position
				depth++
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
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if !_rules[ruleCOMMA]() {
						goto l127
					}
					goto l128
				l127:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
				}
			l128:
				if buffer[position] != rune(']') {
					goto l125
				}
				position++
				depth--
				add(ruleArray, position126)
			}
			return true
		l125:
			position, tokenIndex, depth = position125, tokenIndex125, depth125
			return false
		},
		/* 41 RegularExpression <- <('/' (!'/' .)+ '/' ('i' / 'l' / 'm' / 's' / 'u')*)> */
		func() bool {
			position129, tokenIndex129, depth129 := position, tokenIndex, depth
			{
				position130 := position
				depth++
				if buffer[position] != rune('/') {
					goto l129
				}
				position++
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('/') {
						goto l133
					}
					position++
					goto l129
				l133:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
				}
				if !matchDot() {
					goto l129
				}
			l131:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					{
						position134, tokenIndex134, depth134 := position, tokenIndex, depth
						if buffer[position] != rune('/') {
							goto l134
						}
						position++
						goto l132
					l134:
						position, tokenIndex, depth = position134, tokenIndex134, depth134
					}
					if !matchDot() {
						goto l132
					}
					goto l131
				l132:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
				}
				if buffer[position] != rune('/') {
					goto l129
				}
				position++
			l135:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					{
						position137, tokenIndex137, depth137 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l138
						}
						position++
						goto l137
					l138:
						position, tokenIndex, depth = position137, tokenIndex137, depth137
						if buffer[position] != rune('l') {
							goto l139
						}
						position++
						goto l137
					l139:
						position, tokenIndex, depth = position137, tokenIndex137, depth137
						if buffer[position] != rune('m') {
							goto l140
						}
						position++
						goto l137
					l140:
						position, tokenIndex, depth = position137, tokenIndex137, depth137
						if buffer[position] != rune('s') {
							goto l141
						}
						position++
						goto l137
					l141:
						position, tokenIndex, depth = position137, tokenIndex137, depth137
						if buffer[position] != rune('u') {
							goto l136
						}
						position++
					}
				l137:
					goto l135
				l136:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
				}
				depth--
				add(ruleRegularExpression, position130)
			}
			return true
		l129:
			position, tokenIndex, depth = position129, tokenIndex129, depth129
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
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if !_rules[ruleArray]() {
						goto l148
					}
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if !_rules[ruleObject]() {
						goto l149
					}
					goto l147
				l149:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if !_rules[ruleRegularExpression]() {
						goto l150
					}
					goto l147
				l150:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					{
						position151 := position
						depth++
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							{
								position154 := position
								depth++
								{
									position155, tokenIndex155, depth155 := position, tokenIndex, depth
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
									position, tokenIndex, depth = position155, tokenIndex155, depth155
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
								depth--
								add(ruleBoolean, position154)
							}
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							{
								position158 := position
								depth++
								if !_rules[ruleInteger]() {
									goto l157
								}
								{
									position159, tokenIndex159, depth159 := position, tokenIndex, depth
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
										position162, tokenIndex162, depth162 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l162
										}
										position++
										goto l161
									l162:
										position, tokenIndex, depth = position162, tokenIndex162, depth162
									}
									goto l160
								l159:
									position, tokenIndex, depth = position159, tokenIndex159, depth159
								}
							l160:
								depth--
								add(ruleFloat, position158)
							}
							goto l152
						l157:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if !_rules[ruleInteger]() {
								goto l163
							}
							goto l152
						l163:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if !_rules[ruleString]() {
								goto l164
							}
							goto l152
						l164:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							{
								position165 := position
								depth++
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
								depth--
								add(ruleNullValue, position165)
							}
						}
					l152:
						depth--
						add(ruleScalarType, position151)
					}
				}
			l147:
				depth--
				add(ruleType, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
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
			position198, tokenIndex198, depth198 := position, tokenIndex, depth
			{
				position199 := position
				depth++
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('$') {
						goto l201
					}
					position++
					{
						position202 := position
						depth++
					l203:
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if !_rules[ruleVariableName]() {
								goto l204
							}
							{
								position205 := position
								depth++
								if buffer[position] != rune('.') {
									goto l204
								}
								position++
								depth--
								add(ruleDOT, position205)
							}
							goto l203
						l204:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
						}
						if !_rules[ruleVariableName]() {
							goto l201
						}
						depth--
						add(ruleVariableNameSequence, position202)
					}
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					{
						position206 := position
						depth++
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
						depth--
						add(ruleSKIPVAR, position206)
					}
				}
			l200:
				depth--
				add(ruleVariable, position199)
			}
			return true
		l198:
			position, tokenIndex, depth = position198, tokenIndex198, depth198
			return false
		},
		/* 79 VariableNameSequence <- <((VariableName DOT)* VariableName)> */
		nil,
		/* 80 VariableName <- <(Identifier ('[' _ VariableIndex _ ']')?)> */
		func() bool {
			position208, tokenIndex208, depth208 := position, tokenIndex, depth
			{
				position209 := position
				depth++
				if !_rules[ruleIdentifier]() {
					goto l208
				}
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('[') {
						goto l210
					}
					position++
					if !_rules[rule_]() {
						goto l210
					}
					{
						position212 := position
						depth++
						if !_rules[ruleExpression]() {
							goto l210
						}
						depth--
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
					position, tokenIndex, depth = position210, tokenIndex210, depth210
				}
			l211:
				depth--
				add(ruleVariableName, position209)
			}
			return true
		l208:
			position, tokenIndex, depth = position208, tokenIndex208, depth208
			return false
		},
		/* 81 VariableIndex <- <Expression> */
		nil,
		/* 82 Block <- <(_ (COMMENT / FlowControlWord / StatementBlock) SEMI? _)> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				if !_rules[rule_]() {
					goto l214
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					{
						position218 := position
						depth++
						if !_rules[rule_]() {
							goto l217
						}
						if buffer[position] != rune('#') {
							goto l217
						}
						position++
					l219:
						{
							position220, tokenIndex220, depth220 := position, tokenIndex, depth
							{
								position221, tokenIndex221, depth221 := position, tokenIndex, depth
								if buffer[position] != rune('\n') {
									goto l221
								}
								position++
								goto l220
							l221:
								position, tokenIndex, depth = position221, tokenIndex221, depth221
							}
							if !matchDot() {
								goto l220
							}
							goto l219
						l220:
							position, tokenIndex, depth = position220, tokenIndex220, depth220
						}
						depth--
						add(ruleCOMMENT, position218)
					}
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					{
						position223 := position
						depth++
						{
							position224, tokenIndex224, depth224 := position, tokenIndex, depth
							{
								position226 := position
								depth++
								{
									position227 := position
									depth++
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
									depth--
									add(ruleBREAK, position227)
								}
								{
									position228, tokenIndex228, depth228 := position, tokenIndex, depth
									if !_rules[rulePositiveInteger]() {
										goto l228
									}
									goto l229
								l228:
									position, tokenIndex, depth = position228, tokenIndex228, depth228
								}
							l229:
								depth--
								add(ruleFlowControlBreak, position226)
							}
							goto l224
						l225:
							position, tokenIndex, depth = position224, tokenIndex224, depth224
							{
								position230 := position
								depth++
								{
									position231 := position
									depth++
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
									depth--
									add(ruleCONT, position231)
								}
								{
									position232, tokenIndex232, depth232 := position, tokenIndex, depth
									if !_rules[rulePositiveInteger]() {
										goto l232
									}
									goto l233
								l232:
									position, tokenIndex, depth = position232, tokenIndex232, depth232
								}
							l233:
								depth--
								add(ruleFlowControlContinue, position230)
							}
						}
					l224:
						depth--
						add(ruleFlowControlWord, position223)
					}
					goto l216
				l222:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					{
						position234 := position
						depth++
						{
							position235, tokenIndex235, depth235 := position, tokenIndex, depth
							{
								position237 := position
								depth++
								if !_rules[ruleSEMI]() {
									goto l236
								}
								depth--
								add(ruleNOOP, position237)
							}
							goto l235
						l236:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							if !_rules[ruleAssignment]() {
								goto l238
							}
							goto l235
						l238:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							{
								position240 := position
								depth++
								{
									position241, tokenIndex241, depth241 := position, tokenIndex, depth
									{
										position243 := position
										depth++
										{
											position244 := position
											depth++
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
											depth--
											add(ruleUNSET, position244)
										}
										if !_rules[ruleVariableSequence]() {
											goto l242
										}
										depth--
										add(ruleDirectiveUnset, position243)
									}
									goto l241
								l242:
									position, tokenIndex, depth = position241, tokenIndex241, depth241
									{
										position246 := position
										depth++
										{
											position247 := position
											depth++
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
											depth--
											add(ruleINCLUDE, position247)
										}
										if !_rules[ruleString]() {
											goto l245
										}
										depth--
										add(ruleDirectiveInclude, position246)
									}
									goto l241
								l245:
									position, tokenIndex, depth = position241, tokenIndex241, depth241
									{
										position248 := position
										depth++
										{
											position249 := position
											depth++
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
											depth--
											add(ruleDECLARE, position249)
										}
										if !_rules[ruleVariableSequence]() {
											goto l239
										}
										depth--
										add(ruleDirectiveDeclare, position248)
									}
								}
							l241:
								depth--
								add(ruleDirective, position240)
							}
							goto l235
						l239:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							{
								position251 := position
								depth++
								if !_rules[ruleIfStanza]() {
									goto l250
								}
							l252:
								{
									position253, tokenIndex253, depth253 := position, tokenIndex, depth
									{
										position254 := position
										depth++
										if !_rules[ruleELSE]() {
											goto l253
										}
										if !_rules[ruleIfStanza]() {
											goto l253
										}
										depth--
										add(ruleElseIfStanza, position254)
									}
									goto l252
								l253:
									position, tokenIndex, depth = position253, tokenIndex253, depth253
								}
								{
									position255, tokenIndex255, depth255 := position, tokenIndex, depth
									{
										position257 := position
										depth++
										if !_rules[ruleELSE]() {
											goto l255
										}
										if !_rules[ruleOPEN]() {
											goto l255
										}
									l258:
										{
											position259, tokenIndex259, depth259 := position, tokenIndex, depth
											if !_rules[ruleBlock]() {
												goto l259
											}
											goto l258
										l259:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
										}
										if !_rules[ruleCLOSE]() {
											goto l255
										}
										depth--
										add(ruleElseStanza, position257)
									}
									goto l256
								l255:
									position, tokenIndex, depth = position255, tokenIndex255, depth255
								}
							l256:
								depth--
								add(ruleConditional, position251)
							}
							goto l235
						l250:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							{
								position261 := position
								depth++
								{
									position262 := position
									depth++
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
									depth--
									add(ruleLOOP, position262)
								}
								{
									position263, tokenIndex263, depth263 := position, tokenIndex, depth
									if !_rules[ruleOPEN]() {
										goto l264
									}
								l265:
									{
										position266, tokenIndex266, depth266 := position, tokenIndex, depth
										if !_rules[ruleBlock]() {
											goto l266
										}
										goto l265
									l266:
										position, tokenIndex, depth = position266, tokenIndex266, depth266
									}
									if !_rules[ruleCLOSE]() {
										goto l264
									}
									goto l263
								l264:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
									{
										position268 := position
										depth++
										{
											position269 := position
											depth++
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
											depth--
											add(ruleCOUNT, position269)
										}
										{
											position270, tokenIndex270, depth270 := position, tokenIndex, depth
											if !_rules[ruleInteger]() {
												goto l271
											}
											goto l270
										l271:
											position, tokenIndex, depth = position270, tokenIndex270, depth270
											if !_rules[ruleVariable]() {
												goto l267
											}
										}
									l270:
										depth--
										add(ruleLoopConditionFixedLength, position268)
									}
									if !_rules[ruleOPEN]() {
										goto l267
									}
								l272:
									{
										position273, tokenIndex273, depth273 := position, tokenIndex, depth
										if !_rules[ruleBlock]() {
											goto l273
										}
										goto l272
									l273:
										position, tokenIndex, depth = position273, tokenIndex273, depth273
									}
									if !_rules[ruleCLOSE]() {
										goto l267
									}
									goto l263
								l267:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
									{
										position275 := position
										depth++
										{
											position276 := position
											depth++
											if !_rules[ruleVariableSequence]() {
												goto l274
											}
											depth--
											add(ruleLoopIterableLHS, position276)
										}
										{
											position277 := position
											depth++
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
											depth--
											add(ruleIN, position277)
										}
										{
											position278 := position
											depth++
											{
												position279, tokenIndex279, depth279 := position, tokenIndex, depth
												if !_rules[ruleCommand]() {
													goto l280
												}
												goto l279
											l280:
												position, tokenIndex, depth = position279, tokenIndex279, depth279
												if !_rules[ruleVariable]() {
													goto l274
												}
											}
										l279:
											depth--
											add(ruleLoopIterableRHS, position278)
										}
										depth--
										add(ruleLoopConditionIterable, position275)
									}
									if !_rules[ruleOPEN]() {
										goto l274
									}
								l281:
									{
										position282, tokenIndex282, depth282 := position, tokenIndex, depth
										if !_rules[ruleBlock]() {
											goto l282
										}
										goto l281
									l282:
										position, tokenIndex, depth = position282, tokenIndex282, depth282
									}
									if !_rules[ruleCLOSE]() {
										goto l274
									}
									goto l263
								l274:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
									{
										position284 := position
										depth++
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
										depth--
										add(ruleLoopConditionBounded, position284)
									}
									if !_rules[ruleOPEN]() {
										goto l283
									}
								l285:
									{
										position286, tokenIndex286, depth286 := position, tokenIndex, depth
										if !_rules[ruleBlock]() {
											goto l286
										}
										goto l285
									l286:
										position, tokenIndex, depth = position286, tokenIndex286, depth286
									}
									if !_rules[ruleCLOSE]() {
										goto l283
									}
									goto l263
								l283:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
									{
										position287 := position
										depth++
										if !_rules[ruleConditionalExpression]() {
											goto l260
										}
										depth--
										add(ruleLoopConditionTruthy, position287)
									}
									if !_rules[ruleOPEN]() {
										goto l260
									}
								l288:
									{
										position289, tokenIndex289, depth289 := position, tokenIndex, depth
										if !_rules[ruleBlock]() {
											goto l289
										}
										goto l288
									l289:
										position, tokenIndex, depth = position289, tokenIndex289, depth289
									}
									if !_rules[ruleCLOSE]() {
										goto l260
									}
								}
							l263:
								depth--
								add(ruleLoop, position261)
							}
							goto l235
						l260:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							if !_rules[ruleCommand]() {
								goto l214
							}
						}
					l235:
						depth--
						add(ruleStatementBlock, position234)
					}
				}
			l216:
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if !_rules[ruleSEMI]() {
						goto l290
					}
					goto l291
				l290:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
				}
			l291:
				if !_rules[rule_]() {
					goto l214
				}
				depth--
				add(ruleBlock, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
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
			position296, tokenIndex296, depth296 := position, tokenIndex, depth
			{
				position297 := position
				depth++
				{
					position298 := position
					depth++
					if !_rules[ruleVariableSequence]() {
						goto l296
					}
					depth--
					add(ruleAssignmentLHS, position298)
				}
				{
					position299 := position
					depth++
					if !_rules[rule_]() {
						goto l296
					}
					{
						position300, tokenIndex300, depth300 := position, tokenIndex, depth
						{
							position302 := position
							depth++
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
							depth--
							add(ruleAssignEq, position302)
						}
						goto l300
					l301:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position304 := position
							depth++
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
							depth--
							add(ruleStarEq, position304)
						}
						goto l300
					l303:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position306 := position
							depth++
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
							depth--
							add(ruleDivEq, position306)
						}
						goto l300
					l305:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position308 := position
							depth++
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
							depth--
							add(rulePlusEq, position308)
						}
						goto l300
					l307:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position310 := position
							depth++
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
							depth--
							add(ruleMinusEq, position310)
						}
						goto l300
					l309:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position312 := position
							depth++
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
							depth--
							add(ruleAndEq, position312)
						}
						goto l300
					l311:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position314 := position
							depth++
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
							depth--
							add(ruleOrEq, position314)
						}
						goto l300
					l313:
						position, tokenIndex, depth = position300, tokenIndex300, depth300
						{
							position315 := position
							depth++
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
							depth--
							add(ruleAppend, position315)
						}
					}
				l300:
					if !_rules[rule_]() {
						goto l296
					}
					depth--
					add(ruleAssignmentOperator, position299)
				}
				{
					position316 := position
					depth++
					if !_rules[ruleExpressionSequence]() {
						goto l296
					}
					depth--
					add(ruleAssignmentRHS, position316)
				}
				depth--
				add(ruleAssignment, position297)
			}
			return true
		l296:
			position, tokenIndex, depth = position296, tokenIndex296, depth296
			return false
		},
		/* 88 AssignmentLHS <- <VariableSequence> */
		nil,
		/* 89 AssignmentRHS <- <ExpressionSequence> */
		nil,
		/* 90 VariableSequence <- <((Variable COMMA)* Variable)> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
			l321:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if !_rules[ruleVariable]() {
						goto l322
					}
					if !_rules[ruleCOMMA]() {
						goto l322
					}
					goto l321
				l322:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
				}
				if !_rules[ruleVariable]() {
					goto l319
				}
				depth--
				add(ruleVariableSequence, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 91 ExpressionSequence <- <((Expression COMMA)* Expression)> */
		func() bool {
			position323, tokenIndex323, depth323 := position, tokenIndex, depth
			{
				position324 := position
				depth++
			l325:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l326
					}
					if !_rules[ruleCOMMA]() {
						goto l326
					}
					goto l325
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
				if !_rules[ruleExpression]() {
					goto l323
				}
				depth--
				add(ruleExpressionSequence, position324)
			}
			return true
		l323:
			position, tokenIndex, depth = position323, tokenIndex323, depth323
			return false
		},
		/* 92 Expression <- <(_ ExpressionLHS ExpressionRHS? _)> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				if !_rules[rule_]() {
					goto l327
				}
				{
					position329 := position
					depth++
					{
						position330 := position
						depth++
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if !_rules[ruleType]() {
								goto l332
							}
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							if !_rules[ruleVariable]() {
								goto l327
							}
						}
					l331:
						depth--
						add(ruleValueYielding, position330)
					}
					depth--
					add(ruleExpressionLHS, position329)
				}
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					{
						position335 := position
						depth++
						{
							position336 := position
							depth++
							if !_rules[rule_]() {
								goto l333
							}
							{
								position337, tokenIndex337, depth337 := position, tokenIndex, depth
								{
									position339 := position
									depth++
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
									depth--
									add(ruleExponentiate, position339)
								}
								goto l337
							l338:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position341 := position
									depth++
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
									depth--
									add(ruleMultiply, position341)
								}
								goto l337
							l340:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position343 := position
									depth++
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
									depth--
									add(ruleDivide, position343)
								}
								goto l337
							l342:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position345 := position
									depth++
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
									depth--
									add(ruleModulus, position345)
								}
								goto l337
							l344:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position347 := position
									depth++
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
									depth--
									add(ruleAdd, position347)
								}
								goto l337
							l346:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position349 := position
									depth++
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
									depth--
									add(ruleSubtract, position349)
								}
								goto l337
							l348:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position351 := position
									depth++
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
									depth--
									add(ruleBitwiseAnd, position351)
								}
								goto l337
							l350:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position353 := position
									depth++
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
									depth--
									add(ruleBitwiseOr, position353)
								}
								goto l337
							l352:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position355 := position
									depth++
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
									depth--
									add(ruleBitwiseNot, position355)
								}
								goto l337
							l354:
								position, tokenIndex, depth = position337, tokenIndex337, depth337
								{
									position356 := position
									depth++
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
									depth--
									add(ruleBitwiseXor, position356)
								}
							}
						l337:
							if !_rules[rule_]() {
								goto l333
							}
							depth--
							add(ruleOperator, position336)
						}
						if !_rules[ruleExpression]() {
							goto l333
						}
						depth--
						add(ruleExpressionRHS, position335)
					}
					goto l334
				l333:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
				}
			l334:
				if !_rules[rule_]() {
					goto l327
				}
				depth--
				add(ruleExpression, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
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
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				if !_rules[rule_]() {
					goto l364
				}
				{
					position366 := position
					depth++
					{
						position367, tokenIndex367, depth367 := position, tokenIndex, depth
						if !_rules[ruleIdentifier]() {
							goto l367
						}
						{
							position369 := position
							depth++
							if buffer[position] != rune(':') {
								goto l367
							}
							position++
							if buffer[position] != rune(':') {
								goto l367
							}
							position++
							depth--
							add(ruleSCOPE, position369)
						}
						goto l368
					l367:
						position, tokenIndex, depth = position367, tokenIndex367, depth367
					}
				l368:
					if !_rules[ruleIdentifier]() {
						goto l364
					}
					depth--
					add(ruleCommandName, position366)
				}
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if !_rules[rule__]() {
						goto l370
					}
					{
						position372, tokenIndex372, depth372 := position, tokenIndex, depth
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
						position, tokenIndex, depth = position372, tokenIndex372, depth372
						if !_rules[ruleCommandFirstArg]() {
							goto l374
						}
						goto l372
					l374:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
						if !_rules[ruleCommandSecondArg]() {
							goto l370
						}
					}
				l372:
					goto l371
				l370:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
				}
			l371:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if !_rules[rule_]() {
						goto l375
					}
					{
						position377 := position
						depth++
						{
							position378 := position
							depth++
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
							depth--
							add(ruleASSIGN, position378)
						}
						if !_rules[ruleVariable]() {
							goto l375
						}
						depth--
						add(ruleCommandResultAssignment, position377)
					}
					goto l376
				l375:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
				}
			l376:
				depth--
				add(ruleCommand, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 101 CommandName <- <((Identifier SCOPE)? Identifier)> */
		nil,
		/* 102 CommandFirstArg <- <(Variable / Type)> */
		func() bool {
			position380, tokenIndex380, depth380 := position, tokenIndex, depth
			{
				position381 := position
				depth++
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if !_rules[ruleVariable]() {
						goto l383
					}
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if !_rules[ruleType]() {
						goto l380
					}
				}
			l382:
				depth--
				add(ruleCommandFirstArg, position381)
			}
			return true
		l380:
			position, tokenIndex, depth = position380, tokenIndex380, depth380
			return false
		},
		/* 103 CommandSecondArg <- <Object> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				if !_rules[ruleObject]() {
					goto l384
				}
				depth--
				add(ruleCommandSecondArg, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 104 CommandResultAssignment <- <(ASSIGN Variable)> */
		nil,
		/* 105 Conditional <- <(IfStanza ElseIfStanza* ElseStanza?)> */
		nil,
		/* 106 IfStanza <- <(IF ConditionalExpression OPEN Block* CLOSE)> */
		func() bool {
			position388, tokenIndex388, depth388 := position, tokenIndex, depth
			{
				position389 := position
				depth++
				{
					position390 := position
					depth++
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
					depth--
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
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !_rules[ruleBlock]() {
						goto l392
					}
					goto l391
				l392:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
				}
				if !_rules[ruleCLOSE]() {
					goto l388
				}
				depth--
				add(ruleIfStanza, position389)
			}
			return true
		l388:
			position, tokenIndex, depth = position388, tokenIndex388, depth388
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
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					{
						position406 := position
						depth++
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
						depth--
						add(ruleNOT, position406)
					}
					goto l405
				l404:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
				}
			l405:
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					{
						position409 := position
						depth++
						if !_rules[ruleAssignment]() {
							goto l408
						}
						if !_rules[ruleSEMI]() {
							goto l408
						}
						if !_rules[ruleConditionalExpression]() {
							goto l408
						}
						depth--
						add(ruleConditionWithAssignment, position409)
					}
					goto l407
				l408:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					{
						position411 := position
						depth++
						if !_rules[ruleCommand]() {
							goto l410
						}
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if !_rules[ruleSEMI]() {
								goto l412
							}
							if !_rules[ruleConditionalExpression]() {
								goto l412
							}
							goto l413
						l412:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
						}
					l413:
						depth--
						add(ruleConditionWithCommand, position411)
					}
					goto l407
				l410:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					{
						position415 := position
						depth++
						if !_rules[ruleExpression]() {
							goto l414
						}
						{
							position416 := position
							depth++
							{
								position417, tokenIndex417, depth417 := position, tokenIndex, depth
								{
									position419 := position
									depth++
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
									depth--
									add(ruleMatch, position419)
								}
								goto l417
							l418:
								position, tokenIndex, depth = position417, tokenIndex417, depth417
								{
									position420 := position
									depth++
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
									depth--
									add(ruleUnmatch, position420)
								}
							}
						l417:
							depth--
							add(ruleMatchOperator, position416)
						}
						if !_rules[ruleRegularExpression]() {
							goto l414
						}
						depth--
						add(ruleConditionWithRegex, position415)
					}
					goto l407
				l414:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					{
						position421 := position
						depth++
						{
							position422 := position
							depth++
							if !_rules[ruleExpression]() {
								goto l402
							}
							depth--
							add(ruleConditionWithComparatorLHS, position422)
						}
						{
							position423, tokenIndex423, depth423 := position, tokenIndex, depth
							{
								position425 := position
								depth++
								{
									position426 := position
									depth++
									if !_rules[rule_]() {
										goto l423
									}
									{
										position427, tokenIndex427, depth427 := position, tokenIndex, depth
										{
											position429 := position
											depth++
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
											depth--
											add(ruleEquality, position429)
										}
										goto l427
									l428:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position431 := position
											depth++
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
											depth--
											add(ruleNonEquality, position431)
										}
										goto l427
									l430:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position433 := position
											depth++
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
											depth--
											add(ruleGreaterEqual, position433)
										}
										goto l427
									l432:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position435 := position
											depth++
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
											depth--
											add(ruleLessEqual, position435)
										}
										goto l427
									l434:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position437 := position
											depth++
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
											depth--
											add(ruleGreaterThan, position437)
										}
										goto l427
									l436:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position439 := position
											depth++
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
											depth--
											add(ruleLessThan, position439)
										}
										goto l427
									l438:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position441 := position
											depth++
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
											depth--
											add(ruleMembership, position441)
										}
										goto l427
									l440:
										position, tokenIndex, depth = position427, tokenIndex427, depth427
										{
											position442 := position
											depth++
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
											depth--
											add(ruleNonMembership, position442)
										}
									}
								l427:
									if !_rules[rule_]() {
										goto l423
									}
									depth--
									add(ruleComparisonOperator, position426)
								}
								if !_rules[ruleExpression]() {
									goto l423
								}
								depth--
								add(ruleConditionWithComparatorRHS, position425)
							}
							goto l424
						l423:
							position, tokenIndex, depth = position423, tokenIndex423, depth423
						}
					l424:
						depth--
						add(ruleConditionWithComparator, position421)
					}
				}
			l407:
				depth--
				add(ruleConditionalExpression, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
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
