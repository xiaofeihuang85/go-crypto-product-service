package model

type ProductResponse struct {
	ProductID                string `json:"product_id"`
	DisplayName              string `json:"display_name"`
	BaseCurrencyID           string `json:"base_currency_id"`
	BaseName                 string `json:"base_name"`
	QuoteCurrencyID          string `json:"quote_currency_id"`
	QuoteName                string `json:"quote_name"`
	Status                   string `json:"status"`
	Price                    string `json:"price"`
	PricePercentageChange24h string `json:"price_percentage_change_24h"`
	BaseIncrement            string `json:"base_increment"`
	QuoteIncrement           string `json:"quote_increment"`
}
