# Contributing

## Testing

You can test by running against any of the examples in the `./fixtures` directory.

```sh
go run ./cmd/reqfields ./fixtures/ex1
```

## golangci-lint

If you are testing `reqfields` with `golangci-lint`, keep in mind that you'll need to clear its cache to get updated results whenever you change `reqfields`. You can do that with:

```sh
golangci-lint cache clean
```

## Future Work

- [ ] Support validation on structs across (1st+3rd party) packages
- [ ] Support validation on inline structs
- [ ] Support unnamed parameters
- [ ] Add [formal analyzer-based tests](https://pkg.go.dev/golang.org/x/tools/go/analysis#hdr-Testing_an_Analyzer)

## Inspiration

- [`go/analysis` docs](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Using go/analysis to write a custom linter](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
