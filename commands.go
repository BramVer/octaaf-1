package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"octaaf/models"
	"octaaf/scrapers"
	"strconv"
	"strings"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/go-redis/cache"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

func changelog(message *tgbotapi.Message) error {
	if settings.Version == "" {
		return reply(message, "Current version not found, check the changelog here: "+GitUri+"/tags")
	}

	return reply(message, fmt.Sprintf("%v/tags/%v", GitUri, settings.Version))
}

func all(message *tgbotapi.Message) error {
	members := Redis.SMembers(fmt.Sprintf("members_%v", message.Chat.ID)).Val()

	if len(members) == 0 {
		return reply(message, "I'm afraid I can't do that.")
	}

	// used to load the usernames in goroutines
	var wg sync.WaitGroup
	var response string
	// Get the members' usernames
	for _, member := range members {
		memberID, err := strconv.Atoi(member)

		if err != nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			user, err := getUsername(memberID, message.Chat.ID)
			if err == nil {
				response += fmt.Sprintf("@%v ", user.User.UserName)
			}
		}()
	}

	wg.Wait()
	return reply(message, MDEscape(fmt.Sprintf("%v %v", response, message.CommandArguments())))
}

func remind(message *tgbotapi.Message) error {
	errorMessage := "Malformed message, please send something like `/remind_me 1 hour het is de schuld van de sossen`"
	arr := strings.Split(message.Text, " ")

	if len(arr) < 4 {
		return reply(message, errorMessage)
	}

	delay, err := strconv.Atoi(arr[1])

	if err != nil || delay < 1 {
		return reply(message, errorMessage)
	}

	var deadline time.Time

	now := time.Now()

	unit := arr[2]

	switch unit {
	case "minute", "minutes":
		deadline = now.Add(time.Minute * time.Duration(delay)).UTC()
	case "hour", "hours":
		deadline = now.Add(time.Hour * time.Duration(delay)).UTC()
	case "day", "days":
		deadline = now.Add(time.Hour * time.Duration(delay) * 24).UTC()
	case "week", "weeks":
		deadline = now.Add(time.Hour * time.Duration(delay) * 24 * 7).UTC()
	default:
		return reply(message, "Unknown time format")
	}

	var remindMessage string

	// Parse the alertmessage from the string array to 1 string, skip stuff like "/alert 1 hour"
	for i := 3; i < len(arr); i++ {
		remindMessage += fmt.Sprintf("%v ", arr[i])
	}

	reminder := models.Reminder{
		ChatID:    message.Chat.ID,
		UserID:    message.From.ID,
		MessageID: message.MessageID,
		Message:   remindMessage,
		Deadline:  deadline,
		Executed:  false}

	go startReminder(reminder)
	return reply(message, "Reminder saved!")
}

func sendRoll(message *tgbotapi.Message) error {
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
	return reply(message, roll)
}

func count(message *tgbotapi.Message) error {
	return reply(message, fmt.Sprintf("%v", message.MessageID))
}

func whoami(message *tgbotapi.Message) error {
	return reply(message, fmt.Sprintf("%v", message.From.ID))
}

func m8Ball(message *tgbotapi.Message) error {

	if len(message.CommandArguments()) == 0 {
		return reply(message, "Oi! You have to ask question hé 🖕")
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
	return reply(message, answers[roll])
}

func sendBodegem(message *tgbotapi.Message) error {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	_, err := Octaaf.Send(msg)
	return err
}

func where(message *tgbotapi.Message) error {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := scrapers.GetLocation(argument)

	if !found {
		return reply(message, "This place does not exist 🙈🙈🙈🤔🤔�")
	}

	msg := tgbotapi.NewLocation(message.Chat.ID, location.Lat, location.Lng)
	msg.ReplyToMessageID = message.MessageID
	_, err := Octaaf.Send(msg)
	return err
}

func what(message *tgbotapi.Message) error {
	query := message.CommandArguments()
	resp, err := http.Get(fmt.Sprintf("https://api.duckduckgo.com/?q=%v&format=json&no_html=1&skip_disambig=1", query))
	if err != nil {
		return reply(message, "Just what is this? 🤔")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return reply(message, "Just what is this? 🤔")
	}

	result := gjson.Get(string(body), "AbstractText").String()

	if len(result) == 0 {
		return reply(message, fmt.Sprintf("What is this %v you speak of? 🤔", Markdown(query, mdbold)))
	}

	return reply(message, fmt.Sprintf("%v: %v", Markdown(query, mdbold), result))
}

func weather(message *tgbotapi.Message) error {
	weather, found := scrapers.GetWeatherStatus(message.CommandArguments())
	if !found {
		return reply(message, "No data found 🙈🙈🙈🤔🤔🤔")
	} else {
		return reply(message, "*Weather:* "+weather)
	}
}

func search(message *tgbotapi.Message) error {
	if len(message.CommandArguments()) == 0 {
		return reply(message, "What do you expect me to do? 🤔🤔🤔🤔")
	}

	url, found := scrapers.Search(message.CommandArguments(), message.Command() == "search_nsfw")

	if found {
		return reply(message, MDEscape(url))
	}

	return reply(message, "I found nothing 😱😱😱")
}

func sendStallman(message *tgbotapi.Message) error {

	image, err := scrapers.GetStallman()

	if err != nil {
		return reply(message, "Stallman went bork? 🤔🤔🤔🤔")
	}

	return reply(message, image)
}

func sendImage(message *tgbotapi.Message) error {
	var images []string
	var err error
	key := fmt.Sprintf("images_%v", message.Chat.ID)
	if message.Command() != "more" {
		if len(message.CommandArguments()) == 0 {
			return reply(message, fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
		}

		images, err = scrapers.GetImages(message.CommandArguments(), message.Command() == "img_sfw")
		if err != nil {
			return reply(message, "Something went wrong!")
		}

		Cache.Set(&cache.Item{
			Key:        key,
			Object:     images,
			Expiration: 0,
		})
	} else {
		if err := Cache.Get(key, &images); err != nil {
			return reply(message, "I can't fetch them for you right now.")
		}

		// Randomly order images for a different /more
		for i := range images {
			j := rand.Intn(i + 1)
			images[i], images[j] = images[j], images[i]
		}
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

		img, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Errorf("Unable to load image %v; error: ", url, err)
			continue
		}

		err = reply(message, img)

		if err == nil {
			return nil
		}
	}

	return reply(message, "I did not find images for the query: `"+message.CommandArguments()+"`")
}

func xkcd(message *tgbotapi.Message) error {
	image, err := scrapers.GetXKCD()

	if err != nil {
		return reply(message, "Failed to parse XKCD image")
	}

	return reply(message, image)
}

func doubt(message *tgbotapi.Message) error {
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, "assets/doubt.jpg")
	msg.ReplyToMessageID = message.MessageID
	_, err := Octaaf.Send(msg)
	if err != nil {
		log.Error("Unable to send /doubt: ", err)
	}
	return err
}

