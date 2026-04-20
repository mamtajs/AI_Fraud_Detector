import json
from pathlib import Path

import numpy as np
from sklearn.linear_model import LogisticRegression

FEATURES = [
    "amount",
    "is_international",
    "hour_of_day",
    "account_age_days",
    "transactions_last_hour",
    "email_age_days",
]


def build_dataset(n: int = 3000, seed: int = 42):
    rng = np.random.default_rng(seed)

    amount = rng.uniform(5, 5000, size=n)
    is_international = rng.integers(0, 2, size=n)
    hour_of_day = rng.integers(0, 24, size=n)
    account_age_days = rng.uniform(1, 3650, size=n)
    transactions_last_hour = rng.integers(0, 30, size=n)
    email_age_days = rng.uniform(1, 3650, size=n)

    X = np.column_stack(
        [
            amount,
            is_international,
            hour_of_day,
            account_age_days,
            transactions_last_hour,
            email_age_days,
        ]
    )

    risk_signal = (
        0.0012 * amount
        + 1.2 * is_international
        + 0.11 * np.where((hour_of_day <= 5) | (hour_of_day >= 23), 1, 0)
        - 0.0009 * account_age_days
        + 0.22 * transactions_last_hour
        - 0.0008 * email_age_days
    )

    logits = risk_signal - 2.2
    probabilities = 1 / (1 + np.exp(-logits))
    y = (rng.random(n) < probabilities).astype(int)
    return X, y


def main():
    X, y = build_dataset()
    model = LogisticRegression(max_iter=2000)
    model.fit(X, y)

    payload = {
        "feature_order": FEATURES,
        "intercept": float(model.intercept_[0]),
        "coefficients": {
            name: float(weight) for name, weight in zip(FEATURES, model.coef_[0])
        },
        "thresholds": {
            "approve_below": 0.35,
            "review_below": 0.7,
        },
    }

    output_path = (
        Path(__file__).resolve().parents[1]
        / "go-risk-engine"
        / "model"
        / "weights.json"
    )
    output_path.write_text(json.dumps(payload, indent=2))
    print(f"Exported model weights to {output_path}")


if __name__ == "__main__":
    main()
