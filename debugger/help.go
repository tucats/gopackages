package debugger

import "github.com/tucats/gopackages/app-cli/tables"

var helpText = [][]string{
	{"break at", "Halt execution at a given line number"},
	{"continue", "Resume execution of the program"},
	{"exit", "Exit the debugger"},
	{"help", "display this help text"},
	{"print", "Print the value of an expression"},
	{"set", "Set a variable to a value"},
	{"show symbols", "Display the current symbol table"},
	{"show line", "Display the current program line"},
	{"step", "Execute the next line of the program"},
}

func Help() error {
	table, err := tables.New([]string{"Command", "Description"})
	for _, helpItem := range helpText {
		err = table.AddRow(helpItem)
	}
	if err == nil {
		err = table.SetOrderBy("Command")
		_ = table.Print("text")
	}
	return err
}
