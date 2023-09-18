package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/blocklessnetwork/b7s-attributes/env"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagPrefix string
		flagLimit  uint
		flagIgnore []string
		flagStrict bool
	)

	pflag.StringVar(&flagPrefix, "prefix", attributes.Prefix, "prefix node attributes environment variables have")
	pflag.UintVarP(&flagLimit, "limit", "l", attributes.Limit, "number of node attributes to use")
	pflag.StringSliceVarP(&flagIgnore, "ignore", "i", []string{}, "environment variables to skip")
	pflag.BoolVar(&flagStrict, "strict", true, "stop execution if there are too many attributes")

	pflag.Parse()

	attributes, err := env.ReadAttributes(flagPrefix, flagLimit, flagIgnore)
	if err != nil {
		if !errors.Is(err, env.ErrTooManyAttributes) {
			log.Printf("could not retrieve attributes from the environment: %s", err)
			return failure
		}

		if flagStrict {
			log.Printf("too many arguments set. remove some or rerun with strict mode turned off")
			return failure
		}
	}

	if len(attributes) == 0 {
		return success
	}

	log.Printf("%v attributes retrieved", len(attributes))

	for name, value := range attributes {
		fmt.Printf("%v: %v\n", name, value)
	}

	return success
}
