package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	UID			string	`json:"uiid"`
	ID 			int		`json:"id"`
}

type ResumeTemplate struct {
	Title		string		`json:"title"`
	Version		float64		`json:"version"`
	ID 			int			`json:"id"`
}

type resumeData struct {
	//ID 			int 		`json:"id"`
	//FirstName	string		`json:"first_name"`
	//LastName	string		`json:"last_name"`
	/*
	--- Personal Details ---
		wanted job title
		upload image
		first name
		last name
		email 
		phone (country code)
		country 
		city
		driving license
		address - Not recommended
		postal code - Not recommended
		date of birth
	
	--- Employment History --- [multiple] collection/ list/ slice
		Job Title 
		Employer
		start date
		end date
		City
		Job Description (Look into adding gramarly API, OpenAI feature here BTW)

	--- Education ---
		School
		Degree
		Start date
		End Date
		City
		Description ( gramarly ... )

	--- Websited & Social Links ---
		Label
		Link

	--- Skills ---
		Title ( Use Tagify to add skills, OpenAI Feature; genearate posible skills bastes job history  )

	--- Hobbies ---
		description

	--- Professional Summary ---
		description ( OpenAI Feature )

	--- Courses ---
	--- References ---
	--- Custom Section ---

	*/

}



type Payload struct {
	Model       		string `json:"model"`
	Prompt      		string `json:"prompt"`
	MaxTokens   		int    `json:"max_tokens"`
	Temperature 		float64`json:"temperature"`
}

type OpenaiResponse struct {
	ID					string 		`json:"id"`
	Object				string 		`json:"object"`
	Created 			int    		`json:"created"`
	Model 				string 		`json:"model"`
	Choices 			[]Choice	`json:"choices"`
	Usage 				Usage  		`json:"usage"`

}
type Choice struct {
	Text		  		string 	`json:"text"`
	Index 				int		`json:"index"`
	Logprobs	 		int 	`json:"logprobs"`
	FinishReason 		string	`json:"finish_reason"`
}
type Usage struct {
	PromptTokens		int 	`json:"prompt_tokens"`
	CompletionTokens	int 	`json:"completion_tokens"`
	TotalTokens 		int 	`json:"total_tokens"`
} 

func init(){
	error := godotenv.Load()
	if error != nil{
		log.Fatalln("Exception loading .env : ", error)
	}
}

func NewCompletionsRequest(writer http.ResponseWriter, request *http.Request) { 
	//apiOrg = os.Getenv("API_ORG")
	writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(http.StatusOK)
	apiKey := os.Getenv("API_KEY")
	var payloadRequest Payload
	error := json.NewDecoder(request.Body).Decode(&payloadRequest)
	if error != nil{
		log.Fatalln("An error occured :  ", error)
		return 
	}
	
	bearerToken := "Bearer"+" "+apiKey
	data := Payload{
		Model: "text-davinci-003",
		Prompt: payloadRequest.Prompt,
		MaxTokens: 300,
		Temperature: 0.7,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Exception: ", err)
		return
	}

	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
	if err != nil {
		log.Fatalln("Exception: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Exception: ", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	var result OpenaiResponse
	json.Unmarshal([]byte(responseBody), &result)
	json.NewEncoder(writer).Encode(result)
}

func main(){
	router := mux.NewRouter();

	router.HandleFunc("/api/v1/message", NewCompletionsRequest).Methods("POST")

	fmt.Println("Server started on port 8000: ")
	error := http.ListenAndServe(":8000", router)
	if error != nil {
		log.Fatalln("An error occured listening to server: ", error);
		return
	}
}


