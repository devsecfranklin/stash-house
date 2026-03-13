/*
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
*/

package logging

import (
	"fmt"
	"log"
	"os"
)

const (
	LRED   = "\033[1;31m"
	LGREEN = "\033[1;32m"
	LBLUE  = "\033[1;34m"
	LPURP  = "\033[1;35m"
	NC     = "\033[0m" // No Color
)

func Log_header(msg string) {
	fmt.Printf("\n%s# --- %s %s\n", LPURP, msg, NC)
}

func Log_info(msg string) {
	fmt.Printf("%s%s%s\n", LBLUE, msg, NC)
}

func Log_success(msg string) {
	fmt.Printf("%s%s%s\n", LGREEN, msg, NC)
}

func Log_error(msg string) {
	fmt.Printf("%sERROR: %s%s\n", LRED, msg, NC)
	os.Exit(1) // Exit on critical errors during setup
}

func Log_fatal(msg string) { // Added for graceful server shutdown logging
	fmt.Printf("%sFATAL: %s%s\n", LRED, msg, NC)
	log.Fatal(msg)
}

// should this be removed since it doesn't match
func CheckError(e error) {
    if e != nil {
        panic(e)
    }
}
