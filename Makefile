BINARY := ccpk
PREFIX := $(HOME)/.local/bin

.PHONY: all build install uninstall clean

all: build

build:
	go build -o $(BINARY) ./cmd/

install:
	mkdir -p $(PREFIX)
	install -m 755 $(BINARY) $(PREFIX)/$(BINARY)
	@echo "installed to $(PREFIX)/$(BINARY)"
	@case ":$$PATH:" in \
		*":$(PREFIX):"*) ;; \
		*) echo "warning: $(PREFIX) is not in your PATH" ;; \
	esac

uninstall:
	rm -f $(PREFIX)/$(BINARY)
	@echo "removed $(PREFIX)/$(BINARY)"

clean:
	rm $(BINARY)
