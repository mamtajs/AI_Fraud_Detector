package com.example.fraud.api;

import java.util.List;

public record ScoreResult(
        String transactionId,
        double riskScore,
        String decision,
        List<String> reasons
) {}
