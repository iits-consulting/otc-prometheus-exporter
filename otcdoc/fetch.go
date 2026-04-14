package otcdoc

import (
	"fmt"
	"io"
	"net/http"
)

func fetchBytes(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: HTTP %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}
	return body, nil
}

// FetchDocumentationSource fetches and parses the RST file for ds.
// The OTC namespace (e.g. "SYS.NAT") is parsed from the RST content itself.
// Sources without a GithubRawUrl return an empty page without error.
func FetchDocumentationSource(ds DocumentationSource) (DocumentationPage, error) {
	if ds.GithubRawUrl == "" {
		return DocumentationPage{}, nil
	}
	body, err := fetchBytes(ds.GithubRawUrl)
	if err != nil {
		return DocumentationPage{}, err
	}
	page, err := ParseDocumentationPageFromRstBytes(body)
	if err != nil {
		return page, fmt.Errorf("parse %s: %w", ds.GithubRawUrl, err)
	}
	return page, nil
}

// FetchMarkdownMetrics fetches and parses a Huawei Cloud metric catalog Markdown file.
func FetchMarkdownMetrics(url string) ([]MetricDocumentationEntry, error) {
	body, err := fetchBytes(url)
	if err != nil {
		return nil, err
	}
	return ParseHuaweiMarkdownMetrics(body), nil
}
