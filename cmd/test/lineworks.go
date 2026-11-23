package main

import (
	"context"
	"fmt"
	"growlink/internal/application/usecase/lineworks_usecase"
	"growlink/internal/infrastructure/external/lineworks"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	setting, ok, err := lineworks.Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if !ok {
		os.Exit(2)
	}
	client, err := lineworks.NewClient(setting)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	uc := lineworks_usecase.New(client, nil)
	if err := uc.SendMessage(context.Background(), []string{"ryota.kashio@cap"}, "test"); err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}
	os.Exit(0)
}
