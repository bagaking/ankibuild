.PHONY: bundle example test

example:
	go run .

test:
	go test ./...

bundle:
	$(MAKE) -C bundle -f Makefile clean
	$(MAKE) -f bundle/Makefile
