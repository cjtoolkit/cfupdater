package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Cli struct {
	sync.Mutex
	topName    []string
	bottomName map[string][]string
	commands   map[string]*Command

	longest int
	args    []string

	clearQueue []func()
}

const (
	base    = 10
	bitSize = 64
	spacer  = "     "
)

func (c *Cli) start(logger *log.Logger, appName string) {
	c.Lock()
	defer c.Unlock()

	logger.SetPrefix(fmt.Sprintf("%s: ", appName))

	if len(c.args) > 1 {
		if c.args[1] == "--help" {
			fmt.Println("hello")
			c.help(logger)
		} else {
			c.exec(logger, appName)
		}
	} else {
		c.info()
	}
}

func (c *Cli) exec(logger *log.Logger, appName string) {
	cmdName := c.args[1]

	if c.commands[cmdName] == nil {
		logger.Fatalf("'%s' is not a reconised command name", cmdName)
	}

	logger.SetPrefix(fmt.Sprintf("%s: %s: ", appName, cmdName))

	command := c.commands[cmdName]

	argumentsStrs := []string{}
	argumentsCount := 0
	for _, argOpt := range c.args[2:] {
		if len(argOpt) >= 2 && argOpt[:2] == "--" {
			opt := argOpt[2:]
			if command.optionLookUp[opt] == nil {
				logger.Fatalf("Does not support this option '--%s'", opt)
			}
			option := command.optionLookUp[opt]
			*(option.ptrBool) = true
		} else {
			argumentsStrs = append(argumentsStrs, argOpt)
			argumentsCount++
		}
	}

	var arguments []*Argument
	arguments = append(command.mandatoryArgumentOrdered, command.optionalArgumentOrdered...)

	for key, argument := range arguments {
		if key >= argumentsCount {
			if argument.mandatory {
				logger.Fatalf("'%s' (%d) is a required argument", argument.name, key+1)
			} else {
				break
			}
		}
		c.execPopulateArgument(logger, argument.ptr, argumentsStrs[key])
	}

	context := command.context

	if context, ok := context.(CommandPreInterface); ok {
		context.CommandPre(logger)
	}

	context.CommandExec(logger)

	if context, ok := context.(CommandPostInterface); ok {
		context.CommandPost(logger)
	}
}

func (c *Cli) execPopulateArgument(logger *log.Logger, ptr interface{}, argumentValue string) {
	switch ptr := ptr.(type) {
	case *string:
		*ptr = argumentValue
	case *int64:
		n, err := strconv.ParseInt(argumentValue, base, bitSize)
		if err != nil {
			logger.Fatal(err.Error())
		}
		*ptr = n
	case *uint64:
		n, err := strconv.ParseUint(argumentValue, base, bitSize)
		if err != nil {
			logger.Fatal(err.Error())
		}
		*ptr = n
	case *float64:
		f, err := strconv.ParseFloat(argumentValue, bitSize)
		if err != nil {
			logger.Fatal(err.Error())
		}
		*ptr = f
	}
}

func (c *Cli) help(logger *log.Logger) {
	if len(c.args) <= 2 {
		c.info()
		return
	}
	cmdName := c.args[2]

	if c.commands[cmdName] == nil {
		logger.Fatalf("'%s' is not a reconised command name", cmdName)
	}

	command := c.commands[cmdName]

	fmt.Println("Command Name        :", cmdName)
	fmt.Println("Command Description :", command.description)
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Print(spacer, cmdName)

	if len(command.optionNames) > 0 {
		fmt.Print(" [options]")
	}

	for _, argumentName := range command.mandatoryArgumentNames {
		fmt.Print(" ", argumentName)
	}

	if len(command.optionalArgumentNames) > 0 {
		fmt.Print("[")
		for _, argumentName := range command.optionalArgumentNames {
			fmt.Print("", argumentName)
		}
		fmt.Print("]")
	}

	fmt.Println()

	c.helpArg(command.mandatoryArgumentNames, command.mandatoryArgumentLookUp, "Mandatory Argument:")
	c.helpArg(command.optionalArgumentNames, command.optionalArgumentLookUp, "Optional Argument:")

	if len(command.optionNames) <= 0 {
		return
	}

	fmt.Println()

	fmt.Println("Options:")

	optionNames := command.optionNames
	sort.Strings(optionNames)

	for _, optionName := range optionNames {
		option := command.optionLookUp[optionName]
		spaces := spacer
		addSpaces := cli.longest - option.nameLen
		for i := 0; i < addSpaces; i++ {
			spaces += " "
		}
		fmt.Print(spacer, "--", optionName, spaces, option.description)
		fmt.Println()
	}
}

