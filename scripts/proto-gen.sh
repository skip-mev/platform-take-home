#!/bin/bash
set -e

echo "Generating proto code..."

# Get the root directory of the project
ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR/proto"

# Find all unique directories containing .proto files
proto_dirs=$(find . -name '*.proto' -print0 2>/dev/null | xargs -0 -n1 dirname 2>/dev/null | sort | uniq)

if [ -z "$proto_dirs" ]; then
    echo "No .proto files found in proto directory"
    exit 1
fi

for dir in $proto_dirs; do
    for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
        if grep go_package "$file" &>/dev/null; then
            buf generate --template buf.gen.yaml "$file" || { echo "Failed to generate code for $file"; exit 1; }
        fi
    done
done

echo "Proto code generation completed successfully."
