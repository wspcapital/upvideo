package keywordtool

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


func (this *Service) GetInstantSuggestions(keyword string) (uniq []string, err error) {

	TIMEOUT := time.Second * 5
	client  := new(http.Client)
	um      := make(map[string]struct{}) // unique map

	client.Timeout = TIMEOUT

	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		if err = getRequest(client, keyword, string(r), um); err != nil {
			return
		}
	}

	uniq = make([]string, 0, len(um))
	for t := range um {
		uniq = append(uniq, t)
	}

	return

}

func setRequestHeaders(req *http.Request) {
	// headers
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("accept-encoding", "gzip, deflate, sdch")
	req.Header.Set("accept-language", "en-EN,en;q=0.8")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64)"+
		" AppleWebKit/537.36 (KHTML, like Gecko)"+
		" Chrome/50.0.2661.75 Safari/537.36")
	req.Header.Set("accept", "*/*")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("authority", "www.google.co.uk")
	req.Header.Set("referer", "https://www.google.co.uk/")
}

func addRequestURLQueryParams(req *http.Request, keyword, line string) {
	var q = req.URL.Query()
	q.Add("client", "psy-ab")
	q.Add("hl", "en")
	q.Add("gs_rn", "64")
	q.Add("gs_ri", "psy-ab")
	q.Add("cp", "5")
	q.Add("gs_id", "1h")
	q.Add("q", keyword+" "+line)
	q.Add("xhr", "t")
	req.URL.RawQuery = q.Encode()
}

func getRequest( client *http.Client, keyword string, line string, um map[string]struct{}, ) ( err error) {

	const URL = "https://www.google.co.uk/complete/search"

	var get *http.Request
	if get, err = http.NewRequest(http.MethodGet, URL, nil); err != nil {
		return
	}

	setRequestHeaders(get)
	addRequestURLQueryParams(get, keyword, line)

	var resp *http.Response
	if resp, err = client.Do(get); err != nil {
		return
	}
	defer resp.Body.Close() // close body

	var rc io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var r io.Reader
		if r, err = gzip.NewReader(resp.Body); err != nil {
			return
		}
		rc = ioutil.NopCloser(r)
	case "deflate":
		rc = flate.NewReader(resp.Body)
	default:
		rc = resp.Body
	}
	defer rc.Close() // close decompresser

	var jsonResp []interface{}
	if err = json.NewDecoder(rc).Decode(&jsonResp); err != nil {
		return
	}

	if len(jsonResp) < 3 {
		err = fmt.Errorf("invalid response array length %d", len(jsonResp))
		return
	}

	var ary, ok = jsonResp[1].([]interface{})
	if ok == false {
		err = fmt.Errorf("invalid json response root type %T", jsonResp[1])
		return
	}

	for _, line := range ary {
		var ln, ok = line.([]interface{})
		if ok == false {
			err = fmt.Errorf("invalid json element type: %T", line)
			return
		}
		if len(ln) < 1 {
			err = fmt.Errorf("invalid element list length %d", len(ln))
			return
		}
		var text string
		if text, ok = ln[0].(string); ok == false {
			err = fmt.Errorf("expected string, got %T", ln[0])
			return
		}

		text = strings.Replace(text, "<b>", "", -1)
		text = strings.Replace(text, "</b>", "", -1)
		text = strings.Replace(text, "\n", " ", -1)

		um[text] = struct{}{}

	}

	return
}

