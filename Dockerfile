# Builder stage
FROM golang:1.18 AS build

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY src/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minecraft-player-stats .

#######################
# Runtime stage

FROM gcr.io/distroless/base

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/minecraft-player-stats .

EXPOSE 8080

ENTRYPOINT ["./minecraft-player-stats"]
