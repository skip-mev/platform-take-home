echo "Generating proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &>/dev/null; then
      buf generate --template buf.gen.yaml $file
    fi
  done
done

# Generate external protocol buffers
# echo "Generating cosmwasm protos"
# buf generate buf.build/cosmwasm/wasmd

rm -r ../api/types
mv github.com/skip-mev/platform-take-home/api/types ../api
rm -r github.com

# move proto files to the right places
