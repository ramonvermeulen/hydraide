# =============================================================================
# 📄 License Notice – HydrAIDE Intellectual Property (© 2025 Trendizz.com Kft.)
# =============================================================================
#
# This file is part of the HydrAIDE system and is protected by a custom,
# restrictive license. All rights reserved.
#
# ▸ This source is licensed for the exclusive purpose of building software that
#
#	interacts directly with the official HydrAIDE Engine.
#
# ▸ Redistribution, modification, reverse engineering, or reuse of any part of
#
#	this file outside the authorized HydrAIDE environment is strictly prohibited.
#
# ▸ You may NOT use this file to build or assist in building any:
#
#	– alternative engines,
#	– competing database or processing systems,
#	– protocol-compatible backends,
#	– SDKs for unauthorized runtimes,
#	– or any AI/ML training dataset or embedding extraction pipeline.
#
# ▸ This file may not be used in whole or in part for benchmarking, reimplementation,
#
#	architectural mimicry, or integration with systems that replicate or compete
#	with HydrAIDE’s features or design.
#
# By accessing or using this file, you accept the full terms of the HydrAIDE License.
# Violations may result in legal action, including injunctions or claims for damages.
#
# 🔗 License: https://github.com/hydraide/hydraide/blob/main/LICENSE.md
# ✉ Contact: hello@trendizz.com
# =============================================================================

# =============================================================================
# 📄 Makefile – HydrAIDE Proto Compiler
# =============================================================================
#
# This Makefile provides useful targets for generating gRPC client code from
# .proto definitions. It supports Go out-of-the-box and allows optional
# generation for Python, Node.js, Rust, Java, and C# if tools are installed.
#
# 🧠 Note:
# - Go SDK is pre-generated under ./generated/go
# - All other languages must be generated manually
#
# 🔐 Safe to run in CI/CD – missing tools will not break execution
#
# 🔗 Need help? → https://grpc.io/docs/
#
# =============================================================================

.PHONY: build proto clean proto-python proto-node proto-rust proto-java proto-csharp help

# -----------------------------------------------------------------------------
# 🧪 build – Regenerate Go code + tidy dependencies
# -----------------------------------------------------------------------------
# - Runs protoc with Go plugins
# - Ensures Go dependencies are updated
build: proto
	@echo "✅ Go dependencies updated"
	go mod tidy
	go get -u all

# -----------------------------------------------------------------------------
# 🛠️ proto – Compile .proto files to Go (no go get)
# -----------------------------------------------------------------------------
# - Generates .pb.go and .pb.grpc.go files to ./generated/go
# - Uses source-relative paths for imports
proto:
	@echo "🛠️  Generating Go gRPC files to ./generated/hydraidepbgo"
	protoc --proto_path=proto \
		--go_out=./generated/hydraidepbgo --go_opt=paths=source_relative \
		--go-grpc_out=./generated/hydraidepbgo --go-grpc_opt=paths=source_relative \
		proto/hydraide.proto || echo "⚠️  Go proto generation failed. Check protoc-gen-go & protoc-gen-go-grpc"

# -----------------------------------------------------------------------------
# 🧹 clean – Delete all generated proto output
# -----------------------------------------------------------------------------
# - Deletes all contents in the ./generated folders
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf generated/hydraidepbgo* generated/hydraidepbpy/* generated/hydraidepbjs/* generated/hydraidepbrs/* generated/hydraidepbjv/* generated/hydraidepbcs/*

# -----------------------------------------------------------------------------
# 🔹 proto-python – Generate Python client bindings (if grpc_tools available)
# -----------------------------------------------------------------------------
# Output: ./generated/python
proto-python:
	@echo "🐍 Generating Python gRPC files..."
	@command -v python3 >/dev/null 2>&1 || { echo "⚠️  Python not installed – skipping"; exit 0; }
	@protoc --proto_path=proto \
		--python_out=./generated/hydraidepbpy \
		--pyi_out=./generated/hydraidepbpy \
		--grpc_python_out=./generated/hydraidepbpy \
		proto/hydraide.proto || echo "⚠️  Python proto generation failed."

# -----------------------------------------------------------------------------
# 🔹 proto-node – Generate Node.js client bindings (requires grpc_tools_node_protoc_plugin)
# -----------------------------------------------------------------------------
# Output: ./generated/node
proto-node:
	@echo "🟨 Generating Node.js gRPC files..."
	@command -v protoc-gen-grpc >/dev/null 2>&1 || { echo "⚠️  Node.js gRPC plugin not found – skipping"; exit 0; }
	@protoc --proto_path=proto \
		--js_out=import_style=commonjs,binary:generated/hydraidepbjs \
		--grpc_out=generated/hydraidepbjs \
		proto/hydraide.proto || echo "⚠️  Node.js proto generation failed."

# -----------------------------------------------------------------------------
# 🔹 proto-rust – Generate Rust proto files (requires protoc-gen-prost)
# -----------------------------------------------------------------------------
# Output: ./generated/rust
proto-rust:
	@echo "🦀 Generating Rust proto files..."
	@command -v protoc-gen-prost >/dev/null 2>&1 || { echo "⚠️  protoc-gen-prost not installed – skipping"; exit 0; }
	@protoc --proto_path=proto \
		--prost_out=./generated/hydraidepbrs \
		proto/hydraide.proto || echo "⚠️  Rust proto generation failed."

# -----------------------------------------------------------------------------
# 🔹 proto-java – Generate Java proto files
# -----------------------------------------------------------------------------
# Output: ./generated/java
proto-java:
	@echo "☕ Generating Java proto files..."
	@protoc --proto_path=proto \
		--java_out=./generated/hydraidepbjv \
		--grpc-java_out=./generated/hydraidepbjv \
		proto/hydraide.proto || echo "⚠️  Java proto generation failed."

# -----------------------------------------------------------------------------
# 🔹 proto-csharp – Generate C# (.NET) proto files
# -----------------------------------------------------------------------------
# Output: ./generated/csharp
proto-csharp:
	@echo "🎯 Generating C# proto files..."
	@protoc --proto_path=proto \
		--csharp_out=./generated/hydraidepbcs \
		--grpc_out=./generated/hydraidepbcs \
		proto/hydraide.proto || echo "⚠️  C# proto generation failed."

# -----------------------------------------------------------------------------
# 📋 help – List all available make targets
# -----------------------------------------------------------------------------
help:
	@echo "📦 HydrAIDE Proto Makefile – Available commands:"
	@echo ""
	@echo "🔨 build          – Compile proto to Go and tidy dependencies"
	@echo "🧠 proto          – Only generate Go bindings"
	@echo "🐍 proto-python   – Generate Python gRPC code (if tools exist)"
	@echo "🟨 proto-node     – Generate Node.js gRPC code (if tools exist)"
	@echo "🦀 proto-rust     – Generate Rust proto files (requires protoc-gen-prost)"
	@echo "☕ proto-java     – Generate Java gRPC bindings"
	@echo "🎯 proto-csharp   – Generate C#/.NET gRPC bindings"
	@echo "🧹 clean          – Remove all generated proto code"
	@echo ""
	@echo "🧭 Notes:"
	@echo " - No plugins? No problem. Targets will skip gracefully."
	@echo " - Generated code goes into ./generated/<language>"
	@echo " - Go SDK is already pre-generated."
