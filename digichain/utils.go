package digichain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (d *DigiChainClient) PostCall(data []byte) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}
	// Create a request with POST method and set the request body
	req, err := http.NewRequest("POST", d.rpc, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}
	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (d *DigiChainClient) GetCall() ([]byte, error) {
	resp, err := http.Get(d.rpc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// func main() {
// 	// Example usage:
// 	postDataUrl := "https://example.com/api/post"
// 	postDataBody := []byte(`{"key": "value"}`)
// 	err := postData(postDataUrl, postDataBody)
// 	if err != nil {
// 		fmt.Println("Error posting data:", err)
// 	}

// 	getDataUrl := "https://example.com/api/get"
// 	responseData, err := getData(getDataUrl)
// 	if err != nil {
// 		fmt.Println("Error getting data:", err)
// 	} else {
// 		fmt.Println("Response data:", string(responseData))
// 	}
// }
