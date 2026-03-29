# Synthify

Graph extraction and exploration workspace for structured document analysis.

## Local development

Start the local stack with Docker Compose:

```bash
docker compose up --build
```

Default host endpoints:

- frontend: `http://localhost:5173`
- backend: `http://localhost:8080`
- gcs emulator: `http://localhost:4443`
- firebase auth emulator: `http://localhost:9099`

If those host ports are already in use, override them before startup:

```bash
BACKEND_PORT=18080 VITE_API_BASE_URL=http://localhost:18080 docker compose up --build
```

Notes:

- If `frontend/` or `backend/` is still missing, the corresponding container waits instead of crashing.
- BigQuery and graph queries run in mock mode in the initial local stack.
- The backend serves `GET /healthz` for a basic liveness check.
