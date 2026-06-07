# mini-go-server

A personal project to learn Go by building a minimal HTTP server from scratch — no frameworks, just raw TCP sockets and the standard library.

## What it does

- Listens on port `8080`
- Parses HTTP GET requests manually over raw TCP
- Serves static HTML files from the `static/` directory
- Uses a worker pool (3 goroutines) to handle concurrent connections

## How to run

```bash
go run main.go
```

Then open your browser at `http://localhost:8080`.

### Available routes

| Route | File served |
|-------|-------------|
| `/` | `static/index.html` |
| `/path` | `static/path.html` |
| `/admin` | `static/admin.html` (returns 403) |
| anything else | `static/404.html` |
