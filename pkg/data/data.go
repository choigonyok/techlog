package data

import (
	"encoding/base64"
	"strings"

	"github.com/choigonyok/techlog/pkg/model"
	"github.com/google/uuid"
)

// CreateRandomString returns random 36 bits unique string
func CreateRandomString() string {
	u := uuid.New()
	return u.String()
}

// EncodeBase64 encodes target string to base64 format string
func EncodeBase64(target string) string {
	t := []byte(target)
	return base64.RawStdEncoding.EncodeToString(t)
}

// DecodeBase64 decodes base64 format string to original
func DecodeBase64(target string) string {
	transedTarget, _ := base64.RawStdEncoding.DecodeString(target)
	return string(transedTarget)
}

// MarshalCommentToDatabaseFmt manipulates model.Comment not to get error by letter like ', ", \, `
func MarshalCommentToDatabaseFmt(comment model.Comment) model.Comment {
	comment.Text = strings.ReplaceAll(strings.ReplaceAll(comment.Text, `\`, `\\`), `'`, `\'`)
	comment.WriterID = strings.ReplaceAll(strings.ReplaceAll(comment.WriterID, `\`, `\\`), `'`, `\'`)
	return comment
}

// UnmarshalCommentToDatabaseFmt changes stored manipulated model.Comment to original
func UnmarshalCommentToDatabaseFmt(comments []model.Comment) []model.Comment {
	for _, comment := range comments {
		comment.Text = strings.ReplaceAll(strings.ReplaceAll(comment.Text, `\'`, `'`), `\`, `\\`)
		comment.WriterID = strings.ReplaceAll(strings.ReplaceAll(comment.WriterID, `\'`, `'`), `\`, `\\`)
	}
	return comments
}

// MarshalCommentToDatabaseFmt manipulates model.Comment not to get error by letter like ', ", \, `
func MarshalReplyToDatabaseFmt(reply model.Reply) model.Reply {
	reply.Text = strings.ReplaceAll(strings.ReplaceAll(reply.Text, `\`, `\\`), `'`, `\'`)
	reply.WriterID = strings.ReplaceAll(strings.ReplaceAll(reply.WriterID, `\`, `\\`), `'`, `\'`)
	return reply
}

// UnmarshalCommentToDatabaseFmt changes stored manipulated model.Comment to original
func UnmarshalReplyToDatabaseFmt(replies []model.Reply) []model.Reply {
	for _, reply := range replies {
		reply.Text = strings.ReplaceAll(strings.ReplaceAll(reply.Text, `\'`, `'`), `\`, `\\`)
		reply.WriterID = strings.ReplaceAll(strings.ReplaceAll(reply.WriterID, `\'`, `'`), `\`, `\\`)
	}
	return replies
}

// MarshalPostToDatabaseFmt manipulates model.Post not to get error by letter like ', ", \, `
func MarshalPostToDatabaseFmt(post model.Post) model.Post {
	post.Tags = RemoveWhiteSpace(post.Tags)
	post.Title = RemoveWhiteSpace(post.Title)

	post.Tags = strings.ToUpper(strings.ReplaceAll((strings.ReplaceAll(post.Tags, `\`, `\\`)), `'`, `\'`))
	post.Text = strings.ReplaceAll(strings.ReplaceAll(post.Text, `\`, `\\`), `'`, `\'`)
	post.Title = strings.ReplaceAll(strings.ReplaceAll(post.Title, `\`, `\\`), `'`, `\'`)
	post.WriteTime = strings.ReplaceAll(strings.ReplaceAll(post.WriteTime, `\`, `\\`), `'`, `\'`)

	return post
}

// UnMarshalPostDatabaseFmt changes stored manipulated model.Post to original
func UnMarshalPostDatabaseFmt(post model.Post) model.Post {
	post.Tags = strings.ReplaceAll(strings.ReplaceAll(post.Tags, `\'`, `'`), `\`, `\\`)
	post.Text = strings.ReplaceAll(strings.ReplaceAll(post.Text, `\'`, `'`), `\`, `\\`)
	post.Title = strings.ReplaceAll(strings.ReplaceAll(post.Title, `\'`, `'`), `\`, `\\`)
	post.WriteTime = strings.ReplaceAll(strings.ReplaceAll(post.WriteTime, `\'`, `'`), `\`, `\\`)

	return post
}

// RemoveWhiteSpace manipulates target string to have only 1 white space among words without all leading/trailing white space
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
