# Synthify

Knowledge graph extraction and exploration workspace.

## Local Development

Start the local development stack with Docker Compose:

```bash
docker compose up --build
```

This starts:

- `frontend` on `http://localhost:5173`
- `backend` on `http://localhost:8080`
- `gcs` emulator on `http://localhost:4443`
- `firebase-auth` emulator on `http://localhost:9099`

Notes:

- If `frontend/` or `backend/` is not implemented yet, the corresponding container waits instead of failing.
- BigQuery and graph queries run in mock mode in the initial compose setup.
