# CSVX

`csvx` is a package with a CSV struct scanner for Go. It is licensed under the permissive Apache 2 license, so it can be used in open and closed source projects alike.


## Implemented Types

- [x] string
- [x] int
- [x] int64
- [x] int32
- [x] int16
- [x] int8
- [x] uint
- [x] uint64
- [x] uint32
- [x] uint16
- [x] uint8
- [x] float64
- [x] float32
- [x] bool ("true", "yes", "1", "1.0" = true, "false", "no", "0", "0.0" = false)

- [x] Pointer types to above underlying types, e.g. `*string` (empty string and `"null"` result in `nil` being set on the Go struct)

## Performance

The struct scanner uses `reflect` quite heavily, so this library will not be as fast as writing a specific parser for the struct. However, for the vast majority of cases, the performance hit will be acceptable and the development speed increase and simple client code will be worth it!