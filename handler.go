package urlshort

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type PathUrl []struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

// MapHandler devolverá un http.HandlerFunc (que también
// implementa http.Handler) que intentará asignar cualquier
// rutas (claves en el mapa) a su URL correspondiente (valores
// a la que apunta cada tecla del mapa, en formato de cadena).
// Si la ruta no se proporciona en el mapa, entonces la reserva
// Se llamará a http.Handler en su lugar.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		//http.Redirect(w, r, "http://www.google.com", 301)

		uri := request.RequestURI

		sliceuri := strings.Split(uri, "/")

		firstPath := "/" + sliceuri[1]

		url := pathsToUrls[firstPath]

		if url != "" {
			http.Redirect(w, request, url, 301)
		}

		fallback.ServeHTTP(w, request)
	})

}

func YAMLHandler(yaml string, fallback http.Handler) http.HandlerFunc {

	binaryYAML := getContentFile(yaml)

	yamlSliceMap := parseYAML(binaryYAML)

	pathsToUrls := buildMap(yamlSliceMap)

	return MapHandler(pathsToUrls, fallback)
}

func JSONHandler(jsonFile string, fallback http.Handler) http.HandlerFunc {

	binaryJSON := getContentFile(jsonFile)

	var j PathUrl

	err := json.Unmarshal(binaryJSON, &j)

	pathsToUrls := map[string]string{}

	for _, v := range j {
		pathsToUrls[v.Path] = v.URL
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return MapHandler(pathsToUrls, fallback)

}

func getContentFile(file string) []byte {
	fileRoute := "../" + file

	content, err := ioutil.ReadFile(fileRoute)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//fmt.Printf("%T %v \n", content, content)

	return content
}

func parseYAML(b []byte) []map[string]string {

	m := []map[string]string{}

	err := yaml.Unmarshal(b, &m)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return m

}

func buildMap(sm []map[string]string) map[string]string {

	var pathMap = map[string]string{}

	for _, m := range sm {
		pathMap[m["path"]] = m["url"]
	}

	return pathMap
}
