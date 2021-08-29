package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func sendLine(customerID string, sharedKey string, line string) error {
	// Basic info
	logType := "tracee"
	method := "POST"
	contentType := "application/json"
	resource := "/api/logs"
	contentLen := len(line)

	// Get current time in corect ISO format
	rfc1123date := strings.Replace(
		time.Now().UTC().Format(time.RFC1123),
		"UTC", "GMT", 1,
	)

	// Calculate Signature Hash
	sharedKeyBytes, err := base64.StdEncoding.DecodeString(sharedKey)
	if err != nil {
		return err
	}
	bytesToHash := []byte(
		fmt.Sprintf("%s\n%d\n%s\nx-ms-date:%s\n%s",
			method,
			contentLen,
			contentType,
			rfc1123date,
			resource,
		),
	)
	hash := hmac.New(sha256.New, []byte(sharedKeyBytes))
	hash.Write(bytesToHash)
	signatureHash := base64.StdEncoding.EncodeToString(
		hash.Sum(nil),
	)
	signature := fmt.Sprintf("SharedKey %s:%s", customerID, signatureHash)

	// Set HTTP Headers and URL
	url := fmt.Sprintf(
		"https://%s.ods.opinsights.azure.com/%s?api-version=2016-04-01",
		customerID,
		resource[1:],
	)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(line)))
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"content-type":  []string{contentType},
		"Authorization": []string{signature},
		"Log-Type":      []string{logType},
		"x-ms-date":     []string{rfc1123date},
	}

	// Send request and check result
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Printf("POST Successful at %s\n", rfc1123date)
	} else {
		fmt.Printf("ERORR status code: %d\n", resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
	}

	return nil
}

func main() {
	log.SetPrefix("azupload: ")
	log.SetFlags(0)

	// Get Azure environment variables
	customerID, ok := os.LookupEnv("WORKSPACE_ID")
	if !ok {
		fmt.Println("Missing environmanet variable WORKSPACE_ID")
		os.Exit(1)
	}
	sharedKey, ok := os.LookupEnv("WORKSPACE_SHARED_KEY")
	if !ok {
		fmt.Println("Missing environmanet variable WORKSPACE_SHARED_KEY")
		os.Exit(1)
	}

	// Parse every line from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := sendLine(customerID, sharedKey, scanner.Text())
		if err != nil {
			fmt.Println("ERROR: ", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
