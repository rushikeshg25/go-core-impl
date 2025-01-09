consistent-hashing:
	@echo "Building consistent-hashing..."
	cd consistent-hashing && go run .

thread-pool:
	@echo "Building thread-pool..."
	cd thread-pool && make go run .

websockets:
	@echo "Building websockets..."
	cd websockets && make go run .

.PHONY: consistent-hashing thread-pool websockets