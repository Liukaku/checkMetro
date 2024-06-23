package metro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	structs "github.com/liukaku/checkMetro/cmd/util"
)

func ApiRequest(stopId string, apiUrl string) ([]byte, error){
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

func GetStops(meteroUrl string) []byte {
	apiKey := os.Getenv("API_KEY")
	numberOfStops := "1000"
	getUrl := fmt.Sprintf("%s%s", meteroUrl, numberOfStops)

	fmt.Println(getUrl)

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

	var fileStruct structs.Result
	
	json.Unmarshal(readFile, &fileStruct)

	backToJson, err := json.Marshal(fileStruct)
	if err != nil {
		panic(err)
	}

	os.WriteFile("stops.json", backToJson, 0644)
	return backToJson

}

func RequestHandler(stopId string, apiUrl string) (*bool, error){
	// func RequestHandler() (*bool, error){
		fmt.Println("Lambda started")
	

		
		apiResponse, err := ApiRequest(stopId, apiUrl)
		fmt.Println("fetched")
	
		if err != nil {
			panic(err)
		}
	
		location, err := handleApiResponse(apiResponse)
	
		fmt.Println(location)
	
		if err != nil {
			panic(err)
		}
	
		success := false
	
		if location == "Droylsden" {
			success = true
		}
	
		fmt.Println(success)
	
		return &success, nil
	}

	
func handleApiResponse(responseBody []byte) (string, error){
	var resString string
	var dat structs.MetroStop
	// var apiErr map[string]interface{}

	json.Unmarshal(responseBody, &resString)

	err := json.Unmarshal([]byte(resString), &dat)

	if err != nil {
		return "error", err
	} 
	fmt.Println(dat)
	return dat.StationLocation, nil
}
