package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
)

var n int

var chCont chan int //control de sincronizacion

var s []string

func main() {
	//valores a esperar
	n = 9
	chCont = make(chan int, 1) //canal asincrono
	chCont <- 0
	//nodo servidor
	ln, err := net.Listen("tcp", "localhost:9800")
	if err != nil {
		log.Printf("aqu5")
		panic(err)
	}
	defer ln.Close()

	//manejo de multiples conexiones
	for {
		con, err := ln.Accept()
		if err != nil {
			log.Printf("aqu55")
			panic(err)
		}
		go manejadorConexiones(con) //trabajar de forma concurrente
	}
}

func manejadorConexiones(con net.Conn) {
	//cada servicio tiene una logica
	//cada manejador que va atender una conexion entrante aplica esa logica
	//Aplicamos la lógica del servicio

	buffIn := bufio.NewReader(con)
	msg, _ := buffIn.ReadString('\n')
	msg = strings.TrimSpace(msg)
	defer con.Close()
	//num, _ := strconv.Atoi(msg)

	fmt.Printf("Llegó el valor %d\n", msg)

	//evaluamos la cantidad de numeros que van pasando por el servicio
	cont := <-chCont

	s = append(s, msg)
	//fmt.Println("Llegó el valor %d\n", s)
	//fmt.Println("Llegó el valor %d\n", cont)

	cont++
	//evaluar si llegó la todos los valores
	if cont == n {
		//fmt.Printf("El valor final es %d\n", msg)
		fmt.Println(s)
		proba := AlgorithmTree(s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8])
		probability := fmt.Sprintf("%f", proba)
		enviar_respuesta(probability)
		cont = 0
		s = nil
	}
	chCont <- cont //cont actualizado
}

func enviar_respuesta(value string) {
	con, err := net.Dial("tcp", "localhost:9801")
	if err != nil {
		log.Printf("aqu7")
		panic(err)
	}
	defer con.Close()
	fmt.Fprint(con, value)
}

//Decision tree -> return distric and probability
func AlgorithmTree(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito string) float64 {

	// Load dataset
	response, err := http.Get("https://raw.githubusercontent.com/Concurrente/Backend/main/dataset/casos_cem_2020_lima.csv")
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("aqu4")
		panic(err)
	}
	fileBytes := bytes.NewReader(body)

	data_cem, err := base.ParseCSVToInstancesFromReader(fileBytes, true)
	if err != nil {
		log.Printf("aqu4")
		panic(err)
	}

	// Discretise the dataset with Chi-Merge
	filt := filters.NewChiMergeFilter(data_cem, 0.999)
	for _, a := range base.NonClassFloatAttributes(data_cem) {
		filt.AddAttribute(a)
	}
	filt.Train()
	data_cemF := base.NewLazilyFilteredInstances(data_cem, filt)

	// Create a 80-20 training-test split
	trainData, testData := base.InstancesTrainTestSplit(data_cemF, 0.80)

	//
	// First up, use ID3
	//
	tree := trees.NewID3DecisionTree(0.6)
	// (Parameter controls train-prune split.)

	// Train the ID3 tree
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// Create a new, empty DenseInstances
	newInst := base.NewDenseInstances()

	// Create some Attributes
	attrs := base.ResolveAllAttributes(testData)

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for x, a := range attrs {
		i := a.GetAttribute()
		b := newInst.AddAttribute(i)
		newSpecs[x] = b
	}

	newInst.AddClassAttribute(newInst.AllAttributes()[8])
	//fmt.Println(newInst)
	newInst.Extend(1)

	//Values to evaluate

	//District to evaluate
	entry_dsitrtict := distrito

	//Values to evaluate
	newInst.Set(newSpecs[0], 0, newSpecs[0].GetAttribute().GetSysValFromString(peruana))
	newInst.Set(newSpecs[1], 0, newSpecs[1].GetAttribute().GetSysValFromString(edad))
	newInst.Set(newSpecs[2], 0, newSpecs[2].GetAttribute().GetSysValFromString(embarazada))
	newInst.Set(newSpecs[3], 0, newSpecs[3].GetAttribute().GetSysValFromString(hijos))
	newInst.Set(newSpecs[4], 0, newSpecs[4].GetAttribute().GetSysValFromString(casada))
	newInst.Set(newSpecs[5], 0, newSpecs[5].GetAttribute().GetSysValFromString(estudia))
	newInst.Set(newSpecs[6], 0, newSpecs[6].GetAttribute().GetSysValFromString(trabaja))
	newInst.Set(newSpecs[7], 0, newSpecs[7].GetAttribute().GetSysValFromString(seguro))

	///////////////////////////////////////////////////////////////

	// Generate predictions
	predictions, err := tree.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Evaluate prediction
	fmt.Println("ID3 Performance (information gain)")
	cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))

	//Predict probabilities
	predictions_proba, err := tree.PredictProba(newInst)
	if err != nil {
		panic(err)
	}
	fmt.Println(predictions_proba)

	var dist string
	var valuer float64

	//Evaluate if the district is in the array
	founded := false
	for _, v := range predictions_proba {
		if entry_dsitrtict == v.ClassValue {
			dist = v.ClassValue
			valuer = v.Probability
			fmt.Println("Distrito: ", dist)
			fmt.Println("Probabilidad: ", valuer)
			founded = true
		}
	}
	if founded == false {
		dist = entry_dsitrtict
		valuer = 0.00
		fmt.Println("Distrito: ", dist)
		fmt.Println("Probabilidad: ", valuer)
	}

	return valuer

}
