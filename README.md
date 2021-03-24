# go-required

A Go linter that warns of missing required struct fields at compile-time.

## Installation

```sh
go get github.com/colinking/go-required
```

## Usage

```sh
reqfields <pkg | file>
```

## Example

```sh
go run ./cmd/reqfields ./fixtures/ex1
```

## Future Work

- [ ] Support validation on structs across 1st-party packages
- [ ] Support validation on structs in 3rd-party packages
- [ ] Support validation on inline structs
- [ ] VSCode plugin? See: https://golangci-lint.run/contributing/new-linters/
- [ ] Benchmark test?
- [ ] Unnamed parameters?

## Inspiration

- [Using go/analysis to write a custom linter](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
