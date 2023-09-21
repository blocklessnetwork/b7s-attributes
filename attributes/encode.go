package attributes

import (
	"bytes"
	"fmt"
)

// Encode will encode the list of attributes to binary format.
func Encode(attributes []Attribute) ([]byte, error) {

	data, err := serialize(attributes)
	if err != nil {
		return nil, fmt.Errorf("could not serialize attributes: %w", err)
	}

	compressed, err := compress(data)
	if err != nil {
		return nil, fmt.Errorf("could not compress attributes: %w", err)
	}

	return compressed, nil
}

func serialize(attributes []Attribute) ([]byte, error) {

	var buf bytes.Buffer
	for i, attr := range attributes {

		_, err := buf.WriteString(fmt.Sprintf("%s=%s", attr.Name, attr.Value))
		if err != nil {
			return nil, fmt.Errorf("could not write attribute data: %w", err)
		}

		// If this is the last attribute, no need for the separator.
		if i == len(attributes)-1 {
			break
		}

		err = buf.WriteByte(binaryRecordSeparator)
		if err != nil {
			return nil, fmt.Errorf("could not write separator: %w", err)
		}
	}

	return buf.Bytes(), nil
}
