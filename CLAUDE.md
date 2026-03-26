# CLAUDE.md — portfolio-gateway (API Gateway)

> Extiende: `SHARED-CLAUDE.md`
> Dominio: `api.jcrlabs.net` | Namespace K8s: `portfolio`

## Qué es esto

API Gateway centralizado del portfolio. Routing declarativo vía YAML,
middleware chain completa, health check agregado y métricas Prometheus.

## Stack

- **Language**: Go 1.24
- **HTTP**: net/http stdlib + httputil.ReverseProxy
- **Rate limit**: golang.org/x/time/rate (100 rps/IP)
- **Métricas**: prometheus/client_golang (histograma latencia + counter requests)
- **Logging**: slog (stdlib JSON)
- **Config**: routing.yaml + godotenv

## Estructura

```
cmd/main.go                  # entrypoint — wiring + http.ListenAndServe
internal/
├── config/config.go         # env vars + LoadRoutes(routing.yaml)
├── gateway/
│   ├── proxy.go             # NewMux — registra rutas con strip prefix
│   ├── middleware.go        # RequestID, Logger, CORS, RateLimit, Recover, Metrics
│   └── health.go            # /healthz — consulta /api/health de cada backend
└── metrics/
    └── metrics.go           # Prometheus histogram + counter + promhttp.Handler
routing.yaml                 # config declarativa de rutas
deploy/helm/                 # Chart + values-test.yaml + values-prod.yaml
.github/workflows/ci.yml     # CI/CD pipeline
```

## Routing declarativo (routing.yaml)

```yaml
routes:
  - name: inventory
    prefix: /api/inventory
    backend: http://inventory-back.taller-inventario.svc.cluster.local:8080
```

El gateway hace strip del prefix antes de hacer forward.

## Middleware chain

```
Request → RequestID → Logger(slog) → CORS → RateLimit(100rps/IP)
        → Recover → Metrics(histogram) → ReverseProxy
```

## Endpoints propios

- `GET /healthz` — estado agregado de todos los backends registrados
- `GET /metrics` — Prometheus scrape endpoint

## Qué NO hacer

- No Gin/Echo — net/http es suficiente para un proxy stateless
- No base de datos — el gateway es stateless por diseño
- No autenticación — cada backend maneja su propia auth
- No circuit breaker — el portfolio no requiere esa complejidad ahora
