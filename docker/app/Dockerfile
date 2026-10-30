FROM golang:1.27
LABEL MAINTAINER="ryota.kashio"

RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl bash unzip \
    && rm -rf /var/lib/apt/lists/*

# air をインストール
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Go モジュールをビルド時にキャッシュ（ボリュームマウント対象の /app 外に保存されるため上書きされない）
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

EXPOSE 1323

# ソースはボリュームマウントで提供、air でホットリロード起動
CMD ["air", "-c", ".air.toml"]
