# Growlink frontend (SvelteKit SPA)

## Development

```bash
bun install
bun run dev
```

- Watches Tailwind → `static/css/output.css` (login page)
- Watches Vite/SvelteKit → `static/app/` (authenticated app)

Run alongside the Go server (`docker compose up` starts `web` + `app`).

## Production build

```bash
bun run build
```
