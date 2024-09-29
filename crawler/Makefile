.PHONY: build run backup restore

build:
	go build -o out/crawler

run:
	./out/crawler

backup:
	pg_dump -h localhost -U myuser -d mydatabase -F c -b -v -f mydatabase.backup

restore:
	pg_restore -h localhost -U myuser -d mydatabase -v mydatabase.backup
