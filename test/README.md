# Testing

## Node

```sh
npm uninstall htmllint htmllint-cli
npm install @linthtml/linthtml --save-dev # https://github.com/linthtml/linthtml/blob/develop/doc/docs/user-guide/installation-and-usage.md
npx linthtml --init
npm audit fix
npx linthtml 'yourfile.html' # you can run LintHTML on any file or directory like this
npx linthtml 'src/**/*.html'
```

### OpenBSD

You can do application testing but there is no `docker` in this build env.

```sh
doas pkg_add node
doas npm install -g npm@latest # upgrade npm
npm -v
npm audit fix # fix sec vulns
```

## Testing Golang

```sh
go install github.com/kisielk/errcheck@latest
~/go/bin/errcheck ./...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint help linters
golangci-lint run
```
