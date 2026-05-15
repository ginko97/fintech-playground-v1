# Fintech Playground V1

A high-performance, resilient fintech ledger system built with **Go** and **PostgreSQL**. This project demonstrates industry-standard practices for handling financial transactions, focusing on data integrity, idempotency, and system resilience.

## System Architecture

The project follows **Clean Architecture** principles, separating business logic from infrastructure and external dependencies.

### Payment Flow & Resilience
The following diagram illustrates how the system handles external bank communication and manages the "Uncertain" state during network failures.

```mermaid
sequenceDiagram
    participant C as Client (API/Mobile)
    participant U as Transaction Usecase
    participant R as Ledger Repository (DB)
    participant G as Payment Gateway (Bank)

    Note over C,G: Week 2: External Orchestration
    C->>U: ProcessPayment(tx)
    U->>R: Create (Status: PENDING)
    R-->>U: Success
    
    U->>G: Authorize(tx)
    
    alt Bank Success
        G-->>U: PaymentResponse (SUCCESS)
        U->>R: UpdateStatus (SUCCESS)
        U-->>C: Transaction Completed
    else Bank Timeout (Context Deadline Exceeded)
        G--XU: Timeout
        U->>R: UpdateStatus (UNCERTAIN)
        Note over U,R: Resilience: No data lost during failure!
        U-->>C: Error: Status Uncertain
    end
```