package attributes

import (
	"bytes"
	"fmt"
	"strings"
)

func DecodeAttributes(data []byte) ([]Attribute, error) {

	decompressed, err := decompress(data)
	if err != nil {
		return nil, fmt.Errorf("could not decompress attributes: %w", err)
	}

	attributes, err := deserializeAttributes(decompressed)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize attributes: %w", err)
	}

	return attributes, nil
}

func deserializeAttributes(data []byte) ([]Attribute, error) {

	records := bytes.Split(data, []byte{binaryRecordSeparator})
	out := make([]Attribute, 0, len(records))

	for _, rec := range records {

		line := string(rec)

		fields := strings.SplitN(line, "=", 2)

		if len(fields) != 2 {
			return nil, fmt.Errorf("unexpected attribute format (line: %s)", line)
		}

		attr := Attribute{
			Name:  fields[0],
			Value: fields[1],
		}
		out = append(out, attr)
	}

	return out, nil
}
