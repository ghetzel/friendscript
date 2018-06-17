package friendscript

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/ghetzel/friendscript/commands/core"
	cmdfile "github.com/ghetzel/friendscript/commands/file"
	cmdfmt "github.com/ghetzel/friendscript/commands/fmt"
	cmdutils "github.com/ghetzel/friendscript/commands/utils"
	cmdvars "github.com/ghetzel/friendscript/commands/vars"
	"github.com/ghetzel/friendscript/scripting"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var MaxReaderWait = time.Duration(5) * time.Second

const DefaultEnvironmentName = `friendscript`

type InteractiveContext struct {
	Command string
	Line    string
}

type InteractiveHandlerFunc func(ctx *InteractiveContext, environment *Environment) ([]string, error)
type ContextHandlerFunc func(ctx *scripting.Context, isCompleted bool)

type Environment struct {
	Name            string
	modules         map[string]Module
	script          *scripting.Friendscript
	stack           []*scripting.Scope
	replHandlers    map[string]InteractiveHandlerFunc
	contextHandlers []ContextHandlerFunc
	chlock          sync.Mutex
}

// Create a new scripting environment.
func NewEnvironment() *Environment {
	environment := &Environment{
		Name:         DefaultEnvironmentName,
		stack:        make([]*scripting.Scope, 0),
		modules:      make(map[string]Module),
		replHandlers: make(map[string]InteractiveHandlerFunc),
	}

	environment.RegisterModule(scripting.UnqualifiedModuleName, core.New(environment))
	environment.RegisterModule(`fmt`, cmdfmt.New(environment))
	environment.RegisterModule(`file`, cmdfile.New(environment))
	environment.RegisterModule(`utils`, cmdutils.New(environment))
	environment.RegisterModule(`vars`, cmdvars.New(environment))

	return environment
}

// Registers a command module to the given prefix.  If a module with the same prefix already exists,
// it will be replaced with the given module.  Prefixes will be stripped of spaces and converted to
// snake_case.  If an empty prefix is given, the default UnqualifiedModuleName ("core") will be used.
func (self *Environment) RegisterModule(prefix string, module Module) {
	prefix = strings.TrimSpace(prefix)
	prefix = stringutil.Underscore(prefix)

	if prefix == `` {
		prefix = scripting.UnqualifiedModuleName
	}

	self.modules[prefix] = module
}

// Removes a registered module at the given prefix.
func (self *Environment) UnregisterModule(prefix string) {
	delete(self.modules, prefix)
}

// Retrieve a copy of the currently registered modules.
func (self *Environment) Modules() map[string]Module {
	modules := make(map[string]Module)

	for name, module := range self.modules {
		modules[name] = module
	}

	return modules
}

// Registers a function to handle a specific REPL command.  If command is an empty string, the function will be called
// for each command entered into the REPL.
func (self *Environment) RegisterCommandHandler(command string, handler InteractiveHandlerFunc) error {
	if _, ok := self.replHandlers[command]; ok {
		return fmt.Errorf("A handler is already registered for the interactive command '%s'", command)
	} else {
		self.replHandlers[command] = handler
		return nil
	}
}

// Registers a handler that will receive updates on execution context and state as the script is running.
// Will return an integer that can be used to remove the handler at a later point.
func (self *Environment) RegisterContextHandler(handler ContextHandlerFunc) int {
	self.chlock.Lock()
	defer self.chlock.Unlock()

	self.contextHandlers = append(self.contextHandlers, handler)
	return len(self.contextHandlers)
}

// Remove the context handler with the given ID.
func (self *Environment) UnregisterContextHandler(id int) {
	self.chlock.Lock()
	defer self.chlock.Unlock()

	self.contextHandlers = append(self.contextHandlers[:id], self.contextHandlers[id+1:]...)
}

