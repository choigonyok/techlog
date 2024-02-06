package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/choigonyok/techlog/pkg/model"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubToken *string

const (
	repositoryName = "techlog-posts"
	githubUserName = "choigonyok"
)

// SyncGithubToken synchronizes input token value with githubToken string pointer
func SyncGithubToken(tokenValue string) {
	githubToken = &tokenValue
}

// GetPostsFromGithubRepo returns every post datas in github techlog-posts repo
func GetPostsFromGithubRepo() []model.Post {
	posts := []model.Post{}
	post := model.Post{}
	var postData string

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, directoryContent, _, err := client.Repositories.GetContents(ctx, githubUserName, repositoryName, "/", nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for _, v := range directoryContent {
		fileName, _ := strings.CutPrefix(v.GetDownloadURL(), "https://raw.githubusercontent.com/"+githubUserName+"/"+repositoryName+"/main/")
		fileName = strings.ReplaceAll(fileName, "%20", " ")
		fileName = strings.ReplaceAll(fileName, "=", "%3D")
		fileName = strings.ReplaceAll(fileName, "&", "%26")
		fileName = strings.ReplaceAll(fileName, "+", "%2B")
		fileName = strings.ReplaceAll(fileName, "#", "%23")
		fileName = strings.ReplaceAll(fileName, "/", "%2F")
		fileName = strings.ReplaceAll(fileName, "?", "%3F")
		fileName = strings.ReplaceAll(fileName, `"`, "%22")
		fileName = strings.ReplaceAll(fileName, `'`, "%27")
		fileName = strings.ReplaceAll(fileName, `~`, "%7E")
		fileName = strings.ReplaceAll(fileName, `!`, "%21")

		response, err := http.Get("https://raw.githubusercontent.com/" + githubUserName + "/" + repositoryName + "/main/" + fileName)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		postData = string(body)

		id, afterID, _ := strings.Cut(postData, "\n")
		id, _ = strings.CutPrefix(id, "[ID: ")
		id, _ = strings.CutSuffix(id, "]")
		post.ID, _ = strconv.Atoi(id)

		tags, afterTags, _ := strings.Cut(afterID, "\n")
		tags, _ = strings.CutPrefix(tags, "[Tags: ")
		tags, _ = strings.CutSuffix(tags, "]")
		post.Tags = tags

		title, afterTitle, _ := strings.Cut(afterTags, "\n")
		title, _ = strings.CutPrefix(title, "[Title: ")
		title, _ = strings.CutSuffix(title, "]")
		post.Title = title

		writeTime, afterWriteTime, _ := strings.Cut(afterTitle, "\n")
		writeTime, _ = strings.CutPrefix(writeTime, "[WriteTime: ")
		writeTime, _ = strings.CutSuffix(writeTime, "]")
		post.WriteTime = writeTime

		imageNames, afterImageNames, _ := strings.Cut(afterWriteTime, "\n")
		imageNames, _ = strings.CutPrefix(imageNames, "[ImageNames: ")
		imageNames, _ = strings.CutSuffix(imageNames, "]")
		post.ThumbnailPath = imageNames

		subtitle, content, _ := strings.Cut(afterImageNames, "\n")

		// Deprecated: After previous posts are modified to contain empty subtitle
		// There are some posts not contain subtitle
		// because they are written before subtitle feature is implemented.
		// So GetPostsFromGithubRepo should check subtitle is exist or not
		subtitle, isExist := strings.CutPrefix(subtitle, "[Subtitle: ")
		if isExist {
			subtitle, _ = strings.CutSuffix(subtitle, "]")
			post.Subtitle = subtitle
			post.Text = content
		} else {
			post.Text = afterImageNames
		}

		posts = append(posts, post)
	}
	return posts
}

// PushCreatedPost commits and pushes newly created post body data
func PushCreatedPost(post model.Post, imageNames []string, isUpdate bool) error {
	postTitleIncludeExtension := post.Title
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "%20", " ")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "%", "%25")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "=", "%3D")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "&", "%26")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "+", "%2B")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "#", "%23")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, "?", "%3F")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, `"`, "%22")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, `'`, "%27")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, `~`, "%7E")
	postTitleIncludeExtension = strings.ReplaceAll(postTitleIncludeExtension, `!`, "%21")
	postTitleIncludeExtension = postTitleIncludeExtension + ".md"

	images := ""
	if imageNames == nil {
		images = post.ThumbnailPath
	} else {
		images = strings.Join(imageNames, " ")
	}

	var (
		commitMessage string
		fileContent   = `[ID: ` + strconv.Itoa(post.ID) + `]` + `
[Tags: ` + post.Tags + `]` + `
[Title: ` + post.Title + `]` + `
[WriteTime: ` + post.WriteTime + `]` + `
[ImageNames: ` + images + `]` + `

` + post.Text
	)

	switch isUpdate {
	case false:
		commitMessage = "New(" + strconv.Itoa(post.ID) + "): " + postTitleIncludeExtension
	case true:
		commitMessage = "Update(" + strconv.Itoa(post.ID) + "): Create " + postTitleIncludeExtension
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(fileContent),
	}

	content, resp, err := client.Repositories.CreateFile(ctx, githubUserName, repositoryName, postTitleIncludeExtension, opt)
	if err != nil {
		fmt.Println("[Error StatusCode]:", resp.StatusCode)
		return err
	} else {
		fmt.Printf("[New Post Commit] committer: %s, message: %s\n", *content.Commit.Committer, *content.Commit.Message)
	}
	return nil
}

func PushDeletedPost(postTitle string, postID int, isUpdate bool) error {
	postTitleIncludeExtension := postTitle + ".md"
	var (
		commitMessage string
		sha           string
	)
	switch isUpdate {
	case false:
		commitMessage = "Remove(" + strconv.Itoa(postID) + "): " + postTitleIncludeExtension
	case true:
		commitMessage = "Update(" + strconv.Itoa(postID) + "): Remove " + postTitleIncludeExtension
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, files, _, _ := client.Repositories.GetContents(ctx, githubUserName, repositoryName, "", nil)
	for _, file := range files {
		if string(file.GetName()) == postTitleIncludeExtension {
			sha = string(file.GetSHA())
		}
	}

	opt := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		SHA:     &sha,
	}

	_, _, err := client.Repositories.DeleteFile(ctx, githubUserName, repositoryName, postTitleIncludeExtension, opt)
	if err != nil {
		return err
	}
	return nil
}

func PushUpdatedPost(afterPost model.Post) error {
	var (
		PostTitleIncludeExtension = afterPost.Title + ".md"
		commitMessage             = "Update(" + strconv.Itoa(afterPost.ID) + "): " + PostTitleIncludeExtension
		// commitMessage = "Test: Push to this repo when new post is updated"
		fileContent = `[ID: ` + strconv.Itoa(afterPost.ID) + `]` + `
		[Tags: ` + afterPost.Tags + `]` + `
		[Title: ` + afterPost.Title + `]` + `
		[WriteTime: ` + afterPost.WriteTime + `]` + `
		[ImageNames: ` + afterPost.ThumbnailPath + `]` + `
		
		` + afterPost.Text
	)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fileContents, _, _, _ := client.Repositories.GetContents(ctx, githubUserName, repositoryName, PostTitleIncludeExtension, nil)

	opt := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		SHA:     fileContents.SHA,
		Content: []byte(fileContent),
	}

	fmt.Println("string SHA: ", *fileContents.SHA)

	_, _, err := client.Repositories.UpdateFile(ctx, githubUserName, repositoryName, PostTitleIncludeExtension, opt)
	if err != nil {
		return err
	}
	return nil
}
