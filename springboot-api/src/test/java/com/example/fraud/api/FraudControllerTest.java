package com.example.fraud.api;

import com.example.fraud.FraudApiApplication;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.client.RestClientTest;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;

import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.BDDMockito.given;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest(classes = FraudApiApplication.class)
@AutoConfigureMockMvc
class FraudControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private RiskEngineClient riskEngineClient;

    @Autowired
    private ObjectMapper objectMapper;

    @Test
    void scoreReturnsDecision() throws Exception {
        given(riskEngineClient.score(any())).willReturn(
                new ScoreResult("txn-1", 0.82, "DECLINE", List.of("high transaction amount"))
        );

        FraudRequest request = new FraudRequest(
                "txn-1",
                3200,
                "USD",
                "NG",
                false,
                true,
                2,
                5,
                11,
                3
        );

        mockMvc.perform(post("/api/v1/fraud/score")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.transactionId").value("txn-1"))
                .andExpect(jsonPath("$.decision").value("DECLINE"));
    }
}
