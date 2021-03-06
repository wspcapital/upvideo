package main


import (
	"github.com/jochasinga/requests"
	"fmt"
	"encoding/json"
	"time"
)

type AutoGenerated []interface{}
var AsciiLowercase = [...]string{"a", "b", "c", "d", "e", "f", "g","h", "i", "j", "k", "l", "m", "n","o", "p", "q", "r", "s", "t", "u","v", "w", "x", "y", "z"}


func removeDuplicatesUnordered(elements []string) []string {
    encountered := map[string]bool{}

    // Create a map of all unique elements.
    for v:= range elements {
        encountered[elements[v]] = true
    }

    // Place all keys from the map into a slice.
    result := []string{}
    for key, _ := range encountered {
        result = append(result, key)
    }
    return result
}


func GoogleInstant(keyword string) (keywords []string) {

	// Set the Content-Type.
	setHeaders := func(r *requests.Request) {
	        r.Header.Add("pragma", "no-cache")
	        //r.Header.Add("accept-encoding", "gzip, deflate, sdch")
	        r.Header.Add("accept-language", "en-EN,en;q=0.8")
	        r.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.75 Safari/537.36")
	        r.Header.Add("accept", "*/*")
	        r.Header.Add("cache-control", "no-cache")
	        r.Header.Add("authority", "www.google.co.uk")
	        r.Header.Add("referer", "https://www.google.co.uk/")
	}


 	for _, letter := range AsciiLowercase {

		SetParams := func(r *requests.Request) {
		        r.Params.Add("client", "psy-ab")
		        r.Params.Add("hl", "en")
		        r.Params.Add("gs_rn", "64")
		        r.Params.Add("gs_ri", "psy-ab")
		        r.Params.Add("cp", "5")
		        r.Params.Add("gs_id", "1h")
		        r.Params.Add("q", fmt.Sprintf("%s %s", keyword, letter))
		        r.Params.Add("xhr", "t")
		}

		res, err := requests.Get("https://www.google.co.uk/complete/search", SetParams, setHeaders)
		if err != nil {
			continue
		}

		if res.StatusCode == 200 {
			var r []interface{}
			json.Unmarshal( res.Bytes() , &r)

			for _, key := range r[1:] {
				keywords = append(keywords, key.(string) )  // here need to append keywords
			}

			fmt.Println( res.String() )
			fmt.Println("\n")
		} else {
			fmt.Println("statuscode :", res.StatusCode)
		}

		time.Sleep(1 * time.Second)
	}

	keywords = removeDuplicatesUnordered(keywords)
	return keywords
}


func main() {
	keywords := GoogleInstant("love")
	for _, key := range keywords{
		fmt.Println(key)
	}
}