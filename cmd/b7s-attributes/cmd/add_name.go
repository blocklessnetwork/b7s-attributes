package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"

	"github.com/blocklessnetwork/b7s-attributes/w3s"
)

// TODO: Sequence number management. Currently it's up to the caller to know the correct sequence number.
// Sequence number for the IPNS record needs to be updated when the path changes.

var flagsAddName struct {
	key      string
	validity time.Duration
	cache    time.Duration
	sequence uint64
}

func runAddName(_ *cobra.Command, args []string) error {

	flags := flagsAddName
	if flags.key == "" {
		return errors.New("key is required")
	}

	key, err := readPrivateKey(flags.key)
	if err != nil {
		return fmt.Errorf("could not read private key: %w", err)
	}

	ipfsPath := normalizeIPFSPath(args[0])
	path := path.FromString(ipfsPath)

	validity := time.Now().Add(flags.validity)
	caching := flags.cache

	record, err := ipns.NewRecord(key, path, flags.sequence, validity, caching)
	if err != nil {
		return fmt.Errorf("could not create IPNS record: %w", err)
	}

	encoded, err := encodeKey(key)
	if err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	id, err := w3s.CreateIPNSName(record, encoded)
	if err != nil {
		return fmt.Errorf("could not create w3name record: %w", err)
	}

	fmt.Printf("%v\n", id)

	return nil
}

// normalizeIPFSPath will make sure the path is in '/ipfs/<cid>' format.
func normalizeIPFSPath(path string) string {
	if strings.HasPrefix(path, "/ipfs") {
		return path
	}

	path = strings.TrimPrefix(path, "/")
	out := fmt.Sprintf("/ipfs/%s", path)

	return out
}

func encodeKey(priv crypto.PrivKey) (string, error) {

	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("could not decode peer ID from public key: %w", err)
	}

	name := ipns.NameFromPeer(id)

	return name.String(), nil
}
