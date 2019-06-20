# Kahla API Client written in Go

<https://wiki.aiursoft.com/ReadDoc/Kahla/Auth.md>

## Rebuild

```bash
cd kahla/
protoc *.proto -I ./ -I $GOPATH/src/github.com/ganlvtech/go-rest-client/protoc-gen-gorestclient/rest/ --gorestclient_out=.
```

## Demo

See `examples/` and `kahla/*_test.go`.

## License

MIT License
