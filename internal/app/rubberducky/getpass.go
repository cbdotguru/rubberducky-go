package rubberducky

import (
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// GetPass TODO some comments
func GetPass() string {
	fmt.Println("We all know by now you have to unlock your magical key with a password, please enter your password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()
	return password
}
