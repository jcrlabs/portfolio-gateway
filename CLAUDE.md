# CLAUDE.md — portfolio-gateway

> Extiende: `SHARED-CLAUDE.md`
> Dominio prod: `portfolio-api.jcrlabs.net` | test: `portfolio-api-test.jcrlabs.net`
> Namespace K8s: `portfolio` (prod) / `portfolio-test` (test)

## Qué es esto

API Gateway centralizado del portfolio. Routing declarativo vía `routing.yaml`,
middleware chain completa (RequestID, Logger, CORS, RateLimit, Recover, Metrics),
health check agregado y métricas Prometheus.

## Stack

- Go 1.24 · net/http stdlib · httputil.ReverseProxy · prometheus/client_golang
- Config: `routing.yaml` + godotenv

## Estructura

```
cmd/main.go                  # entrypoint — wiring + http.ListenAndServe
internal/
├── config/config.go         # env vars + LoadRoutes(routing.yaml)
├── gateway/
│   ├── proxy.go             # NewMux — registra rutas con strip prefix
│   ├── middleware.go        # RequestID, Logger, CORS, RateLimit, Recover, Metrics
│   └── health.go            # GET /healthz — consulta /api/health de cada backend
└── metrics/
    └── metrics.go           # Prometheus histogram + counter
routing.yaml                 # rutas declarativas
deploy/helm/                 # values-test.yaml (portfolio-api-test) + values-prod.yaml (portfolio-api)
.github/workflows/           # ci.yml (lint+test+build) | cd.yml (deploy)
```

## Dominios

| Entorno | URL |
|---------|-----|
| Prod    | `portfolio-api.jcrlabs.net` |
| Test    | `portfolio-api-test.jcrlabs.net` |

## CORS

Origins permitidos: `home.jcrlabs.net`, `home-test.jcrlabs.net`, `localhost:3000`

## CI local

Ejecutar **antes de cada commit** para evitar que lleguen errores a GitHub Actions:

```bash
gofmt -l .                      # no debe mostrar ficheros
go vet ./...
golangci-lint run --timeout=5m
go test -race ./...
CGO_ENABLED=0 go build -o /dev/null ./cmd/main.go
```

## Git

- Ramas: `feature/`, `bugfix/`, `hotfix/`, `release/` — sin prefijos adicionales
- Commits: convencional (`feat:`, `fix:`, `chore:`, etc.) — sin mencionar herramientas externas ni agentes en el mensaje
- PRs: título y descripción propios del cambio — sin mencionar herramientas externas ni agentes
- Comentarios y documentación: redactar en primera persona del equipo — sin atribuir autoría a herramientas

## Qué NO hacer

- No Gin/Echo — net/http es suficiente para un proxy stateless
- No base de datos — el gateway es stateless
- No autenticación — cada backend maneja su propia auth
