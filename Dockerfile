# Builder stage
FROM golang:1.20 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download


COPY src/ .



# Debug build
FROM build AS debug
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o minecraft-player-stats-debug .

# Production build
FROM build AS production
RUN CGO_ENABLED=0 GOOS=linux go build -a -o minecraft-player-stats .

#######################
# Debug stage

FROM golang:latest AS debug-runtime

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY /go/bin/dlv /app/dlv
COPY --from=debug /app/minecraft-player-stats-debug /app/minecraft-player-stats
COPY --from=build . /app/src



EXPOSE 8080 2345

ENTRYPOINT ["/app/dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/minecraft-player-stats"]

#######################
# Runtime stage

FROM gcr.io/distroless/base AS production-runtime

WORKDIR /app

COPY --from=production /app/minecraft-player-stats .

EXPOSE 8080

ENTRYPOINT ["./minecraft-player-stats"]
