package gitconfig_test

import (
	"fmt"
	"testing"

	"github.com/denormal/go-gitconfig"
)

func ExampleNew(t *testing.T) {
	// attempt to load git configuration from within the current directory
	config, err := gitconfig.New()
	if err != nil {
		fmt.Printf(
			"error encountered loading git configuration: %s",
			err.Error(),
		)
	} else {
		email := config.Get("user.email")
		if email != nil {
			fmt.Printf("the git user's email is %q\n", email)
		}

		fmt.Println("the core.* git configuration properties are:")
		for _, property := range config.Find("core.*") {
			fmt.Printf("\t%s=%s\n", property.Name(), property)
		}
	}
} // ExampleNew()
