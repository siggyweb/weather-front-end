package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)

type Weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Latitude       float64 `json:"lat"`
		Longitude      float64 `json:"lon"`
		Timezone       string  `json:"tz_id"`
		LocaltimeEpock int     `json:"localtime_epoch"`
		LocalTime      string  `json:"localtime"`
	}
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float32 `json:"temp_c"`
		Tempf       float32 `json:"temp_f"`
		IsDay       int     `json:"is_day"`
		Condition   struct {
			Description string `json:"text"`
			// Icon        string `json:"icon"`
			// Code        int    `json:"code"`
		}
	}
}

// display is used to recursively print response data struct.
func (w Weather) display(flag bool) ([]byte, error) {
	if flag {
		bytes, err := json.MarshalIndent(w, "", "\t")
		if err != nil {
			return nil, errors.New("could not marshall weather JSON")
		}
		return bytes, nil
	} else {
		str := fmt.Sprintf("The weather in %s is %s, with a temperature of %.01fC.\n", w.Location.Name, w.Current.Condition.Description, w.Current.TempC)
		return []byte(str), nil
	}
}

func GetWeather(l string, w *Weather) error {
	baseUrl := "https://api.weatherapi.com/v1/current.json"
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return fmt.Errorf("request creation failed, %v", err)
	}

	req.Header.Set("accept", "application/json")
	q := req.URL.Query()
	q.Add("q", l)
	q.Add("key", string("secret"))
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed, %v", err)
	}
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(w)

	return nil
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  // Your server-side functionality
  w := Weather{}
  GetWeather("London", &w)
  pretty, err := w.display(true)
	if err == nil {
		return string(pretty))
  } else {
    return "failed"
  } 
}

func main() {
  // Make the handler available for Remote Procedure Call
  lambda.Start(handler)
}
