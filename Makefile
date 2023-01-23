test:
	@echo "Testing Mongo Schema Go"
	go test -coverprofile=coverage.out ./...

clean:
	@echo "Clearing coverage"
	rm -r ./coverage.out