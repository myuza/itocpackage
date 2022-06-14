package itocpackage

import (
	"fmt"

	"github.com/spf13/viper"
)

//SetupConfigEnvirontment is function
func SetupConfigEnvirontment(configFile string) {
	//get information about configuration from toml file
	//assign to correlated structs
	viper.SetConfigName(configFile)
	//viper untuk release
	viper.AddConfigPath("/home/itoc/config") // optionally look for config in home with services as sub directory
	viper.AddConfigPath("/home/apps/config") // development look for config in home with services as sub directory
	//viper untuk development
	viper.AddConfigPath(".") // optionally look for config in the working directory
	// viper.AddConfigPath("../")  // optionally look for config in home with services as sub directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		fmt.Println("[ERROR] Error while reading configuration file")
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

//GetDBconnection is function
func GetDBconnection() *DBconnection {
	var dbConn DBconnection

	dbConn.Protocol = viper.Get("database.itoc.protocol").(string)
	dbConn.Server = viper.Get("database.itoc.server").(string)
	dbConn.IPAddress = viper.Get("database.itoc.ipaddress").(string)
	dbConn.Port = viper.Get("database.itoc.port").(string)
	dbConn.DBName = viper.Get("database.itoc.dbname").(string)
	dbConn.Username = viper.Get("database.itoc.username").(string)
	dbConn.Password = viper.Get("database.itoc.password").(string)

	return &dbConn
}

//GetAMQconnection is function
func GetAMQconnection() *AMQconnection {
	var amqConn AMQconnection

	amqConn.Protocol = viper.Get("api.activemq.protocol").(string)
	amqConn.IPAddress = viper.Get("api.activemq.ipaddress").(string)
	amqConn.Port = viper.Get("api.activemq.port").(string)

	return &amqConn
}

//GetRabbitMQConnection is function
func GetRabbitMQConnection() string {
	rmqProtocol := viper.Get("api.rabbitmq.protocol").(string)
	rmqIP := viper.Get("api.rabbitmq.ipaddress").(string)
	rmqPort := viper.Get("api.rabbitmq.port").(string)
	rmqUser := viper.Get("api.rabbitmq.username").(string)
	rmqPass := viper.Get("api.rabbitmq.password").(string)
	rabbitMQConn := rmqProtocol + "://" + rmqUser + ":" + rmqPass + "@" + rmqIP + ":" + rmqPort
	return rabbitMQConn
}

//GetAMQqueue is function
func GetAMQqueue() *AMQqueue {
	var queue AMQqueue

	queue.Exchange = viper.Get("services.escalationrule.exchange").(string)
	queue.Queue = viper.Get("services.escalationrule.queue").(string)
	queue.Routing = viper.Get("services.escalationrule.routing").(string)
	queue.Source = viper.Get("services.escalationrule.queuesource").(string)
	queue.Email = viper.Get("services.escalationrule.queueemail").(string)
	queue.SMS = viper.Get("services.escalationrule.queuesms").(string)
	queue.Telegram = viper.Get("services.escalationrule.queuetelegram").(string)
	queue.Cinderella = viper.Get("services.escalationrule.queuecinderella").(string)
	queue.TTL = viper.Get("services.escalationrule.x-message-ttl").(int64)

	return &queue
}

//GetAMQsender is function
func GetAMQsender() *AMQsender {
	var sender AMQsender

	sender.SMSInfra = viper.Get("services.escalationrule.sendersmsinfra").(string)
	sender.SMSAlarm = viper.Get("services.escalationrule.sendersmsalarm").(string)
	sender.Telegram = viper.Get("services.escalationrule.sendertelegram").(string)
	sender.Cinderella = viper.Get("services.escalationrule.sendercinderella").(string)
	sender.Email = viper.Get("services.escalationrule.senderemail").(string)

	return &sender
}

//TelegramKangITOC is function
func TelegramKangITOC() *TelegramKItoc {
	var telegram TelegramKItoc
	telegram.Logdir = viper.Get("services.telegram.kangitoc.logdir").(string)
	return &telegram
}
