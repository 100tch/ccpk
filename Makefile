BINARY := ccpk
PREFIX := $(HOME)/.local/bin
BUILD_DIR := build

.PHONY: all build install uninstall clean

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/

install:
	mkdir -p $(PREFIX)
	install -m 755 $(BUILD_DIR)/$(BINARY) $(PREFIX)/$(BINARY)
	@echo "installed to $(PREFIX)/$(BINARY)"
	@case ":$$PATH:" in \
		*":$(BINDIR):"*) ;; \
		*) echo "warning: $(BINDIR) is not in your PATH" ;; \
	esac

uninstall:
	rm -f $(PREFIX)/$(BINARY)
	@echo "removed $(PREFIX)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)
