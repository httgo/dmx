
bench:
	@go test --bench .

test:
	@go test ./...

cover:
	@go test -coverprofile=coverage.out

cover-html: cover
	@go tool cover -html=coverage.out

.PHONY: test
