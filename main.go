package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*estructura
type consulta struct {
	peruana    string
	embarazada string
	hijos      string
	trabaja    string
	edad       string
	casada     string
	estudia    string
	seguro     string
	distrito   string
}
*/

func consultaPorcentaje(resp http.ResponseWriter, req *http.Request) {
	log.Println("Ingreso a Consulta")

	//1.-Forma de recuperar parametros de entrada
	peruana := req.FormValue("peruana")
	embarazada := req.FormValue("embarazada")
	hijos := req.FormValue("hijos")
	trabaja := req.FormValue("trabaja")
	edad := req.FormValue("edad")
	casada := req.FormValue("casada")
	estudia := req.FormValue("estudiav")
	seguro := req.FormValue("seguro")
	distrito := req.FormValue("distrito")

	log.Println(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito)

	resp.Header().Set("Content-Type", "application/json")

	//Respuesta de Porcentaje
	porcentaje := AlgorithmTree(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito) //Funcion
	//serializar
	jsonBytes, _ := json.MarshalIndent(porcentaje, "", " ")
	io.WriteString(resp, string(jsonBytes))
}

func mostrarInicio(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Contect-Type", "text/html")
	io.WriteString(resp, `
		<html>
		<body><h2>Mi API de Consulta</h2></body>
		</html>
	`)
}

func handleRequest() {
	//manejo de los contextos de los métodos
	http.HandleFunc("/home", mostrarInicio)
	http.HandleFunc("/consulta", consultaPorcentaje)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	handleRequest()
}
