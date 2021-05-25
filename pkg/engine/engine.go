package engine

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/girishg4t/go-notifier/pkg/libs"
	jsoniter "github.com/json-iterator/go"
	gomail "gopkg.in/mail.v2"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

func Process(c libs.Configuration) {
	switch strings.ToLower(c.Channel) {
	case "http":
		err := sendNotification(c)
		log.Println(err)
	case "smtp":
		err := sendEmail(c)
		log.Println(err)
	}
}

func sendNotification(c libs.Configuration) error {
	jiraBody, _ := j.Marshal(c.Meta)
	req, err := http.NewRequest(http.MethodPost, c.Config.Url, bytes.NewBuffer(jiraBody))
	if err != nil {
		return err
	}
	if c.Config.Accept != "" {
		req.Header.Add("Accept", c.Config.Accept)
	}
	if c.Config.Authorization != "" {
		req.Header.Add("Authorization", c.Config.Authorization)
	}
	if c.Config.From != "" {
		req.Header.Add("From", c.Config.From)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Println(buf.String())
		return errors.New("not able to send the notification")
	}

	return nil
}

func sendEmail(c libs.Configuration) error {
	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", c.Meta["from"].(string))

	recepients := c.Meta["to"].([]interface{})

	//TODO: need to correct this it is just a hack
	to := strings.Split(recepients[0].(string), " ")
	addresses := make([]string, len(to))
	for i, receipient := range to {
		addresses[i] = receipient
	}

	// Set E-Mail receivers
	m.SetHeader("To", addresses...)

	// Set E-Mail subject
	m.SetHeader("Subject", c.Meta["subject"].(string))

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", c.Meta["text"].(string))

	port, _ := strconv.Atoi(c.Config.Port)
	// Settings for SMTP server
	d := gomail.NewDialer(c.Config.Host, port, c.Config.Username, c.Config.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return nil
}
