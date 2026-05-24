# Fintech Playground V1

A high-performance, resilient fintech ledger system built with **Go** and **PostgreSQL**. This project demonstrates industry-standard practices for handling financial transactions, focusing on data integrity, idempotency, and system resilience.

## System Architecture

The project follows **Clean Architecture** principles, separating business logic from infrastructure and external dependencies.

### Payment Flow & Resilience
The following diagram illustrates how the system handles external bank communication, manages the "Uncertain" state during network failures, and validates requests safely through the HTTP edge.

```mermaid
sequenceDiagram
    participant C as Client (API/Mobile)
    participant H as Transaction Handler (Gin)
    participant U as Transaction Usecase
    participant R as Ledger Repository (DB)
    participant G as Payment Gateway (Bank)

    Note over C,G: Week 3: HTTP API Gateway Edge
    C->>H: POST /v1/payments (JSON Payload)
    Note over H: Validate Payload<br/>Parse Account UUID
    
    H->>U: ProcessPayment(ctx, tx)
    Note over U,G: Week 2: External Orchestration
    U->>R: Create (Status: pending)
    
    Note over R: Idempotency Check:<br/>Inserts gatekeeper key
    R-->>U: Success (Stored)
    
    U->>G: Authorize(ctx, tx)
    
    alt Bank Success
        G-->>U: PaymentResponse (success)
        U->>R: UpdateStatus(id, success)
        U-->>H: Return updated processedTx
        H-->>C: HTTP 200 {"id": "...", "status": "success"}
    else Bank Timeout (Context Deadline Exceeded)
        G--XU: Timeout
        U->>R: UpdateStatus(id, uncertain)
        Note over U,R: Resilience: No data lost during failure!
        U-->>H: Return updated processedTx (uncertain)
        H-->>C: HTTP 500 {"error": "Payment processing failed", "details": "..."}
    end
```