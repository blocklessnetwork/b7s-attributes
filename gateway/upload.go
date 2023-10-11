package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/ipfs/go-cid"
)

func Upload(gatewayURL *url.URL, f fs.File, filename string) (cid.Cid, error) {

	var buf bytes.Buffer
	mp := multipart.NewWriter(&buf)

	part, err := mp.CreateFormFile("file", filename)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not create form file: %w", err)
	}

	_, err = io.Copy(part, f)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not copy file to form: %w", err)
	}

	err = mp.Close()
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not close multipart form: %w", err)
	}

	address := gatewayURL
	address.Path = "/api/v1/attributes"

	req, err := http.NewRequest("POST", address.String(), &buf)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", mp.FormDataContentType())

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not execute HTTP request: %w", err)
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not read response payload: %w", err)
	}

	var gatewayUploadResponse struct {
		CID string `json:"cid"`
	}
	err = json.Unmarshal(payload, &gatewayUploadResponse)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not unmarshal response data (response: %s): %w", payload, err)
	}

	if gatewayUploadResponse.CID == "" {
		return cid.Cid{}, fmt.Errorf("CID not found in response (response: %s)", payload)
	}

	uploadedCID, err := cid.Decode(gatewayUploadResponse.CID)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("could not decode CID: %w", err)
	}

	return uploadedCID, nil
}
