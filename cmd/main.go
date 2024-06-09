package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type MetroStop struct {
	StationLocation string `json:"StationLocation"`
	AtcoCode        string `json:"AtcoCode"`
	Direction       string `json:"Direction"`
	Dest0           string `json:"Dest0"`
}

type Result struct {
	Value        map[string]string 
}

func RequestHandler(ctx context.Context,) (*bool, error){
// func RequestHandler() (*bool, error){
	fmt.Println("Lambda started")
	
	err := godotenv.Load()

	if err != nil {
		handleError(err)
	}

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

func handleError(err error) (error) {
	fmt.Println(err)
	return err
}

func main(){
	lambda.Start(RequestHandler)
	// RequestHandler()
}