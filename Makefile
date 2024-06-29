.PHONY: bundle example

example:
	go run .

bundle:
	$(MAKE) -C bundle -f Makefile clean
	$(MAKE) -f bundle/Makefile