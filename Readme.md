# Go Minecraft Stat Stream

This is a simple Go application that fetches player statistics from a Minecraft server and stores them in a Redis database. It offers an API to query player statistics.

## Features

- Fetches player statistics from a Minecraft server using RCON or JSON 
- Stores player statistics in a Redis database.
- Offers an API to query player statistics with optional grouping, sorting, filtering, and limiting results.
- Allows users to query player statistics based on player name and/or specific stat type.
- Caches player UUID to username mappings using Mojang API
- Supports querying summarized stats, current players, and stat types

## Usage

### API Endpoints
1. `GET /getstats`
   - Returns the latest fetched player stats
2. `GET /playerstats?playerName={playerName}`
   - Returns stats for the specified player by their username
3. `GET /currentplayers`
   - Returns a list of the current player names loaded in Redis
4. `GET /summarizedstats?statType={statType}`
   - Returns summarized stats for the specified stat type (e.g., kills, deaths, etc.)
   - Also returns the specific stat types that make up the summarized stats
   - ` curl "https://localhost:8080/summarizedstats?statType=planks"`   
5. `GET /allstattypes`
   - Returns a list of all available stat types (unique) stored in Redis
#### Get Player Stats

- **Endpoint:** `/playerstats`
- **Method:** `GET`
- **Query Parameters:**
  - `playername` (optional): The name of the player whose stats are to be retrieved.
  - `playernames` (optional): return and compare multiple players
  - `stattype` (optional): The specific stat type to retrieve. If not provided, all stats for the player are retrieved. Use dashes (-) instead of colons (:) in the stat type.
  - `groupby` (optional): Group results by specific criteria. Currently supported: `stattype`.
  - `sort` (optional): Sort results by value. Supported values: `asc` (ascending) and `desc` (descending).
  - `top` (optional): Limit the number of results to the top N items.
  - `category` (optional): Filter stats by a specific category.
- **Example Usage:**
  - Get all stats for player `pvpNJ`: `http://localhost:8080/playerstats?playername=pvpNJ`
  - Get specific stat `minecraft:mined:minecraft:chest` for player `pvpNJ`: `http://localhost:8080/playerstats?playername=pvpNJ&stattype=minecraft-mined-minecraft-chest`
  - Get all stats for player `pvpNJ`, grouped by `stattype` and sorted in descending order: `http://localhost:8080/playerstats?playername=pvpNJ&groupby=stattype&sort=desc`
  - Get the top 10 mined items for player `pvpNJ`, sorted in descending order: `http://localhost:8080/playerstats?playername=pvpNJ&sort=desc&top=10&category=mined`
  - Get stats for specific players: `http://localhost:8080/playerstats?playernames=Player1,Player2&category=mined&top=5`
  - Get stats for top ten most crafted items: `http://localhost:8080/playerstats?category=minecraft:crafted&top=10&sort=desc`


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
- `playerstats.go`: The endpoint for querying player stats with grouping, sorting, filtering, and limiting functionality.
## License

This project is licensed under the MIT License.
