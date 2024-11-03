package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const indexFilename = "index.html"

func main() {
	parameters := []parameter{
		newStringParam("v", "docsify version", "{{v}}", "latest"),
		newStringParam("t", "title string", "{{t}}", "Dokumentation"),
		newStringParam("l", "html lang param", "{{l}}", "de"),
		newStringParam("sp", "search placeholder", "{{sp}}", "Suche"),
		newStringParam("spe", "empty search result text", "{{spe}}", "Keine Ergebnisse gefunden"),
		newIntParam("sie", "search index expiration time in millis", "'{{sie}}'", 3600000),
		newStringParam("sin", "search index namespace", "{{sin}}", "docsify-base-namespace"),
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

type parameter interface {
	getName() string
	getFilePlaceholder() string
	getValue() string
}

type stringParameter struct {
	commandLineFlag string
	name            string
	filePlaceholder string
	defaultValue    string
	value           *string
}

func (p stringParameter) getName() string {
	return p.name
}

func (p stringParameter) getFilePlaceholder() string {
	return p.filePlaceholder
}

func (p stringParameter) getValue() string {
	return *p.value
}

func newStringParam(commandLineFlag, name, filePlaceholder, defaultValue string) stringParameter {
	valueFlag := flag.String(commandLineFlag, defaultValue, name)
	return stringParameter{
		commandLineFlag: commandLineFlag,
		name:            name,
		filePlaceholder: filePlaceholder,
		defaultValue:    defaultValue,
		value:           valueFlag,
	}
}

type intParameter struct {
	commandLineFlag string
	name            string
	filePlaceholder string
	defaultValue    int
	value           *int
}

func (p intParameter) getName() string {
	return p.name
}

func (p intParameter) getFilePlaceholder() string {
	return p.filePlaceholder
}

func (p intParameter) getValue() string {
	return strconv.Itoa(*p.value)
}

func newIntParam(commandLineFlag, name, filePlaceholder string, defaultValue int) intParameter {
	valueFlag := flag.Int(commandLineFlag, defaultValue, name)
	return intParameter{
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
		fileContent = strings.Replace(fileContent, p.getFilePlaceholder(), p.getValue(), -1)
		log.Printf("Set value of '%v': %v\n", p.getName(), p.getValue())
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
