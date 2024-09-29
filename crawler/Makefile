.PHONY: build run backup restore

# Define the source directory
SRC_DIR := src

# Define the output binary name
BINARY := out/crawler

# development
dev-crawl-repo: build
	./$(BINARY) crawl --repo

dev-crawl-readme: build
	./$(BINARY) crawl --readme

# production
up: build run

build:
	go build -o $(BINARY) $(SRC_DIR)/*.go

run:
	./$(BINARY) crawl --repo

# Clean up the binary
clean:
	rm -f $(BINARY)

backup:
	pg_dump -h localhost -U myuser -d mydatabase -F c -b -v -f mydatabase.backup

restore:
	pg_restore -h localhost -U myuser -d mydatabase -v mydatabase.backup
