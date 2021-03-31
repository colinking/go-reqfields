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

## VSCode

To use this linter in VSCode, add the following to your settings to configure `golangci-lint` as your linter:

```json
"go.lintTool": "golangci-lint",
"go.lintFlags": ["--fast"],
```

Next, download this repo and run `go generate ./...`. This will generate a Go plugin that `golangci-lint` will use.

Then, add a `.golangci.yml` to your repo with the following configuration:

```yaml
linters-settings:
  custom:
    required:
      # Make sure to update this path to point at your local copy of `colinking/go-required`:
      path: ./cmd/plugin/main.so
      description: Compile-time warnings for required fields.
      original-url: github.com/colinking/go-required

linters:
  enable:
    - required
```

Once you reload VSCode, you should now see lint warnings!

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
