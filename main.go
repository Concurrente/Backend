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

type ConsultaBono struct {
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
	con, _ := net.Dial("tcp", "localhost:9800")
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
	var usuario = ConsultaBono{}
	peruana := req.FormValue("peruana")
	embarazada := req.FormValue("embarazada")
	hijos := req.FormValue("hijos")
	trabaja := req.FormValue("trabaja")
	edad := req.FormValue("edad")
	casada := req.FormValue("casada")
	estudia := req.FormValue("estudiav")
	seguro := req.FormValue("seguro")
	distrito := req.FormValue("distrito")

	usuario.Peruana = peruana
	usuario.Embarazada = embarazada
	usuario.Hijos = hijos
	usuario.Trabaja = trabaja
	usuario.Edad = edad
	usuario.Casada = casada
	usuario.Estudia = estudia
	usuario.Seguro = seguro
	usuario.Distrito = distrito

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

	log.Println(usuario.Peruana, usuario.Embarazada, usuario.Hijos, usuario.Trabaja, usuario.Edad, usuario.Casada, usuario.Estudia, usuario.Seguro, usuario.Distrito)

	//res.Header().Set("Content-Type", "application/json")

	//Respuesta de Porcentaje
	enviar(usuario.Peruana)
	enviar(usuario.Embarazada)
	enviar(usuario.Hijos)
	enviar(usuario.Trabaja)
	enviar(usuario.Edad)
	enviar(usuario.Casada)
	enviar(usuario.Estudia)
	enviar(usuario.Seguro)
	enviar(usuario.Distrito)
	//porcentaje := AlgorithmTree(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito) //Funcion

	//poner en modo escucha, recepcion
	ln, _ := net.Listen("tcp", "localhost:9801") //ln -> listen

	defer ln.Close()

	con2, _ := ln.Accept() //acepta la conexion

	defer con2.Close()

	bufferIn := bufio.NewReader(con2)

	probabilidad, _ := bufferIn.ReadString('\n')

	probabilidad = strings.TrimSpace(probabilidad)
	//fmt.Printf("La probabilidad es %d\n", msg)
	fmt.Println("Probabilidad: ", probabilidad)

	//serializar
	jsonBytes, _ := json.MarshalIndent(probabilidad, "", " ")
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
	handleRequest()
}
