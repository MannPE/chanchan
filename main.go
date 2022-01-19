package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MannPE/chanchan/boardSettings"
)

type Thread struct {
	Id int	 `json:"no"`
	Comment string `json:"com"`
	Subtitle string  `json:"sub"`
}

type Page struct{
	Page int `json:"page"`
	Threads []Thread `json:"threads"`
}

type ImageRequest struct {
	board string
	keywords []string
	url string
	threadId int
}


const ThreadRefreshSeconds = 20
var ImageChannel chan ImageRequest


func main() {
	defer close(ImageChannel)
	res := boardSettings.GetSFW4chanBoards()
	ImageChannel = make(chan ImageRequest)

	for k, v := range res {
		fmt.Println(k, " -> ", v)
		go BoardWorker(k, v)
	}


	go handleImages()
	for {}
}

func handleImages() {
	for {
			newImgReq := <-ImageChannel
			fmt.Println("New Image:", newImgReq)

			urlParts :=  strings.Split(newImgReq.url, "/") 
			imgName := urlParts[len(urlParts)- 1]
			finalName := fmt.Sprintf("../downloads/%s/%s - %d - %s", newImgReq.keywords[0], newImgReq.board, newImgReq.threadId, imgName)

			_, err := os.ReadFile(finalName)
			if(err == nil ) {
				continue
			}

			res, err := http.Get(fmt.Sprintf("https:%s", newImgReq.url))
			imgBytes, err := io.ReadAll(res.Body)


			os.Mkdir(fmt.Sprintf("../downloads/%s",newImgReq.keywords[0]), 0777)
			err = os.WriteFile(finalName, imgBytes, 0666)

			if err != nil  {
				fmt.Println("Error with image", finalName, "::: ", err)
			} else {
				fmt.Println("Successfully downloaded")
			}
			time.Sleep(time.Second / 2)
    }
}


func BoardWorker(boardName string, keywords []string) {
	catalogUrl := fmt.Sprintf("https://a.4cdn.org/%s/catalog.json", boardName)
	catRequest, err := http.Get(catalogUrl)
	if err != nil {
		// handle error
	}
	defer catRequest.Body.Close()
	catRaw, err := io.ReadAll(catRequest.Body)

	var catJson [] Page
	
	if err := json.Unmarshal(catRaw, &catJson); err != nil {
		fmt.Println("process failed for board", boardName, err)
	}

	// fmt.Println("RES", catJson)

	for _, page := range catJson {
		threads := page.Threads
		for _, thread := range threads {
			validKeys := make([]string, 0)
			for _, keyword := range keywords {
				// TODO : create function that looks for words simultaneously when going over the thread titles
				lowerCom, lowerSub := strings.ToLower(thread.Comment), strings.ToLower(thread.Subtitle)
				if strings.Contains(lowerCom, keyword) || strings.Contains(lowerSub, keyword) {
					validKeys = append(validKeys, keyword)
				}
			} 
			if len(validKeys) > 0 {
				fmt.Println("Thread ", len(validKeys),  thread.Comment, "\n")
				go ThreadWorker(boardName, thread.Id , validKeys)
			}
		}
	}

}

func ThreadWorker (boardName string, threadId int, keywords []string) {
	threadUrl := fmt.Sprintf("https://boards.4chan.org/%s/thread/%d", boardName, threadId)
	fmt.Println("Started thread worker", threadUrl)
	req, _ := http.Get(threadUrl)
	body, _ := io.ReadAll(req.Body)

	imgRE := regexp.MustCompile(`(\/\/i(?:s|)\d*\.(?:4cdn|4chan)\.org\/\w+\/(\d+\.(?:jpg|png|gif|webm)))`)
	images := imgRE.FindAllString(string(body), -1)
	fmt.Println("matches: ", len(images))

	for _, imgUrl := range images {
		imgReq := ImageRequest {
			board: boardName,
			keywords: keywords,
			threadId: threadId,
			url: imgUrl,
		}
		ImageChannel <-imgReq
	}
}
