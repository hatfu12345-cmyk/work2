# ---------- stage 1: builder ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

# сначала зависимости — слой кэшируется
COPY go.mod ./
RUN go mod download github.com/lib/pq

# потом исходники
COPY . .

# статический бинарник без debug-инфы
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ./main.go

# ---------- stage 2: runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/app /app/app

EXPOSE 8080
ENTRYPOINT ["/app/app"]