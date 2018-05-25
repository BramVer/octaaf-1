package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"octaaf/models"
	"octaaf/scrapers"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	humanize "github.com/dustin/go-humanize"
	"github.com/o1egl/govatar"
	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

func sendRoll(message *tgbotapi.Message) {
	rand.Seed(time.Now().UnixNano())
	roll := strconv.Itoa(rand.Intn(9999999999-1000000000) + 1000000000)
	points := [9]string{"👌 Dubs", "🙈 Trips", "😱 Quads", "🤣😂 Penta", "👌👌🤔🤔😂😂 Hexa", "🙊🙉🙈🐵 Septa", "🅱️Octa", "💯💯💯 El Niño"}
	var dubscount int8 = -1

	for i := len(roll) - 1; i > 0; i-- {
		if roll[i] == roll[i-1] {
			dubscount++
		} else {
			break
		}
	}

	if dubscount > -1 {
		roll = points[dubscount] + " " + roll
	}
	reply(message, roll)
}

func count(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.MessageID))
}

func whoami(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.From.ID))
}

func m8Ball(message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		reply(message, "Oi! You have to ask question hé 🖕")
		return
	}

	answers := [20]string{"👌 It is certain",
		"👌 It is decidedly so",
		"👌 Without a doubt",
		"👌 Yes definitely",
		"👌 You may rely on it",
		"👌 As I see it, yes",
		"👌 Most likely",
		"👌 Outlook good",
		"👌 Yes",
		"👌 Signs point to yes",
		"☝ Reply hazy try again",
		"☝ Ask again later",
		"☝ Better not tell you now",
		"☝ Cannot predict now",
		"☝ Concentrate and ask again",
		"🖕 Don't count on it",
		"🖕 My reply is no",
		"🖕 My sources say no",
		"🖕 Outlook not so good",
		"🖕 Very doubtful"}
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(19)
	msg := tgbotapi.NewMessage(message.Chat.ID, answers[roll])
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func sendBodegem(message *tgbotapi.Message) {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func where(message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := scrapers.GetLocation(argument)

	if !found {
		reply(message, "This place does not exist 🙈🙈🙈🤔🤔�")
		return
	}

	msg := tgbotapi.NewLocation(message.Chat.ID, location.Lat, location.Lng)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func what(message *tgbotapi.Message) {
	query := message.CommandArguments()
	resp, err := http.Get(fmt.Sprintf("https://api.duckduckgo.com/?q=%v&format=json&no_html=1&skip_disambig=1", query))
	if err != nil {
		reply(message, "Just what is this? 🤔")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reply(message, "Just what is this? 🤔")
		return
	}

	result := gjson.Get(string(body), "AbstractText").String()

	if len(result) == 0 {
		reply(message, fmt.Sprintf("What is this %v you speak of? 🤔", Markdown(query, mdbold)))
		return
	}

	reply(message, fmt.Sprintf("%v: %v", Markdown(query, mdbold), result))
}

func weather(message *tgbotapi.Message) {
	weather, found := scrapers.GetWeatherStatus(message.CommandArguments())
	if !found {
		reply(message, "No data found 🙈🙈🙈🤔🤔🤔")
	} else {
		reply(message, "*Weather:* "+weather)
	}
}

func sendAvatar(message *tgbotapi.Message) {
	img, err := govatar.GenerateFromUsername(govatar.MALE, message.From.UserName)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	bytes := tgbotapi.FileBytes{Name: "avatar.png", Bytes: buf.Bytes()}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func bol(message *tgbotapi.Message) {
	bolURL := "https://www.bol.com/nl/nieuwsbrieven.html?country=BE"
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	resp, _ := client.Get(bolURL)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	token := "bogusTokenValue"

	doc.Find(".newsletter_subscriptions input").Each(func(i int, node *goquery.Selection) {
		name, found := node.Attr("name")
		if found && name == "token" {
			token, _ = node.Attr("value")
		}
	})

	data := url.Values{
		"emailAddress":          {message.CommandArguments()},
		"subscribedNewsLetters": {"DAGAANBIEDINGEN", "SOFT_OPTIN", "HARD_OPTIN", "B2B"},
		"token":                 {token},
		"updateNewsletters":     {"Voorkeuren+opslaan"}}

	req, _ := http.NewRequest("POST", bolURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	client.Do(req)

	reply(message, fmt.Sprintf("Succesfully subscribed *%v* to the bol.com mailing lists!", message.CommandArguments()))
}

func search(message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		reply(message, "What do you expect me to do? 🤔🤔🤔🤔")
		return
	}

	url, found := scrapers.Search(message.CommandArguments(), message.Command() == "search_nsfw")

	if found {
		reply(message, MDEscape(url))
		return
	}

	reply(message, "I found nothing 😱😱😱")
}

func sendStallman(message *tgbotapi.Message) {

	image, err := scrapers.GetStallman()

	if err != nil {
		reply(message, "Stallman went bork? 🤔🤔🤔🤔")
		return
	}

	bytes := tgbotapi.FileBytes{Name: "stally.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func sendImage(message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		reply(message, fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
		return
	}

	images, err := scrapers.GetImages(message.CommandArguments(), message.Command() == "img_sfw")

	if err != nil {
		reply(message, "Something went wrong!")
		return
	}

	timeout := time.Duration(2 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	for _, url := range images {

		res, err := client.Get(url)

		if err != nil {
			continue
		}

		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)

		if err != nil {
			continue
		}

		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, e := Octaaf.Send(msg)

		if e == nil {
			return
		}
	}

	reply(message, "I did not find images for the query: `"+message.CommandArguments()+"`")
}

func xkcd(message *tgbotapi.Message) {
	image, err := scrapers.GetXKCD()

	if err != nil {
		reply(message, "Failed to parse XKCD image")
	}

	bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func doubt(message *tgbotapi.Message) {
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, "assets/doubt.jpg")
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func quote(message *tgbotapi.Message) {
	// Fetch a random quote
	if message.ReplyToMessage == nil {
		quote := models.Quote{}

		err := DB.Where("chat_id = ?", message.Chat.ID).Order("random()").Limit(1).First(&quote)

		if err != nil {
			reply(message, "Something went wrong while fetching a quote!")
			return
		}

		if quote == (models.Quote{}) {
			reply(message, "No quotes have been saved yet.")
			return
		}

		config := tgbotapi.ChatConfigWithUser{
			ChatID:             message.Chat.ID,
			SuperGroupUsername: "",
			UserID:             quote.UserID}

		user, userErr := Octaaf.GetChatMember(config)

		if userErr != nil {
			reply(message, quote.Quote)
		} else {
			msg := fmt.Sprintf("\"%v\"", Markdown(quote.Quote, mdquote))
			msg += fmt.Sprintf(" \n    ~@%v", MDEscape(user.User.UserName))
			reply(message, msg)
		}

		return
	}

	// Unable to store this quote
	if message.ReplyToMessage.Text == "" {
		reply(message, "No text found in the comment. Not saving the quote!")
		return
	}

	err := DB.Save(&models.Quote{
		Quote:  message.ReplyToMessage.Text,
		UserID: message.ReplyToMessage.From.ID,
		ChatID: message.Chat.ID})

	if err != nil {
		reply(message, "Unable to save the quote...")
		return
	}

	reply(message, "Quote successfully saved!")
}

func nextLaunch(message *tgbotapi.Message) {
	res, err := http.Get("https://launchlibrary.net/1.3/launch?next=5&mode=verbose")

	if err != nil {
		reply(message, "Unable to fetch launch data")
		return
	}

	defer res.Body.Close()

	launchJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		reply(message, "Unable to fetch launch data")
		return
	}

	launches := gjson.Get(string(launchJSON), "launches").Array()

	var msg = "*Next 5 launches:*"

	layout := "January 2, 2006 15:04:05 MST"

	for index, launch := range launches {
		log.Printf("Doing dink")
		whenStr := launch.Get("net").String()
		when, err := time.Parse(layout, whenStr)

		msg += fmt.Sprintf("\n*%v*: %v", index+1, MDEscape(launch.Get("name").String()))

		if err != nil {
			msg += fmt.Sprintf("\n	  %v", Markdown(whenStr, mdcursive))
		} else {
			msg += fmt.Sprintf("\n	  %v", Markdown(humanize.Time(when), mdcursive))
		}

		vods := launch.Get("vidURLs").Array()

		if len(vods) > 0 {
			msg += "\n    " + MDEscape(vods[0].String())
		}
	}

	log.Printf(msg)

	reply(message, msg)
}

func issues(message *tgbotapi.Message) {
	res, err := http.Get("https://api.github.com/repos/bartwillems/Octaaf/issues?state=open")

	if err != nil {
		reply(message, "Unable to fetch open issues")
		return
	}

	defer res.Body.Close()

	issuesJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		reply(message, "Unable to fetch open issues")
		return
	}

	issues := gjson.ParseBytes(issuesJSON)

	var msg = "*Octaaf issues:*"

	var count int

	issues.ForEach(func(key, value gjson.Result) bool {
		count++
		msg += fmt.Sprintf("\n*%v: %v*", count, MDEscape(value.Get("title").String()))
		msg += fmt.Sprintf("\n    *url:* %v", Markdown(value.Get("url").String(), mdcursive))
		msg += fmt.Sprintf("\n    *creator:* %v", Markdown(value.Get("user.login").String(), mdcursive))
		return true
	})

	reply(message, msg)
}

func mcaffee(message *tgbotapi.Message) {
	res, err := http.Get("https://diegorod.github.io/WillMcAfeeEatHisOwnDick/")

	var msg = "Will McAffee eat his down dick?"
	if err != nil {
		msg += "\nI don't know!"
		reply(message, msg)
		return
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		msg += "\nI don't know!"
		reply(message, msg)
		return
	}

	dickBool := doc.Find("#isDickOnTheMenu").First().Text()
	percDiff := doc.Find("#percDiff").First().Text()
	currBpi := doc.Find("#currBpi").First().Text()

	msg += dickBool
	msg += fmt.Sprintf(" The current bitcoin price index (%s) is %s", percDiff, currBpi)

	if dickBool == "Yes!" {
		msg += " below the current price target."
	} else {
		msg += " above the current price target."
	}

	reply(message, msg)
}

func kaliRank(message *tgbotapi.Message) {
	if message.Chat.ID != KaliID {
		reply(message, "You are not allowed!")
		return
	}

	kaliRank := []models.MessageCount{}
	err := DB.Order("diff DESC").Limit(5).All(&kaliRank)

	if err != nil {
		reply(message, "Unable to fetch the kali rankings")
		return
	}

	var msg = "*Kali rankings:*"
	for index, rank := range kaliRank {
		msg += fmt.Sprintf("\n`#%v:` *%v messages*   _~%v_", index+1, rank.Diff, rank.CreatedAt.Format("Monday, 2 January 2006"))
	}

	reply(message, msg)
}

func iasip(message *tgbotapi.Message) {
	server := "http://159.89.14.97:6969"

	res, err := http.Get(server)
	if err != nil {
		reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
		return
	}

	reply(message, string(body))
}
