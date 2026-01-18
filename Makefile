GO      ?= go
BINARY  ?= pt
PKG     := ./...
MAIN    := ./cmd/praetor
OUTDIR  := bin

.PHONY: all build run test tidy clean install lint precommit

all: build

build:
	@mkdir -p $(OUTDIR)
	$(GO) build -o $(OUTDIR)/$(BINARY) $(MAIN)

run:
	$(GO) run ./cmd/praetor

test:
	$(GO) test -v $(PKG)

tidy:
	$(GO) mod tidy

clean:
	@rm -rf $(OUTDIR)

install:
	$(GO) install $(MAIN)

lint:
	$(GO) vet $(PKG)
	golangci-lint run

precommit: run test lint
