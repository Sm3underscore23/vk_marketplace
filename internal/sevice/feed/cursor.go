package feed

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marketplace/internal/models"
	"strings"
)

func generateFeedURI(ctx context.Context, baseURl string, uriParams models.FeedURIParams, cursor string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s/feed?%s=%s&%s=%s&%s=%s",
		baseURl,
		models.URIParamCursor, cursor,
		models.URIParamSortBy, uriParams.SortBy,
		models.URIParamOrder, uriParams.Order,
	))

	if uriParams.PriceMin > 0 {
		builder.WriteString(fmt.Sprintf("&%s=%v",
			models.URIParamPriceMin, uriParams.PriceMin))
	}

	if uriParams.PriceMax > 0 {
		builder.WriteString(fmt.Sprintf("&%s=%v",
			models.URIParamPriceMax, uriParams.PriceMax))
	}

	if !uriParams.CreatedAfter.IsZero() {
		builder.WriteString(fmt.Sprintf("&%s=%T",
			models.URIParamCreatedAfter, uriParams.CreatedAfter))
	}

	builder.WriteString(fmt.Sprintf("&%s=%v",
		models.URIParamLimit, uriParams.Limit))

	return builder.String()
}

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
