bin_dir := "bin"

run-client:
	go run ./cmd/client

run-server:
	go run ./cmd/server

build: build-client build-server

build-client:
	mkdir -p {{bin_dir}}
	go build -o ./{{bin_dir}}/les-perissables-client ./cmd/client

build-server:
	mkdir -p {{bin_dir}}
	go build -o ./{{bin_dir}}/les-perissables-server ./cmd/server

clean:
	rm -rf ./{{bin_dir}}

test:
	go test ./...

lint:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

security-scan:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
