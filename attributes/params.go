package attributes

import (
	"compress/zlib"
	"errors"
)

const (
	Prefix               = "B7S_"
	Limit                = 50
	AttributesDataLength = 1024 // 1KiB

	binaryRecordSeparator byte = 0x1F // ASCII 'Unit Separator'

	defaultCompression = zlib.BestCompression
)

var (
	ErrNoSignature = errors.New("no signature")
)
