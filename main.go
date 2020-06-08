//main file to handle a POST request with the payload and return a json formatted response
package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//parse json with struct
//structs for POST request
type request struct {
	Jobs jobs `json:"jobs"`
}

type jobs struct {
	BuildBaseAMI buildBaseAMI `json:"Build base AMI"`
}

type buildBaseAMI struct {
	Builds []builds `json:"Builds"`
}

type builds struct {
	RuntimeSeconds string `json:"runtime_seconds"`
	BuildDate      string `json:"build_date"`
	Result         string `json:"result"`
	Output         string `json:"output"`
}

//structs for response json
type response struct {
	Latest latest `json:"latest"`
}

type latest struct {
	BuildDate  string `json:"build_date"`
	AMIId      string `json:"ami_id"`
	CommitHash string `json:"commit_hash"`
}

//structs for error response
type errorResponse struct {
	Error errorMsg `json:"error"`
}

type errorMsg struct {
	ErrorMsg string `json:"error_msg"`
}

//function to handle error responses as json
func errorHandler(w http.ResponseWriter, msg string) {
	errorBody := &errorResponse{}
	w.WriteHeader(http.StatusBadRequest)
	errorBody.Error.ErrorMsg = msg
	retVal, _ := json.MarshalIndent(errorBody, "", "\t")
	w.Write(retVal)
}

//function to handle accepting post request and returning json response
func parseJSON(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
        //initialize request and response struct
        buildList := &request{}
        responseList := response{}

        //set response content-type header to json
        w.Header().Set("Content-Type", "application/json")

		//read the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, "Unable to process request body")
			return
		}
		//unmarshal the body containing json content
		err = json.Unmarshal(body, buildList)
		if err != nil {
			errorHandler(w, "Unable to parse json body")
			return
		}
		//for loop for BuildDate string to int conversion
		currentDate := int64(0)
		for _, v := range buildList.Jobs.BuildBaseAMI.Builds {
			buildDate, err := strconv.ParseInt(v.BuildDate, 10, 64)
			if err != nil {
				errorHandler(w, "Unable to parse build date")
				return
			}
			if buildDate > currentDate {
				currentDate = buildDate
				outputList := strings.Split(v.Output, " ")
				if len(outputList) != 4 {
					errorHandler(w, "Incorrect json output")
					return
				}
				responseList.Latest.BuildDate = v.BuildDate
				responseList.Latest.AMIId = outputList[2]
				responseList.Latest.CommitHash = outputList[3]
			}
		}
		//marshal the json content
		response, err := json.MarshalIndent(responseList, "", "\t")
		if err != nil {
			errorHandler(w, "Internal error - Unable to marshal json response")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorHandler(w,"Only POST request is allowed")
        return
	}
}

//main function using Chi to test locally
func main() {
	log.Println("Starting server")
	router := chi.NewRouter()
	router.Post("/builds", parseJSON)
	err := http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
