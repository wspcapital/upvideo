package rapidtags

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Service struct {
}

func (this *Service) GetTags(title string) (items []string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://rapidtags.io/api/index.php?tool=tag-generator&input="+url.QueryEscape(strings.TrimSpace(title)), nil)

	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.75 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://rapidtags.io/generator/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("AlexaToolbar-ALX_NS_PH", "AlexaToolbar/alx-4.0.3")
	req.Header.Set("Cache-Control", "no-cache")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Rapidtags service not reachable: \n%v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Rapidtags service respond with statusCode: \n%v\n", res.Status)
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			fmt.Printf("Response body gzip error: \n%v\n", err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = res.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Response body read error: \n%v\n", err)
		return nil, err
	}

	var jsonRes []interface{}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		fmt.Printf("Response body parse error: \n%v\n", err)
		return nil, err
	}

	for _, tag := range jsonRes[:5] {
		if t, ok := tag.(int); ok {
			items = append(items, strconv.Itoa(t))
		}
		if t, ok := tag.(string); ok {
			items = append(items, t)
		}
	}

	return
}

func NewService() *Service {
	return &Service{}
}
