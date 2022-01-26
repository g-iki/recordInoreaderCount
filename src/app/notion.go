package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetNotionDatePageID() string {

	t := time.Now()
	strDate := t.Format("2006/01/02")
	// NOTES: リクエストボディがシンプルなのでとりあえずこの書き方をするが、長く複雑になるなら構造体定義する
	b := []byte(`
		{
			"filter": {
				"property": "日付",
				"title": {
					"equals": "` + strDate + `"
				}
			},
			"page_size": 1
		}
	`)

	dbID := os.Getenv("DATE_DB_ID")
	url := "https://api.notion.com/v1/databases/" + dbID + "/query"
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(b),
	)
	if err != nil {
		fmt.Printf("Error, " + err.Error())
		return ""
	}

	secret := os.Getenv("NOTION_SECRET")
	version := os.Getenv("NOTION_API_VERSION")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Notion-Version", version)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error, " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error, " + err.Error())
		return ""
	}

	data := new(NotionQueryResponse)

	if err := json.Unmarshal(byteArray, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return ""
	}

	if len(data.Results) == 1 {
		return data.Results[0].ID
	}
	return ""
}

func UpdateNotionDatePage(
	pageID string,
	target string,
	readLaterCnt string,
	unreadCnt string,
) error {
	var readLater, unread string
	switch target {
	case "1":
		readLater = "Inoreader後で読む_0時"
		unread = "Inoreader未読_0時"
	case "2":
		readLater = "Inoreader後で読む_8時"
		unread = "Inoreader未読_8時"
	case "3":
		readLater = "Inoreader後で読む_16時"
		unread = "Inoreader未読_16時"
	default:
		return errors.New("no target error")
	}

	b := []byte(`
		{
			"properties": {
				"` + readLater + `": {
					"number": ` + readLaterCnt + `
				},
				"` + unread + `": {
					"number": ` + unreadCnt + `
				}
			}
		}
	`)
	url := "https://api.notion.com/v1/pages/" + pageID
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("Error, " + err.Error())
		return err
	}

	secret := os.Getenv("NOTION_SECRET")
	version := os.Getenv("NOTION_API_VERSION")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Notion-Version", version)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error, " + err.Error())
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error, " + err.Error())
		return err
	}

	data := new(NotionUpdateResponse)
	if err := json.Unmarshal(byteArray, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return err
	}

	return nil
}

// JsonToGo (https://mholt.github.io/json-to-go/)による自動生成
// 今回はIDさえ取得できれば良い
type NotionQueryResponse struct {
	Object  string `json:"object"`
	Results []struct {
		Object         string    `json:"object"`
		ID             string    `json:"id"`
		CreatedTime    time.Time `json:"created_time"`
		LastEditedTime time.Time `json:"last_edited_time"`
		Parent         struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Archived   bool   `json:"archived"`
		URL        string `json:"url"`
		Properties struct {
			Recipes struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Relation []struct {
					ID string `json:"id"`
				} `json:"relation"`
			} `json:"Recipes"`
			CostOfNextTrip struct {
				ID      string `json:"id"`
				Type    string `json:"type"`
				Formula struct {
					Type   string `json:"type"`
					Number int    `json:"number"`
				} `json:"formula"`
			} `json:"Cost of next trip"`
			LastOrdered struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				Date struct {
					Start string      `json:"start"`
					End   interface{} `json:"end"`
				} `json:"date"`
			} `json:"Last ordered"`
			InStock struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Checkbox bool   `json:"checkbox"`
			} `json:"In stock"`
		} `json:"properties"`
	} `json:"results"`
	HasMore    bool        `json:"has_more"`
	NextCursor interface{} `json:"next_cursor"`
}

type NotionUpdateResponse struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	Parent         struct {
		Type       string `json:"type"`
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Icon struct {
		Type  string `json:"type"`
		Emoji string `json:"emoji"`
	} `json:"icon"`
	Cover struct {
		Type     string `json:"type"`
		External struct {
			URL string `json:"url"`
		} `json:"external"`
	} `json:"cover"`
	Archived   bool   `json:"archived"`
	URL        string `json:"url"`
	Properties struct {
		InStock struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Checkbox bool   `json:"checkbox"`
		} `json:"In stock"`
		Name struct {
			ID    string `json:"id"`
			Type  string `json:"type"`
			Title []struct {
				Type string `json:"type"`
				Text struct {
					Content string      `json:"content"`
					Link    interface{} `json:"link"`
				} `json:"text"`
				Annotations struct {
					Bold          bool   `json:"bold"`
					Italic        bool   `json:"italic"`
					Strikethrough bool   `json:"strikethrough"`
					Underline     bool   `json:"underline"`
					Code          bool   `json:"code"`
					Color         string `json:"color"`
				} `json:"annotations"`
				PlainText string      `json:"plain_text"`
				Href      interface{} `json:"href"`
			} `json:"title"`
		} `json:"Name"`
	} `json:"properties"`
}
