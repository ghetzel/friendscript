# Friendscript

Friendscript is an imperative scripting language that is familiar and easy to use in a variety
of applications that wish to provide scripting support.  It is intended to be a lightweight,
embeddable language that authors can extend in many ways to make adding simple scripting capabilities
to your program easy.

## Uses for Friendscript

- Extending your applications by providing a simple yet robust scripting environment that can be used directly by your users or as an intermediate format for storing complex actions.

- Create a DSL (domain specific language) that you can use to expose your application's functionality in a scripting context.

- Creating a safe alternative to a fully-featured scripting language (e.g.: Lua, Python, Ruby) by having tight control over which language features and functionality is exposed to the end user.

## Language Overview

Read the language introduction and overview [here](docs/README.md)

## Usage Examples

- [Embedding _Friendscript_ in a simple command line application](examples/command-line/main.go)

## "But Grimace, why did you invent a completely new language?  Aren't there so many other powerful scripting langauges out there, and isn't doing so considered a little....unhinged?!"

Short answer "Yes.", long answer "No, with a but..."

Truth be told, my original intent here was largely organized around being curious about language design, the level of difficulty in creating the syntax and logical structure of a language, and all-around just challenging myself to take the risk.

That being said, I think there might actually be something to this.  In any given scenario, there are literally dozens of options to choose from with respect to scripting languages, especially embedded ones.  Lua, Tcl, and GNU Guile all come to mind; and if you're willing to roll your sleeves up a little (or find someone else that already has), Python and Ruby are perfectly valid options.

But I find something interesting about the idea of a soup-to-nuts Golang implementation with a clear, unambiguous API in service of the specific mission of providing familiar scripting constructs to end users, while also giving developers a high degree of flexibility in choosing which language features are available.  I think that last bit is important, because its not always desireable to provide the full breadth of a Turing complete environment; for example in the context of something that _might_ be intended for use in a configuration context.  Configuration is a tricky beast because there are times it requires a middleground between "static, pure data" (e.g.: YAML) and Ruby or Python (powerful, but now your config files can SSH into servers.)  This is just one example where Friendscript might be that middleground.

I think there are some interesting possibilities in this space-- check Friendscript out; perhaps you might too!
