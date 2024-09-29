.PHONY: build run backup restore

# development
dev-crawl-repo: build
	./out/crawler crawl --repo

dev-crawl-readme: build
	./out/crawler crawl --readme

# production
up: build run

build:
	go build -o out/crawler

run:
	./out/crawler

backup:
	pg_dump -h localhost -U myuser -d mydatabase -F c -b -v -f mydatabase.backup

restore:
	pg_restore -h localhost -U myuser -d mydatabase -v mydatabase.backup
