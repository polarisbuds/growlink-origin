FROM oven/bun:1.2.18

WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock* ./

RUN bun install

WORKDIR /app

CMD ["sh", "-c", "cd frontend && bun install && bun run dev"]
