# CLAUDE.md — inventory-back

> Extiende: `SHARED-CLAUDE.md`
> Dominio: `invent-back.jcrlabs.net` | Namespace K8s: `taller-inventario`

## Qué es esto

Backend del sistema de inventario. JWT auth (HS256), RBAC (Admin/Manager/Viewer),
CRUD productos con upload de imágenes a MinIO, dashboard de estadísticas.

## Stack

- **Language**: Go 1.24
- **Framework**: Gin
- **ORM**: GORM + PostgreSQL 17
- **Storage**: MinIO (S3-compatible) con presigned URLs
- **Auth**: JWT HS256 — access token 15 min (Bearer), refresh token 7 d (httpOnly cookie)
- **Deploy**: Namespace `taller-inventario` (prod) / `taller-inventario-test` (test)

## Estructura

```
cmd/main.go              # entrypoint — config + wiring + router
internal/
├── config/config.go     # env vars (godotenv)
├── db/db.go             # GORM connect + AutoMigrate
├── domain/models.go     # User, Product, Category, RefreshToken
├── repository/          # GORM repositories (user, product, category)
├── service/             # auth, product, minio, category
├── handler/             # HTTP handlers (auth, product, category, health)
└── middleware/          # JWT, RBAC, RateLimit, CORS, Logger, RequestID
migrations/001_init.sql  # referencia SQL (AutoMigrate lo gestiona)
deploy/helm/             # Chart + values-test.yaml + values-prod.yaml
.github/workflows/       # CI/CD pipeline
```

## RBAC

| Acción | Admin | Manager | Viewer |
|--------|-------|---------|--------|
| GET productos/categorías | ✅ | ✅ | ✅ |
| POST/PUT productos/categorías | ✅ | ✅ | ❌ |
| DELETE | ✅ | ❌ | ❌ |
| Upload imagen | ✅ | ✅ | ❌ |
| Dashboard stats | ✅ | ✅ | ✅ |

## Flujo de request

```
Request → CORS → RequestID → Logger → RateLimit (auth) → JWT → RoleCheck → Handler
```

## Qué NO hacer

- No GraphQL — REST es suficiente
- No Redis para rate limiting — in-memory suficiente para el tráfico esperado
- No soft deletes — delete es delete
- No caché — inventario necesita datos en tiempo real
- No gin-contrib/cors — CORS propio para control exacto de headers de seguridad
