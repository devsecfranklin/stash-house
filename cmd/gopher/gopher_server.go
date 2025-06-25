package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Define the host and port for our Gopher server.
// Port 70 is the standard for Gopher, but requires root/admin privileges.
// We use a high port like 7070 for easy development.
const (
	HOST = "localhost"
	PORT = "7070"
	ADDR = HOST + ":" + PORT
)

// The main function sets up the listener and accepts connections.
func main() {
	// Start listening for TCP connections on our address.
	listener, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Ensure the listener is closed when main() exits.
	defer listener.Close()

	log.Printf("Gopher server listening on gopher://%s", ADDR)

	// The main loop to accept incoming connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		// Handle each connection in a new goroutine for concurrency.
		go handleConnection(conn)
	}
}

// handleConnection processes a single Gopher client request.
func handleConnection(conn net.Conn) {
	// The Gopher protocol is simple: after we send the response, we close the connection.
	defer conn.Close()
	log.Printf("Serving request from %s", conn.RemoteAddr().String())

	// Read the "selector" sent by the client. For the root menu, this is usually empty.
	reader := bufio.NewReader(conn)
	selector, err := reader.ReadString('\n')
	if err != nil {
		// This can happen if the client disconnects, it's not a critical error.
		log.Printf("Could not read from client: %v", err)
		return
	}
	selector = strings.TrimSpace(selector)

	// Route the request based on the selector.
	switch selector {
	case "":
		serveRootMenu(conn)
	case "/about":
		serveTextFile(conn, "This is a Gopher server created for a CTF challenge!\nGood luck, hacker.\n")
	case "/files":
		serveFilesMenu(conn)
	case "/logs/2024-01-18.log":
		serveTextFile(conn, "12:01 - Server started.\n12:05 - User 'admin' logged in.\n13:30 - User 'guest' accessed /files.\n14:00 - Hmm, I see a strange request for a file that doesn't exist... something like /d33p_s3cr3t_p4th?\n")
	case "/d33p_s3cr3t_p4th":
		// This is a hidden path, not listed in any menu!
		serveTextFile(conn, "Congratulations! You found the hidden path!\nHere is your flag: CTF{G0ph3r_Pr0t0c0l_Is_R3tr0_C00l}\n")
	default:
		serveNotFound(conn, selector)
	}
}

// --- Response Handlers ---

// Gopher Item Format: <ItemType><DisplayText><TAB><Selector><TAB><Host><TAB><Port>
// The response must end with a single period (.) on a new line.

func serveRootMenu(conn net.Conn) {
	// ItemType 'i' is for informational text.
	// ItemType '1' is for a menu/subdirectory.
	// ItemType '0' is for a text file.
	fmt.Fprintf(conn, "iWelcome to the Gopher Hole!\t\t\t\r\n")
	fmt.Fprintf(conn, "i---------------------------\t\t\t\r\n")
	fmt.Fprintf(conn, "0About this server\t/about\t%s\t%s\r\n", HOST, PORT)
	fmt.Fprintf(conn, "1Browse Files\t/files\t%s\t%s\r\n", HOST, PORT)
	fmt.Fprintf(conn, "i \t\t\t\r\n")
	fmt.Fprintf(conn, "iHint: Not all paths are listed in menus.\t\t\t\r\n")
	fmt.Fprintf(conn, ".\r\n") // End of response
}

func serveFilesMenu(conn net.Conn) {
	fmt.Fprintf(conn, "iDirectory listing for /files\t\t\t\r\n")
	fmt.Fprintf(conn, "i-----------------------------\t\t\t\r\n")
	fmt.Fprintf(conn, "0README.txt\t/about\t%s\t%s\r\n", HOST, PORT) // Re-using the about page
	fmt.Fprintf(conn, "1logs\t/logs\t%s\t%s\r\n", HOST, PORT) // This will lead to a 404, a potential rabbit hole!
	fmt.Fprintf(conn, "0access_log_2024-01-18.log\t/logs/2024-01-18.log\t%s\t%s\r\n", HOST, PORT)
	fmt.Fprintf(conn, ".\r\n")
}

func serveTextFile(conn net.Conn, content string) {
	// For a text file response, you just send the content directly.
	fmt.Fprintf(conn, content)
	fmt.Fprintf(conn, ".\r\n")
}

func serveNotFound(conn net.Conn, selector string) {
	// ItemType '3' is for an error.
	fmt.Fprintf(conn, "3Error: Path not found: '%s'\t\t\t\r\n", selector)
	fmt.Fprintf(conn, ".\r\n")
}
