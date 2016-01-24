package cli

import (
	"fmt"
	"log"
	"unicode/utf8"
)

type Command struct {
	name        string
	nameLen     int
	description string

	context CommandInterface
	logger  *log.Logger
	cli     *Cli

	mandatoryArgumentNames   []string
	mandatoryArgumentOrdered []*Argument
	mandatoryArgumentLookUp  map[string]*Argument

	optionalArgumentNames   []string
	optionalArgumentOrdered []*Argument
	optionalArgumentLookUp  map[string]*Argument

	optionNames  []string
	optionLookUp map[string]*Option
}

func newCmd(context CommandInterface, cli *Cli, logger *log.Logger) *Command {
	cmd := &Command{}

	cmd.context = context
	cmd.logger = logger
	cmd.cli = cli

	cmd.mandatoryArgumentNames = []string{}
	cmd.mandatoryArgumentOrdered = []*Argument{}
	cmd.mandatoryArgumentLookUp = map[string]*Argument{}

	cmd.optionalArgumentNames = []string{}
	cmd.optionalArgumentOrdered = []*Argument{}
	cmd.optionalArgumentLookUp = map[string]*Argument{}

	cmd.optionNames = []string{}
	cmd.optionLookUp = map[string]*Option{}

	return cmd
}

func (c *Command) Init(name, description string) {
	c.name = name
	c.nameLen = utf8.RuneCountInString(name)
	c.description = description
}

func (c *Command) argumentCheck(name, description string, ptr interface{}, lookup map[string]*Argument) {
	switch {
	case ptr == nil:
		c.logger.Panic("'ptr' cannot be nil")
	case !reAlpha.MatchString(name):
		c.logger.Panicf("'%s' must be a-z", name)
	case description == "":
		c.logger.Panic("description cannot be empty")
	case lookup[name] != nil:
		c.logger.Panicf("'%s' has already been taken", name)
	}

	switch ptr.(type) {
	case *string: // pass
	case *int64: // pass
	case *uint64: // pass
	case *float64: // pass
	default:
		c.logger.Panic("ptr is not a valid data type, valid data type are *string, *int64, *uint64 and *float64")
	}
}

// ptr must be *string, *int64, *uint64 and *float64
func (c *Command) AddMandatoryArgument(name, description string, ptr interface{}) {
	c.logger.SetPrefix(fmt.Sprintf("CLI: Command Register: %T: Add Mandatory Argument: %s:", c.context, name))

	c.argumentCheck(name, description, ptr, c.mandatoryArgumentLookUp)

	arg := newArgument(true, name, description, ptr)
	c.cli.updateLenght(arg.nameLen)

	c.mandatoryArgumentNames = append(c.mandatoryArgumentNames, name)
	c.mandatoryArgumentOrdered = append(c.mandatoryArgumentOrdered, arg)
	c.mandatoryArgumentLookUp[name] = arg
}

// ptr must be *string, *int64, *uint64 and *float64
func (c *Command) AddOptionalArgument(name, description string, ptr interface{}) {
	c.logger.SetPrefix(fmt.Sprintf("CLI: Command Register: %T: Add Optional Argument: %s:", c.context, name))

	c.argumentCheck(name, description, ptr, c.optionalArgumentLookUp)

	arg := newArgument(false, name, description, ptr)
	c.cli.updateLenght(arg.nameLen)

	c.optionalArgumentNames = append(c.optionalArgumentNames, name)
	c.optionalArgumentOrdered = append(c.optionalArgumentOrdered, arg)
	c.optionalArgumentLookUp[name] = arg
}

func (c *Command) AddOption(name, description string, ptrBool *bool) {
	c.logger.SetPrefix(fmt.Sprintf("CLI: Command Register: %T: Add Option Argument: %s:", c.context, name))

	switch {
	case ptrBool == nil:
		c.logger.Panic("'ptrBool' cannot be nil")
	case !reAlpha.MatchString(name):
		c.logger.Panicf("'%s' must be a-z", name)
	case description == "":
		c.logger.Panic("description cannot be empty")
	case c.optionLookUp[name] != nil:
		c.logger.Panicf("'%s' has already been taken", name)
	}

	opt := newOption(name, description, ptrBool)
	c.cli.updateLenght(opt.nameLen)

	c.optionNames = append(c.optionNames, name)
	c.optionLookUp[name] = opt
}

type CommandInterface interface {
	CommandSetup(*Command)
	CommandExec(*log.Logger)
}

type CommandPreInterface interface {
	CommandInterface
	CommandPre(*log.Logger)
}

type CommandPostInterface interface {
	CommandInterface
	CommandPost(*log.Logger)
}
