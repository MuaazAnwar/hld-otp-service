# HLD OTP Service

A high-level design (HLD) for a scalable and resilient One-Time Password (OTP) service used for user authentication (e.g., login, payments) over SMS.

---

## 🧩 Problem Statement

Design an OTP service that can:

- Handle high volumes of OTP requests from users.
- Deliver OTPs via SMS reliably and quickly.
- Prevent abuse (e.g., bot spamming, brute forcing).
- Support fallback to secondary SMS providers when the primary is unavailable.
- Be scalable, secure, and resilient against failures at both system and external dependency levels.

Constraints:
- OTPs must expire after a short time (e.g., 5 minutes).
- Users should not be able to request OTPs more than a limited number of times in a short window.
- External SMS providers may fail intermittently or have rate limits.

---

## ✅ Solution Overview

The OTP service is split into independent components:

- **API Gateway**: Handles incoming requests, applies rate limiting (IP/UserID), and routes to OTP service.
- **OTP Generator**: Creates a secure OTP and stores it with TTL.
- **Queue**: Decouples OTP generation from delivery for async handling.
- **Worker**: Picks tasks from the queue and sends OTP using a primary SMS provider; if that fails, switches to a backup provider using a Redis-based health check.
- **Redis**: Stores rate limiting counters and tracks SMS provider health.
- **SMS Gateways**: External APIs like Twilio or Nexmo (mocked in this design).

---

## 🗺️ Architecture Diagram

```mermaid
graph TD
    %% Styling
    classDef primary fill:#4a90e2,stroke:#333,stroke-width:2px,color:white,font-weight:bold
    classDef secondary fill:#50c878,stroke:#333,stroke-width:2px,color:white,font-weight:bold
    classDef storage fill:#f5a623,stroke:#333,stroke-width:2px,color:white,font-weight:bold
    classDef external fill:#d0021b,stroke:#333,stroke-width:2px,color:white,font-weight:bold
    classDef subgraphStyle fill:#f8f9fa,stroke:#333,stroke-width:1px,color:#333

    %% Main Nodes
    Client[Client]

    %% API Layer
    subgraph APILayer[API Layer]
        style APILayer subgraphStyle
        APIGateway[API Gateway]
    end

    %% Service Layer
    subgraph ServiceLayer[Service Layer]
        style ServiceLayer subgraphStyle
        OTPService[OTP Service]
        Worker[Worker]
    end

    %% Storage Layer
    subgraph StorageLayer[Storage Layer]
        style StorageLayer subgraphStyle
        Redis[(Redis)]
        Database[(Database)]
        Queue[(Queue)]
    end

    %% External Layer
    subgraph ExternalLayer[External Services]
        style ExternalLayer subgraphStyle
        SMSProviderPrimary[SMS Provider Primary]
        SMSProviderSecondary[SMS Provider Secondary]
    end

    %% Main Flow
    Client -->|Request OTP| APIGateway
    APIGateway -->|Rate limit check| Redis
    APIGateway -->|Forward request| OTPService

    OTPService -->|Store OTP| Database
    OTPService -->|Queue task| Queue

    Queue -->|Process task| Worker
    Worker -->|Check status| Redis

    Worker -->|Send OTP| SMSProviderPrimary
    Worker -.->|Fallback| SMSProviderSecondary

    SMSProviderPrimary -->|Update status| Redis
    SMSProviderSecondary -->|Update status| Redis

    %% Apply styles
    class Client,APIGateway primary
    class OTPService,Worker secondary
    class Redis,Database,Queue storage
    class SMSProviderPrimary,SMSProviderSecondary external

    %% Add flow direction
    linkStyle default stroke:#666,stroke-width:2px,fill:none