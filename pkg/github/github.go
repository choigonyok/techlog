package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

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
func GetPostsFromGithubRepo() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, directoryContent, _, err := client.Repositories.GetContents(ctx, githubUserName, repositoryName, "/", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var text string
	for _, v := range directoryContent {
		response, err := http.Get(v.GetDownloadURL())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		text = string(body)
	}
	fmt.Println(text)
}

// PushCreatedPost commits and pushes newly created post body data
func PushCreatedPost(postID, postTitle, postContent string) {
	var (
		commitMessage     = "Docs(): Create new post"
		fullPostFileTitle = postID + "-" + postTitle + ".md"
	)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// f, _ := os.OpenFile(postTitle, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	// f.WriteString(postContent)
	// defer f.Close()

	// fileInfo, _ := f.Stat()
	// fileSize := fileInfo.Size()
	// fileContent := make([]byte, fileSize)

	// f.Seek(0, 0)
	// f.Read(fileContent)

	opt := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(postContent),
	}

	content, resp, err := client.Repositories.CreateFile(ctx, githubUserName, repositoryName, fullPostFileTitle, opt)
	if err != nil {
		fmt.Println("[Error StatusCode]:", resp.StatusCode)
		fmt.Println("[Error Message]:", err.Error())
	} else {
		fmt.Printf("[New Post Commit] committer: %s, message: %s\n", *content.Commit.Committer, *content.Commit.Message)
	}
}
