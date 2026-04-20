package scoring

import (
	"testing"

	"fraud-detection-platform/go-risk-engine/internal/model"
)

func TestRiskyTransactionGetsReviewOrDecline(t *testing.T) {
	s := &Scorer{
		weights: model.ModelWeights{
			FeatureOrder: []string{
				"amount",
				"is_international",
				"hour_of_day",
				"account_age_days",
				"transactions_last_hour",
				"email_age_days",
			},
			Intercept: -1.9,
			Coefficients: map[string]float64{
				"amount":                 0.0011,
				"is_international":       1.0,
				"hour_of_day":            -0.008,
				"account_age_days":       -0.0007,
				"transactions_last_hour": 0.19,
				"email_age_days":         -0.0006,
			},
			Thresholds: model.Thresholds{
				ApproveBelow: 0.35,
				ReviewBelow:  0.7,
			},
		},
	}

	resp := s.Score(model.ScoreRequest{
		TransactionID:        "txn-1",
		Amount:               3200,
		Currency:             "USD",
		Country:              "NG",
		DeviceTrusted:        false,
		International:        true,
		HourOfDay:            2,
		AccountAgeDays:       10,
		TransactionsLastHour: 9,
		EmailAgeDays:         5,
	})

	if resp.Decision == "APPROVE" {
		t.Fatalf("expected risky transaction not to be approved")
	}
	if resp.RiskScore <= 0 {
		t.Fatalf("expected positive risk score")
	}
}
