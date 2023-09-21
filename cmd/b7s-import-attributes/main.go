package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
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
		flagImport     string
		flagSigningKey string
	)

	pflag.StringVarP(&flagImport, "import", "i", "", "file to import attributes from")
	pflag.StringVarP(&flagSigningKey, "key", "k", "", "key to be used for signing (optional")

	pflag.Parse()

	log.SetFlags(0)

	in, err := os.Open(flagImport)
	if err != nil {
		log.Printf("could not open import file: %s", err)
		return failure
	}

	att, err := attributes.Import(in)
	if err != nil {
		log.Printf("could not import attestation: %s", err)
		return failure
	}

	data, _ := json.Marshal(att)

	fmt.Printf("%s\n", string(data))

	return success
}
