package api

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type BingResponse struct {
    Result []struct {
        URL string `json:"img"`
    } `json:"result"` 
}

// Exported function to handle web requests
func BingSearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    if query == "" {
        http.Error(w, "Query parameter is required", http.StatusBadRequest)
        return
    }

    result, err := SearchBing(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// Exported function to perform Bing search
func SearchBing(query string) (BingResponse, error) {
    url := fmt.Sprintf("https://horrid-api.vercel.app/images?page=7&query=%s", query)
    resp, err := http.Get(url)
    if err != nil {
        return BingResponse{}, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return BingResponse{}, err
    }

    var result BingResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return BingResponse{}, err
    }

    return result, nil
}

// Private function
func SearchBingInline(query string) (BingResponse, error) {
    url := fmt.Sprintf("https://horrid-api.vercel.app/images?page=40&query=%s", query)
    resp, err := http.Get(url)
    if err != nil {
        return BingResponse{}, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return BingResponse{}, err
    }

    var result BingResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return BingResponse{}, err
    }

    return result, nil
}
