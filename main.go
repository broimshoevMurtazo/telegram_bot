package main

import (
	// "crypto/internal/edwards25519/field"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gin-gonic/gin"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovary()

	r.Run(":3434")
}

func Recovary() {
	ReadUser()
	botrezilt, err := tgbotapi.NewBotAPI("6678092429:AAEepBh0V4AsimWpj9tF20Ezu7f7xMW0Mq8")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("conectin Done")

	update := tgbotapi.NewUpdate(0)

	retriverezult, error := botrezilt.GetUpdatesChan(update)
	if error != nil {
		fmt.Printf("error: %v\n", error)
	}

	isEditingPassword := false

	for update := range retriverezult {
		if update.Message.IsCommand() {
			if update.Message.Command() == "reset" {
				isEditingPassword = true
				for _, item := range SignUpSlice {
					if item.TelegramLogin == update.Message.Chat.UserName {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "enter new password")
						botrezilt.Send(msg)
					}
				}
			}
		} else {
			if isEditingPassword {
				for _, item := range SignUpSlice {
					if item.TelegramLogin == update.Message.Chat.UserName {
						item.Password = update.Message.Text
						// newpassword := field.TelegramLogin
						// SignUpSlice[i].Password
					}
				}
			}
		}
	}
}

func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		Writer()
	}
}
func Writer() {
	maesheldata, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", maesheldata, 0644)
}
func ReadUser() {
	readbyte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readbyte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
