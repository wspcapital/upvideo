package checklinks

import (
    _"fmt"
    _"log"
    "net/http"
)

func CheckDisabledUrl(domain string) bool {

    resp, err := http.Get(domain)
    if err != nil {
        //log.Fatal(err)
        return false
    }

    // Print the HTTP Status Code and Status Name
    //fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

    if resp.StatusCode == 404 {
        return true
    }

    return false
}