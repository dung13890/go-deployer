package handlers

import (
	"flag"
	"fmt"
	"os"
)

var (
	configFile    string
	showHelp      bool
	globalCommand = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	deployCommand command
	pingCommand   command
	copyCommand   command
)

type command struct {
	Name        string
	Description string
	Example     string
	Flag        *flag.FlagSet
	StringFlags map[string]*string
	BoolFlags   map[string]*bool
}

func init() {
	deployCommand = command{
		Name:        "deploy",
		Description: "deploy with flag branch or tag into servers",
		Example:     "[truck deploy -t=1.0.0]",
		Flag:        flag.NewFlagSet("deploy", flag.ExitOnError),
	}
	deployCommand.StringFlags = map[string]*string{
		"tag":    deployCommand.Flag.String("t", "", "Tags to build ex:[truck deploy -t=1.0.0]"),
		"branch": deployCommand.Flag.String("b", "", "Branches to build ex:[truck deploy -b=master]"),
	}
	deployCommand.BoolFlags = map[string]*bool{
		"help": deployCommand.Flag.Bool("h", false, "Display this help message for deploy"),
	}

	pingCommand = command{
		Name:        "ping",
		Description: "Testing connection into servers",
		Example:     "[truck ping]",
		Flag:        flag.NewFlagSet("ping", flag.ExitOnError),
	}
	pingCommand.BoolFlags = map[string]*bool{
		"help": pingCommand.Flag.Bool("h", false, "Display this help message for ping"),
	}

	copyCommand = command{
		Name:        "copy",
		Description: "Copy file into servers",
		Example:     "[truck copy]",
		Flag:        flag.NewFlagSet("copy", flag.ExitOnError),
	}
	copyCommand.BoolFlags = map[string]*bool{
		"help": copyCommand.Flag.Bool("h", false, "Display this help message for ping"),
	}

	globalCommand.BoolVar(&showHelp, "h", false, "Display this help message")
	globalCommand.StringVar(&configFile, "f", "", "Custom path to ./config[.yml]")
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		f.VisitAll(func(o *flag.Flag) {
			fmt.Fprintf(os.Stderr,
				"	%v	%v\n",
				fillColor("-"+o.Name, colorGreen),
				o.Usage,
			)
		})
	}
}

func printHelp() {
	fmt.Printf("%s\n	truck [flags] [<commands>] [args ...]\n\n%s\n",
		fillColor("Usage:", colorYellow),
		fillColor("Flags:", colorYellow),
	)
	setupFlags(globalCommand)
	globalCommand.Usage()
	fmt.Printf("\n%s\n",
		fillColor("Commands:", colorYellow),
	)
	// Ping command
	fmt.Printf("	%s	%s	%s\n",
		fillColor(pingCommand.Name, colorGreen),
		pingCommand.Description,
		fillColor(pingCommand.Example, colorCyan),
	)
	// Coppy Command
	fmt.Printf("	%s	%s	%s\n",
		fillColor(copyCommand.Name, colorGreen),
		copyCommand.Description,
		fillColor(copyCommand.Example, colorCyan),
	)

	// Deploy Command
	fmt.Printf("	%s	%s	%s\n",
		fillColor(deployCommand.Name, colorGreen),
		deployCommand.Description,
		fillColor(deployCommand.Example, colorCyan),
	)
}

func printCommand(c command) {
	setupFlags(c.Flag)
	fmt.Printf("%s\n	%s\n\n%s\n	%s	%s\n\n%s\n",
		fillColor("Usage:", colorYellow),
		c.Example,
		fillColor("Commands:", colorYellow),
		fillColor(c.Name, colorGreen),
		c.Description,
		fillColor("Flags:", colorYellow),
	)
	c.Flag.Usage()
}

func Run() {
	globalCommand.Parse(os.Args[1:])
	if len(os.Args) < 2 || showHelp {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		deployCommand.Flag.Parse(os.Args[2:])
	case "ping":
		pingCommand.Flag.Parse(os.Args[2:])
	case "copy":
		copyCommand.Flag.Parse(os.Args[2:])
	default:
		printHelp()
		os.Exit(1)
	}

	if deployCommand.Flag.Parsed() {
		if *deployCommand.BoolFlags["help"] {
			printCommand(deployCommand)
			os.Exit(1)
		}
		// TODO
		fmt.Printf("tag=%v", *deployCommand.StringFlags["tag"])
		fmt.Printf("branch=%v", *deployCommand.StringFlags["branch"])
	}
	if pingCommand.Flag.Parsed() {
		if *pingCommand.BoolFlags["help"] {
			printCommand(pingCommand)
			os.Exit(1)
		}
		// TODO
		fmt.Printf("%v", pingCommand)
	}
	if copyCommand.Flag.Parsed() {
		if *copyCommand.BoolFlags["help"] {
			printCommand(copyCommand)
			os.Exit(1)
		}
		// TODO
		fmt.Printf("%v", copyCommand)
	}
}
