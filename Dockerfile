# syntax=docker/dockerfile:1

# --- Etape de compilation ---
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Telechargement des dependances (cache optimise)
COPY go.mod go.sum ./
RUN go mod download

# Compilation du binaire
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/forum .

# --- Etape d'execution ---
FROM alpine:3.20

WORKDIR /app

# Certificats racine pour les appels HTTPS (API TMDB)
RUN apk add --no-cache ca-certificates

# Binaire et ressources servies au runtime
COPY --from=builder /app/forum ./forum
COPY views ./views
COPY static ./static

EXPOSE 8080

CMD ["./forum"]
