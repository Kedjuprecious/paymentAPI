package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"regexp"


	"github.com/joho/godotenv"
)


type Response struct {
	Reference string `json:"reference"`
	Ussd_code  string `json:"ussd_code"` 
}

type CheckStatus struct {
	Status string `json:"status"`
}

func main() {
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API key not found in .env file")
		return
	}

	// Print API Key to verify it's being loaded for debugging
	fmt.Println("API Key Loaded:", apiKey)

	// Get user input
	var number, amount, reference, description string

	for {
		fmt.Print("Enter your number:\n237")
		fmt.Scanln(&number)
		if isValidNumber(number) {
			break
		}
		fmt.Println("Invalid number. It must be exactly 9 digits and contain only numbers.")
	}
	number = "237" + number

	for {
		fmt.Println("Enter the amount:")
		fmt.Scanln(&amount)
		if isValidAmount(amount) {
			break
		}
		fmt.Println("Invalid amount. It must contain only numbers.")
	}

	fmt.Println("Enter a description:")
	fmt.Scanln(&description)

	fmt.Println("Enter a reference:")
	fmt.Scanln(&reference)

	// Calling POST function
	referenceID := PostRequest(number, reference, description, amount, apiKey)
	if referenceID == "" {
		fmt.Println("Transaction failed or no reference returned")
		return
	}

	// Wait for transaction to process
	time.Sleep(20 * time.Second)

	// Calling GET function
	GetStatus(apiKey, referenceID)
}

// Functions to check validity of number
func isValidNumber(number string) bool {
	match, _ := regexp.MatchString(`^\d{9}$`, number)
	return match
}

// Function to check validity of amount
func isValidAmount(amount string) bool {
	match, _ := regexp.MatchString(`^\d+$`, amount)
	return match
}

// Function to make a POST request to initiate a transaction
func PostRequest(number, reference, description, amount, apiKey string) string {
	payload := map[string]interface{}{
		"from":              number,
		"amount":            amount,
		"description":       description,
		"external_reference": reference,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return ""
	}

	url := "https://demo.campay.net/api/collect/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Print HTTP status code for debugging
	fmt.Println("HTTP Status Code:", resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	// Printing raw body for debugging
	fmt.Println("Raw Response:", string(body))

	// Parse JSON response
	var transaction Response
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return ""
	}

	
	fmt.Println("Reference:", transaction.Reference)
	fmt.Println("U_code:", transaction.Ussd_code) 

	return transaction.Reference // Return the reference
}

// Function to check transaction status
func GetStatus(apiKey, reference string) {
	client := &http.Client{}
	url1 := fmt.Sprintf("https://demo.campay.net/api/transaction/%s/", reference) // Insert reference here

	req1, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return
	}

	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Token "+apiKey)

	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the body response for degugging
	fmt.Println("Raw Status Response:", string(body1))

	var status CheckStatus
	err = json.Unmarshal(body1, &status)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println("Transaction Status:", status.Status)
}



// NOt using functions
// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// )


// type Response struct {
// 	Reference string `json:"reference"`
// 	Ussd_code string `json:"ussd_code"`

// }

// type Check_Status struct {
// 	Status string `json:"status"`
// } 


// func main () {
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("Error loadind .env file")
// 		return
// 	}

// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		fmt.Println("API key not found in .env file")
// 		return
// 	}

// 	var number, amount, reference, description string

	

// 	fmt.Print("Enter your number:\n 237")
//     fmt.Scanln(&number)
// 	number = "237" + number;
	
// 	fmt.Println("Enter the amount")
//     fmt.Scanln(&amount)

// 	fmt.Println("Enter a description")
// 	fmt.Scanln(&description)

// 	fmt.Println("Enter a references")
// 	fmt.Scanln(&reference)

// 	payload := map[string]interface{}{
// 		"from": number,
// 		"amount": amount,
// 		"description": description,
// 		"external_reference": reference,
// 	}
// 	jsonData, err := json.Marshal(payload)
// 	// fmt.Printf(string(jsonData))
// 	if err != nil {
// 		fmt.Println("Error encoding JSON:", err)
// 		return
// 	}

// 	// endpoint, the door through which we are passing trough to access the campay services - its address
// 	url := "https://demo.campay.net/api/collect/"
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("Enter creating request:", err)
// 		return
// 	}
//     // fmt.Println(req)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Token "+apiKey)
	
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return
// 	}

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response:", err)
// 		return
// 	}

// 	var transaction Response
// 	err = json.Unmarshal(body, &transaction)
// 	if err != nil {
// 		fmt.Println("Enter decoding JSON response:", err)
// 		return
// 	}

// 	fmt.Println("Reference:", transaction.Reference)
// 	fmt.Println("Ussd_code:", transaction.Ussd_code)


// 	// Writing code for GET request

	

// 	url1 := fmt.Sprintf("https://demo.campay.net/api/transaction/%s/",transaction.Reference)
// 	req1, err := http.NewRequest("GET", url1, nil)

// 	req1.Header.Set("Content-Type", "application/json")
// 	req1.Header.Set("Authorization", "Token "+apiKey)
	
// 	time.Sleep(20 * time.Second)
	
// 	resp1, err := client.Do(req1)
// 	if err != nil{
// 		fmt.Println("Error making GET request", err)
// 		return
// 	}
// 	defer resp1.Body.Close()

// 	body1, err := io.ReadAll(resp1.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response:", err)
// 		return
// 	}

// 	var status Check_Status
// 	err = json.Unmarshal(body1, &status)
// 	if err != nil {
// 		fmt.Println("Enter decoding JSON response:", err)
// 		return
// 	}

// 	fmt.Println("Status:", status.Status)


// }