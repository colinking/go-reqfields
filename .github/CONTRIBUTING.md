# Contributing

## Testing

You can test by running against any of the examples in the `./fixtures` directory.

```sh
go run ./cmd/reqfields ./fixtures/ex1
```

## Future Work

- [ ] Support validation on structs across 1st-party packages
- [ ] Support validation on structs in 3rd-party packages
- [ ] Support validation on inline structs
- [ ] Support unnamed parameters

## Inspiration

- [Using go/analysis to write a custom linter](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
