package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"
)

const (
	w3nameAPI = "https://name.web3.storage/name"
)

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

	id, err := w3nameCreate(record, encoded)
	if err != nil {
		return fmt.Errorf("could not create w3name record: %w", err)
	}

	fmt.Printf("%v\n", id)

	return nil
}

func encodeKey(priv crypto.PrivKey) (string, error) {

	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("could not decode peer ID from public key: %w", err)
	}

	name := ipns.NameFromPeer(id)

	return name.String(), nil
}

// TODO: What is the record encoding method?
// TODO: What is the key encoding method?
// TODO: What is the ID returned used for?
func w3nameCreate(rec *ipns.Record, key string) (string, error) {

	cli := &http.Client{}

	payload, err := ipns.MarshalRecord(rec)
	if err != nil {
		return "", fmt.Errorf("could not marshal IPNS record: %w", err)
	}

	enc := base64.StdEncoding.EncodeToString(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", w3nameAPI, key), strings.NewReader(enc))
	if err != nil {
		return "", fmt.Errorf("could not create HTTP request: %w", err)
	}

	token := os.Getenv(apiToken)
	if token == "" {
		return "", fmt.Errorf("web3storage auth token not set")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute HTTP request: %w", err)
	}
	defer res.Body.Close()

	payload, err = io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response payload: %w", err)
	}

	var w3nameResponse struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(payload, &w3nameResponse)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response data: %w", err)
	}

	if w3nameResponse.ID == "" {
		return "", errors.New("received empty ID")
	}

	return w3nameResponse.ID, nil
}
