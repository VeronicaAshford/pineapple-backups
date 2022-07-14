package src

import (
	"fmt"
	"github.com/schollz/progressbar"
	"math/rand"
	"os"
	"sf/src/boluobao"
	"sf/src/config"
	"strconv"
	"time"
)

type AutoGenerated struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Index    int    `json:"index"`
	IsVip    bool   `json:"is_vip"`
	VolumeID string `json:"volume_id"`
	Content  string `json:"content"`
}

func GetCatalogue(BookData Books) bool {
	response := boluobao.GetCatalogueDetailedById(BookData.NovelID)
	for _, data := range response.Data.VolumeList {
		fmt.Println("start download volume: ", data.Title)
		bar := progressbar.New(len(data.ChapterList))
		for _, Chapter := range data.ChapterList {
			if Chapter.OriginNeedFireMoney == 0 {
				GetContent(len(data.ChapterList), BookData, strconv.Itoa(Chapter.ChapID), bar)
			} else {
				fmt.Println("this chapter is VIP and need fire money, skip it")
			}
		}
	}
	return true
}

func GetContent(ChapLength int, BookData Books, ChapterId string, bar *progressbar.ProgressBar) {
	if err := bar.Add(1); err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}
	response := boluobao.GetContentDetailedByCid(ChapterId)
	if response.Status.HTTPCode != 200 {
		if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" {
			fmt.Println(response.Status.Msg)
			os.Exit(0)
		} else {
			fmt.Println(response.Status.Msg)
			GetContent(ChapLength, BookData, ChapterId, bar)
		}
	} else {
		if f, err := os.OpenFile(config.Var.SaveFile+"/"+BookData.NovelName+".txt",
			os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			defer func(f *os.File) {
				err = f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(f)
			if _, ok := f.WriteString("\n\n\n" +
				response.Data.Title + ":" + response.Data.AddTime + "\n" +
				response.Data.Expand.Content + "\n" + BookData.AuthorName,
			); ok != nil {
				fmt.Println(ok)
			}
		} else {
			fmt.Println(err)
		}
	}
	//fmt.Printf(" %d/%d \r", response.Data.ChapOrder, ChapLength)
}
