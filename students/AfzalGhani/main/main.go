package main

import (
	"io/ioutil"
	"net/http"
)

func main() {

	mux := defaulServerMux()
	ans, err := ioutil.ReadFile(storage/data.yaml")

	//fmt.Printf("the value is %v", string(ans))

	if err != nil {
		panic(err)
	}
	maphandler, _ := YAMLHandler(ans, mux)

	http.ListenAndServe(":8181", maphandler)

}

func defaulServerMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
	defer r.Body.Close()
}
