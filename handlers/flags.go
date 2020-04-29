package handlers

import (
	"flag"
	"fmt"
	"os"
)

var (
	configFile    string
	isShowHelp    bool
	deployTag     string
	deployBranch  string
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
	pingCommand   = flag.NewFlagSet("ping", flag.ExitOnError)
)

func init() {
	flag.BoolVar(&isShowHelp, "h", false, "Show help")
	flag.StringVar(&configFile, "f", "", "Custom path to ./config[.yml]")
	deployCommand.StringVar(&deployTag, "t", "", "Tags to build ex:[truck deploy -t=1.0.0]")
	deployCommand.StringVar(&deployBranch, "b", "", "Branches to build ex:[truck deploy -b=hotfix/yourticket]")
}

func showHelp() {
	fmt.Printf("%s %s\n%s\n",
		FillColor("Usage:", ColorYellow),
		FillColor("truck [flags] [commands] [args ...]", ColorGreen),
		FillColor("Flags:", ColorYellow),
	)
	flag.PrintDefaults()
	showCommand()
}

func showCommand() {
	fmt.Printf("\n%s\n %s\n",
		FillColor("Commands:", ColorYellow),
		FillColor("deploy	deploy with flag branch or tag", ColorGreen),
	)
	deployCommand.PrintDefaults()
	fmt.Printf("\n %s\n",
		FillColor("ping	Testing connection into servers", ColorGreen),
	)
}

func Run() {
	flag.Parse()
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		deployCommand.Parse(os.Args[2:])

	default:
		showCommand()
		os.Exit(1)
	}

	if deployTag == "" {
		showCommand()
		os.Exit(1)
	}

	fmt.Println(deployTag)
}
