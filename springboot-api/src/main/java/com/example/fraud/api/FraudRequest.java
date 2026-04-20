package com.example.fraud.api;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

public record FraudRequest(
        @NotBlank String transactionId,
        @Positive double amount,
        @NotBlank String currency,
        @NotBlank String country,
        @NotNull Boolean deviceTrusted,
        @NotNull Boolean international,
        @Min(0) @Max(23) int hourOfDay,
        @Min(0) int accountAgeDays,
        @Min(0) int transactionsLastHour,
        @Min(0) int emailAgeDays
) {}
