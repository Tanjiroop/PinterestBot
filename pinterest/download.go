package pinterest

import (	
	"regexp"
	"strings"	
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func ExtractURL(message string) string {
    pattern := regexp.MustCompile(`https?://\S+`)
    match := pattern.FindString(message)
    return match
}

func DownloadSend(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	chk := message.Text
	if strings.HasPrefix(chk, "/") {
		return nil
	}
	pattern := regexp.MustCompile(`https://pin\.it/(?P<url>[\w]+)`)
	if !pattern.MatchString(chk) {
		return nil
	}

	link, _ := ExtractURL(chk)
	photo := gotgbot.InputMediaPhoto{			
		Media: gotgbot.InputFileByURL(link),
	}
	_, err := b.SendPhoto(ctx.EffectiveChat.Id, photo.Media, &gotgbot.SendPhotoOpts{})
	if err != nil {
		message.Reply(b, "Failed to Send Photo", &gotgbot.SendMessageOpts{})
		return err
	}
	return nil
}
