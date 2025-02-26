APP_NAME=agentos

build:
	go build -o $(APP_NAME)

run:
	go run cmd/main.go $(ARGS)

clean:
	rm -f $(APP_NAME)
