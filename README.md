# GymmTracker

GymmTracker is a script for fetching GymmBoxx occupancy status and then storing it in a PostgreSQL database for further analysis. It comes with a cron job scheduler that, by default, runs at 3-hour intervals from 8AM to 8PM.

### Prerequisites

- [Go 1.18.3+](https://go.dev/)
- [PostgreSQL 14.4+](https://www.postgresql.org/)

### Usage

1. Edit the database URL field in [.env.example](.env.example) and rename the file to .env
2. Build the executable with `go build ./cmd/gymmtracker.go`
3. Run the executable with `./gymmtracker.exe`

### Command flags
- `stdout` - Enable/disable logging to console. Default: `true`
- `logfile` - Set file to output log to. Default: `logs/client.log`