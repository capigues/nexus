package main

import (
	"fmt"
	"os"
)

// func debug(format string, v ...interface{}) {
// 	if settings.Debug {
// 		format = fmt.Sprintf("[debug] %s\n", format)
// 		log.Output(2, fmt.Sprintf(format, v...))
// 	}
// }

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

func verifyNexusFolderExists() {
	_, err := os.Stat(os.Getenv("NEXUS_HOME"))

	if os.IsNotExist(err) {
		err := os.Mkdir(os.Getenv("NEXUS_HOME"), 0755)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Creating NEXUS dir")
		return
	}
}

func main() {
	// initializing env vars for saving api
	if err := os.Setenv("NEXUS_HOME", os.Getenv("HOME")+"/.nexus"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Setenv("NEXUS_SERVERS_PATH", os.Getenv("NEXUS_HOME")+"/servers.json")

	// verifying folder to hold info exists
	verifyNexusFolderExists()

	// loading existing saved apis from file
	servers := &ModelServers{}
	servers.Load()

	// initialize root 'nexus' command and subcommands
	cmd, err := NewRootCmd(servers, os.Stdout, os.Args[1:])
	if err != nil {
		warning("%+v", err)
		os.Exit(1)
	}

	// executing command from user
	if err := cmd.Execute(); err != nil {
		// debug("%+v", err)
		os.Exit(0)
	}
}

// GENERAL PLAN FOR IMPLEMENTING QUICKLY

//  1. FIGURE OUT STORAGE
//  	- Decide and find examples of how I will store the LLM APIs added to the cli
//  	- Write functions (and test functions) that can write LLM API struct, edit it, delete it and read it from the file

//  2. WRITE BOILERPLATE FUNCTIONS
// 		- Write the boilder plate function for adding, removing, and updating LLM APIs

//  3. NARROW DOWN ON ADD
// 		- nexus add
