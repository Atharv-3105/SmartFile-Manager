package search

import (
	"bytes"
	"encoding/binary"
	"fmt"
)


//Function to convert from BLOB to a slice of float32
func BlobToFloat32Slice(blob []byte) ([]float32, error) {
	if len(blob) % 4 != 0{
		return nil, fmt.Errorf("[SEARCH] invalid blob length %d (not multiple of 4)", len(blob))
	}

	vec := make([]float32, len(blob)/4)

	reader := bytes.NewReader(blob)
	if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
		return nil, err
	}

	return vec, nil
}