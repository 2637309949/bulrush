# all depend on middleware
all: middleware

# $@ mean end product `middleware`
middleware:
	@go build -o middleware
	@./run 1 true $@
	@./run 5 true $@
	@./run 10 true $@
	@./run 15 true $@
	@./run 20 true $@
	@./run 30 true $@
	@./run 50 true $@
	@echo

# means the word "all, middleware" doesn't represent a file name in this Makefile
.PHONY: all middleware