func (self *Environment) EvaluateFile(path string, scope ...*scripting.Scope) (*scripting.Scope, error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()

		return self.EvaluateReader(file, scope...)
	} else {
		return nil, err
	}
}

func (self *Environment) EvaluateReader(reader io.Reader, scope ...*scripting.Scope) (*scripting.Scope, error) {
	var data []byte
	var errchan = make(chan error)

	go func() {
		if d, err := ioutil.ReadAll(reader); err == nil {
			data = d
			errchan <- nil
		} else {
			errchan <- err
		}
	}()

	select {
	case err := <-errchan:
		if err == nil {
			return self.EvaluateString(string(data), scope...)
		} else {
			return nil, err
		}
	case <-time.After(MaxReaderWait):
		return nil, fmt.Errorf("Failed to read Friendscript after %v", MaxReaderWait)
	}
}

func (self *Environment) EvaluateString(data string, scope ...*scripting.Scope) (*scripting.Scope, error) {
	if script, err := scripting.Parse(data); err == nil {
		return self.Evaluate(script, scope...)
	} else {
		return nil, err
	}
}

func (self *Environment) Evaluate(script *scripting.Friendscript, scope ...*scripting.Scope) (*scripting.Scope, error) {
	var rootScope *scripting.Scope

	if len(scope) > 0 && scope[0] != nil {
		rootScope = scope[0]
	} else {
		rootScope = scripting.NewScope(nil)
	}

	self.script = script
	self.pushScope(rootScope)

	for _, block := range script.Blocks() {
		if err := self.evaluateBlock(block); err != nil {
			return self.Scope(), err
		}
	}

	return self.Scope(), nil
}

