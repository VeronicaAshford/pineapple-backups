package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"os"
	"path"
	"strconv"
)

type BookInits struct {
	BookID      string
	ShowBook    bool
	Locks       *threading.GoLimit
	EpubSetting *epub.Epub
}
type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  string
	NovelCover string
	AuthorName string
	CharCount  string
	SignStatus string
}

func InitEpubFile() *epub.Epub {
	AddImage := true
	var epubSetting *epub.Epub
	switch command.Command.AppType {
	case "sfacg":
		epubSetting = epub.NewEpub(config.Current.NewBooks["novel_name"]) // set epub setting and add section
		epubSetting.SetAuthor(config.Current.NewBooks["author_name"])     // set author
		if !config.Exist(config.Current.CoverPath) {
			if reader := request.Request(config.Current.NewBooks["novel_cover"]); reader == nil {
				fmt.Println("download cover failed!")
				AddImage = false
			} else {
				_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
			}
		}
	case "cat":
		epubSetting = epub.NewEpub(config.APP.Hbooker.BookInfo.BookName) // set epub setting and add section
		epubSetting.SetAuthor(config.APP.Hbooker.BookInfo.AuthorName)    // set author
		if !config.Exist(config.Current.CoverPath) {
			if reader := request.Request(config.APP.Hbooker.BookInfo.Cover); reader == nil {
				fmt.Println("download cover failed!")
				AddImage = false
			} else {
				_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
			}
		}
		if AddImage {
			_, _ = epubSetting.AddImage(config.Current.CoverPath, "")
			epubSetting.SetCover(path.Join("../images", config.Current.NewBooks["novel_name"]+".jpg"), "")
		}
	}
	return epubSetting
}

func SettingBooks(bookId string) (*Catalogue, error) {
	var err error
	switch command.Command.AppType {
	case "sfacg":
		BookInfo := boluobao.API.Book.NovelInfo(bookId)
		if BookInfo != nil {
			config.Current.NewBooks = map[string]string{
				"novel_name":  tools.RegexpName(BookInfo.NovelName),
				"novel_id":    strconv.Itoa(BookInfo.NovelID),
				"novel_cover": BookInfo.NovelCover,
				"author_name": BookInfo.AuthorName,
				"char_count":  strconv.Itoa(BookInfo.CharCount),
				"mark_count":  strconv.Itoa(BookInfo.MarkCount),
			}
			tools.Mkdir(path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"]))
			config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.Current.NewBooks["novel_name"])
			config.Current.CoverPath = path.Join(config.Vars.ConfigName, config.Vars.CoverFile, config.Current.NewBooks["novel_name"]+".jpg")
		} else {
			return nil, fmt.Errorf("book_id:%v is invalid", bookId)
		}
	case "cat":
		config.APP.Hbooker.BookInfo, err = config.APP.Hbooker.Client.API.GetBookInfo(bookId)
		if err != nil {
			return nil, err
		}
		tools.Mkdir(path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName))
		config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.APP.Hbooker.BookInfo.BookName)
		config.Current.CoverPath = path.Join(config.Vars.ConfigName, config.Vars.CoverFile, config.APP.Hbooker.BookInfo.BookName+".jpg")
	}
	return BookDetailed(), nil

}

func BookDetailed() *Catalogue {

	var briefIntroduction string
	if command.Command.AppType == "sfacg" {
		briefIntroduction = fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
			config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_id"], config.Current.NewBooks["author_name"],
			config.Current.NewBooks["char_count"],
		)
		file.Open(
			path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_name"]+".txt"),
			briefIntroduction,
			"w",
		)
	} else {
		briefIntroduction = fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
			config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookID, config.APP.Hbooker.BookInfo.AuthorName,
			config.APP.Hbooker.BookInfo.TotalWordCount,
		)

		file.Open(
			path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookName+".txt"),
			briefIntroduction,
			"w",
		)
	}
	fmt.Println(briefIntroduction)
	return &Catalogue{EpubSetting: InitEpubFile()}
}
