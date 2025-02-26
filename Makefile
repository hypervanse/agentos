APP_NAME ?= agentos

build:
	go build -o $(APP_NAME) ./cmd/main.go

run: build
	sudo ./$(APP_NAME) $(ARGS)

clean:
	rm -f $(APP_NAME)
