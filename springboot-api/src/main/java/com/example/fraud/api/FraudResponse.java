package com.example.fraud.api;

import java.time.Instant;
import java.util.List;

public record FraudResponse(
        String transactionId,
        double amount,
        String decision,
        double riskScore,
        List<String> reasons,
        Instant evaluatedAt
) {}
