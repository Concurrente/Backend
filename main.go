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

func enviar(consulta string) {
	con, _ := net.Dial("tcp", "localhost:9800")
	defer con.Close()
	fmt.Fprint(con, consulta)
}

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

	log.Println(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito)

	resp.Header().Set("Content-Type", "application/json")

	//Respuesta de Porcentaje
	enviar(peruana)
	enviar(embarazada)
	enviar(hijos)
	enviar(trabaja)
	enviar(edad)
	enviar(casada)
	enviar(estudia)
	enviar(seguro)
	enviar(distrito)
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
