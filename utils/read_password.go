package utils

import (
	"fmt"
	"golang.org/x/term"
	"strings"
	"syscall"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytePassword)), nil
}
