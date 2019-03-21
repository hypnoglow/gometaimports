package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

var (
	version      = "unknown"
	templateDir  = "templates"
	templateName = "index.html.tmpl"
)

func main() {
	var (
		port       int
		configPath string
	)
	flag.IntVar(&port, "port", 8080, "Port to start the application on")
	flag.StringVar(&configPath, "config", "", "Path to gometaimports config file")
	flag.Parse()

	if configPath == "" {
		flag.PrintDefaults()
		log.Fatal("config path cannot be empty")
	}

	templatePath := filepath.Join(templateDir, templateName)

	gen, err := NewGenerator(configPath, templatePath)
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	log.Printf("gometaimports version %v", version)

	log.Printf("Start gometaimports server on :%d", port)
	defer log.Print("Stop gometaimports server")

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), ImportsHandler{
		Generator: gen,
		Logger:    stdLogger{},
	})

	log.Println(err)
}
