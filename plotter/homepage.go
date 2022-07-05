package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func homepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Serving homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}
