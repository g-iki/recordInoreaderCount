package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	app "getInoreaderCount/src/app"

	"github.com/joho/godotenv"
)

func main() {
	exe, _ := os.Executable()
	path := filepath.Dir(exe)
	err := godotenv.Load(path + "/.env")
	fmt.Println(path + "/.env")
	if err != nil {
		fmt.Println("Failed to get environment value.")
		os.Exit(0)
	}

	f := flag.String("no", "0", "update target flag")
	flag.Parse()
	updateTarget := *f

	readLaterCnt, unReadCnt := app.GetInoreaderInfo()

	pageID := app.GetNotionDatePageID()
	if pageID == "" {
		fmt.Println("Failed to get Notion date page.")
		return
	}
	app.UpdateNotionDatePage(
		pageID,
		updateTarget,
		readLaterCnt,
		unReadCnt,
	)
}
