package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/database"
	img "github.com/choigonyok/techlog/pkg/image"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/choigonyok/techlog/pkg/time"
	"github.com/gin-gonic/gin"
)

// CreatePost creates new post with client input data
func CreatePost(c *gin.Context) {
	VerifyAdminUser(c)
	// statusCode := VerifyAdminUser(c)
	// if statusCode == 401 {
	// 	resp.Response401(c)
	// 	return
	// }

	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	post := model.Post{}
	image := model.PostImage{}
	imageNames := []string{}

	formData, err := c.MultipartForm()
	if err != nil {
		resp.Response500(c, err)
		return
	}

	// handle post data
	postDatas := formData.Value["data"]
	postData := []byte(postDatas[0])
	err = json.Unmarshal(postData, &post)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	post.WriteTime = time.GetCurrentTimeByFormat("2006-01-02")
	postID, err := svcMaster.CreatePost(post)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	// handle image data
	imageDatas := formData.File["file"]

	for i, v := range imageDatas {
		image.ImageName = data.CreateRandomString() + v.Filename[strings.LastIndex(v.Filename, "."):]
		if i == 0 {
			image.Thumbnail = "true"
		} else {
			image.Thumbnail = "false"
		}

		image.PostID = postID

		err = img.Upload(v, image.ImageName)
		if err != nil {
			fmt.Println(err.Error())
		}

		err = svcMaster.StoreImage(image)
		if err != nil {
			err := rollBackSavedImageByImageName(image.ImageName)
			resp.Response500(c, err)
			return
		}
		imageNames = append(imageNames, image.ImageName)
	}
	post.ID = postID
	// err = github.PushCreatedPost(post, imageNames, false)
	// if err != nil {
	// 	resp.Response500(c, err)
	// }
}

// rollBackSavedImageByImageName deletes saved image by file name
func rollBackSavedImageByImageName(imageName string) error {
	return os.Remove("assets/" + imageName)
}

// DeletePostByPostID deletes post and images by post id
func DeletePostByPostID(c *gin.Context) {
	VerifyAdminUser(c)
	// statusCode := VerifyAdminUser(c)
	// if statusCode == 401 {
	// 	resp.Response401(c)
	// 	return
	// }

	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	postID := c.Param("postid")

	// post, err := svc.GetPostByID(postID)
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }

	imageNames, err := svcMaster.DeletePostByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	for _, v := range imageNames {
		if err := img.Remove(v); err != nil {
			resp.Response500(c, err)
			return
		}
	}

	// err = github.PushDeletedPost(post.Title, post.ID, false)
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }
}

// GetPost returns post data including post body
func GetPost(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	postID := c.Param("postid")

	post, err := svcSlave.GetPostByID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	err = resp.ResponseDataWith200(c, post)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// UpdatePost updates title, tags, body of post
func UpdatePostByPostID(c *gin.Context) {
	VerifyAdminUser(c)
	// statusCode := VerifyAdminUser(c)
	// if statusCode == 401 {
	// 	resp.Response401(c)
	// 	return
	// }

	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	postID := c.Param("postid")

	beforePost, _ := svcSlave.GetPostByID(postID)

	afterPost := model.Post{}
	err := c.ShouldBindJSON(&afterPost)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	afterPost = data.MarshalPostToDatabaseFmt(afterPost)

	if afterPost.ID, err = strconv.Atoi(postID); err != nil {
		resp.Response500(c, err)
		return
	} else {
		err = svcMaster.UpdatePost(afterPost)
	}

	if err != nil {
		resp.Response500(c, err)
		return
	}

	if beforePost.Title == afterPost.Title {
		if beforePost.ThumbnailPath != afterPost.ThumbnailPath {
			imageNames := strings.Split(beforePost.ThumbnailPath, " ")
			for _, v := range imageNames {
				img.Remove(v)
			}
			// img.Upload()
		}
		// err = github.PushUpdatedPost(afterPost)
		// if err != nil {
		// 	resp.Response500(c, err)
		// 	return
		// }
	}
	// else {
	// 	if err := github.PushDeletedPost(beforePost.Title, beforePost.ID, true); err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// 	if err := github.PushCreatedPost(afterPost, nil, true); err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// }
}

func UpdatePostImagesByPostID(c *gin.Context) {
	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	postID := c.Param("postid")
	images := []model.PostImage{}

	c.ShouldBindJSON(&images)

	if err := svcMaster.DeletePostImagesByPostID(postID); err != nil {
		resp.Response500(c, err)
		return
	}
	if err := svcMaster.CreatePostImagesByPostID(postID, images); err != nil {
		resp.Response500(c, err)
		return
	}
}

// GetPostCards returns every posts data without post body
func GetPosts(c *gin.Context) {
	tag := c.Query("tag")
	if tag != "ALL" {
		cards, err := getPostsByTag(tag)
		if err != nil {
			resp.Response500(c, err)
			return
		} else {
			resp.ResponseDataWith200(c, cards)
			return
		}
	}
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)

	cards, err := svcSlave.GetPosts()
	if err != nil {
		resp.Response500(c, err)
		return
	}
	err = resp.ResponseDataWith200(c, cards)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// getEveryCardByTag returns posts data by tag without post body
func getPostsByTag(tag string) ([]model.PostCard, error) {

	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)

	return svcSlave.GetPostsByTag(tag)
}

// GetTags returns every stored tags
func GetTags(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)

	tags, err := svcSlave.GetTags()
	if err != nil {
		resp.Response500(c, err)
		return
	}
	if resp.ResponseDataWith200(c, tags) != nil {
		resp.Response500(c, err)
	}
}

// GetThumbnailByPostID returns post thumbnail image file
func GetThumbnailByPostID(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	postID := c.Param("postid")

	thumbnailName, err := svcSlave.GetThumbnailNameByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	image, err := img.Download(thumbnailName)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	defer image.Body.Close()

	c.Header("Content-Type", *image.ContentType)
	io.Copy(c.Writer, image.Body)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

func GetPostImageByImageID(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	imageID := c.Param("imageid")

	imageName, err := svcSlave.GetPostImageNameByImageID(imageID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	image, err := img.Download(imageName)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	defer image.Body.Close()

	c.Header("Content-Type", *image.ContentType)
	io.Copy(c.Writer, image.Body)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

func GetImagesByPostID(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	postID := c.Param("postid")
	images, err := svcSlave.GetImagesByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	resp.ResponseDataWith200(c, images)
}
