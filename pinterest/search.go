package pinterest

import (
"net/http"
"encoding/json"   
"fmt"
"io/ioutil"		
"strings"	
"github.com/PaulSonOfLars/gotgbot/v2"
"github.com/PaulSonOfLars/gotgbot/v2/ext"	
)


func searchPinterest(query string) ([]string, error) {
 url := fmt.Sprintf("https://horridapi.onrender.com/pinterest?query=%s", query)
 resp, err := http.Get(url)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 body, err := ioutil.ReadAll(resp.Body)
 if err != nil {
  return nil, err
 }

 var result struct {
  Data []struct {
   URL string json:"url"
  } json:"data"
 }
 err = json.Unmarshal(body, &result)
 if err != nil {
  return nil, err
 }

 var urls []string
 for _, item := range result.Data {
  urls = append(urls, item.URL)
 }

 return urls, nil
}


func FindImage(b *gotgbot.Bot, ctx *ext.Context) error {    
    message := ctx.Message
    query := strings.Replace(m, "/h", "", -1)
    urls, err := searchPinterest(query)
    if err != nil {       
      fmt.Println(err)
      message.Reply(b, "Image not found", &gotgbot.SendMessageOpts{})
      return nil       
    media := ([]gotgbot.InputMedia, 0)
    for _, fuck := range urls {
      media = append(media, gotgbot.InputMediaPhoto{
          Media: fuck,
      })
    b.SendMediaGroup(      
      ctx.EffectiveUser.Id,
      media,
      &gotgbot.SendMediaGroupOpts{},
    )
    return nil
   
