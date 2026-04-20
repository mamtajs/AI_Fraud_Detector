package model

type ScoreRequest struct {
	TransactionID        string  `json:"transactionId"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	Country              string  `json:"country"`
	DeviceTrusted        bool    `json:"deviceTrusted"`
	International        bool    `json:"international"`
	HourOfDay            int     `json:"hourOfDay"`
	AccountAgeDays       int     `json:"accountAgeDays"`
	TransactionsLastHour int     `json:"transactionsLastHour"`
	EmailAgeDays         int     `json:"emailAgeDays"`
}

type ScoreResponse struct {
	TransactionID string   `json:"transactionId"`
	RiskScore     float64  `json:"riskScore"`
	Decision      string   `json:"decision"`
	Reasons       []string `json:"reasons"`
}

type ModelWeights struct {
	FeatureOrder []string           `json:"feature_order"`
	Intercept    float64            `json:"intercept"`
	Coefficients map[string]float64 `json:"coefficients"`
	Thresholds   Thresholds         `json:"thresholds"`
}

type Thresholds struct {
	ApproveBelow float64 `json:"approve_below"`
	ReviewBelow  float64 `json:"review_below"`
}
