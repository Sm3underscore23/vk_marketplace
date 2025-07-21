package feed

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"marketplace/internal/models"
)

func generateCursor(aesgcm cipher.AEAD, lastAdData models.LastAdData) (string, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	plaintext, err := json.Marshal(lastAdData)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	full := append(nonce, ciphertext...)

	cursor := base64.URLEncoding.EncodeToString(full)

	return cursor, nil

}

func parseCursor(aesgcm cipher.AEAD, cursor string) (models.LastAdData, error) {
	var lastAdData models.LastAdData

	full, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return lastAdData, err
	}
	if len(full) < 12 {
		return lastAdData, models.ErrorShortCursor
	}

	nonce := full[:12]
	ciphertext := full[12:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return lastAdData, err
	}

	err = json.Unmarshal(plaintext, &lastAdData)
	if err != nil {
		return models.LastAdData{}, err
	}

	return lastAdData, nil
}
