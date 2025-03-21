# Makefile – Hydra Proto Compiler
# ==================================
# This Makefile provides useful targets for generating Go code from .proto files.
#
# 📌 Default Behavior:
# - If you are using Go, the pre-generated .pb.go files are available under ./generated/go.
# - If you need to generate files for another language (Python, Rust, Node.js, etc.), see the "Custom Generation" section below.

# Define phony targets (not actual files)
.PHONY: build proto clean

# --------------------------------------
# build – Compile the proto file and update Go dependencies
# --------------------------------------
# This target:
# - Runs protoc with Go + gRPC plugin to generate .pb.go and .pb.grpc.go files
# - Ensures source-relative paths are used for clean import paths
# - Updates all Go module dependencies (recommended after regeneration)
build: proto
	@echo "✅ Go dependencies updated"
	go mod tidy
	go get -u all

# --------------------------------------
# proto – Only compile the proto files (no go get)
# --------------------------------------
# Use this if you just want to regenerate the .pb.go files
# without updating Go dependencies.
proto:
	@echo "🛠️ Generating Go gRPC files to ./generated/go"
	protoc --proto_path=proto \
		--go_out=./generated/go --go_opt=paths=source_relative \
		--go-grpc_out=./generated/go --go-grpc_opt=paths=source_relative \
		proto/hydra.proto

# --------------------------------------
# clean – Optional: implement file cleanup logic here (e.g., remove generated files)
# --------------------------------------
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf generated/go/*

# --------------------------------------
# 🔹 Custom Proto Compilation for Other Languages 🔹
# --------------------------------------
# If you need to generate client bindings for **other programming languages**,
# use the commands below and adjust paths as necessary.

# Example for Python:
# protoc --proto_path=proto --python_out=generated/python proto/hydra.proto

# Example for Node.js:
# protoc --proto_path=proto --js_out=import_style=commonjs,binary:generated/node proto/hydra.proto

# Example for Rust:
# protoc --proto_path=proto --rust_out=generated/rust proto/hydra.proto

# Example for Java:
# protoc --proto_path=proto --java_out=generated/java proto/hydra.proto

# --------------------------------------
# 🚀 Notes:
# - The **Go SDK** is already pre-generated and available under ./generated/go.
# - If you need another language, **you must generate the bindings manually** using protoc.
# - This Makefile does NOT auto-generate bindings for non-Go languages (to avoid unnecessary dependencies).
#
# Need help? See https://grpc.io/docs/