package scripting

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var rxInterpolate = regexp.MustCompile(`({[^}]+})`)
var placeholderVarName = `_`
var maxInterpolateSequences = 64

// This represents the name of a module whose commands do not need to be qualified with
// a "name::" prefix.
var UnqualifiedModuleName = `core`

type tracer int

type Scope struct {
	Environment    Commandable
	SkipPreclear   bool
	parent         *Scope
	data           map[string]any
	isolatedReads  bool
	isolatedWrites bool
	mostRecentKey  string
	evallock       sync.Mutex
	ctx            *Context
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent: parent,
		data:   make(map[string]any),
	}
}

func NewEphemeralScope(parent *Scope) *Scope {
	scope := NewScope(nil)
	scope.isolatedReads = false
	scope.isolatedWrites = true
	return scope
}

func NewIsolatedScope(parent *Scope) *Scope {
	scope := NewScope(parent)
	scope.isolatedReads = true
	scope.isolatedWrites = true
	return scope
}

func (self *Scope) Env() Commandable {
	if self.Environment != nil {
		return self.Environment
	} else if self.parent != nil {
		return self.parent.Env()
	} else {
		return new(NullCommand)
	}
}

func (self *Scope) Level() int {
	if self.parent == nil {
		return 0
	} else {
		return self.parent.Level() + 1
	}
}

func (self *Scope) MostRecentValue() any {
	if self.mostRecentKey == `` {
		return nil
	}

	return self.Get(self.mostRecentKey)
}

// Sets the scope evaluation lock and store the given context.
func (self *Scope) LockContext(ctx *Context) {
	self.evallock.Lock()
	self.ctx = ctx
}

// Clears the evaluation context and clears the lock.
func (self *Scope) Unlock() {
	self.ctx = nil
	self.evallock.Unlock()
}

// Returnt the current evaluation context (if any, may be nil).
func (self *Scope) EvalContext() *Context {
	if self.ctx != nil {
		return self.ctx
	} else if self.parent != nil {
		return self.parent.EvalContext()
	} else {
		return nil
	}
}

func (self *Scope) String() string {
	if data, err := json.MarshalIndent(self.Data(), ``, `  `); err == nil {
		return string(data)
	} else {
		return err.Error()
	}
}

func (self *Scope) Data() map[string]any {
	var output = make(map[string]any)

	maputil.Walk(self.data, func(value any, path []string, isLeaf bool) error {
		if resolvable, ok := value.(Resolvable); ok {
			maputil.DeepSet(output, path, resolvable.Resolve())
		} else if typeutil.IsArray(value) {
			maputil.DeepSet(output, path, value)
			return maputil.SkipDescendants
		} else if isLeaf {
			maputil.DeepSet(output, path, value)
		}

		return nil
	})

	return output
}

func (self *Scope) Declare(key string) {
	if key == `` || key == placeholderVarName {
		return
	}

	var e emptyValue
	key = self.prepVariableName(key)

	// log.Infof("DECL scope(%d)[%v]", self.Level(), key)
	maputil.DeepSet(self.data, strings.Split(key, `.`), e)
}

func (self *Scope) Set(key string, value any) {
	key = self.prepVariableName(key)
	scope := self.OwnerOf(key)
	scope.set(key, value)
	self.mostRecentKey = key
}

func (self *Scope) Get(key string, fallback ...any) any {
	value, _ := self.get(key, fallback...)

	// the emptyValue type is used by the "declare" statement to put a non-nil placeholder
	// value in a scope for the purpose of occupying they key.  When used as a value outside
	// of this package, it should be nil.
	if IsEmpty(value) {
		return nil
	}

	return value
}

// Returns the scope that "owns" the given key.  This works by first checking for an
// already-set key in the current scope.  If none exists, the parent scope
// is consulted for non-nil values (and so on, all the way up the scope chain).
//
// If none of the ancestor scopes have a non-nil value at the given key, the current
// scope becomes the owner of the key and will be returned.
func (self *Scope) OwnerOf(key string) *Scope {
	if self.isolatedWrites || self.IsLocal(key) {
		return self
	} else {
		_, scope := self.get(key)
		return scope
	}
}

func (self *Scope) IsLocal(key string) bool {
	if _, ok := maputil.DeepGet(self.data, strings.Split(key, `.`), tracer(0)).(tracer); ok {
		return false
	}

	return true
}

func (self *Scope) set(key string, value any) {
	if key == `` || key == placeholderVarName {
		return
	}

	if IsEmpty(value) {
		value = new(emptyValue)
	} else if v, err := exprToValue(value); err == nil {
		value = v
	} else {
		log.Panicf("Cannot set %v: %v", key, err)
	}

	value = intIfYouCan(value)
	value = mapifyStruct(value)

	// fmt.Printf("SSET scope(%d)[%v] = %T(%v)\n", self.Level(), key, value, value)
	//
	// for _, st := range log.StackTrace(3) {
	// 	if strings.Contains(st.String(), `friendscript`) {
	// 		fmt.Printf("  " + st.String() + "\n")
	// 	}
	// }

	maputil.DeepSet(self.data, strings.Split(key, `.`), value)
}

func (self *Scope) get(key string, fallback ...any) (any, *Scope) {
	key = self.prepVariableName(key)

	v := maputil.DeepGet(self.data, strings.Split(key, `.`))

	if !IsEmpty(v) {
		// return *copies* of compound types
		if typeutil.IsMap(v) {
			v = maputil.DeepCopyStruct(v)
		} else if typeutil.IsArray(v) {
			v = sliceutil.Sliceify(v)
		}

		// fmt.Printf("SGET scope(%d)[%v] -> %T(%v)\n", self.Level(), key, v, v)
		return v, self
	} else if self.parent != nil && !self.isolatedReads {
		if v, scope := self.parent.get(key, fallback...); v != nil {
			// fmt.Printf("SGET scope(%d)[%v] -> PARENT(%T)\n", self.Level(), key, v)
			return v, scope
		}
	}

	if len(fallback) > 0 && fallback[0] != nil {
		// fmt.Printf("SGET scope(%d)[%v] -> %T(%v) FALLBACK\n", self.Level(), key, fallback[0], fallback[0])
		return fallback[0], self
	} else {
		// fmt.Printf("SGET scope(%d)[%v] -> nil FALLBACK\n", self.Level(), key)
		return new(emptyValue), self
	}
}

func (self *Scope) Interpolate(in string) string {
	for range maxInterpolateSequences {
		if match := rxutil.Match(rxInterpolate, in); match != nil {
			seq := match.Group(1)
			seq = stringutil.Unwrap(seq, `{`, `}`)

			value := self.Get(seq)

			if IsEmpty(value) {
				value = ``
			}

			in = match.ReplaceGroup(1, typeutil.String(value))
		} else {
			break
		}
	}

	return in
}

func (self *Scope) prepVariableName(key string) string {
	key = strings.TrimPrefix(key, `$`)

	return key
}

func IsEmpty(in any) bool {
	if in == nil {
		return true
	} else if _, ok := in.(*emptyValue); ok {
		return true
	} else if _, ok := in.(emptyValue); ok {
		return true
	}

	return false
}

type Commandable interface {
	ExecuteCommand(command *Command) (string, any, error)
}

type NullCommand struct {
	Commandable
}

func (n *NullCommand) ExecuteCommand(command *Command) (string, any, error) {
	return ``, nil, fmt.Errorf("no command interface found")
}
