package internal

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateUUID() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		panic(err)
	}

	// Set version and variant bits (UUIDv4)
	u[6] = (u[6] & 0x0F) | 0x40 // Version 4
	u[8] = (u[8] & 0x3F) | 0x80 // Variant

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(u[0:4]),
		hex.EncodeToString(u[4:6]),
		hex.EncodeToString(u[6:8]),
		hex.EncodeToString(u[8:10]),
		hex.EncodeToString(u[10:16]),
	)
}
