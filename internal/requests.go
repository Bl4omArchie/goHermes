package hermes

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/html"
)


type HermesNetwork struct {
	Client 	*http.Client
	Logger	*SlogWrapper
}


func NewHermes(logger *SlogWrapper, timeout time.Duration) *HermesNetwork {
	return &HermesNetwork{
		Client: &http.Client{
			Timeout: timeout,
		},
		Logger: logger,
	}
}


func (hn *HermesNetwork) GetPageContent(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		hn.Logger.Log.Error("create request failed", "url", url, "err", err)
		return "", err
	}

	resp, err := hn.Client.Do(req)
	if err != nil {
		hn.Logger.Log.Error("http request failed", "url", url, "err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code %d", resp.StatusCode)
		hn.Logger.Log.Error("bad response", "url", url, "status", resp.StatusCode)
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		hn.Logger.Log.Error("read body failed", "url", url, "err", err)
		return "", err
	}

	return string(data), nil
}


func (hn *HermesNetwork) GetParsedPageContent(ctx context.Context, url string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		hn.Logger.Log.Error("create request failed", "url", url, "err", err)
		return nil, err
	}

	resp, err := hn.Client.Do(req)
	if err != nil {
		hn.Logger.Log.Error("http request failed", "url", url, "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code %d", resp.StatusCode)
		hn.Logger.Log.Error("bad response", "url", url, "status", resp.StatusCode)
		return nil, err
	}

	return html.Parse(resp.Body)
}


func (hn *HermesNetwork) DownloadDocumentReturnHash(ctx context.Context, url, filePath string) (string, error) {
	data, err := hn.GetPageContent(ctx, url)
	if err != nil {
		hn.Logger.Log.Error("download failed", "url", url, "err", err)
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		hn.Logger.Log.Error("mkdir failed", "path", filePath, "err", err)
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		hn.Logger.Log.Error("file create failed", "path", filePath, "err", err)
		return "", err
	}
	defer file.Close()

	if _, err := file.Write([]byte(data)); err != nil {
		hn.Logger.Log.Error("write failed", "path", filePath, "err", err)
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		hn.Logger.Log.Error("seek failed", "path", filePath, "err", err)
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		hn.Logger.Log.Error("hash failed", "err", err)
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
