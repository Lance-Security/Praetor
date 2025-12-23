GO      ?= go
BINARY  ?= pt
PKG     := ./...
MAIN    := ./main.go
OUTDIR  := bin

.PHONY: all build run test tidy clean

all: build

build:
	@mkdir -p $(OUTDIR)
	$(GO) build -o $(OUTDIR)/$(BINARY) $(MAIN)

run:
	$(GO) run $(MAIN)

test:
	$(GO) test $(PKG)

tidy:
	$(GO) mod tidy

clean:
	@rm -rf $(OUTDIR)