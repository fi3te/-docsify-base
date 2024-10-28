package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const indexFilename = "index.html"

func main() {
	parameters := []parameter{
		newParam("v", "docsify version", "{{v}}", "latest"),
		newParam("t", "title string", "{{t}}", "Dokumentation"),
		newParam("l", "html lang param", "{{l}}", "de"),
		newParam("sp", "search placeholder", "{{sp}}", "Suche"),
		newParam("spe", "empty search result text", "{{spe}}", "Keine Ergebnisse gefunden"),
	}
	portFlag := flag.Int("p", 80, "port to listen on")
	flag.Parse()

	updatePlaceholders(indexFilename, parameters)

	port := *portFlag
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Printf("Starting server on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

type parameter struct {
	commandLineFlag string
	name            string
	filePlaceholder string
	defaultValue    string
	value           *string
}

func newParam(commandLineFlag, name, filePlaceholder, defaultValue string) parameter {
	valueFlag := flag.String(commandLineFlag, defaultValue, name)
	return parameter{
		commandLineFlag: commandLineFlag,
		name:            name,
		filePlaceholder: filePlaceholder,
		defaultValue:    defaultValue,
		value:           valueFlag,
	}
}

func updatePlaceholders(filename string, placeholders []parameter) {
	fileContent := readFile(filename)
	log.Println("Updating placeholders in file content...")
	for _, p := range placeholders {
		fileContent = strings.Replace(fileContent, p.filePlaceholder, *p.value, -1)
		log.Printf("Set value of '%v': %v\n", p.name, *p.value)
	}
	writeFile(filename, fileContent)
}

func readFile(filename string) string {
	log.Printf("Reading file '%s'...\n", filename)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func writeFile(filename string, content string) {
	log.Printf("Writing file '%s'...\n", filename)
	bytes := []byte(content)
	err := os.WriteFile(filename, bytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
