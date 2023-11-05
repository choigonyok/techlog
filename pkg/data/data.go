package data

import (
	"encoding/base64"
	"strings"

	"github.com/choigonyok/techlog/pkg/model"
	"github.com/google/uuid"
)

func CreateRandomString() string {
	u := uuid.New()
	return u.String()
}

func EncodeBase64(target string) string {
	t := []byte(target)
	return base64.RawStdEncoding.EncodeToString(t)
}

func MarshalPostToDatabaseFmt(post model.Post) model.Post {
	post.Tags = RemoveWhiteSpace(post.Tags)
	post.Title = RemoveWhiteSpace(post.Title)
	post.Tags = strings.ToUpper(strings.ReplaceAll(post.Tags, `'`, `\'`))

	post.Text = strings.ReplaceAll(post.Text, `'`, `\'`)
	post.Title = strings.ReplaceAll(post.Title, `'`, `\'`)
	post.WriteTime = strings.ReplaceAll(post.WriteTime, `'`, `\'`)

	return post
}

func UnMarshalDatabaseFormat(tag, text string) string {
	return strings.ReplaceAll(text, `\'`, `'`)
}

func RemoveWhiteSpace(target string) string {
	target = strings.TrimSpace(target)
	var before string
	var found bool
	var str []string
	for {
		before, target, found = strings.Cut(target, " ")
		if found {
			str = append(str, before)
			target = strings.TrimSpace(target)
		} else {
			str = append(str, before)
			break
		}
	}
	result := strings.Join(str, " ")
	return result
}
