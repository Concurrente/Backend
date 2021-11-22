package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

var nodo string = "localhost:9800"
var nodo2 string = "localhost:9700"
var nodo3 string = "localhost:9600"

type ConsultaPorcentaje struct {
	Peruana    string `json:"peruana"`
	Embarazada string `json:"embarazada"`
	Hijos      string `json:"hijos"`
	Trabaja    string `json:"trabaja"`
	Edad       string `json:"edad"`
	Casada     string `json:"casada"`
	Estudia    string `json:"estudia"`
	Seguro     string `json:"seguro"`
	Distrito   string `json:"distrito"`
}

func enviar(consulta string) {
	con, err := net.Dial("tcp", nodo2)
	if err != nil {
		log.Printf("aqui10")
		panic(err)
	}
	defer con.Close()
	fmt.Fprint(con, consulta)
}

func consultaPorcentaje(res http.ResponseWriter, req *http.Request) {
	log.Println("Ingreso a Consulta")

	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	log.Println("Llamada al endpoint /knn")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	res.Header().Set("Access-Control-Expose-Headers", "Authorization")
	//1.-Forma de recuperar parametros de entrada
	//req.Header().Set("Content-type", "application/json")
	var consultaPorcentaje = ConsultaPorcentaje{}

	err := json.NewDecoder(req.Body).Decode(&consultaPorcentaje)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		peruana := "1"
		embarazada := "1"
		hijos := "0"
		trabaja := "1"
		edad := "25"
		casada := "1"
		estudia := "0"
		seguro := "1"
		distrito := "San Juan de Lurigancho"
	*/

	log.Println(consultaPorcentaje.Peruana, consultaPorcentaje.Embarazada, consultaPorcentaje.Hijos, consultaPorcentaje.Trabaja, consultaPorcentaje.Edad, consultaPorcentaje.Casada, consultaPorcentaje.Estudia, consultaPorcentaje.Seguro, consultaPorcentaje.Distrito)

	//res.Header().Set("Content-Type", "application/json")

	//Respuesta de Porcentaje
	enviar(consultaPorcentaje.Peruana)
	enviar(consultaPorcentaje.Embarazada)
	enviar(consultaPorcentaje.Hijos)
	enviar(consultaPorcentaje.Trabaja)
	enviar(consultaPorcentaje.Edad)
	enviar(consultaPorcentaje.Casada)
	enviar(consultaPorcentaje.Estudia)
	enviar(consultaPorcentaje.Seguro)
	enviar(consultaPorcentaje.Distrito)
	//porcentaje := AlgorithmTree(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito) //Funcion

	//poner en modo escucha, recepcion
	ln, err := net.Listen("tcp", "localhost:9801") //ln -> listen
	if err != nil {
		log.Printf("aqui")
		panic(err)
	}
	defer ln.Close()

	con2, err := ln.Accept() //acepta la conexion
	if err != nil {
		log.Printf("aqui2")
		panic(err)
	}
	defer con2.Close()

	bufferIn := bufio.NewReader(con2)

	probabilidad, _ := bufferIn.ReadString('\n')

	probabilidad = strings.TrimSpace(probabilidad)
	//fmt.Printf("La probabilidad es %d\n", msg)
	fmt.Println("Probabilidad: ", probabilidad)

	//serializar
	jsonBytes, err := json.MarshalIndent(probabilidad, "", " ")
	if err != nil {
		log.Printf("aqu4")
		panic(err)
	}
	io.WriteString(res, string(jsonBytes))

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
	//go server()
	handleRequest()
}
