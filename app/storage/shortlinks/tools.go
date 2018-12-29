package shortlinks

import (
    "net/http"
)

func CheckDisabledUrl(domain string) bool {

    resp, err := http.Get(domain)
    if err != nil {
        return false
    }

    //if StatusCode == 404 this result like link has bin disabled.
    
    if resp.StatusCode == 404 {
        return true
    }

    return false
}