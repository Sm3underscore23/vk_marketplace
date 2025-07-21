package validator

import (
	"context"
	"io"
	"marketplace/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (v *validatorWrapper) CheckImage(ctx context.Context, imageURL string) error {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, imageURL, nil)
	if err != nil {
		return models.ErrorImageInvalidURL
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.ErrorImageNotAccessible
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.ErrorImageNonSuccessStatusCode
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return models.ErrorImageInvalidContentType
	}

	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return models.ErrorImageInvalidContentLengthHeader
		}
		if size > v.maxImageSizeBytes {
			return models.ErrorImageTooLarge
		}
		return nil
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return models.ErrorImageGETRequestFailed
	}

	resp, err = client.Do(req)
	if err != nil {
		return models.ErrorImageGETRequestFailed
	}
	defer resp.Body.Close()

	limitedReader := io.LimitReader(resp.Body, v.maxImageSizeBytes+1)
	n, err := io.Copy(io.Discard, limitedReader)
	if err != nil {
		return models.ErrorImageReadFailed
	}

	if n > v.maxImageSizeBytes {
		return models.ErrorImageTooLarge
	}

	return nil
}
