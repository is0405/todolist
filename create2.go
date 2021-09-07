package main

import (
	//"net/url"
	"net/http"
	"fmt"
	"os"
	//"strings"
	"io/ioutil"
	"bytes"
	//"encoding/json"	
)


func main() {
	// c_url := "http://localhost:10001/recipes/3"
	// c_url := "https://ec2-3-23-42-55.us-east-2.compute.amazonaws.com/recipes"
	c_url := "http://localhost:10001/todo/login"

	jsonStr :=`{
      "mail": "is0405@eee.co.jp",
      "password": "aaaaa"
}`

	
	req, err := http.NewRequest( "POST", c_url, bytes.NewBuffer([]byte(jsonStr)) )
	// req, err := http.NewRequest( "GET", c_url, nil )
	// req, err := http.NewRequest( "PATCH", c_url, bytes.NewBuffer([]byte(jsonStr)) )
	// req, err := http.NewRequest( "DELETE", c_url, nil )
	
	if err != nil {
		fmt.Println( "Error:http" )
		fmt.Println( err )
		os.Exit( 0 )
	}

	//req.Header.Set( "Content-Type", form_type )
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvZmZpY2VfaWQiOjEsInJvbGUiOiJvZmZpY2UiLCJleHAiOjE2MDUzNTE5ODQsImlzcyI6ImNvbS5jYXJlY29uIn0.1PMXMq7dhtGBICJWKRzkf8eRvinZGkOAdtqfYx9zENI"
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", token)
	client := new( http.Client )
	
	resp, err := client.Do( req )

	if resp != nil {
		defer resp.Body.Close()
		var byteArray, _ = ioutil.ReadAll( resp.Body )
		fmt.Println( string( byteArray ) )
	} else {
		fmt.Println( err )
	}
}
