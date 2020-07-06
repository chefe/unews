run:
	go run .

docker:
	docker build -t unews .
	docker run --rm --name unews -p 8080:8080 unews:latest
