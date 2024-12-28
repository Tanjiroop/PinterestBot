package pinterest

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
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

func FindImage(b *gotgbot.Bot, ctx *ext.Context) error {
    message := ctx.Message
    query := strings.Replace(message.Text, "/h", "", -1)
    urls, err := searchPinterest(query)
    if err != nil {
        fmt.Println(err)
        message.Reply(b, "Image not found", &gotgbot.SendMessageOpts{})
        return err
    }

    media := make([]gotgbot.InputMedia, 0)
    for _, item := range urls.Data { 
        media = append(media, gotgbot.InputMediaPhoto{
						Media: gotgbot.InputFileByURL(item.URL),
        })        
    }

    _, err = b.SendMediaGroup(
        ctx.EffectiveUser.Id,
        media,
        &gotgbot.SendMediaGroupOpts{},
    )
    return err 
}
