package arbitrage

import (
	"io"
	"os"
)

func ReadFile(filename string) ([]byte, error) {
	handle, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	bytes, err := io.ReadAll(handle)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
