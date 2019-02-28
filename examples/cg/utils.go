package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const AxisFormat = "2006-01-02T15:04:05.000Z"

func PathExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CommandLine(args []string) ([]string, error) {
	// command line
	if len(args) < 1 {
		Header()
		fmt.Print("Cmd line:")
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		inputs := strings.Fields(strings.Replace(input, "\n", "", -1))
		args = append(args, inputs...)
	}
	return args, nil
}

// Check if error exist
// Exit if err != nil
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Header() {
	fmt.Println(`
   ____   U  ___ u  __  __     ____    U  ___ u 
U /"___|   \/"_ \/U|' \/ '|uU /"___|u   \/"_ \/ 
\| | u     | | | |\| |\/| |/\| |  _ /   | | | | 
 | |/__.-,_| |_| | | |  | |  | |_| |.-,_| |_| | 
  \____|\_)-\___/  |_|  |_|   \____| \_)-\___/  
 _// \\      \\   <<,-,,-.    _)(|_       \\    
(__)(__)    (__)   (./  \.)  (__)__)     (__) 
`)
}

// Print ng help and exit
func Help() {
	fmt.Println(`comgo is a comtrade file parser.

Usage:
	detail:  cg [-f] filepath [-d]
	parse:   cg [-f] filepath [-c] channel No.

options:
	-f	--file		 cfg file path
	-h	--help		 information about the commands
	-c	--channel	 channel No. to save
	-d	--detail	 provide analog channel names
	-v	--version	 print netgo version`)
	os.Exit(1)
}

// Print ng version and exit
func Version() {
	fmt.Println(`comgo version cg[0.0.1]`)
	os.Exit(1)
}
