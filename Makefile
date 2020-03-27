.PHONY: test
test: vet
	go test -v pkg/*

.PHONY: vet
vet:
	go vet pkg/*

# TODO(gibi): which go linter is the industry standard?
