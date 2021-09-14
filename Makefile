deps:
	@go get

lint:
	@golangci-lint run -v

vet:
	@go vet ./...

test:
	@go test -coverprofile=cover.out -v ./...

cov:
	@go tool cover -html=cover.out

tidy:
	@go mod tidy

build:
	@go build .

doc:
	@$(MAKE) -C doc all

release:
	@goreleaser $(GORELEASER_ARGS)

snapshot: GORELEASER_ARGS= --rm-dist --snapshot
snapshot: release

.PHONY: deps lint vet test cov tidy build doc tidy release snapshot