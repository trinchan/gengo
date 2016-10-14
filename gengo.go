package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trinchan/gengo/sign"
)

var logger *log.Logger

const (
	// SandboxBaseURL is the Gengo sandbox base URL
	SandboxBaseURL = "http://api.sandbox.gengo.com/v2"
	// ProductionBaseURL is the Gengo production base URL
	ProductionBaseURL = "http://api.gengo.com/v2"
)

// SetLogger let's library users supply a logger, so that api debugging
// can be logged along with the application's debugging info.
func SetLogger(l *log.Logger) {
	logger = l
}

// Client is a client for connecting to the Gengo APi
type Client struct {
	PublicKey    string
	PrivateKey   string
	BaseURL      string
	RoundTripper http.RoundTripper
	signer       sign.Signer
}

// New creates a new Gengo Client with the given keys and base URL
func New(publicKey, privateKey, baseURL string) *Client {
	return &Client{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		BaseURL:      baseURL,
		RoundTripper: http.DefaultTransport,
		signer:       sign.NewHMACSigner([]byte(privateKey)),
	}
}

// NewFromEnv creates a new Gengo Client from environment variables
// GENGO_PUBLIC_KEY = public key
// GENGO_PRIVATE_KEY = private key
// Set GENGO_PRODUCTION to anything non-empty to use production URL
func NewFromEnv() *Client {
	publicKey, privateKey, sandbox := os.Getenv("GENGO_PUBLIC_KEY"), os.Getenv("GENGO_PRIVATE_KEY"), os.Getenv("GENGO_PRODUCTION") == ""
	if publicKey == "" {
		log.Fatal("Public key not set (check GENGO_PUBLIC_KEY environment variable)")
	}
	if privateKey == "" {
		log.Fatal("Private key not set (check GENGO_PRIVATE_KEY environment variable)")
	}
	baseURL := SandboxBaseURL
	if !sandbox {
		baseURL = ProductionBaseURL
	}
	return New(publicKey, privateKey, baseURL)
}

// SetRoundTripper allows a custom HTTP RoundTripper to be used
func (c *Client) SetRoundTripper(rt http.RoundTripper) {
	c.RoundTripper = rt
}

const (
	// OPStatOK marks a successful API request
	OPStatOK = "ok"
	// OPStatError marks an errored API request
	OPStatError = "error"
)

type response struct {
	OPStat   string          `json:"opstat"`
	Response json.RawMessage `json:"response,omitempty"`
	Error    ErrorResponse   `json:"err,omitempty"`
}

// ErrorResponse describes the Gengo error response
type ErrorResponse struct {
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (c *Client) urlEncoded(method, path string, params url.Values, resp interface{}) error {
	req, err := http.NewRequest(method, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	vals := c.vals(nil)
	for key, param := range params {
		for _, p := range param {
			vals.Add(key, p)
		}
	}
	req.URL.RawQuery = vals.Encode()
	return c.do(req, resp)
}

func (c *Client) formEncoded(method, path string, body io.Reader, resp interface{}) error {
	vals := c.vals(body)
	req, err := http.NewRequest(method, c.BaseURL+path, strings.NewReader(vals.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req, resp)
}

func (c *Client) get(path string, params url.Values, resp interface{}) error {
	return c.urlEncoded(http.MethodGet, path, params, resp)
}

func (c *Client) delete(path string, params url.Values, resp interface{}) error {
	return c.urlEncoded(http.MethodDelete, path, params, resp)
}

func (c *Client) post(path string, body io.Reader, resp interface{}) error {
	return c.formEncoded(http.MethodPost, path, body, resp)
}

func (c *Client) put(path string, body io.Reader, resp interface{}) error {
	return c.formEncoded(http.MethodPut, path, body, resp)
}

func (c *Client) multipart(path string, body *bytes.Buffer, writer *multipart.Writer, resp interface{}) error {
	ts := strconv.Itoa(int(time.Now().Unix()))
	writer.WriteField("api_key", c.PublicKey)
	writer.WriteField("api_sig", c.signer.Sign(ts))
	writer.WriteField("ts", ts)
	err := writer.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return c.do(req, resp)
}

func (c *Client) vals(body io.Reader) url.Values {
	ts := strconv.Itoa(int(time.Now().Unix()))
	vals := url.Values{}
	vals.Add("api_key", c.PublicKey)
	vals.Add("api_sig", c.signer.Sign(ts))
	vals.Add("ts", ts)
	if body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		vals.Add("data", buf.String())
	}
	return vals
}
func (c *Client) do(req *http.Request, resp interface{}) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "Gengo Go Library; Version 0.0.1; https://www.gengo.com")
	re, err := c.RoundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	defer re.Body.Close()
	dec := json.NewDecoder(re.Body)
	r := new(response)
	err = dec.Decode(r)
	if err != nil {
		return err
	}
	if r.OPStat != OPStatOK {
		return r.Error
	}
	if len(r.Response) > 0 {
		err = json.Unmarshal(r.Response, resp)
	}
	return err
}
