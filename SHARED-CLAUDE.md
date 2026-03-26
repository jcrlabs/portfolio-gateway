# CLAUDE.md — 00 Shared Principles (jcrlabs)

> Base para todos los proyectos de jcrlabs.

## Principio #1: CÓDIGO SIMPLE

La solución más simple que funcione correctamente es siempre la correcta.

- `net/http` nativo > Gin (si no necesitas router avanzado)
- In-memory > Redis (si los datos caben en memoria)
- `slog` stdlib > librerías de logging externas

## Principio #2: Código mínimo

```
back/
├── cmd/main.go
├── internal/
└── Dockerfile
```

## Dominios

| Proyecto | Dominio prod | Dominio test |
|----------|-------------|---------------|
| Landing | jcrlabs.net | test.jcrlabs.net |
| Gateway | api.jcrlabs.net | api-test.jcrlabs.net |
| Inventory | invent-back.jcrlabs.net | invent-back-test.jcrlabs.net |
| Blog | blog.jcrlabs.net | blog-test.jcrlabs.net |
| Chat | chat.jcrlabs.net | chat-test.jcrlabs.net |

## Versiones

| Tech | Versión |
|------|---------|
| Go | 1.24 |
| Next.js | 16.2 |
| React | 19 |
| Node.js | 22 LTS |
| PostgreSQL | 17 |
| Tailwind CSS | 4 |

## K8s Deploy (todos los proyectos)

```
Namespace prod:  portfolio
Namespace test:  portfolio-test
Docker:          multi-stage → distroless (Go) / node:22-alpine (Node)
Ingress:         nginx + cert-manager TLS
HPA:             min 1-2, max 4-5
```

## Seguridad mínima

- Security headers: HSTS, X-Frame-Options, X-Content-Type-Options
- CORS: lista explícita de origins
- Rate limiting por IP
- Secrets: Kubernetes Sealed Secrets

## Health check

Todos los proyectos exponen:
```
GET /api/health → { status: "ok", version: "x.y.z", uptime: 12345 }
```
