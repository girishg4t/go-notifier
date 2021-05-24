package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

var t = new(template.Template)

func main() {
	top, err := ioutil.ReadFile("sample.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	b := []byte(`{"fields": {"issuetype":{"name":"Task"},"project":{"key": "ALT"}, "summary": "Making it work using yaml"}}`)
	var f map[string]interface{}
	json.Unmarshal(b, &f)

	template.Must(t.New("top").Parse(string(top)))

	topBuf := new(bytes.Buffer)
	_ = t.ExecuteTemplate(topBuf, "top", f)

	yamlString := topBuf.String()
	y := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(yamlString), &y)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	SendNotification(y["meta"])
}

func SendNotification(i interface{}) error {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	jiraBody, err := json.Marshal(&i)
	fmt.Println(err)
	req, err := http.NewRequest(http.MethodPost, "https://accurics.atlassian.net/rest/api/3/issue",
		bytes.NewBuffer(jiraBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+BasicAuth("cs-girish.talekar@accurics.com", os.Getenv("JIRA_TOKEN")))
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
