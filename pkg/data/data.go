package data

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func CreateRandomString() string {
	u := uuid.New()
	fmt.Println(u.Domain())
	fmt.Println(u.ID())
	fmt.Println(u.MarshalText())
	fmt.Println(u.String())

	return u.String()
}

func EncodeBase64(target string) string {
	t := []byte(target)
	return base64.RawStdEncoding.EncodeToString(t)
}
