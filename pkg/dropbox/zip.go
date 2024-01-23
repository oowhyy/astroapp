package dropbox

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func (client *Client) FetchZip(zipPath string) (*zip.Reader, error) {
	f, err := client.Files.Download(&DownloadInput{
		Path: zipPath,
	})
	if err != nil {
		return nil, fmt.Errorf("db download error: %w", err)
	}
	buff := new(bytes.Buffer)
	size, err := io.Copy(buff, f.Body)
	if err != nil {
		return nil, fmt.Errorf("io copy error: %w", err)
	}
	f.Body.Close()
	br := bytes.NewReader(buff.Bytes())
	reader, err := zip.NewReader(br, size)
	// reader, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, fmt.Errorf("open zip reader error: %w", err)
	}
	return reader, nil
}

func (client *Client) FetchFile(path string) (io.ReadCloser, error) {
	f, err := client.Files.Download(&DownloadInput{
		Path: path,
	})
	if err != nil {
		return nil, fmt.Errorf("db download error: %w", err)
	}
	return f.Body, nil
}
