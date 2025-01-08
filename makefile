consistent-hashing:
	@echo "Building consistent-hashing..."
	cd consistent-hashing && make $@

thread-pool:
	@echo "Building thread-pool..."
	cd thread-pool && make $@

websockets:
	@echo "Building websockets..."
	cd websockets && make $@

.PHONY: all
all: consistent-hashing thread-pool websockets