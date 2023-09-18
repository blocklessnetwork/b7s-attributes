package attributes

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func compress(data []byte) ([]byte, error) {

	var compressed bytes.Buffer
	w, err := zlib.NewWriterLevel(&compressed, defaultCompression)
	if err != nil {
		return nil, fmt.Errorf("could not create zlib writer: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("could not write data to zlib writer: %w", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close zlib writer: %w", err)
	}

	return compressed.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {

	input := bytes.NewReader(data)
	reader, err := zlib.NewReader(input)
	if err != nil {
		return nil, fmt.Errorf("could not create zlib reader: %w", err)
	}

	var out bytes.Buffer
	_, err = io.Copy(&out, reader)
	if err != nil {
		return nil, fmt.Errorf("could not decompress zlib data stream: %w", err)
	}

	return out.Bytes(), nil
}
