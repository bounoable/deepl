test:
	go test -short

e2e-test:
	./scripts/e2e-test $(authKey)

docs:
	@./scripts/docs

.PHONY: test e2e-test docs
