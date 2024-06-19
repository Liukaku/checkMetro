package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/liukaku/checkMetro/cmd/s3"
)

type MetroStop struct {
	StationLocation string `json:"StationLocation"`
	AtcoCode        string `json:"AtcoCode"`
	Direction       string `json:"Direction"`
	Dest0           string `json:"Dest0"`
}

type Result struct {
	ODataContext string      `json:"@odata.context"`
	Value        []MetroStop `json:"value"`
}

func RequestHandler(ctx context.Context,) (*bool, error){
// func RequestHandler() (*bool, error){
	fmt.Println("Lambda started")
	
	stopId := os.Getenv("STOP_ID")
	apiUrl := os.Getenv("API_URL")
	
	apiResponse, err := apiRequest(stopId, apiUrl)
	fmt.Println("fetched")

	if err != nil {
		handleError(err)
	}

	location, err := handleApiResponse(apiResponse)

	fmt.Println(location)

	if err != nil {
		handleError(err)
	}

	success := false

	if location == "Droylsden" {
		success = true
	}

	fmt.Println(success)

	return &success, nil
}

func apiRequest(stopId string, apiUrl string) ([]byte, error){
	postBodyJson, _ := json.Marshal(map[string]string {
		"stop": stopId,
	})
	
	postBodyBuffer := bytes.NewBuffer(postBodyJson)
	
	fmt.Println("fetching")
	response, err := http.Post(apiUrl, "application/json", postBodyBuffer)

	if err != nil {

		return []byte("frick"), err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return []byte("frick"), err
	}

	return responseBody, nil
}

func handleApiResponse(responseBody []byte) (string, error){
	var resString string
	var dat MetroStop
	// var apiErr map[string]interface{}

	json.Unmarshal(responseBody, &resString)

	err := json.Unmarshal([]byte(resString), &dat)

	if err != nil {
		return "error", err
	} 
	fmt.Println(dat)
	return dat.StationLocation, nil
}

func getStops() []byte {
	apiKey := os.Getenv("API_KEY")
	numberOfStops := "1000"
	getUrl := fmt.Sprintf("https://api.tfgm.com/odata/Metrolinks?$top=%s", numberOfStops)

	fmt.Println(getUrl)
	fmt.Println(apiKey)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", getUrl, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	res, err := client.Do(req)
	
	if err != nil {
		panic(err)
	}
	
	fileOpen := res.Body

	readFile, err := io.ReadAll(fileOpen)
	if err != nil {
		panic(err)
	}

	var fileStruct Result
	
	json.Unmarshal(readFile, &fileStruct)

	backToJson, err := json.Marshal(fileStruct)
	if err != nil {
		panic(err)
	}

	os.WriteFile("stops.json", backToJson, 0644)
	return backToJson

}

func handleError(err error) (error) {
	fmt.Println(err)
	return err
}

func main(){
	// lambda.Start(RequestHandler)
	err := godotenv.Load()

	if err != nil {
		handleError(err)
	}

	s3.LoadConfig()
	stopResults := getStops()
	s3.SaveToBucket("metro-stops", stopResults)
	// RequestHandler()
}