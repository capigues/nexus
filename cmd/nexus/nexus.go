package main

import (
	"fmt"
	"log"
	"os"
)

var (
	NEXUS_FILE = "/.nexus"
)

func debug(format string, v ...interface{}) {
	// if settings.Debug {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
	// }
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

func main() {
	fmt.Println("hello world")
	servers := &ModelServers{}
	servers.Load(NEXUS_FILE + "/servers.json")
	fmt.Println(servers)
}

// func main() {
// 	servers := &ModelServers{}

// 	servers.Load(NEXUS_FILE + "/servers.json")

// 	cmd, err := NewRootCmd(servers, os.Stdout, os.Args[1:])
// 	if err != nil {
// 		warning("%+v", err)
// 		os.Exit(1)
// 	}

// 	if err := cmd.Execute(); err != nil {
// 		debug("%+v", err)
// 		os.Exit(1)
// 	}
// 	// nexus.Execute()
// }

// GENERAL PLAN FOR IMPLEMENTING QUICKLY

//  1. FIGURE OUT STORAGE
//  	- Decide and find examples of how I will store the LLM APIs added to the cli
//  	- Write functions (and test functions) that can write LLM API struct, edit it, delete it and read it from the file
//
//  2. WRITE BOILERPLATE FUNCTIONS
// 		- Write the boilder plate function for adding, removing, and updating LLM APIs
//
//  3. NARROW DOWN ON ADD
// 		- nexus add
//
//
