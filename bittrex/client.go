package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const apiUrl = "https://bittrex.com/api/"

type jsonResponse struct {
	Success bool
	Message string
	Result  json.RawMessage
}

type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

func NewClient(apiKey, apiSecret string) (c *Client) {
	return &Client{apiKey, apiSecret, &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) get(method string, params map[string]string, version string) (jsonResp jsonResponse, e error) {
	url := fmt.Sprintf("%sv%s%s", apiUrl, version, method)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	nonce := time.Now().UnixNano()
	q := req.URL.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	q.Set("apikey", c.apiKey)
	q.Set("nonce", fmt.Sprintf("%d", nonce))

	req.URL.RawQuery = q.Encode()

	mac := hmac.New(sha512.New, []byte(c.apiSecret))
	_, err = mac.Write([]byte(req.URL.String()))
	sign := hex.EncodeToString(mac.Sum(nil))

	req.Header.Add("apisign", sign)

	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &jsonResp)
	return jsonResp, err
}
