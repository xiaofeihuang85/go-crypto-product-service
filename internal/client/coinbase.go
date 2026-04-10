package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
)

type CoinbaseClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCoinbaseClient(baseURL string, httpClient *http.Client) *CoinbaseClient {
	return &CoinbaseClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *CoinbaseClient) GetProduct(ctx context.Context, productID string) (model.CoinbaseProduct, error) {
	requestURL := fmt.Sprintf("%s/api/v3/brokerage/market/products/%s", c.baseURL, url.PathEscape(productID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return model.CoinbaseProduct{}, fmt.Errorf("build coinbase request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.CoinbaseProduct{}, fmt.Errorf("request coinbase product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return model.CoinbaseProduct{}, ErrProductNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return model.CoinbaseProduct{}, fmt.Errorf("%w: status=%d body=%s", ErrUnexpectedStatus, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var product model.CoinbaseProduct
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return model.CoinbaseProduct{}, fmt.Errorf("decode coinbase product response: %w", err)
	}

	return product, nil
}

var (
	ErrProductNotFound  = errors.New("coinbase product not found")
	ErrUnexpectedStatus = errors.New("coinbase returned unexpected status")
)
