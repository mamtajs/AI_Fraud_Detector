package com.example.fraud.api;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;

@Component
public class RiskEngineClient {

    private final RestClient restClient;
    private final String riskEngineBaseUrl;

    public RiskEngineClient(
            RestClient restClient,
            @Value("${risk-engine.base-url}") String riskEngineBaseUrl
    ) {
        this.restClient = restClient;
        this.riskEngineBaseUrl = riskEngineBaseUrl;
    }

    public ScoreResult score(FraudRequest request) {
        return restClient.post()
                .uri(riskEngineBaseUrl + "/score")
                .contentType(MediaType.APPLICATION_JSON)
                .body(request)
                .retrieve()
                .body(ScoreResult.class);
    }
}
