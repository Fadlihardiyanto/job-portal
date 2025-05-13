package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

func MapToJSON(m map[string]string) string {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return `{"error": "failed to process error message"}`
	}
	return string(jsonData)
}

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func StringToInt(s string) (int, error) {
	var i int
	err := json.Unmarshal([]byte(s), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}
