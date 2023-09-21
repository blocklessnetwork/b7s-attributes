package attributes

import (
	"compress/zlib"
)

const (
	Prefix = "B7S_"
	Limit  = 50

	binaryRecordSeparator byte = 0x1F // ASCII 'Unit Separator'

	defaultCompression = zlib.BestCompression
)

type Attribute struct {
	Name  string
	Value string
}
