package gateway

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ipfs/boxo/ipns"
)

type attributeCreateRequest struct {
	Name   string `json:"ipnsName"`
	Record string `json:"ipnsRecord"`
}

func CreateIPNSName(gateway *url.URL, rec *ipns.Record, key string) (string, error) {

	addr := gateway
	addr.Path = ipnsNameEndpoint

	payload, err := encodeIPNSRequest(rec, key)
	if err != nil {
		return "", fmt.Errorf("could not encode IPNS request: %w", err)
	}

	req, err := http.NewRequest("POST", addr.String(), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("could not create HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	cli := &http.Client{}
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
		return "", fmt.Errorf("could not unmarshal response data (response: %s): %w", payload, err)
	}

	if w3nameResponse.ID == "" {
		return "", fmt.Errorf("ID not found in response (response: %s)", payload)
	}

	return w3nameResponse.ID, nil
}

func encodeIPNSRequest(rec *ipns.Record, name string) ([]byte, error) {

	payload, err := ipns.MarshalRecord(rec)
	if err != nil {
		return nil, fmt.Errorf("could not marshal IPNS record: %w", err)
	}
	enc := base64.StdEncoding.EncodeToString(payload)

	req := attributeCreateRequest{
		Name:   name,
		Record: enc,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not JSON encode request: %w", err)
	}

	return data, nil
}
