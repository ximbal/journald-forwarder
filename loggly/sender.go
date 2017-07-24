package loggly

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"encoding/json"
	"unicode/utf8"
)


type Response struct {
	Response string `json:"response"`

}

func clearUTF(entry string) string {

	 if !utf8.ValidString(entry) {
        v := make([]rune, 0, len(entry))
        for i, r := range entry {
            if r == utf8.RuneError {
                _, size := utf8.DecodeRuneInString(entry[i:])
                if size == 1 {
                    continue
                }
            }
            v = append(v, r)	
        }
		entry = string(v)
    }

	return entry
}

func GenerateUri(token string, tag string) string {

	uri := ""
	switch tag {
	case "": uri = fmt.Sprintf("https://logs-01.loggly.com/inputs/%s", token)
	default: uri = fmt.Sprintf("https://logs-01.loggly.com/inputs/%s/tag/%s/", token, tag)
	}

	log.Println(fmt.Sprintf("Will send data to %s", uri))

	return uri
}

func SendEvent(entry string, uri string) {
	

	resp, err := http.Post(uri, "application/x-www-form-urlencoded", ioutil.NopCloser(strings.NewReader(clearUTF(entry))) )

	if err != nil {
		log.Println("Unable to send data: %v", err)
		log.Println(entry)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to get a response: %v", err)
		return
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Invalid response: %v", err)
		return
	}
	if response.Response != "ok" {
		log.Println("Response was: %v", response.Response)
	}

}
