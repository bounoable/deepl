test:
	go test -short

e2e-test:
	./scripts/e2e-test $(authKey)

.PHONY: test e2e-test
