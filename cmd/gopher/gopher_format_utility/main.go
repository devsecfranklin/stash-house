/*
gfu - gophermap format utility

`gfu` manipulates gophermaps (gopher menus). It is intended to be used as part
of an automation chain for managing a gopherhole using maps as the main
doctype, easily allowing any document to contain links.

`gfu` can:
- Convert all lines in a file that are not valid gopher links into gopher info
text (item type 'i') lines
- Deconstruct all gopher info text lines in a file back to plain text for easy
editing

There are also plans to include additional features, such as:
- Adding the contents of a header file into the gophermap
- Adding the contents of a footer file into the gophermap *Please note - many
servers already support includes, so the above may not be needed. If you are
interested in this feature, you may want to check your server documentation
first.*

Documentation

For information on using `gfu`, see the help information: gfu --help

Information on building, downloading and installing `gfu` can be found at:
https://tildegit.org/sloum/gfu

Further information on gfu can be found at the `gfu` homepage:
https://rawtext.club/~sloum/gfu.html
*/
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

//version information, set by ldflags with some default fallbacks
var (
	version string = "not set"
	build   string = "not set"
)

//command line flag global variables
var (
	deconstructInfoTextLines bool
	header                   string
	footer                   string
	stdout                   bool
	printVersion             bool
)

//regex global variable that identifies the format of gopher menu lines
var reGopherLines = regexp.MustCompile(`.+\t.*\t.*\t.*`)

func buildInfoText(ln string) string {
	var out strings.Builder
	out.Grow(20 + len(ln))
	out.WriteString("i")
	out.WriteString(ln)
	out.WriteString("\tfalse\tnull.host\t1")
	return out.String()
}

func deconstructInfoText(ln string) string {
	text := strings.SplitN(ln, "\t", 2)
	infotext := text[0]
	if len(infotext) > 1 {
		return infotext[1:]
	}
	return ""
}

//takes a string representing the file path, reads the file, then passes the
//data for processing. returns data as a bytes.Buffer and any error
//information.
func readFile(path string) (outFile bytes.Buffer, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()

	outFile, err = processFile(file)
	if err != nil {
		return
	}

	return outFile, err
}

//takes data as a bytes.Buffer, reads and processes each line according to the
//deconstructInfoTextLines flag. returns processed data as a bytes.Buffer and
//any error information.
func processFile(file io.Reader) (outFile bytes.Buffer, err error) {
	scanner := bufio.NewScanner(file)

	if !deconstructInfoTextLines {
		for scanner.Scan() {
			l := scanner.Text()
			if exp := reGopherLines.MatchString(l); exp {
				outFile.WriteString(l)
			} else {
				outFile.WriteString(buildInfoText(l))
			}
			outFile.WriteString("\n")
		}
	} else {
		for scanner.Scan() {
			l := scanner.Text()
			if exp := reGopherLines.MatchString(l); exp && l[0] == 'i' {
				outFile.WriteString(deconstructInfoText(l))
			} else {
				outFile.WriteString(l)
			}
			outFile.WriteString("\n")
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return outFile, nil
}

//takes a file path as a string and some data as bytes.Buffer. writes the data
//to the specified file. returns any error information.
func writeFile(path string, outFile bytes.Buffer) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(outFile.Bytes())
	return nil
}

//configure command line flags
func init() {
	flag.BoolVar(&deconstructInfoTextLines, "d", false, "Deconstruct a gophermap's info text lines back to plain text")
	flag.StringVar(&header, "head", "", "Path to a file containing header content")
	flag.StringVar(&footer, "foot", "", "Path to a file containing footer content")
	flag.BoolVar(&stdout, "stdout", false, "Instead of writing changes to a file, return them to stdout")
	flag.BoolVar(&printVersion, "v", false, "Output version information")
}

//PrintHelp produces a nice display message when the --help flag is used
func PrintHelp() {
	art := `gfu - gophermap formatting utility

syntax: gfu [flags...] [filepath]

example: gfu -d ~/gopher/phlog/gophermap

default
	Convert plain text lines to gophermap info text (item type 'i') lines
`
	fmt.Fprint(os.Stdout, art)
	flag.PrintDefaults()
}

func main() {
	//process command line flags
	flag.Usage = PrintHelp
	flag.Parse()
	args := flag.Args()

	if printVersion {
		fmt.Printf("gfu - gophermap format utility - version %s build %s\n", version, build)
		os.Exit(0)
	}
	if l := len(args); l != 1 {
		fmt.Fprintf(os.Stderr, "Incorrect number of arguments. Expected 1, got %d\n", l)
		os.Exit(1)
	}
	if header != "" {
		fmt.Println("Header functionality is not built yet, proceeding with general gophermap conversion...")
	}

	if footer != "" {
		fmt.Println("Footer functionality is not built yet, proceeding with general gophermap conversion...")
	}

	//read and process file
	outFile, err := readFile(args[0])
	if err != nil {
		log.Fatalln("Error while reading file -", err)
	}

	//output data to stdout or file
	if stdout {
		fmt.Print(outFile.String())
	} else {
		err = writeFile(args[0], outFile)
		if err != nil {
			log.Fatalln("Error while writing file -", err)
		}
	}

	//the end
	os.Exit(0)
}
