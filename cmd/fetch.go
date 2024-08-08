package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// FetchCommand maneja la ejecución de peticiones fetch
func FetchCommand(attributes string, successContent, errorContent string) string {
	attrMap := ParseAttributes(attributes)
	url := attrMap["url"]
	method := attrMap["method"]
	responseType := attrMap["responseType"]
	resultName := attrMap["result"]
	timerStr := attrMap["timer"]
	limit := 0
	if limitStr, exists := attrMap["limit"]; exists {
		limit, _ = strconv.Atoi(limitStr)
	}

	timer := 0
	if timerStr != "" {
		timer, _ = strconv.Atoi(timerStr)
	}

	var wg sync.WaitGroup
	var result interface{}
	var fetchErr error
	resultChan := make(chan interface{})
	errorChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			errorChan <- err
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorChan <- err
			return
		}

		if responseType == "json" {
			err := json.Unmarshal(body, &result)
			if err != nil {
				errorChan <- err
				return
			}
			if limit > 0 {
				// Limitar la cantidad de resultados
				switch v := result.(type) {
				case []interface{}:
					if len(v) > limit {
						result = v[:limit]
					}
				case map[string]interface{}:
					// Asumimos que queremos limitar la cantidad de elementos en un campo específico
					for key, value := range v {
						if array, ok := value.([]interface{}); ok {
							if len(array) > limit {
								v[key] = array[:limit]
							}
						}
					}
				}
			}
		} else {
			result = string(body)
		}

		resultChan <- result
	}()

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	select {
	case result = <-resultChan:
		Variables[resultName] = result
	case fetchErr = <-errorChan:
		return errorContent
	}

	if fetchErr != nil {
		return errorContent
	}

	// Generar el contenido HTML procesado
	successHTML := ProcessCustomTags(successContent)

	// Incluir el script de actualización si timer está presente
	if timer > 0 {
		successHTML += fmt.Sprintf(`
			<script>
				setInterval(function() {
					fetch("%s", { method: "%s" })
						.then(response => response.json())
						.then(data => {
							// Actualizar los elementos en la página
							document.querySelectorAll("[data-update='%s']").forEach(function(element) {
								var key = element.getAttribute('data-key');
								var value = data[0][key];
								element.innerText = value;
							});
						})
						.catch(error => {
							console.error('Error fetching data:', error);
						});
				}, %d);
			</script>
		`, url, method, resultName, timer*1000)
	}

	return successHTML
}
