package cmd

import "fmt"

const Version = "0.1.0"

func OutputVersion() string {
	return fmt.Sprintf("gs version v%s", Version)
}
