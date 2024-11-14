make:
	@echo "Usage:"
	@echo "  make test"
	@echo "  make build"
	@echo "  make collect"

test:
	@echo "Running tests..."
	@go test -v ./...

build:
	@echo "Building..."
	@mkdir -p bin
	@go build -o bin/ant_watcher cmd/ant_watcher/main.go

collect:
	@output_file="output/ant_watcher.txt"; \
	echo -n > "$$output_file"; \
	find . -type f ! -path "./output/*" ! -path "./docs/*" ! -path "./script/*" ! -path "./bin/*" -exec grep -Iq . {} \; -and -exec cat {} + >> "$$output_file"

.PHONY: test build collect