func (self *Environment) replCompleter(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	// {Text: "users", Description: "Store the username and age"},
	// {Text: "articles", Description: "Store the article text posted by user"},
	// {Text: "comments", Description: "Store the text commented to articles"},

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (self *Environment) REPL() (*scripting.Scope, error) {
	replScope := scripting.NewScope(nil)
	var replErr error

	var options = []prompt.Option{
		prompt.OptionPrefix(self.Name + `> `),
	}

	exec := func(line string) {
		if handled, err := self.evaluateReplBuiltin(line); handled {
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			_, replErr = self.EvaluateString(line, replScope)

			if replErr != nil {
				fmt.Println(replErr.Error())
			}
		}
	}

	repl := prompt.New(exec, self.replCompleter, options...)
	repl.Run()

	return replScope, replErr
}

func (self *Environment) evaluateReplBuiltin(line string) (bool, error) {
	var handler InteractiveHandlerFunc

	cmd, _ := stringutil.SplitPair(strings.TrimSpace(line), ` `)
	ctx := &InteractiveContext{
		Command: cmd,
		Line:    line,
	}

	if h, ok := self.replHandlers[cmd]; ok {
		handler = h
	} else if h, ok := self.replHandlers[``]; ok {
		handler = h
	}

	if handler != nil {
		if output, err := handler(ctx, self); err == nil {
			for _, line := range output {
				fmt.Println(line)
			}
		} else {
			fmt.Println(color.New(color.FgRed).Sprint(`error:`) + ` ` + err.Error())
		}

		return true, nil
	} else {
		return false, nil
	}
}

func (self *Environment) pushScope(scope *scripting.Scope) {
	// if len(self.stack) > 0 {
	// 	log.Debugf("PUSH scope(%d) is masked", self.Scope().Level())
	// } else {
	// 	log.Debugf("PUSH scope(%d) is ROOT", scope.Level())
	// }

	self.stack = append(self.stack, scope)

	if self.script != nil {
		self.script.SetScope(self.Scope())
	}

	// log.Debugf("PUSH scope(%d) is active", self.Scope().Level())
}

func (self *Environment) Scope() *scripting.Scope {
	if len(self.stack) > 0 {
		return self.stack[len(self.stack)-1]
	} else {
		return nil
	}
}

func (self *Environment) popScope() *scripting.Scope {
	if len(self.stack) > 1 {
		top := self.stack[len(self.stack)-1]
		self.stack = self.stack[0 : len(self.stack)-1]

		if self.script != nil {
			self.script.SetScope(self.Scope())
		}

		// log.Debugf("POP  scope(%d) is active", self.Scope().Level())

		return top
	} else if len(self.stack) == 1 {
		return self.stack[0]
	} else {
		log.Fatal("attempted pop on an empty scope stack")
		return nil
	}
}

func (self *Environment) evaluateBlock(block *scripting.Block) error {
	log.Debug(strings.Repeat("-", 70))

	switch block.Type() {
	case scripting.StatementBlock:
		for _, statement := range block.Statements() {
			if err := self.evaluateStatement(statement); err != nil {
				return err
			}
		}

	case scripting.EventHandlerBlock:
		return fmt.Errorf("Not Implemented")

	case scripting.FlowControlWord:
		if levels := block.FlowBreak(); levels > 0 {
			return scripting.NewFlowControl(scripting.FlowBreak, levels)
		} else if levels := block.FlowContinue(); levels > 0 {
			return scripting.NewFlowControl(scripting.FlowContinue, levels)
		} else {
			return fmt.Errorf("invalid flow control statement")
		}
	}

	return nil
}

func (self *Environment) evaluateStatement(statement *scripting.Statement) error {
	switch statement.Type() {
	case scripting.AssignmentStatement:
		return self.evaluateAssignment(statement.Assignment(), false)

	case scripting.DirectiveStatement:
		return self.evaluateDirective(statement.Directive())

	case scripting.ConditionalStatement:
		_, err := self.evaluateConditional(statement.Conditional())
		return err

	case scripting.LoopStatement:
		return self.evaluateLoop(statement.Loop())

	case scripting.CommandStatement:
		_, err := self.evaluateCommand(statement.Command(), false)
		return err

	case scripting.NoOpStatement:
		return nil

	default:
		return fmt.Errorf("Unrecognized statement: %v", statement.Type())
	}
}

func (self *Environment) evaluateAssignment(assignment *scripting.Assignment, forceDeclare bool) error {
	// log.Debugf("ASSN %v", assignment)

	if assignment.Operator.ShouldPreclear() {
		// clear out all the left-hand side variables
		for _, lhs := range assignment.LeftHandSide {
			if forceDeclare {
				self.Scope().Declare(lhs)
			} else {
				self.Scope().Set(lhs, nil)
			}
		}
	}

	// unpack
	if len(assignment.RightHandSide) == 1 {
		if rhs, err := assignment.RightHandSide[0].Value(); err == nil {
			totalLhsCount := len(assignment.LeftHandSide)

			if totalLhsCount > 1 && typeutil.IsArray(rhs) {
				for i, rhs := range sliceutil.Sliceify(rhs) {
					if i < totalLhsCount {
						if result, err := assignment.Operator.Evaluate(
							self.Scope().Get(assignment.LeftHandSide[i]),
							rhs,
						); err == nil {
							self.Scope().Set(assignment.LeftHandSide[i], result)
						} else {
							return err
						}
					}
				}

				return nil
			}
		}
	}

	for i, lhs := range assignment.LeftHandSide {
		if i < len(assignment.RightHandSide) {
			if result, err := assignment.Operator.Evaluate(
				self.Scope().Get(lhs),
				assignment.RightHandSide[i],
			); err == nil {
				self.Scope().Set(lhs, result)
			} else {
				return err
			}
		}
	}

	return nil
}

func (self *Environment) evaluateDirective(directive *scripting.Directive) error {
	switch directive.Type() {
	case scripting.UnsetDirective:
		return fmt.Errorf("'unset' not implemented yet")
	case scripting.IncludeDirective:
		return fmt.Errorf("'include' not implemented yet")
	case scripting.DeclareDirective:
		for _, varname := range directive.VariableNames() {
			self.Scope().Declare(varname)
		}
	}

	return nil
}

func (self *Environment) evaluateCommand(command *scripting.Command, forceDeclare bool) (string, error) {
	modname, name := command.Name()
	log.Debugf("EXEC %v::%v", modname, name)
	ctx := command.SourceContext()
	self.sendContextUpdate(ctx, false)

	if first, rest, err := command.Args(); err == nil {
		// locate the module this command belongs to
		if module, ok := self.modules[modname]; ok {
			// log.Debugf("CMND called %T(%v), %T(%v)", first, first, rest, rest)

			// tell that module to execute the command, giving it the name and arguments
			if result, err := module.ExecuteCommand(name, first, rest); err == nil {
				self.sendContextUpdate(ctx, true)
				// log.Debugf("CMND returned %T(%v)", result, result)

				// if there is an output variable destination, set that in the current scope
				if resultVar := command.OutputName(); resultVar != `` {
					if forceDeclare {
						self.Scope().Declare(resultVar)
					}

					self.Scope().Set(resultVar, result)

					return resultVar, nil
				}

				return ``, nil
			} else {
				ctx.Error = err
			}
		} else {
			ctx.Error = fmt.Errorf("Cannot locate module %q", modname)
		}
	} else {
		ctx.Error = fmt.Errorf("invalid arguments: %v", err)
	}

	self.sendContextUpdate(ctx, true)
	return ``, ctx.Error
}

func (self *Environment) evaluateConditional(conditional *scripting.Conditional) (bool, error) {
	var blocks = make([]*scripting.Block, 0)
	var trueBranch bool
	var conditionScope = scripting.NewScope(self.Scope())
	self.pushScope(conditionScope)
	defer self.popScope()

	switch conditional.Type() {
	case scripting.ConditionWithAssignment:
		assignment, condition := conditional.WithAssignment()

		if err := self.evaluateAssignment(assignment, true); err == nil {
			result := condition.IsTrue()
			blocks, trueBranch = self.evaluateConditionalGetBranch(conditional, result)
		} else {
			return trueBranch, err
		}

	case scripting.ConditionWithCommand:
		command, condition := conditional.WithCommand()

		if _, err := self.evaluateCommand(command, true); err == nil {
			result := condition.IsTrue()
			blocks, trueBranch = self.evaluateConditionalGetBranch(conditional, result)
		} else {
			return trueBranch, err
		}

	case scripting.ConditionWithRegex:
		expression, matchOp, rx := conditional.WithRegex()
		result := matchOp.Evaluate(rx, expression)
		blocks, trueBranch = self.evaluateConditionalGetBranch(conditional, result)

	case scripting.ConditionWithComparator:
		lhs, cmp, rhs := conditional.WithComparator()

		result := cmp.Evaluate(lhs, rhs)
		blocks, trueBranch = self.evaluateConditionalGetBranch(conditional, result)

	default:
		return trueBranch, fmt.Errorf("Unrecognized Conditional type")
	}

	for _, block := range blocks {
		if err := self.evaluateBlock(block); err != nil {
			return trueBranch, err
		}
	}

	return trueBranch, nil
}

func (self *Environment) evaluateConditionalGetBranch(conditional *scripting.Conditional, result bool) ([]*scripting.Block, bool) {
	var blocks = make([]*scripting.Block, 0)
	var trueBranch bool

	if conditional.IsNegated() {
		result = !result
	}

	if result {
		// log.Debugf("IF branch")
		blocks = conditional.IfBlocks()
		trueBranch = true
	} else {
		var tookElifBranch bool

		for _, elif := range conditional.ElseIfConditions() {
			if t, err := self.evaluateConditional(elif); err == nil {
				if t {
					// log.Debugf("ELSE-IF %d branch", ei)
					tookElifBranch = true
					blocks = elif.IfBlocks()
					break
				}
			} else {
				log.Fatal(err)
				return nil, false
			}
		}

		if !tookElifBranch {
			// log.Debugf("ELSE branch")
			blocks = conditional.ElseBlocks()
		}
	}

	return blocks, trueBranch
}

func (self *Environment) evaluateLoop(loop *scripting.Loop) error {
	var i int
	var sourceVar string
	var destVars []string
	var loopScope = scripting.NewScope(self.Scope())

	loopScope.Declare(`index`)

	self.pushScope(loopScope)
	defer self.popScope()

	// if we have an iterator, we have to initialize the values
	if loop.Type() == scripting.IteratorLoop {
		if s, d, err := self.evaluateLoopIterationStart(loop, loopScope); err == nil {
			sourceVar = s
			destVars = d

			log.Debugf("Iterator initialized: %v -> %v", sourceVar, destVars)
		} else {
			return err
		}
	}

	// log.Debugf("LOOP BEGIN")

LoopEval:
	for loop.ShouldContinue() {
		if loop.Type() == scripting.IteratorLoop {
			iterVector := loopScope.Get(sourceVar)

			if typeutil.IsMap(iterVector) {
				remap := make([][]interface{}, 0)
				keys := maputil.StringKeys(iterVector)
				sort.Strings(keys)

				for _, key := range keys {
					remap = append(remap, []interface{}{
						key,
						maputil.Get(iterVector, key),
					})
				}

				iterVector = remap
			}

			if iterLen := sliceutil.Len(iterVector); i < iterLen {
				if iterItem, ok := sliceutil.At(iterVector, i); ok {
					var didSet bool

					if totalLhsCount := len(destVars); totalLhsCount > 1 {
						if typeutil.IsArray(iterItem) {
							for j, rhs := range sliceutil.Sliceify(iterItem) {
								if j < totalLhsCount {
									loopScope.Set(destVars[j], rhs)
									didSet = true
								}
							}
						}
					}

					if !didSet {
						loopScope.Set(destVars[0], iterItem)
					}
				} else {
					return fmt.Errorf("Failed to retrieve iterator item %d", i)
				}
			} else {
				break
			}
		}

		loopScope.Set(`index`, loop.CurrentIndex())

		for _, block := range loop.Blocks() {
			if err := self.evaluateBlock(block); err != nil {
				if fc, ok := err.(*scripting.FlowControlErr); ok {
					if fc.Level <= 0 {
						return fc
					} else if fc.Level == 1 {
						if fc.Type == scripting.FlowContinue {
							continue LoopEval
						} else {
							break LoopEval
						}
					} else {
						fc.Level = fc.Level - 1
						return fc
					}
				} else {
					return err
				}
			}
		}

		i += 1
	}

	return nil
}

func (self *Environment) evaluateLoopIterationStart(loop *scripting.Loop, scope *scripting.Scope) (string, []string, error) {
	destVars, source := loop.IteratableParts()
	var sourceVar string

	if cmd, ok := source.(*scripting.Command); ok {
		// since we totally need the results of the command to iterate on them, if the command
		// didn't specify a result variable, we're going to force it to have one
		if outputName := cmd.OutputName(); outputName == `` {
			cmd.SetOutputNameOverride(scripting.DefaultIteratorCommandResultVariableName)
		}

		if resultVar, err := self.evaluateCommand(cmd, true); err == nil {
			sourceVar = resultVar
		} else {
			return ``, nil, err
		}
	} else if srcvar, ok := source.(string); ok {
		sourceVar = srcvar
	}

	for _, v := range destVars {
		scope.Declare(v)
	}

	return sourceVar, destVars, nil
}

func (self *Environment) sendContextUpdate(ctx *scripting.Context, isDone bool) {
	self.chlock.Lock()
	defer self.chlock.Unlock()

	for _, handler := range self.contextHandlers {
		handler(ctx, isDone)
	}
}
