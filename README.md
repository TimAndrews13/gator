# Gator 🐊

A command-line RSS feed aggregator built in Go. Gator lets you register users, subscribe to RSS feeds, automatically poll and collect posts, and browse the latest content — all from your terminal.

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [PostgreSQL](https://www.postgresql.org/) running locally 

Gator expects a Postgres database reachable at:

```
postgres://postgres:postgres@localhost:5432/gator
```

Make sure you've created the `gator` database before running the app:

```sql
CREATE DATABASE gator;
```

## Installation

```bash
git clone https://github.com/TimAndrews13/gator.git
cd gator
go build -o gator .
```

Or install directly with:

```bash
go install github.com/TimAndrews13/gator@latest
```

## Usage

```bash
./gator <command> [arguments]
```

### User Management

| Command | Arguments | Description |
|---|---|---|
| `register` | `<name>` | Create a new user and set them as current |
| `login` | `<name>` | Switch to an existing user |
| `users` | — | List all registered users |
| `reset` | — | Delete all users (destructive) |

**Example:**
```bash
./gator register alice
./gator login alice
./gator users
```

### Feed Management

These commands require a logged-in user.

| Command | Arguments | Description |
|---|---|---|
| `addfeed` | `<name> <url>` | Add a new RSS feed and follow it |
| `feeds` | — | List all available feeds |
| `follow` | `<url>` | Follow an existing feed by URL |
| `following` | — | List feeds the current user follows |
| `unfollow` | `<url>` | Unfollow a feed by URL |

**Example:**
```bash
./gator addfeed "Hacker News" https://news.ycombinator.com/rss
./gator follow https://news.ycombinator.com/rss
./gator following
./gator unfollow https://news.ycombinator.com/rss
```

### Aggregation

| Command | Arguments | Description |
|---|---|---|
| `agg` | `<duration>` | Start polling feeds at the given interval (e.g. `30s`, `5m`, `1h`) |

The aggregator runs continuously, fetching the next due feed on each tick and saving new posts to the database.

```bash
./gator agg 1m
```

### Browsing Posts

| Command | Arguments | Description |
|---|---|---|
| `browse` | `[limit]` | Show the latest posts for the current user (default: 2) |

```bash
./gator browse
./gator browse 10
```

## Project Structure

```
.
├── main.go                 # Entry point, command registration
├── commands.go             # Command routing and dispatch
├── handler_functions.go    # All command handlers
├── rss_feed.go             # RSS fetching and parsing
├── internal/
│   ├── config/             # Config file read/write (current user)
│   └── database/           # sqlc-generated database queries
└── sql/                    # SQL schema and queries
```

## Dependencies

- [`github.com/lib/pq`](https://github.com/lib/pq) — PostgreSQL driver
- [`github.com/google/uuid`](https://github.com/google/uuid) — UUID generation
- [`sqlc`](https://sqlc.dev/) — Type-safe SQL query generation