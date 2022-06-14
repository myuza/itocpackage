package itocpackage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RequestLog struct {
	StartTime    string
	EndTime      string
	Username     string
	Fullname     string
	ChatroomType string
	RequestType  string
	Request      string
	MessageID    string
	Message      string
	BotAction    string
	BotResponse  string
	Status       string
	ResponseTime string
}

type BotConfig struct {
	BotApiKey  string `json:"botAPIKey"`
	ItocDB     string `json:"itocDB"`
	GrafanaUrl string `json:"grafanaUrl"`
	GrafanaKey string `json:"grafanaKey"`
	UimAPI     string `json:"uimAPI"`
	UimUser    string `json:"uimUser"`
	UimPass    string `json:"uimPass"`
	BotLog     string `json:"botLog"`
}

type BotSpeak struct {
	Speak string   `json:"speak"`
	Word  []string `json:"words"`
}

type CaptureGrafanaConfig struct {
	Title        string `json:"title"`
	Tags         string `json:"tags"`
	Caption      string `json:"caption"`
	DashboardUrl string `json:"dashboardUrl"`
	PanelId      string `json:"panelId"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	Timerange    int    `json:"timerange"`
}

type Config struct {
	BotConfig BotConfig
	BotSpeak  []BotSpeak
}

var config Config

func Initiate() *Config {
	//=============read bot config
	jsonFile, err := os.Open("botconfig.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var botConfig BotConfig
	json.Unmarshal(byteValue, &botConfig)

	//=============read speak config
	jsonFile, err = os.Open("wordingconfig.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)
	var botSpeak []BotSpeak
	json.Unmarshal(byteValue, &botSpeak)

	config.BotConfig = botConfig
	config.BotSpeak = botSpeak
	return &config
}

func GetBotConfig() *Config {
	return &config
}

func Speak(speak string) string {
	var word string
	speaklen := len(config.BotSpeak)
	for i := 0; i < speaklen; i++ {
		if config.BotSpeak[i].Speak == speak {
			wordList := config.BotSpeak[i].Word
			rand.Seed(time.Now().Unix())
			word = wordList[rand.Intn(len(wordList))]
			break
		}
	}
	return word
}

func BotHelp(botReply chan string, requestLog *RequestLog) {
	botReply <- Speak("greetings")
	botReply <- "Baiklah, Terima kasih telah menghubungi Kang ITOC. Akang senang membantu Bapak/Ibu " + requestLog.Username + ". Silahkan ketik /menu untuk mengetahui layanan yang Akang sediakan."
	close(botReply)
	requestLog.BotResponse = "bot_guidance"
}

func BotMenu(botReply chan string, botReplyMsgConf chan tgbotapi.MessageConfig, requestLog *RequestLog) {
	botReply <- Speak("greetings")
	botReply <- "Selamat Datang di Menu Utama bersama Kang ITOC"
	replymsg := "Saya bersama dengan Bapak/Ibu " + requestLog.Username + ", ada yang bisa dibantu?"
	botReply <- replymsg
	close(botReply)

	requestLog.BotResponse = "bot_menu"
}

func BotLocation(botReply chan string, requestLog *RequestLog, update tgbotapi.Update) {
	botReply <- Speak("greetings")
	time.Sleep(time.Second * 1)
	chatid := update.Message.Chat.ID
	chattype := update.Message.Chat.Type
	var chatdesc string
	if chattype == "private" {
		chatdesc = update.Message.Chat.FirstName + " " + update.Message.Chat.LastName
	} else {
		chatdesc = update.Message.Chat.Title
	}
	replymsg := "Saya berada di " + chattype + " bernama " + chatdesc + " dengan ChatID " + strconv.FormatInt(chatid, 10)
	botReply <- replymsg
	close(botReply)
	requestLog.BotResponse = "bot_location"
}

/*
func GetInfraMetrics(botReply chan string, requestLog *RequestLog, msgText string) {
	s := strings.Split(msgText, " ")
	if len(s) <= 2 {
		botReply <- "please provide more argument, example:\n- /infra [qos] [source] [target] \n- /infra [qos] [source]"
		requestLog.BotResponse = "argument_not_completed"
	} else if len(s) == 3 {
		botReply <- Speak("on_check")

		resultFound := false
		qos, source := s[1], s[2]

		listQos := []string{qos}
		listSource := []string{source}
		listTarget := []string{}
		if qos == "QOS_DISK_USAGE_PERC" {
			listTarget = []string{"--alltarget--"}
		} else {
			listTarget = []string{source}
		}

		graphResult := uimconn.GetQosValue(listSource, listQos, listTarget, "latest", "target")

		responseText := ""
		for _, aGraphResult := range graphResult {
			resultFound = true
			datapoints := aGraphResult.Datapoints
			target := aGraphResult.Target
			responseText = responseText + "\ntarget: " + target + "\nvalue: " + strconv.FormatFloat(datapoints[0][0], 'f', 2, 64) + "\n"
		}
		if resultFound {
			botReply <- responseText
			filePhoto := itocconn.CapturePanelInfra(qos, listSource, listTarget, 0, 0)
			caption := "[" + qos + "] last 6 hours trend"
			botReply <- "this_is_photo_msg||" + filePhoto + "||" + caption
			requestLog.BotResponse = "request_sent"
		} else {
			botReply <- Speak("request_not_found")
			requestLog.BotResponse = "request_not_found"
		}

	} else if len(s) == 4 {
		botReply <- Speak("on_check")
		resultFound := true
		qos, source, target := s[1], s[2], s[3]

		listQos := []string{qos}
		listSource := []string{source}
		listTarget := []string{}

		if qos == "QOS_DISK_USAGE_PERC" {
			listTarget = []string{target}
		} else {
			listTarget = []string{"--alltarget--"}
		}

		graphResult := uimconn.GetQosValue(listSource, listQos, listTarget, "latest", "target")

		responseText := qos + "\n"
		for _, aGraphResult := range graphResult {
			datapoints := aGraphResult.Datapoints
			target := aGraphResult.Target
			responseText = responseText + "\nTarget: " + target + "\nValue: " + strconv.FormatFloat(datapoints[0][0], 'f', 2, 64) + "\n"
		}
		if resultFound {
			botReply <- responseText
			requestLog.BotResponse = "request_sent"
		} else {
			requestLog.BotResponse = "request_not_found"
			botReply <- Speak("request_not_found")
		}
	}
	close(botReply)
}

func GetRechargeAmount(botReply chan string, requestLog *RequestLog, msgText string) {
	s := strings.Split(msgText, " ")
	if len(s) <= 1 {
		botReply <- Speak("on_check")
		filePhoto := itocconn.CapturePanelRecharge("today")
		caption := "[RECHARGE] Today"
		botReply <- "this_is_photo_msg||" + filePhoto + "||" + caption
		requestLog.BotResponse = "request_sent"
	}
	close(botReply)
}

func SystemHealthCheck(botReply chan string, requestLog *RequestLog, msgText string) {
	//=============read bot config
	jsonFile, err := os.Open("conf/healthcheck.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var captureGrafanaConfig []CaptureGrafanaConfig
	json.Unmarshal(byteValue, &captureGrafanaConfig)

	botReply <- Speak("on_check")
	for _, captureGrafana := range captureGrafanaConfig {
		filePhoto := itocconn.CapturePanel(captureGrafana.DashboardUrl, captureGrafana.PanelId, captureGrafana.Height, captureGrafana.Height, captureGrafana.Timerange)
		caption := "[" + captureGrafana.Title + "] " + captureGrafana.Caption
		botReply <- "this_is_photo_msg||" + filePhoto + "||" + caption
	}
	requestLog.BotResponse = "request_sent"

	close(botReply)
}

func CaptureHomeDashboard(botReply chan string, requestLog *RequestLog, msgText string) {
	botReply <- Speak("on_check")
	botReply <- "masuk capture"

	filePhoto := itocconn.CaptureWebpage("http://itoc.telkomsel.co.id/APM/dashboard_itoc/index_v2.php")
	caption := "[Success Rate] Home Dashboard"
	botReply <- "this_is_photo_msg||" + filePhoto + "||" + caption

	requestLog.BotResponse = "request_sent"

	close(botReply)
}

func HandlingReportRequest(botReplyMsgConf chan tgbotapi.MessageConfig, botReply chan string, requestLog *RequestLog, msgText string) {
	close(botReply)
	requestLog.BotResponse = "reportname_list"
	var reportList = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("SPM REPORT", "spm_report_list"),
		),
	)
	var msgConf tgbotapi.MessageConfig
	msgConf.ReplyMarkup = reportList
	msgConf.Text = "What kind of report?"
	botReplyMsgConf <- msgConf
	close(botReplyMsgConf)
}

func GetReportType(botReplyEditMsgConf chan tgbotapi.EditMessageReplyMarkupConfig, requestLog *RequestLog, reportname string) {
	switch reportname {
	case "event":
		requestLog.BotResponse = "report_event_type_list"
		var reportList = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("daily", "report_event_hourly_list"),
				tgbotapi.NewInlineKeyboardButtonData("hourly", "report_event_daily_list"),
			),
		)
		var msgConf tgbotapi.EditMessageReplyMarkupConfig
		msgConf.ReplyMarkup = &reportList

		botReplyEditMsgConf <- msgConf
		close(botReplyEditMsgConf)
	default:

	}
}

func GetSPMReportName(botReplyEditMsgConf chan tgbotapi.EditMessageReplyMarkupConfig, requestLog *RequestLog, reportname string) {
	switch reportname {
	case "event":
		requestLog.BotResponse = "report_event_type_list"
		var reportList = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("daily", "report_event_hourly_list"),
				tgbotapi.NewInlineKeyboardButtonData("hourly", "report_event_daily_list"),
			),
		)
		var msgConf tgbotapi.EditMessageReplyMarkupConfig
		msgConf.ReplyMarkup = &reportList

		botReplyEditMsgConf <- msgConf
		close(botReplyEditMsgConf)
	default:

	}
}

func DownloadFile(botReply chan string, botReplyDocConf chan tgbotapi.DocumentConfig, requestLog *RequestLog, command string) {
	//dl_spm_0001
	s := strings.Split(command, "_")
	docOwner := s[1]
	docId := s[2]
	switch docOwner {
	case "spm":
		botReply <- "this may take a while.. please wait"
		close(botReply)
		var docConf tgbotapi.DocumentConfig
		spmreport := itocconn.GetSPMReportByID(docId)
		fmt.Printf("%+v\n", spmreport)
		if &spmreport != nil {
			filePath := SaveFileFromWeb("pdf", spmreport.FileUrl, spmreport.FileName)
			docConf.File = filePath
			docConf.Caption = "happy reading:)"
		}
		botReplyDocConf <- docConf
		close(botReplyDocConf)
	}
}

func SaveFileFromWeb(fileformat string, url string, filename string) string {
	filePath := "doc/" + filename
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("GET", url, nil)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dowload From Web Finish")
	return filePath
}
*/
