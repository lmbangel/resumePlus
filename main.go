package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"os"
	"io/ioutil"
	"bytes"
	"github.com/joho/godotenv"
)

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl https://api.openai.com/v1/completions \
//   -H "Content-Type: application/json" \
//   -H "Authorization: Bearer $OPENAI_API_KEY" \
//   -d '{
//     "model": "text-davinci-003",
//     "prompt": "Say this is a test",
//     "max_tokens": 7,
//     "temperature": 0
//   }'

type Payload struct {
	Model       		string `json:"model"`
	Prompt      		string `json:"prompt"`
	MaxTokens   		int    `json:"max_tokens"`
	Temperature 		float64`json:"temperature"`
}

type openaiResponse struct {
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
	var result openaiResponse
	json.Unmarshal([]byte(responseBody), &result)
	//fmt.Println(result)
	json.NewEncoder(writer).Encode(result)
	//return string(responseBody)

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

