# Go Minecraft Stat Stream

This is a simple Go application that fetches player statistics from a Minecraft server and stores them in a Redis database. It offers an API to query player statistics.

## Features

- Fetches player statistics from a Minecraft server using RCON or JSON.
- Stores player statistics in a Redis database.
- Offers an API to query player statistics with optional grouping and sorting.
- Allows users to query player statistics based on player name and/or specific stat type.

## Usage

### API Endpoints

#### Get Player Stats

- **Endpoint:** `/playerstats`
- **Method:** `GET`
- **Query Parameters:**
  - `playername` (required): The name of the player whose stats are to be retrieved.
  - `stattype` (optional): The specific stat type to retrieve. If not provided, all stats for the player are retrieved.
  - `groupby` (optional): Group results by specific criteria. Currently supported: `stattype`.
  - `sort` (optional): Sort results by value. Supported values: `asc` (ascending) and `desc` (descending).
- **Example Usage:**
  - Get all stats for player `pvpNJ`: `http://localhost:8080/playerstats?playername=pvpNJ`
  - Get specific stat `minecraft:mined:minecraft:chest` for player `pvpNJ`: `http://localhost:8080/playerstats?playername=pvpNJ&stattype=minecraft:mined:minecraft:chest`
  - Get all stats for player `pvpNJ`, grouped by `stattype` and sorted in descending order: `http://localhost:8080/playerstats?playername=pvpNJ&groupby=stattype&sort=desc`

## Configuration

Configure the application using environment variables in the `docker-compose.yml` file:

- `POLLING_INTERVAL`: The interval at which the application polls for new stats (e.g., `5m` for 5 minutes).
- `STATS_SOURCE`: The source of stats, either `rcon` for RCON or `json` for JSON.
- `JSON_STATS_DIRECTORY`: The directory where JSON stat files are located (used only if `STATS_SOURCE` is `json`).
- `REDIS_ADDR`: The address of the Redis server.
- `REDIS_PASSWORD`: The password for the Redis server (if applicable).
- `REDIS_DB`: The Redis database index.
- `SERVER_PORT`: (optional) 

## Running the Application

1. Build the Docker image: `docker-compose build`
2. Run the application: `docker-compose up`

## Development

The application is written in Go and consists of the following components:

- `main.go`: Entry point of the application.
- `polling.go`: Handles the periodic polling of player stats.
- `stats.go`: Fetches player stats and stores them in Redis.
- `mojang.go`: Converts player UUIDs to player names.
- `handler.go`: HTTP handlers for the API endpoints.
- `playerstats.go`: The endpoint for querying player stats with grouping and sorting functionality.

## License

This project is licensed under the MIT License.
