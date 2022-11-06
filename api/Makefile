test-cov-total:
	@go tool cover -func cover.out | tail -n 1 | xargs -n1
test-cov:
	@go test ./... -coverprofile fmtcoverage.html fmt
	@go test ./... -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html
	@google-chrome 'cover.html'
test-cov-html:
	@go test -coverprofile=cover.out ./... && go tool cover -html=cover.out
mock:
	@go generate ./...