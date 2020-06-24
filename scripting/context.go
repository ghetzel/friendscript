package scripting

import (
	"fmt"
	"time"
)

type ContextType string

const (
	BlockContext     ContextType = `block`
	StatementContext             = `statement`
	CommandContext               = `command`
)

type Context struct {
	Type                ContextType
	Label               string
	Script              *Friendscript
	Filename            string
	Parent              *Context
	AbsoluteStartOffset int
	Length              int
	Error               error
	StartedAt           time.Time
	Took                time.Duration
}

func (self *Context) String() string {
	return fmt.Sprintf("[%v] %v %d + %d", self.Type, self.Label, self.AbsoluteStartOffset, self.Length)
}

func (self *Context) Snippet() string {
	if self.Script != nil {
		var src = self.Script.Buffer

		if self.AbsoluteStartOffset >= 0 && self.Length > 0 {
			if self.AbsoluteStartOffset < len(src) {
				var endIndex = self.AbsoluteStartOffset + self.Length

				if endIndex <= len(src) {
					return src[self.AbsoluteStartOffset:endIndex]
				}
			}
		}
	}

	return ``
}
