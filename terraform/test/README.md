# Terratest AWS 

## Setup

Initialize the module and configure. For example: 

```sh
go mod init "github.com/devsecfranklin/customer-demo"
go mod tidy
go mod vendor # ln -s ~/go/pkg/mod
```

## Testing

```sh
go test -v -timeout 30m #all test cases
go test -v -run TestTerraformBasicExample # single test case
```

