# CLAUDE.md — home.jcrlabs.net (Portfolio Landing)

> Extiende: `SHARED-CLAUDE.md` 
> Dominio: `home.jcrlabs.net` | Namespace K8s: `home`

## Qué es esto

Landing principal del portfolio. Muestra todos los proyectos de jcrlabs con estado live, stack y links. Bebe de los demás proyectos vía sus APIs públicas para mostrar datos en tiempo real (ej: últimos posts del blog, uptime del inventory, usuarios activos en el chat).

## Principio rector: SIMPLE

Cada decisión de código debe ir hacia la solución más simple. Si hay dos formas de hacer algo, elegir siempre la más simple. No añadir abstracciones sin necesidad real.

## Stack

- **Framework**: Next.js 16 (App Router, SSG/ISR para el home)
- **Runtime**: Node.js 22 LTS
- **Estilos**: Tailwind CSS v4 — zero config, sin `tailwind.config.js`
- **Fuentes de datos**: fetch estático en build + ISR cada 60s por proyecto
- **Deploy**: K8s `home` namespace, `home.jcrlabs.net`

## Estructura de archivos (mínima)

```
home-front/
├── app/
│   ├── page.tsx          # Home: hero + grid de proyectos
│   ├── layout.tsx        # Root layout con fuente y metadata
│   └── globals.css       # Tailwind @import
├── components/
│   └── ProjectCard.tsx   # Una sola card reutilizable
├── lib/
│   └── projects.ts       # Config estática de todos los proyectos
└── next.config.ts        # Solo lo imprescindible
```

## Proyectos que se muestran (lib/projects.ts)

```typescript
export const PROJECTS = [
  {
    id: "inventory",
    name: "Inventory",
    description: "Sistema de gestión de inventario con MinIO para assets",
    url: "https://inventory.jcrlabs.net",
    apiUrl: "https://inventory.jcrlabs.net/api/health",
    stack: ["Go", "React", "PostgreSQL", "MinIO"],
    status: "live",
    repo: "github.com/jonathanCaamano/inventory-back",
  },
  {
    id: "blog",
    name: "Blog",
    description: "CMS con ingestión RSS automática de fuentes CNCF/K8s",
    url: "https://blog.jcrlabs.net",
    apiUrl: "https://blog.jcrlabs.net/api/posts/latest",
    stack: ["NestJS", "MongoDB", "GraphQL", "Next.js"],
    status: "live",
    repo: "github.com/jonathanCaamano/blog-cms-back",
  },
  {
    id: "dashboard",
    name: "K8s Dashboard",
    description: "Monitorización del cluster en tiempo real",
    url: "https://dashboard.jcrlabs.net",
    apiUrl: "https://dashboard.jcrlabs.net/api/health",
    stack: ["Go", "client-go", "React", "Redis"],
    status: "live",
    repo: "github.com/jonathanCaamano/k8s-dashboard",
  },
  {
    id: "chat",
    name: "Chat",
    description: "Chat en tiempo real con WebSockets y Redis pub/sub",
    url: "https://chat.jcrlabs.net",
    apiUrl: "https://chat.jcrlabs.net/api/health",
    stack: ["Go", "WebSockets", "Redis", "React"],
    status: "live",
    repo: "github.com/jonathanCaamano/chat-back",
  },
  {
    id: "fincontrol",
    name: "FinControl",
    description: "Control de finanzas personales con análisis y alertas",
    url: "https://fincontrol.jcrlabs.net",
    apiUrl: "https://fincontrol.jcrlabs.net/api/health",
    stack: ["NestJS", "PostgreSQL", "React"],
    status: "wip",
    repo: "github.com/jonathanCaamano/fincontrol-back",
  },
]
```

## Data fetching por proyecto

Cada proyecto expone `/api/health` (JSON). El home hace ISR: `revalidate: 60`.

```typescript
// app/page.tsx
export const revalidate = 60

async function getProjectStatus(apiUrl: string) {
  try {
    const res = await fetch(apiUrl, { next: { revalidate: 60 } })
    return res.ok ? "online" : "degraded"
  } catch {
    return "offline"
  }
}
```

## Dominio y K8s deploy

```yaml
# ingress.yaml
host: home.jcrlabs.net

# Namespace: home
# Docker: multi-stage → distroless
# HPA: 1-2 replicas (es un landing estático, no necesita más)
```

## Qué NO hacer

- No Redux, no Zustand — no hay estado global que justifique una librería
- No tRPC — las llamadas son simples fetch a APIs REST
- No Storybook — los componentes son simples, no lo necesitan
- No GraphQL en el home — consume REST de los proyectos
- No tests unitarios de UI — el home es mostly estático
- No `tailwind.config.js` — Tailwind v4 funciona sin él

## Referencia

- Patrones base: `github.com/jonathanCaamano/inventory-front` (mismo stack Next.js)
- Shared K8s patterns: `00-SHARED-PRINCIPLES-CLAUDE.md`
