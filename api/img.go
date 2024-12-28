package api

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type PinterestResponse struct {
    Data []struct {
        URL string `json:"url"`
    } `json:"data"` 
}

func searchPinterest(query string) (PinterestResponse, error) {
    url := fmt.Sprintf("https://horridapi.onrender.com/pinterest?query=%s", query)
    resp, err := http.Get(url)
    if err != nil {
        return PinterestResponse{}, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return PinterestResponse{}, err
    }

    var result PinterestResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return PinterestResponse{}, err
    }

    return result, nil
}
