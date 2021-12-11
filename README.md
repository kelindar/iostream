<p align="center">
<img height="110" src=".github/logo.png" border="0" alt="kelindar/iostream">
<br>
<img src="https://img.shields.io/github/go-mod/go-version/kelindar/iostream" alt="Go Version">
<a href="https://pkg.go.dev/github.com/kelindar/iostream"><img src="https://pkg.go.dev/badge/github.com/kelindar/iostream" alt="PkgGoDev"></a>
<a href="https://goreportcard.com/report/github.com/kelindar/iostream"><img src="https://goreportcard.com/badge/github.com/kelindar/iostream" alt="Go Report Card"></a>
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
<a href="https://coveralls.io/github/kelindar/iostream"><img src="https://coveralls.io/repos/github/kelindar/iostream/badge.svg" alt="Coverage"></a>
</p>

## Simple Binary Stream Reader/Writer

This package contains a set of simple utility reader and writer that can be used to efficiently read/write binary information over the wire.

```go
import "github.com/kelindar/iostream"
```

## Usage

In order to use it, you can simply instantiate either a `Reader` or a `Writer` which require an `io.Reader` or `io.Writer` respectively. Here's a simple example

```go
// Fake stream
stream := bytes.NewBuffer(nil)

// Write some data into the stream...
w := iostream.NewWriter(stream)
w.WriteString("Roman")
w.WriteUint32(36)

// Read the data back...
r := iostream.NewReader(stream)
name, err := r.ReadString()
age, err  := r.ReadUint32()
```

## Contributing

We are open to contributions, feel free to submit a pull request and we'll review it as quickly as we can. This library is maintained by [Roman Atachiants](https://www.linkedin.com/in/atachiants/)

## License

Tile is licensed under the [MIT License](LICENSE.md).
