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

func DecodeBase64(target string) string {
	transedTarget, _ := base64.RawStdEncoding.DecodeString(target)
	return string(transedTarget)
}

func MarshalCommentToDatabaseFmt(comment model.Comment) model.Comment {
	comment.Text = strings.ReplaceAll(strings.ReplaceAll(comment.Text, `\`, `\\`), `'`, `\'`)
	comment.WriterPW = strings.ReplaceAll(strings.ReplaceAll(comment.WriterPW, `\`, `\\`), `'`, `\'`)
	comment.WriterID = strings.ReplaceAll(strings.ReplaceAll(comment.WriterID, `\`, `\\`), `'`, `\'`)
	return comment
}

func UnmarshalCommentToDatabaseFmt(comments []model.Comment) []model.Comment {
	for _, comment := range comments {
		comment.Text = strings.ReplaceAll(strings.ReplaceAll(comment.Text, `\'`, `'`), `\`, `\\`)
		comment.WriterID = strings.ReplaceAll(strings.ReplaceAll(comment.WriterID, `\'`, `'`), `\`, `\\`)
		comment.WriterPW = strings.ReplaceAll(strings.ReplaceAll(comment.WriterPW, `\'`, `'`), `\`, `\\`)
	}
	return comments
}

func MarshalPostToDatabaseFmt(post model.Post) model.Post {
	post.Tags = RemoveWhiteSpace(post.Tags)
	post.Title = RemoveWhiteSpace(post.Title)

	post.Tags = strings.ToUpper(strings.ReplaceAll((strings.ReplaceAll(post.Tags, `\`, `\\`)), `'`, `\'`))
	post.Text = strings.ReplaceAll(strings.ReplaceAll(post.Text, `\`, `\\`), `'`, `\'`)
	post.Title = strings.ReplaceAll(strings.ReplaceAll(post.Title, `\`, `\\`), `'`, `\'`)
	post.WriteTime = strings.ReplaceAll(strings.ReplaceAll(post.WriteTime, `\`, `\\`), `'`, `\'`)

	return post
}

func UnMarshalPostDatabaseFmt(post model.Post) model.Post {
	post.Tags = strings.ReplaceAll(strings.ReplaceAll(post.Tags, `\'`, `'`), `\`, `\\`)
	post.Text = strings.ReplaceAll(strings.ReplaceAll(post.Text, `\'`, `'`), `\`, `\\`)
	post.Title = strings.ReplaceAll(strings.ReplaceAll(post.Title, `\'`, `'`), `\`, `\\`)
	post.WriteTime = strings.ReplaceAll(strings.ReplaceAll(post.WriteTime, `\'`, `'`), `\`, `\\`)

	return post
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
