package handler

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/choigonyok/techlog/internal/usecase"
	"github.com/choigonyok/techlog/pkg/image"
	"github.com/choigonyok/techlog/pkg/model"
	"github.com/choigonyok/techlog/pkg/time"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	gin.HandlerFunc
	usecase *usecase.PostUsecase
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		usecase: usecase.NewPostUsecase(),
	}
}

// GetPosts returns every posts data without post body
func (h *PostHandler) GetPosts(c *gin.Context) {
	tag := c.Query("tag")
	posts, err := h.usecase.GetPosts(tag)
	if err != nil {
		fmt.Println("ERR GETTING POSTS:", err)
		// RETURN
	}
	b, _ := json.Marshal(posts)
	c.Writer.Write(b)
}

// GetPost returns post data including post body
func (h *PostHandler) GetPost(c *gin.Context) {
	postId := c.Param("postId")
	post, err := h.usecase.GetPost(postId)
	if err != nil {
		fmt.Println("ERR GETTING POST:", err)
		// RETURN
		return
	}

	b, _ := json.Marshal(post)
	c.Writer.Write(b)
}

// CreatePost creates new post with client input data
func (h *PostHandler) CreatePost(c *gin.Context) {
	post := model.Post{}
	post.WriteTime = time.GetCurrentTimeByFormat("2006-01-02")

	m, err := c.MultipartForm()
	if err != nil {
		fmt.Println("ERR PARSING MULTIPART FORM: ", err)
		// RETURN
	}
	postDatas := m.Value["data"]
	imageDatas := m.File["file"]

	postData := []byte(postDatas[0])
	json.Unmarshal(postData, &post)

	err = h.usecase.CreatePost(&post, imageDatas)
	if err != nil {
		fmt.Println("ERR CREATING POST: ", err)
		// RETURN
	}

	// PUSH POST TO GITHUB REPO
	// post.ID = postID
	// err = github.PushCreatedPost(post, imageNames, false)
	//
	//	if err != nil {
	//		resp.Response500(c, err)
	//	}
}

// GetTags returns every stored tags and number of posts each tag contains
func (h *PostHandler) GetTags(c *gin.Context) {
	tags, err := h.usecase.GetTags()
	if err != nil {
		fmt.Println("ERR GETTING TAGS:", err)
		// RETURN
		return
	}

	b, _ := json.Marshal(tags)

	c.Writer.Write(b)
}

// GetThumbnailByPostID returns post thumbnail image file
func (h *PostHandler) GetThumbnail(c *gin.Context) {
	postId := c.Param("postId")

	thubmail := h.usecase.GetThumbnail(postId)

	img, err := image.Download(thubmail.ID)
	if err != nil {
		// RETURN
		return
	}
	defer img.Body.Close()

	c.Header("Content-Type", *img.ContentType)
	io.Copy(c.Writer, img.Body)
	if err != nil {
		// RETURN
		return
	}

	b, _ := json.Marshal(thubmail)
	c.Writer.Write(b)
}

// DeletePost deletes post and images by post id
func (h *PostHandler) DeletePost(c *gin.Context) {
	postId := c.Param("postId")

	if err := h.usecase.DeletePost(postId); err != nil {
		fmt.Println("ERR DELETING POST: ", err)
		// RETURN
	}

	// for _, v := range imageNames {
	// 	if err := img.Remove(v); err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// }

	// err = github.PushDeletedPost(post.Title, post.ID, false)
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }
}

// UpdatePost updates title, subtitle, tags, body of post
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postId := c.Param("postId")
	post := model.Post{}
	post.ID = postId
	m, err := c.MultipartForm()
	if err != nil {
		fmt.Println("ERR PARSING MULTIPART FORM: ", err)
		// RETURN
	}
	postDatas := m.Value["data"]
	imageDatas := m.File["file"]

	postData := []byte(postDatas[0])
	json.Unmarshal(postData, &post)

	h.usecase.UpdatePost(&post, imageDatas)

	// if beforePost.Title == afterPost.Title {
	// 	err = github.PushUpdatedPost(afterPost)
	// 	if err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// } else {
	// 	if err := github.PushDeletedPost(beforePost.Title, beforePost.ID, true); err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// 	if err := github.PushCreatedPost(afterPost, nil, true); err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// }

	// if err := svcMaster.UpdatePost(afterPost); err != nil {
	// 	resp.Response500(c, err)
	// }
}

func (h *PostHandler) GetImages(c *gin.Context) {
	postId := c.Param("postId")
	images, err := h.usecase.GetImages(postId)
	if err != nil {
		fmt.Println("ERR GETTING IMAGES :", err)
		// RETURN
		return
	}

	b, _ := json.Marshal(images)
	c.Writer.Write(b)
}

// // isImageEqual compares images before update with images after update
// func isImageEqual(beforeImages string, afterImages string) bool {
// 	cmp := map[string]bool{}
// 	for _, v := range strings.Fields(beforeImages) {
// 		cmp[v] = true
// 	}
// 	for _, v := range strings.Fields(afterImages) {
// 		if !cmp[v] {
// 			return false
// 		}
// 	}
// 	return true
// }

func (h *PostHandler) UpdateImages(c *gin.Context) {
	postId := c.Param("postId")
	images := []model.Image{}

	c.ShouldBindJSON(&images)

	if err := h.usecase.UpdateImages(images, postId); err != nil {
		fmt.Println("ERR UPDATING IMAGES:", err)
		// RETURN
		return
	}
}

func (h *PostHandler) GetImage(c *gin.Context) {
	imageId := c.Param("imageId")

	image := h.usecase.GetImage(imageId)

	b, _ := json.Marshal(image)
	c.Writer.Write(b)

	// image, err := img.Download(imageName)
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }
	// defer image.Body.Close()

	// c.Header("Content-Type", *image.ContentType)
	// io.Copy(c.Writer, image.Body)
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }
}

// func GetEveryPostCount(c *gin.Context) {
// 	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
// 	svcSlave := service.NewService(pvrSlave)
// 	count := svcSlave.GetEveryPostCount()
// 	if count == 0 {
// 		resp.Response500(c, errors.New("there is no post"))
// 	}
// 	resp.ResponseDataWith200(c, struct {
// 		Count string `json:"count"`
// 	}{
// 		Count: strconv.Itoa(count),
// 	})
// }
