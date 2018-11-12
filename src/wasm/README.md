## Build
```bash
$ GOOS=js GOARCH=wasm go build -o main.wasm
```

## Run
### Node
Install Node.js and then
```bash
$ GOOS=js GOARCH=wasm go run -exec="$(go env GOROOT)/misc/wasm/go_js_wasm_exec" .
```

### Life VM
Install Life and then
```bash
# entry point is `run`
life -entry run ./main.wasm
```

But now it is showing 
```bash
$ life -entry run ./main.wasm
proto: tag has unknown wire type: "varint,1,opt,name=version,proto3"
proto: tag has unknown wire type: "varint,1,opt,name=version,proto3"
proto: tag has unknown wire type: "varint,1,opt,name=version,proto3"
proto: tag has unknown wire type: "varint,1,opt,name=version,proto3"
...
```
It is from `github.com/golang/protobuf@v1.2.0/proto/properties.go`
```go
// Parse populates p by parsing a string in the protobuf struct field tag style.
func (p *Properties) Parse(s string) {
	// "bytes,49,opt,name=foo,def=hello!"
	fields := strings.Split(s, ",") // breaks def=, but handled below.
	if len(fields) < 2 {
		fmt.Fprintf(os.Stderr, "proto: tag has too few fields: %q\n", s)
		return
	}

	p.Wire = fields[0]
	switch p.Wire {
	case "varint":
		p.WireType = WireVarint
	case "fixed32":
		p.WireType = WireFixed32
	case "fixed64":
		p.WireType = WireFixed64
	case "zigzag32":
		p.WireType = WireVarint
	case "zigzag64":
		p.WireType = WireVarint
	case "bytes", "group":
		p.WireType = WireBytes
		// no numeric converter for non-numeric types
	default:
		fmt.Fprintf(os.Stderr, "proto: tag has unknown wire type: %q\n", s)
		return
	}

	var err error
	p.Tag, err = strconv.Atoi(fields[1])
	if err != nil {
		return
	}

outer:
	for i := 2; i < len(fields); i++ {
		f := fields[i]
		switch {
		case f == "req":
			p.Required = true
		case f == "opt":
			p.Optional = true
		case f == "rep":
			p.Repeated = true
		case f == "packed":
			p.Packed = true
		case strings.HasPrefix(f, "name="):
			p.OrigName = f[5:]
		case strings.HasPrefix(f, "json="):
			p.JSONName = f[5:]
		case strings.HasPrefix(f, "enum="):
			p.Enum = f[5:]
		case f == "proto3":
			p.proto3 = true
		case f == "oneof":
			p.oneof = true
		case strings.HasPrefix(f, "def="):
			p.HasDefault = true
			p.Default = f[4:] // rest of string
			if i+1 < len(fields) {
				// Commas aren't escaped, and def is always last.
				p.Default += "," + strings.Join(fields[i+1:], ",")
				break outer
			}
		}
	}
}
```
It is an issue about the VM executing the switch statement.   
Node is ok without this issue.