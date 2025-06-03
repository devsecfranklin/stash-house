# games.bitsmasher.net

```sh
/usr/local/go/bin/go mod init games
go mod tidy
```

## Testing

```sh
go install github.com/kisielk/errcheck@latest
$HOME/go/bin/errcheck ./...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint help linters
golangci-lint run
```
