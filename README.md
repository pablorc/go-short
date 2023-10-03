# go-short

A simple url shortener service (mostly to practice some Golang)

# Usage

1. `docker compose up` for Redis
1. `go run .` to run the server

Save new URLs with `localhost:3333/save?url=<your-url>`. It gives back the new URL.

Access the URL to get a redirection to the original one 

