package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/pflag"
)

var commands []*command

type command struct {
	Name        string
	Short       string
	Description string
	Main        func(args []string, fs *pflag.FlagSet) int
}

func main() {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	cmdMap := map[string]*command{}
	for _, cmd := range commands {
		for _, v := range []string{cmd.Name, cmd.Short} {
			if _, seen := cmdMap[v]; seen {
				panic("command already set: " + v)
			}
			cmdMap[v] = cmd
		}
	}

	if len(os.Args) < 2 {
		globalHelp()
		os.Exit(0)
	}

	if os.Args[1] == "help" {
		globalHelp()
		for _, cmd := range commands {
			fmt.Printf("\n### Help for %s:\n\n", cmd.Name)
			z := os.Args[0] + " " + cmd.Name
			cmd.Main([]string{z, "--help"}, pflag.NewFlagSet(z, pflag.ExitOnError))
		}
	} else if cmd, ok := cmdMap[os.Args[1]]; !ok {
		globalHelp()
		os.Exit(2)
	} else {
		args := append([]string{os.Args[0] + " " + os.Args[1]}, os.Args[2:]...)
		fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
		os.Exit(cmd.Main(args, fs))
	}
}

func globalHelp() {
	fmt.Fprintf(os.Stderr, "Usage: %s command [options] [arguments]\n\nDictutil provides low-level utilities to manipulate Kobo dictionaries.\n\nCommands:\n", os.Args[0])
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "  %-20s %s\n", fmt.Sprintf("%s (%s)", cmd.Name, cmd.Short), cmd.Description)
	}
	fmt.Fprintf(os.Stderr, "  %-20s %s\n", "help", "Show help for all commands")
	fmt.Fprintf(os.Stderr, "\nOptions:\n  -h, --help   Show this help text\n")
}
