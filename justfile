tool_dir := "bin"
staticcheck_version := "v0.7.0"
govulncheck_version := "v1.1.4"

run-client:
	go run ./cmd/client

run-server:
	go run ./cmd/server

build: build-client build-server

build-client:
	mkdir -p {{tool_dir}}
	go build -o ./{{tool_dir}}/les-perissables-client ./cmd/client

build-server:
	mkdir -p {{tool_dir}}
	go build -o ./{{tool_dir}}/les-perissables-server ./cmd/server

clean:
	rm -rf ./{{tool_dir}}

test:
	go test ./...

lint: install-tools
	go vet ./...
	./{{tool_dir}}/staticcheck ./...

security-scan: install-tools
	./{{tool_dir}}/govulncheck ./...

install-tools:
	mkdir -p {{tool_dir}} && GOBIN="$(pwd)/{{tool_dir}}" go install honnef.co/go/tools/cmd/staticcheck@{{staticcheck_version}} && GOBIN="$(pwd)/{{tool_dir}}" go install golang.org/x/vuln/cmd/govulncheck@{{govulncheck_version}}
