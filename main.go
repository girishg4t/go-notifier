package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
)

var t = new(template.Template)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Start()
	// top, err := ioutil.ReadFile("./templates/jira.yaml")
	// if err != nil {
	// 	log.Printf("yamlFile.Get err   #%v ", err)
	// }

	// template.Must(t.New("top").Parse(string(top)))

	// b := []byte(`{"fields": {"issuetype":{"name":"Task"},"project":{"key": "ALT"}, "summary": "Making it work using yaml"}}`)
	// var m map[string]interface{}
	// json.Unmarshal(b, &m)

	// topBuf := new(bytes.Buffer)
	// _ = t.ExecuteTemplate(topBuf, "top", m)

	// yamlString := topBuf.String()
	// y := make(map[string]interface{})
	// err = yaml.Unmarshal([]byte(yamlString), &y)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// jiraBody, err := json.Marshal(&y)

	// if err != nil {
	// 	// do error check
	// 	fmt.Println(err)
	// 	return
	// }

	// var c Configuration

	// if err := json.Unmarshal(jiraBody, &c); err != nil {
	// 	panic(err)
	// }

	// SendNotification(c.Services[0].Meta)

}

func SendNotification(i map[string]interface{}) error {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	jiraBody, err := json.Marshal(&i)
	fmt.Println(err)
	req, err := http.NewRequest(http.MethodPost, "https://accurics.atlassian.net/rest/api/3/issue",
		bytes.NewBuffer(jiraBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+BasicAuth("cs-girish.talekar@accurics.com", os.Getenv("SLACK_TOKEN")))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		fmt.Println(buf)
		return errors.New("non-ok response returned from Jira")
	}
	return nil
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
