# In backend/Dockerfile

# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files from the backend subdirectory
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# ---- Stage 2: Final ----
FROM alpine:latest

RUN apk update && apk add --no-cache tzdata

WORKDIR /app

# Copy the built binary from the 'builder' stage
COPY --from=builder /app/main .

# Copy the frontend files
COPY frontend/ ./frontend

EXPOSE 8080

CMD ["./main"]