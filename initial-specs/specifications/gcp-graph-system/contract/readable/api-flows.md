# API Flows and Interactions

This document visualizes the primary interactions between the Frontend, Backend, and external services like Firebase Auth and Stripe.

## 1. Authentication and User Sync

When a user logs in for the first time or returns, the following flow ensures their profile is synchronized.

```mermaid
sequenceDiagram
    participant F as Frontend (React)
    participant A as Firebase Auth
    participant B as Backend (UserService)
    participant DB as BigQuery

    F->>A: Login (Google OAuth)
    A-->>F: ID Token (JWT)
    F->>B: SyncUser (Header: Authorization: Bearer <Token>)
    B->>B: Verify JWT
    B->>DB: Check/Create User Record
    DB-->>B: User Data
    B-->>F: SyncUserResponse (User info, is_new_user)
```

## 2. Workspace Management & Billing (Stripe)

Upgrading a workspace to the 'Pro' plan involves a redirection to Stripe.

```mermaid
sequenceDiagram
    participant F as Frontend
    participant B as Backend (BillingService)
    participant S as Stripe
    participant DB as BigQuery (workspaces)

    F->>B: CreateCheckoutSession(workspace_id)
    B->>S: Create Session API
    S-->>B: Checkout URL
    B-->>F: Returns URL
    F->>S: User completes payment on Stripe
    S-->>B: Webhook (checkout.session.completed)
    B->>B: Verify Signature
    B->>DB: Update workspace plan to 'pro'
```

## 3. Interactive Graph Exploration

The core value of the system is the interactive traversal of the knowledge graph.

```mermaid
sequenceDiagram
    participant F as Frontend (React Flow)
    participant B as Backend (GraphService)
    participant S as Spanner Graph

    F->>B: GetGraph(workspace_id, document_id)
    B->>S: Query initial graph for document
    S-->>B: Nodes & Edges
    B-->>F: GetGraphResponse (Initial View)
    
    Note over F, S: User clicks a node to expand
    
    F->>B: ExpandNeighbors(node_id, depth, edge_types)
    B->>S: Traversal query (GQL)
    S-->>B: Subgraph (Adjacent nodes/edges)
    B-->>F: ExpandNeighborsResponse (Merged into View)
```
