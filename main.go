package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	// Command-line flags.
	modePtr := flag.String("mode", "", "Mode: upload or download")
	filePtr := flag.String("file", "", "For upload: path to the file; for download: destination path")
	cidPtr := flag.String("cid", "", "For download: the CID of the file to retrieve")
	ipfsAPI := flag.String("api", "localhost:5001", "IPFS API endpoint (default localhost:5001)")

	flag.Parse()

	// Create a new IPFS shell connecting to the specified API endpoint.
	ipfsShell := shell.NewShell(*ipfsAPI)

	switch *modePtr {
	case "upload":
		if *filePtr == "" {
			log.Fatalf("Upload mode requires -file=<path-to-file>")
		}

		// Read file into memory.
		data, err := ioutil.ReadFile(*filePtr)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		// Upload file to IPFS (IPFS automatically handles chunking).
		cid, err := ipfsShell.Add(bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Failed to upload to IPFS: %v", err)
		}

		fmt.Printf("File uploaded to IPFS with CID: %s\n", cid)

	case "download":
		if *cidPtr == "" {
			log.Fatalf("Download mode requires -cid=<file-CID>")
		}
		if *filePtr == "" {
			log.Fatalf("Download mode requires -file=<destination-path>")
		}

		// Retrieve file from IPFS using the CID.
		reader, err := ipfsShell.Cat(*cidPtr)
		if err != nil {
			log.Fatalf("Failed to retrieve file from IPFS: %v", err)
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalf("Failed to read data from IPFS response: %v", err)
		}

		// Save the retrieved data to a local file.
		err = ioutil.WriteFile(*filePtr, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

		fmt.Printf("File downloaded from IPFS and saved as: %s\n", *filePtr)

	default:
		fmt.Println("Invalid mode. Use -mode=upload or -mode=download")
		flag.Usage()
		os.Exit(1)
	}
}
