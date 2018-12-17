package keywordtool

import (
	"bitbucket.org/marketingx/upvideo/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Service struct {
	APIParams *config.KeywordtoolParams
}

type keywordtoolRequest struct {
	ApiKey          string `json:"apikey"`
	Keyword         string `json:"keyword"`
	Country         string `json:"country"`
	Language        string `json:"language"`
	Metrics         string `json:"metrics"`
	MetricsLocation string `json:"metrics_location"`
	MetricsLanguage string `json:"metrics_language"`
	MetricsNetwork  string `json:"metrics_network"`
	MetricsCurrency string `json:"metrics_currency"`
	Output          string `json:"output"`
}

type keywordtoolResponse struct {
	Results map[string]interface{} `json:"results"`
}

func (this *Service) GetKeywords(keyword string) (items []string, err error) {
	jsonReq, err := json.Marshal(keywordtoolRequest{
		ApiKey:          this.APIParams.ApiKey,
		Keyword:         keyword,
		Country:         "US",
		Language:        "en",
		Metrics:         "true",
		MetricsLocation: "2840",
		MetricsLanguage: "en",
		MetricsNetwork:  "googlesearchnetwork",
		MetricsCurrency: "USD",
		Output:          "json",
	})

	client := &http.Client{}
	res, err := client.Post("https://api.keywordtool.io/v2/search/suggestions/google", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Printf("Keywordtool service not reachable: \n%v\n", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Keywordtool service respond with statusCode: \n%v\n", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Response body read error: \n%v\n", err)
		return nil, err
	}

	var jsonRes keywordtoolResponse
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		fmt.Printf("Response body parse error: \n%v\n", err)
		return nil, err
	}

	for _, value := range jsonRes.Results {
		for _, value2 := range value.([]interface{}) {
			data := value2.(map[string]interface{})
			if keyword, ok := data["string"]; ok {
				items = append(items, strings.TrimSpace(keyword.(string)))
			}
		}
	}

	return
}

func NewService(params *config.KeywordtoolParams) *Service {
	return &Service{APIParams: params}
}
