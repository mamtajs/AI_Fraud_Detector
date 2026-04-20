# Architecture Notes

## Flow
1. Client submits a transaction payload to the Spring Boot API.
2. Spring Boot validates the request and forwards it to the Go scoring engine.
3. The Go engine loads exported ML weights and computes a fraud probability.
4. Business rules adjust the risk score for edge cases:
   - high amount
   - high transaction velocity
   - international charge from untrusted device
   - odd-hour activity
5. The final result is returned as one of:
   - APPROVE
   - REVIEW
   - DECLINE

## Why split Java and Go?
- **Spring Boot** is great for enterprise APIs, validation, and future integrations.
- **Go** is excellent for lightweight, low-latency inference and high-throughput scoring.

## ML strategy
The hot path does not depend on Python runtime. Instead:
- Python trains a logistic regression model offline
- model coefficients are exported to JSON
- Go performs online inference using those coefficients

This is a practical architecture pattern because it keeps scoring fast and portable.

## Future scale path
- Kafka for streaming transactions
- PostgreSQL for immutable fraud audit records
- Redis for device, account, and velocity lookups
- Prometheus/Grafana for observability
- Kubernetes for deployment
