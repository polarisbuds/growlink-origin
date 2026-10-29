FROM oven/bun:1.2.18

WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock* ./

RUN mkdir -p /app/static/vendor /app/static/app

RUN bun install

WORKDIR /app

ENV CHOKIDAR_USEPOLLING=true
ENV CHOKIDAR_INTERVAL=1000

CMD ["sh", "-c", "cd frontend && bun install && bun run dev:docker"]
