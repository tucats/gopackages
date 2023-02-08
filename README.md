# gopackages

This project contains a set of packages intended to support a command-line tool written in Go.
They allow for the definition of the command line grammar (including type checking on option
values, missing arguments, etc) and a defined action routine called when a subcommand is
processed successfully.

A simple command line tool defines a grammar for the commands and subcommands, and their 
options. It then calls the app package Run() method which handles parsing and execution
control from then on.

The command line tool developer:

* Defines the grammar structure
* Provides action routines for each subcommand that results in execution of the command function.
* Action routines can use support packages for these additional common functions:
  * Query and set configuration values
  * Handle basic messaging to the console
  * Enable debug logging as needed
  * Generate consistently formatted tabular output
  * Generate JSON tabular output
* Support for general expression handling is also provided in the github.com/tucats/gopackages/expressions subpackage.
This is also used to support filtering of table output when using the app-cli/tables package.
* A compiler and bytecode interpreter for a Go-like expression evaluator is included
* A library of built-in runtime functions is provided.

See related repositories at github.com/tucats/csv and github.com/tucats/weather for sample
command line tools that demonstrate using this set of packages.
