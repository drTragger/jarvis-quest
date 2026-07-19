# jarvis-quest

A JARVIS-themed interactive puzzle quest. Vue 3 + Vite frontend, Go backend,
single binary serving both the API and the built SPA.

## Development

```bash
cd backend && go run .        # API on :8080
cd frontend && npm run dev    # dev server, proxies /api to :8080
```

## Building

```bash
make build   # installs frontend deps, builds the SPA into backend/dist, builds the Go binary
make run     # build + run locally
```

## Deployment

Runs as a systemd service (`jarvis-quest.service`) behind a Cloudflare Tunnel.
Pushing to `main` triggers `.github/workflows/deploy.yml` on a self-hosted
runner, which builds and restarts the service.

```bash
make deploy  # pull + build + restart + status (manual deploy)
make logs    # tail service logs
```
