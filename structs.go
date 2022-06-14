package itocpackage

//AMQconnection is struct
type AMQconnection struct {
	Protocol  string
	IPAddress string
	Port      string
}

//AMQqueue is struct
type AMQqueue struct {
	Exchange   string
	Queue      string
	Routing    string
	Source     string
	SMS        string
	Telegram   string
	Email      string
	Cinderella string
	TTL        int64
}

//AMQsender is struct
type AMQsender struct {
	SMSInfra   string
	SMSAlarm   string
	Telegram   string
	Email      string
	Cinderella string
}

//DBconnection is struct
type DBconnection struct {
	Protocol  string
	Server    string
	IPAddress string
	Port      string
	Username  string
	Password  string
	DBName    string
}

//TelegramKItoc
type TelegramKItoc struct {
	Logdir string
}

//SourceTarget is struct
type SourceTarget struct {
	Logdir          string
	IncidentLog     string
	ApplicationList string
}
