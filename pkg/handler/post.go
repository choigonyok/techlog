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
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// 작성된 게시글에 썸네일 추가
func WritePostImageHandler(c *gin.Context) {
	// VerifyAdminUser(c)

	// pvr := database.NewMysqlProvider(database.GetConnector())
	// svc := service.NewService(pvr)
	// image, err := c.MultipartForm()
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }

	// recentID, err := model.GetRecentPostID()
	// if err != nil {
	// 	fmt.Println("ERROR #4 : ", err.Error())
	// 	c.Writer.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// var thumbnailName string
	// everyIMG := imgFile.File["file"]
	// for i, v := range everyIMG {
	// 	dotIndex := strings.LastIndex(v.Filename, ".")
	// 	exetension := v.Filename[dotIndex:]
	// 	if i == 0 {
	// 		thumbnailName = exetension
	// 	}
	// 	err = c.SaveUploadedFile(v, "assets/"+strconv.Itoa(recentID)+"-"+strconv.Itoa(i)+exetension)
	// 	if err != nil {
	// 		c.Writer.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Println("ERROR #34 : ", err.Error())
	// 		return
	// 	}
	// }
	// err = model.UpdatePostImagePath(recentID, strconv.Itoa(recentID)+"-0"+thumbnailName)
	// if err != nil {
	// 	c.Writer.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Println("ERROR #36 : ", err.Error())
	// 	return
	// }

}

// CreatePost creates new post with client input data
func CreatePost(c *gin.Context) {
	VerifyAdminUser(c)

	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	post := model.Post{}
	image := model.PostImage{}

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

	postID, err := svc.CreatePost(post)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	// handle image data
	imageDatas := formData.File["file"]

	for i, v := range imageDatas {
		image.ImageName = data.CreateRandomString() + v.Filename[strings.LastIndex(v.Filename, "."):]
		fmt.Println(image.ImageName)
		if i == 0 {
			image.Thumbnail = "true"
		} else {
			image.Thumbnail = "false"
		}

		image.PostID = postID

		err = c.SaveUploadedFile(v, "assets/"+image.ImageName)
		if err != nil {
			resp.Response500(c, err)
			return
		}
		err = svc.StoreImage(image)
		if err != nil {
			fmt.Println(err.Error())
			err := rollBackSavedImageByImageName(image.ImageName)
			resp.Response500(c, err)
			return
		}
	}
}

// rollBackSavedImageByImageName deletes saved image by file name
func rollBackSavedImageByImageName(imageName string) error {
	return os.Remove("assets/" + imageName)
}

// DeletePostByPostID deletes post and images by post id
func DeletePostByPostID(c *gin.Context) {
	VerifyAdminUser(c)

	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	imageNames, err := svc.DeletePostByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	for _, v := range imageNames {
		os.Remove("assets/" + v)
	}
}

// GetPost returns post data including post body
func GetPost(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	posts, err := svc.GetPostByID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	err = resp.ResponseDataWith200(c, posts)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// UpdatePost updates title, tags, body of post
func UpdatePostByPostID(c *gin.Context) {
	VerifyAdminUser(c)

	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	post := model.Post{}
	err := c.ShouldBindJSON(&post)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	post = data.MarshalPostToDatabaseFmt(post)

	if post.ID, err = strconv.Atoi(postID); err != nil {
		resp.Response500(c, err)
		return
	} else {
		err = svc.UpdatePost(post)
	}
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// GetEveryCard returns every posts data without post body
func GetEveryCard(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	cards, err := svc.GetEveryCard()
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

// GetEveryCardByTag returns posts data by tag without post body
func GetEveryCardByTag(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	m := model.PostTags{}
	if err := c.ShouldBindJSON(&m); err != nil {
		resp.Response500(c, err)
		return
	}
	cards, err := svc.GetEveryCardByTag(m.Tags)
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

// GetTags returns every stored tags
func GetTags(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	tags, err := svc.GetEveryTags()
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
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	thumbnailName, err := svc.GetThumbnailNameByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	image, err := os.Open("assets/" + thumbnailName)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	defer image.Close()
	_, err = io.Copy(c.Writer, image)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}
