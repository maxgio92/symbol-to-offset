PROGRAM := symbol-to-offset
go ?= $(shell command -v go)

$(PROGRAM):
	$(go) build -v -o $(PROGRAM) .

