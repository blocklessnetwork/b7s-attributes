package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/web3-storage/go-w3s-client"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

const (
	apiToken = "WEB3STORAGE_TOKEN"
)

func runUpload(_ *cobra.Command, args []string) error {

	input := args[0]

	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("could not open attributes file: %w", err)
	}
	defer f.Close()

	// Make sure the file is vaid before uploading it.
	att, err := attributes.Import(f)
	if err != nil {
		return fmt.Errorf("could not read attribute file: %w", err)
	}

	err = attributes.Validate(att)
	if err != nil {
		return fmt.Errorf("attribute file is invalid: %w", err)
	}

	// TODO: Document the token env var.
	token := os.Getenv(apiToken)
	if token == "" {
		return fmt.Errorf("web3storage auth token not set")
	}

	client, err := w3s.NewClient(w3s.WithToken(token))
	if err != nil {
		return fmt.Errorf("could not create web3storage client: %w", err)
	}

	// Seek back to start of file to upload it (since we've already read it).
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek back to start of attributes file: %w", err)
	}

	cid, err := client.Put(context.Background(), f)
	if err != nil {
		return fmt.Errorf("could not upload attributes file: %w", err)
	}

	fmt.Printf("%v\n", cid.String())

	return nil
}
