package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	metro "github.com/liukaku/checkMetro/cmd/api"
	"github.com/liukaku/checkMetro/cmd/s3"
	structs "github.com/liukaku/checkMetro/cmd/util"
)

func handleError(err error) (error) {
	fmt.Println(err)
	return err
}

func getStopFromArray(stops []structs.MetroStop, targetStop string)(structs.MetroStop, error){
	fmt.Println(len(stops))

	for i := 0; i < len(stops); i++ {
		if stops[i].StationLocation == targetStop {
			return stops[i], nil
		}
	}
	return structs.MetroStop{}, errors.New("stop not found")
}

func getAndSaveNewTestStop(stops structs.Result, targetStop string){
	stop, err := getStopFromArray(stops.Value, targetStop)

	if err != nil {
		panic(err)
	}

	fmt.Println(stop)

	stopAsBytes, err := json.Marshal(stop)

	if err != nil {
		panic(err)
	}

	s3.SaveToBucket("metro-stops", stopAsBytes, "targetStop.json")
}

func handleRequest()(bool, error){
	err := godotenv.Load()

	if err != nil {
		handleError(err)
	}

	apiUrl := os.Getenv("API_URL")
	metroUrl := os.Getenv("METRO_URL")
	emailKey := os.Getenv("EMAIL_KEY")
	emailTo := os.Getenv("EMAIL_TO")
	emailFrom := os.Getenv("EMAIL_FROM")
	sparkUrl := os.Getenv("SPARK_URL")

	s3.LoadConfig()

	// 1.fetch from s3
	// 2.try that id
	// 3.if true then return
	// 4.else request new ID & save it

	result := s3.GetFromBucket("metro-stops", "targetStop.json")
	readFile, _ := io.ReadAll(result)

	var fileToJson structs.MetroStop
	
	json.Unmarshal(readFile, &fileToJson)
	
	success, _ := metro.RequestHandler(strconv.Itoa(fileToJson.Id), apiUrl)

	if *success {
		fmt.Printf("success no need to get a new list")
		metro.SendEmail(*success, emailKey, emailTo, emailFrom, sparkUrl)
		return true, nil
	}

	// fetch stops return []bytes
	stopResults := metro.GetStops(metroUrl)

	// convert to json
	var stopsJson structs.Result
	json.Unmarshal(stopResults, &stopsJson)

	// determine the new stop struct and save to s3 as individual
	getAndSaveNewTestStop(stopsJson, "Droylsden")

	// save the new whol response to s3
	s3.SaveToBucket("metro-stops", stopResults, "metrostops.json")

	fmt.Println("oh lawdy we updating")

	metro.SendEmail(*success, emailKey, emailTo, emailFrom, sparkUrl)
	return true, nil
}


func main() {
	lambda.Start(handleRequest)
}