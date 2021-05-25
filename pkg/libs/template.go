package libs

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

func GetParsedTemplate(path string, msg map[string]interface{}) (string, error) {
	var t = new(template.Template)
	temp, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return "", err
	}

	template.Must(t.New("temp").Parse(string(temp)))

	tempBuf := new(bytes.Buffer)
	_ = t.ExecuteTemplate(tempBuf, "temp", msg)

	return tempBuf.String(), nil

}