func (c *Cli) helpArg(argumentNames []string, lookup map[string]*Argument, title string) {
	if len(argumentNames) > 0 {
		fmt.Println()

		fmt.Println(title)
		sort.Strings(argumentNames)
		for _, argumentName := range argumentNames {
			arg := lookup[argumentName]
			spaces := spacer
			addSpaces := cli.longest - arg.nameLen
			for i := 0; i < addSpaces; i++ {
				spaces += " "
			}
			fmt.Print(spacer, argumentName, spaces, arg.description)
			fmt.Printf(" (type: %T)", arg.ptr)
			fmt.Println()
		}
	}
}

func (c *Cli) info() {
	fmt.Println("List of avaliable command are:")
	fmt.Println()
	fmt.Println("Use '--help commandName' to get more details for each command")
	fmt.Println()
	topNames := c.topName
	sort.Strings(topNames)

	for _, topName := range topNames {
		fmt.Println(topName)

		names := c.bottomName[topName]
		sort.Strings(names)
		for _, name := range names {
			command := c.commands[name]
			spaces := spacer
			addSpaces := cli.longest - command.nameLen
			for i := 0; i < addSpaces; i++ {
				spaces += " "
			}
			fmt.Print(spacer, name, spaces, command.description)
			fmt.Println()
		}
	}
}

func (c *Cli) updateLenght(lenght int) {
	if lenght > c.longest {
		c.longest = lenght
	}
}

func (cli *Cli) register(command CommandInterface) {
	cli.Lock()
	defer cli.Unlock()

	logger := log.New(os.Stderr, "CLI: Command Register: ", log.LstdFlags)

	if fmt.Sprintf("%T", command)[0] != '*' {
		logger.Panicf("'%T' is not a pointer", command)
	}

	cmd := newCmd(command, cli, logger)

	command.CommandSetup(cmd)

	logger.SetPrefix(fmt.Sprintf("CLI: Command Register: %T: ", command))

	if cmd.name == "" && cmd.description == "" {
		logger.Panicf("'%T' has not been initialised properly", command)
	}

	if cmd.name == "" {
		logger.Panicf("'%T' must have a name", command)
	} else if !reAlphaColon.MatchString(cmd.name) {
		logger.Panicf("Name of '%T' must only contain a-z and :", command)
	}

	cli.updateLenght(cmd.nameLen)

	if cmd.description == "" {
		logger.Panicf("'%T' must have a description", command)
	}

	pos := strings.Index(cmd.name, ":")

	if pos == -1 {
		logger.Panicf("Name of '%T' must have at least one colon", command)
	}

	if cli.commands[cmd.name] != nil {
		logger.Panicf("That name '%s' has been taken", cmd.name)
	}

	cli.commands[cmd.name] = cmd

	topName := cmd.name[:pos]

	if cli.bottomName[topName] == nil {
		cli.topName = append(cli.topName, topName)
		cli.bottomName[topName] = []string{cmd.name}
	} else {
		cli.bottomName[topName] = append(cli.bottomName[topName], cmd.name)
	}
}

func (c *Cli) clear() {
	c.topName = nil
	c.bottomName = nil
	c.commands = nil
	c.args = nil

	queue := c.clearQueue
	c.clearQueue = nil

	for _, fn := range queue {
		fn()
	}
}

func (c *Cli) addToClearQueue(fn func()) {
	if fn == nil {
		return
	}

	if c.clearQueue == nil {
		c.clearQueue = []func(){fn}
	} else {
		c.clearQueue = append(c.clearQueue, fn)
	}
}

var (
	cli = &Cli{
		topName:    []string{},
		bottomName: map[string][]string{},
		commands:   map[string]*Command{},
		args:       os.Args,
	}
	reAlphaColon = regexp.MustCompile(`^[a-z-:]+$`)
	reAlpha      = regexp.MustCompile(`^[a-z-]+$`)
)

func Start(logger *log.Logger, appName string) {
	cli.start(logger, appName)
}

func Register(command CommandInterface) {
	cli.register(command)
}

func Clear() {
	cli.clear()
}

func AddToClearQueue(fn func()) {
	cli.addToClearQueue(fn)
}
