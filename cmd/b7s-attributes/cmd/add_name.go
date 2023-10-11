package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"

	"github.com/blocklessnetwork/b7s-attributes/gateway"
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

	path := path.FromString(args[0])
	err = path.IsValid()
	if err != nil {
		return fmt.Errorf("path is not valid: %w", err)
	}

	validity := time.Now().Add(flags.validity)
	caching := flags.cache

	record, err := ipns.NewRecord(key, path, flags.sequence, validity, caching)
	if err != nil {
		return fmt.Errorf("could not create IPNS record: %w", err)
	}

	err = ipns.Validate(record, key.GetPublic())
	if err != nil {
		return fmt.Errorf("IPNS record is not valid: %w", err)
	}

	ipnsName, err := ipnsNameFromKey(key)
	if err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	id, err := gateway.CreateIPNSName(record, ipnsName)
	if err != nil {
		return fmt.Errorf("could not create w3name record: %w", err)
	}

	fmt.Printf("%s\n", id)

	return nil
}

func ipnsNameFromKey(priv crypto.PrivKey) (string, error) {

	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("could not decode peer ID from public key: %w", err)
	}

	name := ipns.NameFromPeer(id)

	return name.String(), nil
}
