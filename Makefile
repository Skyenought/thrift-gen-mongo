build:
	go build .
	mv thrift-gen-mongo $(GOPATH)/bin/thrift-gen-mongo
test:
	cd examples && make test
format:
	gofumpt -l -w .
	find . -type f \( -name "*.toml" -o -name "*.go" \) -exec sed -i 's/$/\r/' {} +


