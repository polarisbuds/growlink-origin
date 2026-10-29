---
description: Growlink frontend (SvelteKit SPA + Tailwind for login CSS).
globs: "*.ts, *.tsx, *.html, *.css, *.js, *.jsx, *.svelte, package.json"
alwaysApply: false
---

Use **Bun** for package scripts (`bun install`, `bun run dev`, `bun run build`).

- **Svelte 5** + **SvelteKit 2** SPA (`ssr = false`, `adapter-static`) builds to `../static/app`.
- **Vite** is used via SvelteKit (`vite build` / `vite build --watch` in `dev`).
- **Tailwind 3.4** + **DaisyUI 4.12** — app styles via `src/app.css`; login page uses `build:css` → `static/css/output.css`.
- Docker: `web` service runs `bun run dev:docker` (CSS polling via `CHOKIDAR_USEPOLLING` + `dev:kit:poll`). Host `bun run dev` keeps `fs.watch`.
