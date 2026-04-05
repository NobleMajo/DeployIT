package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	bunPath, err := getBunPath(wd)
	if err != nil {
		log.Fatal(err)
	}

	embedDir := filepath.Join(wd, "embed")

	st, err := os.Stat(bunPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fetchbun: notice: Bun CLI ready at %s (%d bytes)", bunPath, st.Size())

	entrypoint := filepath.Join(wd, "ts", "interface.ts")
	outfile := filepath.Join(embedDir, "interface")
	if effectiveGOOS() == "windows" {
		outfile = filepath.Join(embedDir, "interface.exe")
	}

	err = bunCompile(wd, bunPath, entrypoint, outfile, false)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fetchbun: notice: compiled %s -> %s", entrypoint, outfile)
}
