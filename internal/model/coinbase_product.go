package model

type CoinbaseProduct struct {
	ProductID                string `json:"product_id"`
	Price                    string `json:"price"`
	PricePercentageChange24h string `json:"price_percentage_change_24h"`
	BaseIncrement            string `json:"base_increment"`
	QuoteIncrement           string `json:"quote_increment"`
	BaseName                 string `json:"base_name"`
	QuoteName                string `json:"quote_name"`
	Status                   string `json:"status"`
	QuoteCurrencyID          string `json:"quote_currency_id"`
	BaseCurrencyID           string `json:"base_currency_id"`
	DisplayName              string `json:"display_name"`
}
