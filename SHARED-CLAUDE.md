# CLAUDE.md — 00 Shared Principles (jcrlabs)

> Este archivo es la base para todos los proyectos de jcrlabs. Todos los CLAUDE.md de proyecto extienden estas reglas.

## Principio #1: CÓDIGO SIMPLE

**La solución más simple que funcione correctamente es siempre la correcta.**

Antes de añadir cualquier abstracción, librería, o patrón, preguntarse:
- ¿Resuelve un problema real que tengo ahora mismo?
- ¿Podría hacer esto con lo que ya tengo?

Si la respuesta a la primera es no, o a la segunda es sí → no añadir.

Ejemplos concretos:
- `fetch` nativo > axios (sin motivo real para axios)
- Una función `extractTags(item)` > una librería de NLP
- Un `map` en memoria > Redis (si los datos caben en memoria)
- Un `@Cron` de NestJS > una cola de jobs (si no hay retry complejo)
- Un `interface` de TypeScript > Zod (si no hay validación en runtime)

## Principio #2: Código mínimo

Cada proyecto tiene la estructura mínima necesaria. Nada más.

```
back/
├── cmd/main.go (Go) o src/main.ts (NestJS)
├── internal/ o src/modules/
└── Dockerfile
```

## Dominios del portfolio

| Proyecto | Dominio | Stack |
|----------|---------|-------|
| Home (landing) | home.jcrlabs.net | Next.js 16 |
| Inventory | invent-back.jcrlabs.net | Go + React + PostgreSQL + MinIO |
| Blog CMS | blog.jcrlabs.net | NestJS + MongoDB + Next.js |
| K8s Dashboard | dashboard.jcrlabs.net | Go + client-go + React |
| Chat | chat.jcrlabs.net | Go + WebSockets + Redis |
| FinControl | fincontrol.jcrlabs.net | NestJS + PostgreSQL + React |

## Versiones

| Tech | Versión |
|------|---------|
| Go | 1.24 |
| NestJS | 11 |
| Next.js | 16 |
| React | 19 |
| Node.js | 22 LTS |
| PostgreSQL | 17 |
| MongoDB | 8 |
| Redis | 7.4 |
| Tailwind CSS | 4 |

## K8s Deploy standards (todos los proyectos)

```yaml
k8s/
├── deployment.yaml
├── service.yaml
├── ingress.yaml
└── hpa.yaml

# Docker multi-stage → distroless
# GitHub Actions CI/CD
# Namespace = nombre del proyecto
# Ingress Nginx con TLS (cert-manager)
```

## Seguridad (mínimo obligatorio en todos)

- JWT HS256: access token 15min (Bearer header), refresh token 7d (httpOnly cookie)
- bcrypt cost 12 para passwords
- Rate limiting en auth: 10 req/15min/IP
- Security headers: HSTS, X-Frame-Options, X-Content-Type-Options
- CORS: lista explícita de origins, nunca `*`
- Secrets: Kubernetes Sealed Secrets

## API health check (todos los proyectos)

```
GET /api/health → { status: "ok", version: "1.0.0", uptime: 12345 }
```

## Referencia base

- Primer proyecto del portfolio: `github.com/jonathanCaamano/inventory-back`
- Todos los proyectos siguen los mismos patrones que el inventory
