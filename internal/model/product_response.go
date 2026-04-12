package model

type ProductResponse struct {
	ProductID        string `json:"product_id"`
	MarketPair       string `json:"market_pair"`
	ProductName      string `json:"product_name"`
	BaseCurrency     string `json:"base_currency"`
	QuoteCurrency    string `json:"quote_currency"`
	Status           string `json:"status"`
	IsTradingEnabled bool   `json:"is_trading_enabled"`
	Price            string `json:"price"`
	PriceChange24H   string `json:"price_change_24h"`
	CacheStatus      string `json:"cache_status"`
	Source           string `json:"source"`
}
