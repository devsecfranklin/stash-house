/*
# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock connection struct to simulate net.Conn
type mockConn struct {
	readBuf  bytes.Buffer
	writeBuf bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuf.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t_ int) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t_ int) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t_ int) error {
	return nil
}

// Helper function to create a mock connection with a given input and capture output
func newMockConn(input string) *mockConn {
	conn := &mockConn{}
	conn.readBuf.WriteString(input)
	return conn
}

// TestServeRootMenu tests the serveRootMenu function
func TestServeRootMenu(t *testing.T) {
	conn := newMockConn("") // No input needed for serving the root menu
	serveRootMenu(conn)

	expectedOutput := fmt.Sprintf(
		"iWelcome to the Gopher Hole!\t\t\t\r\n" +
			"i---------------------------\t\t\t\r\n" +
			"0About this server\t/about\t%s\t%s\r\n" +
			"1Browse Files\t/files\t%s\t%s\r\n" +
			"i \t\t\t\r\n" +
			"iHint: Not all paths are listed in menus.\t\t\t\r\n" +
			".\r\n", HOST, PORT, HOST, PORT)

	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "Root menu output mismatch")
}

// TestServeFilesMenu tests the serveFilesMenu function
func TestServeFilesMenu(t *testing.T) {
	conn := newMockConn("")
	serveFilesMenu(conn)

	expectedOutput := fmt.Sprintf(
		"iDirectory listing for /files\t\t\t\r\n" +
			"i-----------------------------\t\t\t\r\n" +
			"0README.txt\t/about\t%s\t%s\r\n" +
			"1logs\t/logs\t%s\t%s\r\n" +
			"0access_log_2024-01-18.log\t/logs/2024-01-18.log\t%s\t%s\r\n" +
			".\r\n", HOST, PORT, HOST, PORT, HOST, PORT)

	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "Files menu output mismatch")
}

// TestServeTextFile tests the serveTextFile function
func TestServeTextFile(t *testing.T) {
	conn := newMockConn("")
	content := "Hello, Gopher!\nThis is a test file."
	serveTextFile(conn, content)

	expectedOutput := content + ".\r\n"
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "Text file output mismatch")
}

// TestServeNotFound tests the serveNotFound function
func TestServeNotFound(t *testing.T) {
	conn := newMockConn("")
	selector := "/nonexistent_path"
	serveNotFound(conn, selector)

	expectedOutput := fmt.Sprintf("3Error: Path not found: '%s'\t\t\t\r\n.\r\n", selector)
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "Not found output mismatch")
}

// TestHandleConnectionRoot tests handling of the root selector
func TestHandleConnectionRoot(t *testing.T) {
	conn := newMockConn("\n") // Empty selector for root
	handleConnection(conn)

	expectedOutput := fmt.Sprintf(
		"iWelcome to the Gopher Hole!\t\t\t\r\n" +
			"i---------------------------\t\t\t\r\n" +
			"0About this server\t/about\t%s\t%s\r\n" +
			"1Browse Files\t/files\t%s\t%s\r\n" +
			"i \t\t\t\r\n" +
			"iHint: Not all paths are listed in menus.\t\t\t\r\n" +
			".\r\n", HOST, PORT, HOST, PORT)

	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection root output mismatch")
}

// TestHandleConnectionAbout tests handling of the /about selector
func TestHandleConnectionAbout(t *testing.T) {
	conn := newMockConn("/about\n")
	handleConnection(conn)

	expectedOutput := "This is a Gopher server created for a CTF challenge!\nGood luck, hacker.\n.\r\n"
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection /about output mismatch")
}

// TestHandleConnectionFiles tests handling of the /files selector
func TestHandleConnectionFiles(t *testing.T) {
	conn := newMockConn("/files\n")
	handleConnection(conn)

	expectedOutput := fmt.Sprintf(
		"iDirectory listing for /files\t\t\t\r\n" +
			"i-----------------------------\t\t\t\r\n" +
			"0README.txt\t/about\t%s\t%s\r\n" +
			"1logs\t/logs\t%s\t%s\r\n" +
			"0access_log_2024-01-18.log\t/logs/2024-01-18.log\t%s\t%s\r\n" +
			".\r\n", HOST, PORT, HOST, PORT, HOST, PORT)

	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection /files output mismatch")
}

// TestHandleConnectionLogs tests handling of a specific log file selector
func TestHandleConnectionLogs(t *testing.T) {
	conn := newMockConn("/logs/2024-01-18.log\n")
	handleConnection(conn)

	expectedOutput := "12:01 - Server started.\n12:05 - User 'admin' logged in.\n13:30 - User 'guest' accessed /files.\n14:00 - Hmm, I see a strange request for a file that doesn't exist... something like /d33p_s3cr3t_p4th?\n.\r\n"
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection /logs/2024-01-18.log output mismatch")
}

// TestHandleConnectionHiddenPath tests handling of the hidden path
func TestHandleConnectionHiddenPath(t *testing.T) {
	conn := newMockConn("/d33p_s3cr3t_p4th\n")
	handleConnection(conn)

	expectedOutput := "Congratulations! You found the hidden path!\nHere is your flag: CTF{G0ph3r_Pr0t0c0l_Is_R3tr0_C00l}\n.\r\n"
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection hidden path output mismatch")
}

// TestHandleConnectionNotFound tests handling of an unknown selector
func TestHandleConnectionNotFound(t *testing.T) {
	conn := newMockConn("/unknown_path\n")
	handleConnection(conn)

	expectedOutput := "3Error: Path not found: '/unknown_path'\t\t\t\r\n.\r\n"
	assert.Equal(t, expectedOutput, conn.writeBuf.String(), "HandleConnection not found output mismatch")
}

/*
func TestMain(t *testing.T) {
	// Testing main() for a server application typically involves
	// starting the server in a goroutine and then making a client
	// connection to it. This is more of an integration test.
	// For unit tests, we focus on the individual handler functions.
	// A simple way to test main's fatal error handling is to mock log.Fatalf
	// and trigger an error condition (e.g., failed to listen).

	oldFatalf := log.Fatalf
	defer func() { log.Fatalf = oldFatalf }()

	// Simulate listen failure
	log.Fatalf = func(format string, v ...interface{}) {
		panic(fmt.Sprintf(format, v...)) // Catch fatal errors as panics
	}

	// This part is difficult to test without actually trying to listen on a port
	// that's already in use or privileged without root access.
	// For a true unit test, you would mock net.Listen.
}
*/

