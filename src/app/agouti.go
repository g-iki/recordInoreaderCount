package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
)

func GetInoreaderInfo() (string, string) {
	var param []string
	param = append(param, "--headless")
	param = append(param, "--disable-dev-shm-usage")
	param = append(param, "--no-sandbox")
	param = append(param, "--blink-settings=imagesEnabled=false")
	param = append(param, "--disable-gpu")

	/**
	chrome起動オプションに任意のディレクトリを--user-data-dirの値として設定して起動する
	→Inoreaderにログインする
	→下記--user-data-dirオプションに同じディレクトリを指定する
	*/
	prifileDir := os.Getenv("PROFILE_PATH")
	profileName := os.Getenv("PROFILE_NAME")
	param = append(param,
		"--user-data-dir="+prifileDir)
	param = append(param,
		"--profile-dir="+profileName)

	driver := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args",
			param,
		),
	)

	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()
	page, err := driver.NewPage()
	if err != nil {
		log.Printf("Failed to open page:%v", err)
		return "", ""
	}
	// Inoreaderのページを開く
	if err := page.Navigate("https://www.inoreader.com/"); err != nil {
		log.Printf("Failed to navigate:%v", err)
		return "", ""
	}

	time.Sleep(time.Second * time.Duration(5))
	readLater, _ := page.FindByID("unread_cnt_starred").Text()
	unRead, _ := page.FindByID("unread_cnt_all_items").Text()
	fmt.Println("readLater", readLater)
	fmt.Println("unRead", unRead)
	return readLater, unRead
}
