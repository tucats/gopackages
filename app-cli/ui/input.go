package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Prompt prints a prompt string, and gets input from the console.
// The line endings are removed and the remainder of the input is
// returned as a string.
func Prompt(p string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(p)
	buffer, _ := reader.ReadString('\n')

	//Remove any extra line endings (CRLF or LF)
	buffer = strings.Replace(buffer, "\r\n", "", -1)
	buffer = strings.Replace(buffer, "\n", "", -1)

	return buffer
}

// PromptPassword prompts the user with a string prompt, and then
// allows the user to enter confidential information such as a password
// without it being echoed on the terminal. The value entered is returned
// as a string.
func PromptPassword(p string) string {
	fmt.Print(p)
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)
	fmt.Println() // it's necessary to add a new line after user's input

	return password
}
