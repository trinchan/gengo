package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/trinchan/gengo/lang"
)

const (
	serviceNamespace = "/translate/service"
)

type LanguagePairsResponse struct {
	LanguagePairs []LanguagePairWithPrice
}

func (lp *LanguagePairsResponse) UnmarshalJSON(data []byte) error {
	l := new([]LanguagePairWithPrice)
	err := json.Unmarshal(data, l)
	if err != nil {
		return err
	}
	lp.LanguagePairs = *l
	return nil
}

type LanguagePairWithPrice struct {
	lang.Pair
	Tier
	Currency  string  `json:"currency"`
	UnitPrice Float64 `json:"unit_price"`
}

type LanguagePairsRequest struct {
	Options url.Values
}

type LanguagePairsRequestOption func(*LanguagePairsRequest)

func WithSource(s lang.Code) LanguagePairsRequestOption {
	return func(r *LanguagePairsRequest) {
		r.Options["lc_src"] = []string{string(s)}
	}
}

func NewLanguagePairsRequest(options ...LanguagePairsRequestOption) *LanguagePairsRequest {
	r := &LanguagePairsRequest{Options: url.Values{}}
	for _, option := range options {
		option(r)
	}
	return r
}

func (c *Client) LanguagePairs(req *LanguagePairsRequest) (*LanguagePairsResponse, error) {
	asr := new(LanguagePairsResponse)
	err := c.get(serviceNamespace+"/language_pairs", req.Options, asr)
	return asr, err
}

type LanguagesResponse struct {
	Languages []Language
}

func (l *LanguagesResponse) UnmarshalJSON(b []byte) error {
	langs := new([]Language)
	err := json.Unmarshal(b, langs)
	if err != nil {
		return err
	}
	l.Languages = *langs
	return nil
}

type Language struct {
	UnitType      string `json:"unit_type"`
	Code          string `json:"lc"`
	LocalizedName string `json:"localized_name"`
	Name          string `json:"language"`
}

func (c *Client) Languages() (*LanguagesResponse, error) {
	l := new(LanguagesResponse)
	err := c.get(serviceNamespace+"/languages", nil, l)
	return l, err
}

type QuoteTextRequest struct {
	Jobs []*JobRequest `json:"jobs"`
}

func NewQuoteTextRequest(jobs ...*JobRequest) *QuoteTextRequest {
	return &QuoteTextRequest{Jobs: jobs}
}

type QuoteTextResponse struct {
	Jobs []TextQuote `json:"jobs"`
}

type TextQuote struct {
	Type           string  `json:"type"`
	Credits        Float64 `json:"credits"`
	ETA            int     `json:"eta"`
	UnitCount      Int     `json:"unit_count"`
	DetectedSource string  `json:"lc_src_detected,omitempty"`
	Currency       string  `json:"currency"`
}

func (c *Client) QuoteText(req *QuoteTextRequest) (*QuoteTextResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	qr := new(QuoteTextResponse)
	err = c.post(serviceNamespace+"/quote", bytes.NewReader(b), qr)
	return qr, err
}

type QuoteFileRequest struct {
	Jobs []*FileJobRequest `json:"jobs"`
}

func NewQuoteFileRequest(jobs ...*FileJobRequest) *QuoteFileRequest {
	return &QuoteFileRequest{Jobs: jobs}
}

type QuoteFileResponse struct {
	Jobs []FileQuote `json:"jobs"`
}

type FileQuote struct {
	TextQuote
	Identifier string          `json:"identifier"`
	Error      *FileQuoteError `json:"err,omitempty"`
}

type FileQuoteError struct {
	Name string `json:"filename"`
	Code int    `json:"code"`
	Key  string `json:"key"`
}

func (c *Client) QuoteFile(req *QuoteFileRequest) (*QuoteFileResponse, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for i, fjr := range req.Jobs {
		key := fmt.Sprintf("file_%d", i)
		fjr.FileKey = key
		part, fErr := writer.CreateFormFile(key, filepath.Base(fjr.FilePath))
		if fErr != nil {
			return nil, fErr
		}
		file, fErr := os.Open(fjr.FilePath)
		if fErr != nil {
			return nil, fErr
		}
		_, fErr = io.Copy(part, file)
		if fErr != nil {
			return nil, fErr
		}
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	writer.WriteField("data", string(b))
	qr := new(QuoteFileResponse)
	err = c.multipart(serviceNamespace+"/quote/file", body, writer, qr)
	return qr, err
}
