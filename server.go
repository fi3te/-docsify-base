package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const templateFilename = "index.html.tmpl"
const indexFilename = "index.html"

func main() {
	maxNameLength := 3
	parameters := []parameter{
		newStringParam("v", "latest", "docsify version"),
		newStringParam("t", "Dokumentation", "title string"),
		newStringParam("l", "de", "html lang param"),
		newStringParam("sp", "Suche", "search placeholder"),
		newStringParam("spe", "Keine Ergebnisse gefunden", "empty search result text"),
		newIntParam("sie", 3600000, "search index expiration time in millis"),
		newStringParam("sin", "docsify-base-namespace", "search index namespace"),
	}
	portFlag := flag.Int("p", 80, "port to listen on")
	log.Println("Reading parameters...")
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("%-*s = %s\n", maxNameLength, f.Name, f.Value.String())
	})

	log.Printf("Generating %s...", indexFilename)
	generateHTML(parameters)

	port := *portFlag
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Printf("Starting server on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func generateHTML(parameters []parameter) {
	tp := template.Must(template.ParseFiles(templateFilename))

	outputFile, err := os.OpenFile(indexFilename, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	if err := tp.Execute(outputFile, toMap(parameters)); err != nil {
		log.Fatal(err)
	}
}

type parameter struct {
	name           string
	getValueString func() string
	description    string
}

func newStringParam(name, defaultValue, description string) parameter {
	valueFlag := flag.String(name, defaultValue, description)
	return parameter{
		name: name,
		getValueString: func() string {
			return *valueFlag
		},
		description: description,
	}
}

func newIntParam(name string, defaultValue int, description string) parameter {
	valueFlag := flag.Int(name, defaultValue, description)
	return parameter{
		name: name,
		getValueString: func() string {
			return strconv.Itoa(*valueFlag)
		},
		description: description,
	}
}

func toMap(parameters []parameter) map[string]string {
	m := make(map[string]string)
	for _, p := range parameters {
		m[p.name] = p.getValueString()
	}
	return m
}
