package pinterest

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func InlineSearch(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery.Query
	if query == "" {
		_, err := ctx.InlineQuery.Answer(b, []gotgbot.InlineQueryResult{
			gotgbot.InlineQueryResultArticle{Title: "Not Found"},
		})
		return err
  }
	return nil
}
