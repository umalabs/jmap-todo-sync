# JMAP Todo Sync Example

This is a simple example project demonstrating data synchronization using a simplified JSON Meta Application Protocol (JMAP) approach.

## Features

*   **Backend:** Go server with SQLite database, handling JMAP-like requests for Todo items.
*   **Frontend:** React client interacting with the backend using JMAP-style JSON requests.
*   **Functionality:** Create, read, update, and delete Todo items.

## Setup and Run

**Backend (Go):**

1.  Navigate to the `backend` directory: `cd backend`
2.  Initialize Go modules: `go mod tidy`
3.  Run the backend server: `go run main.go`
    *   The server will start on `http://localhost:8080`.

**Frontend (React):**

1.  Navigate to the `frontend` directory: `cd frontend`
2.  Install dependencies: `npm install` or `yarn install`
3.  Start the frontend development server: `npm start` or `yarn start`
    *   The frontend will be accessible at `http://localhost:3000`.

## Accessing the App

Open your browser and go to `http://localhost:3000`. You should see the Todo application.

## Important Notes

*   **Simplified JMAP:** This is a highly simplified example and not a fully compliant JMAP implementation.
*   **No Authentication:** Authentication and authorization are not implemented for simplicity.
*   **Polling:** The client uses simple polling to check for updates, not efficient push notifications or long-polling.
*   **Minimal Error Handling:** Error handling is basic for demonstration purposes.
*   **Example Todo Capability:**  The `urn:example:params:jmap:todo` capability is used as an example. In a real JMAP environment, you would use registered or well-defined capabilities.

This project is for educational purposes to illustrate the basic principles of JMAP-like data synchronization. For production systems, a full and robust JMAP implementation would be necessary.