package urlshort

import (
	"fmt"
	"net/http"
	"strings"
)

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

		fmt.Println(firstPath)

		url := pathsToUrls[firstPath]

		if firstPath == url {
			http.Redirect(w, request, url, 301)
		}

		fallback.ServeHTTP(w, request)
	})

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	return nil, nil
}
