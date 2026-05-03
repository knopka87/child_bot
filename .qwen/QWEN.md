# Project Context for Qwen Code

## Project Overview
- **Name:** Child Bot
- **Type:** Full-stack application (Go backend + React frontend)
- **Frontend:** React 18 + TypeScript, VKUI, TailwindCSS, Vite, Zustand
- **Backend:** Go, REST API, PostgreSQL, Redis
- **Source path:** `frontend/sourse/src` (note: typo in directory name)

## Key Commands
```bash
# Frontend
cd frontend && npm run dev      # Start dev server
cd frontend && npm run build    # Build for production
cd frontend && npm run lint     # Run linter
cd frontend && npm run typecheck # TypeScript check

# Backend
make build                      # Build Go server
make run                        # Run backend locally
make test                       # Run all tests

# Full stack
make dev                        # Start all services (Docker)
make dev-down                   # Stop all services
```

## Architecture Notes
- Frontend uses VKUI component library
- State management via Zustand
- Routing: react-router-dom v6
- Styling: TailwindCSS + clsx
- API: axios for HTTP requests

## Figma Integration Workflow
- **Figma export:** `frontend/source` — полный React-проект из Figma (референс для сравнения)
- **Actual source code:** `frontend/src` — здесь ведём разработку
- **Goal:** Привести `frontend/src` к виду как в `frontend/source`

### Figma comparison workflow:
1. Reference implementation: `frontend/source/src/`
2. Current codebase: `frontend/src/`
3. Compare component-by-component and update current code to match the Figma design
4. Pay attention to: styles, layout, colors, typography, spacing, component structure

## Conventions
- TypeScript strict mode
- ESLint + Prettier formatting
- Component-based architecture in `frontend/sourse/src/`
- API handlers in Go backend under `api/`
