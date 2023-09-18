package main

import (
	"log"
	"os"
	"reflect"

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

	log.SetFlags(0)

	attrs, err := env.ReadAttributes(flagPrefix, flagIgnore)
	if err != nil {
		log.Printf("could not retrieve attributes from the environment: %s", err)
		return failure

	}

	if len(attrs) == 0 {
		return success
	}

	if uint(len(attrs)) > flagLimit {

		if flagStrict {
			log.Printf("too many attributes set. remove some or re-run with 'strict' mode turned off")
			return failure
		}

		attrs = attrs[:flagLimit]
	}

	log.Printf("%v attributes retrieved", len(attrs))

	encoded, err := attributes.Encode(attrs)
	if err != nil {
		log.Printf("could not encode attributes: %s", err)
		return failure
	}

	decoded, err := attributes.Decode(encoded)
	if err != nil {
		log.Printf("could not decode attributes: %s", err)
		return failure
	}

	same := reflect.DeepEqual(attrs, decoded)
	if !same {
		log.Printf("attributes not equal!")
		log.Printf("orig: %+#v", attrs)
		log.Printf("decoded: %+#v", decoded)
		return failure
	}

	log.Printf("original and decoded data are the same")

	return success
}
