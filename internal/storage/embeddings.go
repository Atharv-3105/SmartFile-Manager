package storage

import(
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
)

//This will store a vector embedding for a file
func InsertEmbedding(
	tx *sql.Tx,
	fileID  int64,
	model  string,
	vector  []float32,
) error {
	if len(vector) == 0 {
		return fmt.Errorf("[DB] embedding vector is empty")
	}

	blob, err := float32SliceToBlob(vector)
	if err != nil {
		return fmt.Errorf("[DB] encode embedding: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO embeddings (file_id, model, vector, dim)
		VALUES	(?, ?, ?, ?)`,
		fileID,model,blob,len(vector),	
	)
	if err != nil {
		return fmt.Errorf("[DB] insert embedding: %w", err)
	}

	return nil
}


//Function to Encode from Float32 into Binary Blob
func float32SliceToBlob(vec []float32) ([]byte, error) {
	buf := new(bytes.Buffer)

	//Use little-endian method
	//It's a method of storing multi-byte data types(i.e INT,FLOAT)
	//in memory where the LSB(least-significant byte) is stored at the Smallest Memory Address.
	if err := binary.Write(buf, binary.LittleEndian, vec); err != nil {
		return nil, err 
	}

	return buf.Bytes(), nil
}