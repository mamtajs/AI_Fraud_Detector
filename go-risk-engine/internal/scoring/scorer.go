package scoring

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"

	"fraud-detection-platform/go-risk-engine/internal/model"
)

type Scorer struct {
	weights model.ModelWeights
}

func NewScorer(weightsPath string) (*Scorer, error) {
	data, err := os.ReadFile(weightsPath)
	if err != nil {
		return nil, fmt.Errorf("read weights: %w", err)
	}

	var w model.ModelWeights
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("unmarshal weights: %w", err)
	}

	return &Scorer{weights: w}, nil
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (s *Scorer) Score(req model.ScoreRequest) model.ScoreResponse {
	features := map[string]float64{
		"amount":                 req.Amount,
		"is_international":       boolToFloat(req.International),
		"hour_of_day":            float64(req.HourOfDay),
		"account_age_days":       float64(req.AccountAgeDays),
		"transactions_last_hour": float64(req.TransactionsLastHour),
		"email_age_days":         float64(req.EmailAgeDays),
	}

	logit := s.weights.Intercept
	contributions := map[string]float64{}
	for _, feature := range s.weights.FeatureOrder {
		contrib := s.weights.Coefficients[feature] * features[feature]
		contributions[feature] = contrib
		logit += contrib
	}

	score := sigmoid(logit)
	reasons := topReasons(contributions, req)

	if req.Amount > 2500 {
		score += 0.08
		reasons = append(reasons, "high transaction amount")
	}
	if req.TransactionsLastHour >= 8 {
		score += 0.07
		reasons = append(reasons, "high transaction velocity")
	}
	if req.International && !req.DeviceTrusted {
		score += 0.12
		reasons = append(reasons, "international charge from untrusted device")
	}
	if score > 0.99 {
		score = 0.99
	}

	decision := "APPROVE"
	if score >= s.weights.Thresholds.ReviewBelow {
		decision = "DECLINE"
	} else if score >= s.weights.Thresholds.ApproveBelow {
		decision = "REVIEW"
	}

	return model.ScoreResponse{
		TransactionID: req.TransactionID,
		RiskScore:     round(score, 4),
		Decision:      decision,
		Reasons:       uniqueStrings(reasons),
	}
}

func topReasons(contributions map[string]float64, req model.ScoreRequest) []string {
	type pair struct {
		Feature string
		Value   float64
	}

	pairs := make([]pair, 0, len(contributions))
	for k, v := range contributions {
		pairs = append(pairs, pair{Feature: k, Value: v})
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Value > pairs[j].Value })

	labels := map[string]string{
		"amount":                 "amount increased fraud risk",
		"is_international":       "international transaction",
		"hour_of_day":            "unusual transaction time",
		"account_age_days":       "newer account profile",
		"transactions_last_hour": "burst of recent transactions",
		"email_age_days":         "newer email profile",
	}

	result := []string{}
	for _, p := range pairs {
		if len(result) == 3 {
			break
		}
		if p.Value > 0 {
			result = append(result, labels[p.Feature])
		}
	}

	if req.HourOfDay <= 5 || req.HourOfDay >= 23 {
		result = append(result, "transaction occurred during odd hours")
	}

	return result
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func boolToFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func round(value float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Round(value*pow) / pow
}
