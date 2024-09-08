package bot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/xid"
)

func downloadFileFromLink(ctx context.Context, fileUrl string, mime string) (string, error) {
	fileType, err := mimeToType(mime)
	if err != nil {
		return "", fmt.Errorf("failed to convert mime: %w", err)
	}

	s := xid.New().String() + fileType
	out, err := os.Create(s)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileUrl, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy body: %w", err)
	}

	return s, nil
}

func mimeToType(mime string) (string, error) {
	switch mime {
	case "audio/mpeg":
		return ".mp3", nil
	case "audio/ogg":
		return ".ogg", nil
	case "audio/vnd.wav":
		return ".wav", nil
	default:
		return "", errors.New("unknown mime type")
	}
}
