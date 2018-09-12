package main
//
import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"log"
	"fmt"
)
//
type CommandData struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
	author *discordgo.User
	channel string
	prefix  string
	guild *discordgo.Guild
}
func (data CommandData) LoadData(session *discordgo.Session, message *discordgo.MessageCreate) {
	rand.Seed(time.Now().UTC().UnixNano())
	data.session = session
	data.message = message
	var channel, _ = data.session.State.Channel(data.message.ChannelID)
	var guild, _ = data.session.State.Guild(channel.GuildID)
	data.guild = guild
	channel, guild = nil, nil
	data.author = message.Author
	data.channel = message.ChannelID
	data.prefix = configuration.Prefix
	data.checkCommand()
}
func (data CommandData) startswith(forcheck string) (bool) {
	var start = strings.Split(data.message.Content, " ")[0]
	if start ==  forcheck {
		return true
	} else {
		return false
	}
}
func (data CommandData) roll() {
	var elements = strings.Split(data.message.Content, " ")
	var first, second, result, elemLen = 0, 100, 0, len(elements)
	if elemLen == 2 {
		second, _ = strconv.Atoi(elements[1])
	} else
	if elemLen == 3 {
		first, _ = strconv.Atoi(elements[1])
		second, _ = strconv.Atoi(elements[2])
	}
	if first > second {
		first, second = second, first
	}
	if first == second {
		result = first
	} else {
		result = rand.Intn(second - first) + first
	}
	data.session.ChannelMessageSend(data.channel, "Result: " + strconv.Itoa(result))
}
func (data CommandData) top() {
	var res, QueryError = configuration.database.Query("SELECT username, points FROM users ORDER BY points DESC LIMIT 10")
	if QueryError != nil {
		log.Print(QueryError.Error())
	}
	var username string
	var points int
	var inline = false
	var counter = 0
	var fields []*discordgo.MessageEmbedField
	for res.Next() {
		res.Scan(&username, &points)
		counter += 1
		if counter > 1 {
			inline = true
		}
		var tmp = &discordgo.MessageEmbedField{
			fmt.Sprintf("%d - %s", counter, username),
			strconv.Itoa(points),
			inline,
		}
		fields = append(fields, tmp)
	}
	embed := &discordgo.MessageEmbed {
		Author: &discordgo.MessageEmbedAuthor{Name:data.author.Username},
		Color: 0x00ff00,
		Fields: fields,
	}
	data.session.ChannelMessageSendEmbed(data.channel, embed)
}
func (data CommandData) coins() {
	var res, QueryError = configuration.database.Query("SELECT discord_id, points FROM users ORDER BY points DESC")
	if QueryError != nil {
		log.Print(QueryError.Error())
	}
	var points int
	var id string
	var counter = 1
	for res.Next() {
		res.Scan(&id, &points)
		if id == data.author.ID{
			break
		}
		counter += 1
	}
	var fields = []*discordgo.MessageEmbedField {
		{
			"Баланс",
			fmt.Sprintf("%d DGC", points),
			true,
		},
		{
			"Место",
			fmt.Sprintf("%d место среди пользователей", counter),
			true,
		},
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{Name:data.author.Username},
		Color: 0x00ff00,
		Fields: fields,
	}
	data.session.ChannelMessageSendEmbed(data.channel, embed)
}
func (data CommandData) throw() {
	var target = data.message.Mentions
	if len(target) == 0 {
		return
	}
	var targetAllInfo, _ = data.session.GuildMember(data.guild.ID, target[0].ID)
	var authorAllInfo, _ = data.session.GuildMember(data.guild.ID, data.author.ID)
	var targetNick, authorNick = targetAllInfo.Nick, authorAllInfo.Nick
	if targetNick == "" {
		targetNick = target[0].Username
	}
	if authorNick == "" {
		authorNick = data.author.Username
	}
	var allEmoji = data.guild.Emojis
	var staticEmoji []*discordgo.Emoji
	for emoji := range allEmoji {
		if allEmoji[emoji].Animated == false {
			staticEmoji = append(staticEmoji, allEmoji[emoji])
		}
	}
	var targetEmoji = staticEmoji[rand.Intn(len(staticEmoji))]
	var emojiString = fmt.Sprintf("<:%s:%s>", targetEmoji.Name, targetEmoji.ID)
	data.session.ChannelMessageSend(data.channel, fmt.Sprintf("**%s** threw %s at **%s**", authorNick, emojiString, targetNick))
}
func (data CommandData) checkCommand() {
	var start = strings.Split(data.message.Content, " ")[0]
	switch start {
	case data.prefix + "roll":
		data.roll()
		break
	case data.prefix + "top":
		data.top()
		break
	case data.prefix + "coins":
		data.coins()
		break
	case data.prefix + "throw":
		data.throw()
		break
	}
}
