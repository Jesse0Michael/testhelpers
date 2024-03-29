#################################################################################
# TEST COMMANDS
#################################################################################
test:
	go test -cover ./... 
	golangci-lint run ./...

test-coverage:
	go test -coverpkg ./... -coverprofile coverage.out ./... && go tool cover -html=coverage.out
