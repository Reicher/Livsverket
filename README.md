# Livsverket

A simple split-stack project with a Go backend and a lightweight, React-inspired frontend tailored for a portrait-oriented mobile experience.

> **Note**
> The execution environment for this challenge does not provide access to the public npm registry. To keep the workflow offline-friendly the frontend ships with a tiny `mini-react` helper that mimics the React hook API used in the UI. No third-party packages are required to run the interface.

## Backend

```bash
go run main.go
```

The backend exposes `GET /api/status` returning a JSON payload confirming the server is available. The service listens on port `8080` by default and honours the `PORT` and `ALLOWED_ORIGIN` environment variables.

## Frontend

```bash
cd frontend
npm install  # optional, no dependencies are required
npm run dev
```

The frontend runs an unbundled development server on port `5173`. It renders a mobile-first layout with a top title bar, a dynamic content area, and a bottom navigation bar. Selecting any bottom navigation button swaps the central content panel. On load, the app attempts to reach the Go backend and shows the connection status beneath the active page text.

## Connecting the stacks

When the frontend loads it calls `http://localhost:8080/api/status`. You can run both services concurrently to validate the integration:

```bash
# Terminal 1
cd frontend
npm run dev

# Terminal 2
go run main.go
```
