package w3s

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ipfs/boxo/ipns"
)

const (
	w3nameAPI = "https://name.web3.storage/name"
)

func CreateIPNSName(rec *ipns.Record, key string) (string, error) {

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

	token := os.Getenv(APITokenEnv)
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
		return "", fmt.Errorf("ID not found in response (response: %s)", payload)
	}

	return w3nameResponse.ID, nil
}
