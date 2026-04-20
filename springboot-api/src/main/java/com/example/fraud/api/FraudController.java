package com.example.fraud.api;

import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/fraud")
public class FraudController {

    private final RiskEngineClient riskEngineClient;

    public FraudController(RiskEngineClient riskEngineClient) {
        this.riskEngineClient = riskEngineClient;
    }

    @GetMapping("/health")
    public HealthResponse health() {
        return new HealthResponse("ok");
    }

    @PostMapping("/score")
    public FraudResponse score(@Valid @RequestBody FraudRequest request) {
        ScoreResult scoreResult = riskEngineClient.score(request);
        return new FraudResponse(
                request.transactionId(),
                request.amount(),
                scoreResult.decision(),
                scoreResult.riskScore(),
                scoreResult.reasons(),
                java.time.Instant.now()
        );
    }

    public record HealthResponse(String status) {}
}
