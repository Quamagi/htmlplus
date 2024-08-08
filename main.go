package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"goweb/cmd"
)

// TemplateData holds the data for the HTML template
type TemplateData struct {
	Content string
}

// CommandHandler handles the execution of custom commands within custom HTML tags
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	readChan := make(chan string)
	processChan := make(chan string)
	errorChan := make(chan error)
	done := make(chan struct{})

	// Goroutine para leer el contenido del archivo
	wg.Add(1)
	go func() {
		defer wg.Done()
		content, err := ioutil.ReadFile("templates/index.html")
		if err != nil {
			errorChan <- err
			close(readChan)
			return
		}
		readChan <- string(content)
		close(readChan)
	}()

	// Goroutine para procesar el contenido
	wg.Add(1)
	go func() {
		defer wg.Done()
		content, ok := <-readChan
		if !ok {
			close(processChan)
			return
		}
		processedContent := cmd.ProcessCustomTags(content)
		processChan <- processedContent
		close(processChan)
	}()

	// Goroutine para escribir el contenido procesado
	wg.Add(1)
	go func() {
		defer wg.Done()
		processedContent, ok := <-processChan
		if !ok {
			return
		}
		tmpl, err := template.New("index").Parse(processedContent)
		if err != nil {
			errorChan <- err
			return
		}

		data := TemplateData{Content: processedContent}
		err = tmpl.Execute(w, data)
		if err != nil {
			errorChan <- err
		}
		close(done)
	}()

	// Goroutine para esperar a que todas las goroutines terminen y manejar errores
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Manejo de errores
	select {
	case err := <-errorChan:
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case <-done:
		// Todo ha terminado correctamente
	}

	// Print all variables before executing the template
	log.Println("Variables before executing the template:")
	for name, value := range cmd.Variables {
		log.Printf("%s: %v\n", name, value)
	}
}

func main() {
	http.HandleFunc("/", CommandHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
