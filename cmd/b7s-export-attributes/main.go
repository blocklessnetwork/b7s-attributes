package main

import (
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
		flagPrefix     string
		flagLimit      uint
		flagIgnore     []string
		flagStrict     bool
		flagExport     string
		flagSigningKey string
	)

	pflag.StringVar(&flagPrefix, "prefix", attributes.Prefix, "prefix node attributes environment variables have")
	pflag.UintVarP(&flagLimit, "limit", "l", attributes.Limit, "number of node attributes to use")
	pflag.StringSliceVarP(&flagIgnore, "ignore", "i", []string{}, "environment variables to skip")
	pflag.BoolVar(&flagStrict, "strict", true, "stop execution if there are too many attributes")
	pflag.StringVarP(&flagExport, "export", "e", "", "file to export attributes to")
	pflag.StringVarP(&flagSigningKey, "key", "k", "", "key to be used for signing (optional")

	pflag.Parse()

	log.SetFlags(0)

	attrs, err := env.ReadAttributes(flagPrefix, flagIgnore)
	if err != nil {
		log.Printf("could not retrieve attributes from the environment: %s", err)
		return failure

	}

	if len(attrs) == 0 {
		log.Printf("no attributes found")
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

	att := attributes.Attestation{
		Attributes: attrs,
	}

	export, err := os.Create(flagExport)
	if err != nil {
		log.Printf("could not open export file: %s", err)
		return failure
	}
	defer export.Close()

	err = attributes.Export(export, att)
	if err != nil {
		log.Printf("could not export attestation: %s", err)
		return failure
	}

	log.Printf("attributes successfully exported")

	return success
}
