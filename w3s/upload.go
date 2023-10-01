package w3s

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

func Upload(f fs.File) (cid.Cid, error) {

	token := os.Getenv(APITokenEnv)
	if token == "" {
		return cid.Cid{}, fmt.Errorf("web3storage auth token not set")
	}

	client, err := w3s.NewClient(w3s.WithToken(token))
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not create web3storage client: %w", err)
	}

	id, err := client.Put(context.Background(), f)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not upload file: %w", err)
	}

	return id, nil
}