func quote(message *tgbotapi.Message) error {
	// Fetch a random quote
	if message.ReplyToMessage == nil {
		quote := models.Quote{}

		err := quote.Search(DB, message.Chat.ID, message.CommandArguments())

		if err != nil {
			log.Errorf("Quote fetch error: %v", err)
			return reply(message, "No quote found boi")
		}

		user, userErr := getUsername(quote.UserID, message.Chat.ID)

		if userErr != nil {
			log.Errorf("Unable to find the username for id '%v' : %v", quote.UserID, userErr)
			return reply(message, quote.Quote)
		} else {
			msg := fmt.Sprintf("\"%v\"", Markdown(quote.Quote, mdquote))
			msg += fmt.Sprintf(" \n    ~@%v", MDEscape(user.User.UserName))
			return reply(message, msg)
		}
	}

	// Unable to store this quote
	if message.ReplyToMessage.Text == "" {
		return reply(message, "No text found in the comment. Not saving the quote!")
	}

	err := DB.Save(&models.Quote{
		Quote:  message.ReplyToMessage.Text,
		UserID: message.ReplyToMessage.From.ID,
		ChatID: message.Chat.ID})

	if err != nil {
		log.Errorf("Unable to save quote '%v', error: %v", message.ReplyToMessage.Text, err)
		return reply(message, "Unable to save the quote...")
	}

	return reply(message, "Quote successfully saved!")
}

func nextLaunch(message *tgbotapi.Message) error {
	res, err := http.Get("https://launchlibrary.net/1.3/launch?next=5&mode=verbose")

	if err != nil {
		return reply(message, "Unable to fetch launch data")
	}

	defer res.Body.Close()

	launchJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return reply(message, "Unable to fetch launch data")
	}

	launches := gjson.Get(string(launchJSON), "launches").Array()

	var msg = "*Next 5 launches:*"

	layout := "January 2, 2006 15:04:05 MST"

	for index, launch := range launches {
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

	return reply(message, msg)
}

func issues(message *tgbotapi.Message) error {
	res, err := http.Get("https://api.github.com/repos/bartwillems/Octaaf/issues?state=open")

	if err != nil {
		return reply(message, "Unable to fetch open issues")
	}

	defer res.Body.Close()

	issuesJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return reply(message, "Unable to fetch open issues")
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

	return reply(message, msg)
}

func kaliRank(message *tgbotapi.Message) error {
	if message.Chat.ID != settings.Telegram.KaliID {
		return reply(message, "You are not allowed!")
	}

	kaliRank := []models.MessageCount{}
	err := DB.Order("diff DESC").Limit(5).All(&kaliRank)

	if err != nil {
		log.Error("Unable to fetch kali rankings: ", err)
		return reply(message, "Unable to fetch the kali rankings")
	}

	var msg = "*Kali rankings:*"
	for index, rank := range kaliRank {
		msg += fmt.Sprintf("\n`#%v:` *%v messages*   _~%v_", index+1, rank.Diff, rank.CreatedAt.Format("Monday, 2 January 2006"))
	}

	return reply(message, msg)
}

func iasip(message *tgbotapi.Message) error {
	server := "http://159.89.14.97:6969"

	res, err := http.Get(server)
	if err != nil {
		log.Error("Unable to fetch IASIP quote: ", err)
		return reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Unable to fetch IASIP quote: ", err)
		return reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
	}

	return reply(message, string(body))
}

func reported(message *tgbotapi.Message) error {
	if message.Chat.ID != settings.Telegram.KaliID {
		return reply(message, "Yeah well, you need to update to Strontbot Enterprise edition for Workgroups to use this feature.")
	}

	reportCount, err := DB.Count(models.Report{})

	if err != nil {
		log.Error("Report fetch error: ", err)
		return reply(message, "I can't seem to be able to count the reports.")
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             message.Chat.ID,
		SuperGroupUsername: "",
		UserID:             settings.Telegram.ReporterID}

	reporter, err := Octaaf.GetChatMember(config)

	if err != nil {
		return reply(message, fmt.Sprintf("So far, %v people have been reported by Dieter", reportCount))
	} else {
		return reply(message, MDEscape(fmt.Sprintf("So far, %v people have been reported by: @%v", reportCount, reporter.User.UserName)))
	}
}
