## Build
```bash
$ GOOS=js GOARCH=wasm go build -o main.wasm
```

## Run
Install Node.js and then
```bash
$ GOOS=js GOARCH=wasm go run -exec="$(go env GOROOT)/misc/wasm/go_js_wasm_exec" .
```