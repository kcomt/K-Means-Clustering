package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type PredictionNonJson struct {
	Edad               int     `json:"edad,omitempty"`
	Peso               float64 `json:"peso,omitempty"`
	Distrito           int     `json:"distrito,omitempty"`
	Tos                int     `json:"tos,omitempty"`
	Fiebre             int     `json:"fiebre,omitempty"`
	DificultadRespirar int     `json:"dificultadRespirar,omitempty"`
	PerdidaOlfato      int     `json:"perdidaOlfato,omitempty"`
	Enfermo            int     `json:"enfermo,omitempty"`
}

type ClusteredGroup struct {
	name string
	data []PredictionNonJson
}

type Clusters struct {
	score    float64
	clusters []ClusteredGroup
}

var arrayOfData = make([]PredictionNonJson, 0, 496)
var arrayOfClusters = make([]Clusters, 0, 10)

func load(i int, noCommasRow []string, dataBuffer chan PredictionNonJson) {
	enfermo, _ := strconv.Atoi(noCommasRow[7])
	if enfermo == 1 {
		edad, _ := strconv.Atoi(noCommasRow[0])
		peso, _ := strconv.ParseFloat(noCommasRow[1], 2)
		tos, _ := strconv.Atoi(noCommasRow[3])
		fiebre, _ := strconv.Atoi(noCommasRow[4])
		dificultadRespirar, _ := strconv.Atoi(noCommasRow[5])
		perdidaOlfato, _ := strconv.Atoi(noCommasRow[6])
		distrito := 0
		switch noCommasRow[2] {
		case "Callao":
			distrito = 1
		case "Ventanilla":
			distrito = 2
		case "Ate":
			distrito = 3
		case "Barranco":
			distrito = 4
		case "Chorrillos":
			distrito = 5
		case "Comas":
			distrito = 6
		case "Jesus Maria":
			distrito = 7
		case "La Molina":
			distrito = 8
		case "La Victoria":
			distrito = 9
		case "Lince":
			distrito = 10
		case "Los Olivos":
			distrito = 11
		case "Lurin":
			distrito = 12
		case "Magdalena del Mar":
			distrito = 13
		case "Miraflores":
			distrito = 14
		case "Pueblo Libre":
			distrito = 15
		case "Puente Piedra":
			distrito = 16
		case "Rimac":
			distrito = 17
		case "San Borja":
			distrito = 18
		case "San Isidro":
			distrito = 19
		case "San Juan de Lurigancho":
			distrito = 20
		case "San Martin de Porres":
			distrito = 21
		case "San Miguel":
			distrito = 22
		case "Santiago de Surco":
			distrito = 23
		case "Surquillo":
			distrito = 24
		case "Villa El Salvador":
			distrito = 25
		default:
			distrito = 0
		}
		data := PredictionNonJson{edad, peso, distrito, tos, fiebre, dificultadRespirar, perdidaOlfato, enfermo}
		dataBuffer <- data
	}
}

func train() {
	dataBuffer := make(chan PredictionNonJson, 496)
	csvfile, err := os.Open("covidPeruDataSet.csv")
	if err != nil {
		log.Fatalln("No se pudo abrir el archivo", err)
	}
	r := csv.NewReader(csvfile)
	for i := 0; i < 1000; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		noCommasRow := strings.Split(record[0], ";")
		go load(i, noCommasRow, dataBuffer)
	}
	for i := 0; i < 496; i++ {
		arrayOfData = append(arrayOfData, <-dataBuffer)
	}
}

func calculateDistance(predict PredictionNonJson, train PredictionNonJson) float64 {
	result1 := math.Pow((float64(predict.Edad) - float64(train.Edad)), 2)
	result2 := math.Pow((predict.Peso - train.Peso), 2)
	result3 := math.Pow((float64(predict.Distrito) - float64(train.Distrito)), 2)
	result4 := math.Pow((float64(predict.Tos) - float64(train.Tos)), 2)
	result5 := math.Pow((float64(predict.Fiebre) - float64(train.Fiebre)), 2)
	result6 := math.Pow((float64(predict.DificultadRespirar) - float64(train.DificultadRespirar)), 2)
	result7 := math.Pow((float64(predict.PerdidaOlfato) - float64(train.PerdidaOlfato)), 2)
	result := math.Sqrt(result1 + result2 + result3 + +result4 + result5 + result6 + result7)

	return result
}

func createClusters(k int) {
	arrayOfRootClusters := make([]PredictionNonJson, 0, 4)
	arrayOfGroups := make([]ClusteredGroup, 0, 4)

	variation := make([]float64, 0, 4)
	contOfDistances := make([]int, 0, 4)

	for i := 0; i < k; i++ {
		name := "Group " + strconv.Itoa(i)
		contOfDistances = append(contOfDistances, 0)
		variation = append(variation, 0.0)
		auxArray := make([]PredictionNonJson, 0, 150)
		auxGroup := ClusteredGroup{name, auxArray}
		arrayOfGroups = append(arrayOfGroups, auxGroup)
	}

	for i := 0; i < k; i++ {
		x := rand.Intn(495)
		rootCluster := arrayOfData[x]
		arrayOfRootClusters = append(arrayOfRootClusters, rootCluster)
	}

	for i := 0; i < 496; i++ {
		arrayOfDistances := make([]float64, 0, 4)
		for j := 0; j < 4; j++ {
			distance := calculateDistance(arrayOfData[i], arrayOfRootClusters[j])
			arrayOfDistances = append(arrayOfDistances, distance)
		}
		indexOfMin := 0
		min := arrayOfDistances[0]
		for j := 0; j < 4; j++ {
			if min > arrayOfDistances[j] {
				min = arrayOfDistances[j]
				indexOfMin = j
			}
		}
		arrayOfGroups[indexOfMin].data = append(arrayOfGroups[indexOfMin].data, arrayOfData[i])
		variation[indexOfMin] = variation[indexOfMin] + min
		contOfDistances[indexOfMin] = contOfDistances[indexOfMin] + 1
	}

	score := 0.0
	for i := 0; i < k; i++ {
		variation[i] = variation[i] / float64(contOfDistances[i])
		score = score + variation[i]
	}
	score = score / float64(k)

	cluster := Clusters{score, arrayOfGroups}
	arrayOfClusters = append(arrayOfClusters, cluster)
}

func findClustersNtimes() {
	N := 10
	k := 4
	for i := 0; i < N; i++ {
		createClusters(k)
	}

	min := arrayOfClusters[0]
	for i := 0; i < len(arrayOfClusters); i++ {
		if min.score > arrayOfClusters[i].score {
			min = arrayOfClusters[i]
		}
	}

	fmt.Println(min.clusters)
}
func main() {
	train()
	findClustersNtimes()
}
