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

You can do application testing, but there is no `docker` binary in this build env.
Docker does not work on OpenBSD.

```sh
doas pkg_add node
doas npm install -g npm@latest # upgrade npm
npm -v
npm audit fix # fix sec vulns
```

## Testing Golang

- Errcheck

```sh
go install github.com/kisielk/errcheck@latest
~/go/bin/errcheck ./...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint help linters
golangci-lint run
```

- check-golang-templates

```sh
git clone git@github.com:gosom/check-golang-templates.git
cd check-golang-template
go mod download
go install
check-golang-templates -folder sites/games/web/templates
```

Explanation of Tests:
main_test.go
Mocking:

Logging: The init() function at the beginning of main_test.go overrides the logging package functions to prevent actual log output during tests. This is crucial for isolating the code under test.

HTTP Requests/Responses: httptest package is used to create mock http.Request and http.ResponseWriter objects, allowing us to simulate client requests and inspect server responses without actually starting an HTTP server.

http.DefaultClient: For the OAuth callback test, the http.DefaultClient is temporarily replaced with a custom RoundTripFunc to mock the HTTP call to the Twitch token endpoint. This allows us to control the response received from Twitch's API.

auth.GenerateRandomState: The GenerateRandomState function from internal/auth is mocked to return a predictable string, ensuring deterministic testing of the OAuth state handling.

Test Cases:

TestHandler: Verifies that the root handler serves the indexPage correctly and sets appropriate cache control headers.

TestOauthHandlerInitialRequest: Checks if the /oauth endpoint correctly displays the authorization page and generates/stores a new OAuth state.

TestOauthHandlerCallbackSuccess: Simulates a successful callback from Twitch, including the exchange of an authorization code for an access token. It verifies the response content and that the state is removed after use.

TestOauthHandlerCallbackError: Tests the scenario where Twitch returns an error during the OAuth callback.

TestOauthHandlerMissingClientID: Checks behavior when environment variables for Twitch client ID are not set.

TestTwitchChatHandler: Verifies the twitchChatHandler correctly renders the chatPage and passes the twitchOauthToken.

gopher_server_test.go
Mocking net.Conn:

A mockConn struct is created that implements the net.Conn interface. This allows us to simulate network connections by controlling what bytes are "read" from the client and capturing what bytes are "written" to the client.

The newMockConn helper function makes it easy to set up test scenarios.

Test Cases:

TestServeRootMenu, TestServeFilesMenu, TestServeTextFile, TestServeNotFound: These tests directly call the individual Gopher response functions and assert that they write the correct Gopher protocol formatted output to the mock connection's buffer.

TestHandleConnectionRoot, TestHandleConnectionAbout, TestHandleConnectionFiles, TestHandleConnectionLogs, TestHandleConnectionHiddenPath, TestHandleConnectionNotFound: These tests simulate full Gopher client requests by providing a "selector" to the handleConnection function via the mockConn's read buffer and then asserting the complete Gopher response generated in the write buffer. This covers the routing logic of the server.

These tests provide good coverage for the core logic of both applications. For more comprehensive testing, you might consider:

Integration Tests: To verify that different components (e.g., main.go's HTTP server and external services like Twitch) work together correctly.

End-to-End Tests: To simulate a user's entire flow through the application.

Table-Driven Tests: For functions with multiple input/output scenarios (especially in gopher_server.go for different selectors).