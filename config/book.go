package config

import (
	"fmt"
	"regexp"
)

func FindID(url string) string {
	current_book_id := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
	if len(current_book_id) > 1 {
		if Vars.AppType == "cat" {
			if len(current_book_id[1]) != 9 { // test if the input is hbooker book id
				fmt.Println("hbooker bookid is 9 characters, please input again:")
			} else {
				return current_book_id[1]
			}
		}
	}
	return ""
}
