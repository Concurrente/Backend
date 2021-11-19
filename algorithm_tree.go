package main

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
)

func AlgorithmTree() {

	// Load dataset
	data_cem, err := base.ParseCSVToInstances("dataset/casos_cem_2020_lima.csv", true)
	if err != nil {
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
	fmt.Println(newInst)
	newInst.Extend(1)

	//Values to evaluate
	VICTIMA_PERUANA := "1"
	EDAD_VICTIMA := "25"
	VICTIMA_GESTANDO := "0"
	TIENE_HIJOS := "0"
	ESTADO_CIVIL_VICTIMA := "1"
	ESTUDIA := "0"
	TRABAJA_VICTIMA := "1"
	TIENE_SEGURO := "1"

	//District to evaluate
	entry_dsitrtict := "San Juan de Lurigancho"

	//Values to evaluate
	newInst.Set(newSpecs[0], 0, newSpecs[0].GetAttribute().GetSysValFromString(VICTIMA_PERUANA))
	newInst.Set(newSpecs[1], 0, newSpecs[1].GetAttribute().GetSysValFromString(EDAD_VICTIMA))
	newInst.Set(newSpecs[2], 0, newSpecs[2].GetAttribute().GetSysValFromString(VICTIMA_GESTANDO))
	newInst.Set(newSpecs[3], 0, newSpecs[3].GetAttribute().GetSysValFromString(TIENE_HIJOS))
	newInst.Set(newSpecs[4], 0, newSpecs[4].GetAttribute().GetSysValFromString(ESTADO_CIVIL_VICTIMA))
	newInst.Set(newSpecs[5], 0, newSpecs[5].GetAttribute().GetSysValFromString(ESTUDIA))
	newInst.Set(newSpecs[6], 0, newSpecs[6].GetAttribute().GetSysValFromString(TRABAJA_VICTIMA))
	newInst.Set(newSpecs[7], 0, newSpecs[7].GetAttribute().GetSysValFromString(TIENE_SEGURO))

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

	//Evaluate if the district is in the array
	founded := false
	for _, v := range predictions_proba {
		if entry_dsitrtict == v.ClassValue {
			dist := v.ClassValue
			valuer := v.Probability
			fmt.Println("Distrito: ", dist)
			fmt.Println("Probabilidad: ", valuer)
			founded = true
		}
	}
	if founded == false {
		dist := entry_dsitrtict
		valuer := 0
		fmt.Println("Distrito: ", dist)
		fmt.Println("Probabilidad: ", valuer)
	}

}
