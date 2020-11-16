test:
	go test -short

integration-test:
	./scripts/integration-test $(authKey)

.PHONY: test
