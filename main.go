package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/barnjamin/algorand-contract-interface/generator"
)

var (
	path = flag.String("path", "", "Path to the directory containing the manifest and contracts")
)

func main() {

	flag.Parse()

	m, err := generator.NewManifest(*path)
	if err != nil {
		log.Fatalf("Failed to parse manifest file: %+v", err)
	}

	ci, err := m.GenerateInterface()
	if err != nil {
		log.Fatalf("Failed to generate contract interface: %+v", err)
	}

	b, err := json.MarshalIndent(ci, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal contract interfacee: %+v", err)
	}

	if err := ioutil.WriteFile(*path+"/"+"asc.json", b, 0666); err != nil {
		log.Fatalf("Failed to write interface file: %+v", err)
	}

	log.Printf("Done")
}
