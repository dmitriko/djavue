package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

const TokenSize = 20

func NewToken() (string, error) {
	b := make([]byte, TokenSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewULIDNow() (string, error) {
	return NewULID(time.Now())
}

func NewULID(tm time.Time) (string, error) {
	id, err := ulid.New(ulid.Timestamp(tm), rand.Reader)
	return fmt.Sprintf("%s", id), err
}
