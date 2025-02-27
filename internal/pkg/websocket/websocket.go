package websocket

/**
 * Architecture Overview:
 *
 * ┌────────────────────┐     ┌────────────────────┐     ┌────────────────────┐
 * │                    │     │                    │     │                    │
 * │   Server           │     │   Connection       │     │   Handler          │
 * │                    │     │                    │     │                    │
 * │  - Manages TCP     │     │  - Manages WS      │     │  - Business Logic  │
 * │    listener        │◄───►│    protocol        │◄───►│  - Message         │
 * │  - Accepts         │     │  - Reads/writes    │     │    processing      │
 * │    connections     │     │    frames          │     │  - State           │
 * │  - Event polling   │     │  - Maintains       │     │    management      │
 * │                    │     │    connection      │     │                    │
 * └────────────────────┘     └────────────────────┘     └────────────────────┘
 *
 * Events flow:
 * 1. Server accepts TCP connection
 * 2. Connection handles WebSocket protocol details
 * 3. Handler receives parsed messages and implements business logic
 * 4. Handler sends messages back through Connection 
 */
