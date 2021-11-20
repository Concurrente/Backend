package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
)

//Decision tree -> return distric and probability
func AlgorithmTree(peruana, embarazada, hijos, trabaja, edad, casada, estudia, seguro, distrito string) (string, float64) {

	// Load dataset
	response, err := http.Get("https://raw.githubusercontent.com/Concurrente/Backend/main/dataset/casos_cem_2020_lima.csv")
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	fileBytes := bytes.NewReader(body)

	data_cem, err := base.ParseCSVToInstancesFromReader(fileBytes, true)

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

	return dist, valuer

}
