import Header from "../Header/Header";
import "./Deletepage.css";
import Footer from "../UI/Footer";
import axios from "axios";
import { useEffect, useState } from "react";
import MDEditor from "@uiw/react-md-editor";
import Writepage from "./Writepage"
import { useNavigate } from "react-router-dom";

const Deletepage = () => {
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/oauth2/auth")
      .then((response) => {
        if (response.status !== 202) {
          navigate("/");
        }
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
          } else {
            console.error(error);
          }
        }
      });
  }, [])

  axios.defaults.withCredentials = true;

  const [changeEvent, setChangeEvent] = useState(false);
  const [md, setMD] = useState("");
  const [titleText, setTitleText] = useState("");
  const [subtitleText, setSubtitleText] = useState("");
  const [tagText, setTagText] = useState("");
  const [dateText, setDateText] = useState("");
  const [bodyText, setBodyText] = useState("");
  const [id, setID] = useState();
  const [isDeleted, setIsDeleted] = useState(false);
  const [toModify, setToModify] = useState(false);
  const [allPost, setAllPost] = useState(false);
  const [postData, setPostData] = useState([]);
  const [isPosts, setIsPosts] = useState(false);
  const [isComments, setIsComments] = useState(false);
  const [isWrite, setIsWrites] = useState(false);
  const [comInfo, setComInfo] = useState([]);
  const [postRequest, setPostRequest] = useState(false);
  const [imageNames, setImageNames] = useState([]);
  const [img, setIMG] = useState([]);
  const [newImages, setNewImages] = useState([]);
  const [preImageLength, setPreImageLength] = useState(0);
  

  const postHandler = () => {
    axios
      .get(process.env.REACT_APP_HOST + "/oauth2/auth")
      .then((response) => {
        if (response.status === 202) {
          requestEditPost();
        }
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
          } else {
            console.error(error);
          }
        }
      })
  };

  const requestEditPost = () => {
    let images = [];
    imageNames.map((item) => (
      item.thumbnail === "1"
        ?
        images.push(item.imageName)
        :
        ""
    ))
    imageNames.map((item) => (
      item.thumbnail === "0"
        ?
        images.push(item.imageName)
        :
        ""
    ))
    const imageString = images.join(" ");

    const postdata = {
      title: titleText,
      tags: tagText,
      writeTime: dateText,
      text: bodyText,
      thumbnailPath: imageString,
      subtitle: subtitleText,
    };

    // console.log("SEND DATA:",JSON.stringify(postdata));

    const formData = new FormData();
    for (let i = 0; i < img.length; i++) {
      formData.append("file", img[i]);
    }
    formData.append('data', JSON.stringify(postdata));

    axios
      .put(process.env.REACT_APP_HOST + "/api/posts/" + id, formData, {
        "Content-type": "multipart/form-data",
        withCredentials: true,

      })
      .then((response) => {
      })
      .catch((error) => {
        console.log(error);
      });

    axios
      .put(process.env.REACT_APP_HOST + "/api/posts/" + id + "/images", imageNames, {
        "Content-type": "multipart/form-data",
        withCredentials: true,
      })
      .then((response) => {
        setToModify(false);
        setIsDeleted(!isDeleted);
        setChangeEvent(!changeEvent);
      })
      .catch((error) => {
        console.log(error);
      });
  }

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts")
      .then((response) => {
        setAllPost(response.data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, [isDeleted]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [changeEvent]);

  useEffect(() => {
    setBodyText(md);
  }, [md]);

  const deleteHandler = (value) => {
    axios
      .get(process.env.REACT_APP_HOST + "/oauth2/auth")
      .then((response) => {
        if (response.status === 202) {
          requestDeletePost(value);
        }
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
          } else {
            console.error(error);
          }
        }
      })
  };

  const requestDeletePost = (value) => {
    axios
      .delete(process.env.REACT_APP_HOST + "/api/posts/" + value, {
        withCredentials: true,
      })
      .then((response) => {
        setPostData(response.data);
        setIsDeleted(!isDeleted);
      })
      .catch((error) => {
        console.error(error);
      });
  }

  const modifyHandler = (value) => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + value)
      .then((response) => {
        setToModify(true);
        setID(value);
        setTitleText(response.data.title);
        setTagText(response.data.tags);
        setDateText(response.data.writetime);
        setMD(response.data.text);
        setSubtitleText(response.data.subtitle);
        setChangeEvent(!changeEvent);
      })
      .catch((error) => {
        console.error(error);
      });

    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + value + "/images")
      .then((response) => {
        setImageNames([...response.data]);
        setPreImageLength(response.data.length);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const titleHandler = (e) => {
    setTitleText(e.target.value);
  };

  const subtitleHandler = (e) => {
    setSubtitleText(e.target.value);
  };

  const tagHandler = (e) => {
    setTagText(e.target.value);
  };

  const dateHandler = (e) => {
    setDateText(e.target.value);
  };

  const isCommentsHandler = () => {
    setIsComments(true);
    setIsPosts(false);
    setPostRequest(!postRequest);
    setIsWrites(false);
  };

  const CommentDeleteHandler = (value) => {
    axios
      .get(process.env.REACT_APP_HOST + "/oauth2/auth")
      .then((response) => {
        if (response.status === 202) {
          requestDeleteComment(value);
        }
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
          } else {
            console.error(error);
          }
        }
      })
  };

  const requestDeleteComment = (value) => {
    axios
      .delete(process.env.REACT_APP_HOST + "/api/comments/" + value +
        "?type=admin")
      .then((response) => {
        setPostRequest(!postRequest);
      })
      .catch((error) => {
        console.error(error);
      });
  }

  const isPostsHandler = () => {
    setIsComments(false);
    setIsPosts(true);
    setToModify(false);
    setIsWrites(false);
  };

  const isWriteHandler = () => {
    setIsComments(false);
    setIsPosts(false);
    setIsWrites(true);
  };

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/comments")
      .then((response) => {
        setComInfo([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [postRequest]);

  const changeThumbnailHandler = (newThumbnailImage) => {
    let originalImages = imageNames
    let beforeThumbnail = imageNames.filter((element) => element.thumbnail === "1")
    let afterThumbnail = imageNames.filter((element) => element.imageName === newThumbnailImage.imageName)
    beforeThumbnail[0].thumbnail = "0";
    afterThumbnail[0].thumbnail = "1";
    setImageNames([...imageNames]);

    // setImageNames(
    //   [...afterThumbnail, ...originalImages.filter((element) => element.imageName !== beforeThumbnail[0].imageName).filter((element) => element.imageName !== afterThumbnail[0].imageName), ...beforeThumbnail]
    // );
  }

  const imgHandler = (e) => {
    setIMG((img) => [...img, ...e.target.files]);
    console.log("e.target.files : ",e.target.files);
    let f = document.getElementById("imgfile").files;
    if (f.length !== 0) {
      for (let i = 0; i < f.length; i++) {
        
        const reader = new FileReader();
        reader.onload = (e) => {
          setNewImages([...newImages, reader.result]);
        }
        reader.readAsDataURL(e.target.files[i]);
        
        let img = {
          id: 0,
          postid: id,
          imageName: f[i].name,
          thumbnail: 0,
        };
        setImageNames([...imageNames, img]);
      }
    }
  };

  const deleteIMGHandler = (value) => {
    setImageNames(
      imageNames.filter((element) => String(element.imageName) !== value.imageName)
    );
    if (value.thumbnail === "1") {
      imageNames.filter((element) => String(element.imageName) !== value.imageName)[0].thumbnail = "1";
    }
    if (value.id !== 0) {
      setPreImageLength(preImageLength - 1);
    }
  
    setIMG(
      img.filter((element) => String(element.name) !== value.imageName)
    );
  };

  console.log("IMG LENGTH : ", img.length);

  return (
    <div>
      <Header />
      <div className="delete-container">
        <div className="delete-main">ADMIN</div>
        <div className="select-container">
          <input
            type="button"
            className="select-button"
            value="WRITE"
            onClick={isWriteHandler}
          />
          <input
            type="button"
            className="select-button"
            value="POSTS"
            onClick={isPostsHandler}
          />
          <input
            type="button"
            className="select-button"
            value="COMMENTS"
            onClick={isCommentsHandler}
          />
        </div>
        {isWrite && <div>
          <Writepage />
        </div>}
        {isPosts && (
          <div>
            {toModify && (
              <div className="modify-container">
                <div className="admin-titletagdate">
                  <input type="text" value={tagText} onChange={tagHandler} />
                </div>
                <div className="admin-titletagdate">
                  <input
                    type="text"
                    value={titleText}
                    onChange={titleHandler}
                  />
                </div>
                <div className="admin-titletagdate">
                  <input
                    type="text"
                    value={subtitleText}
                    onChange={subtitleHandler}
                  />
                </div>
                <div className="image-thumbnail-container">
                  {imageNames.map((item, index) => (
                    item.thumbnail === "1" ?
                      (item.id !== 0 ?
                        <div className="image-thumbnail">
                          <img
                            className="image-thumbnail-image"
                            alt="my"
                            src={process.env.REACT_APP_HOST + "/api/posts/" + id + "/images/" + item.id}
                          />
                          <div className="image-thumbnail-name">
                            {item.imageName}
                          </div>
                        </div>
                        :
                        <div className="image-thumbnail">
                          <img
                            className="image-thumbnail-image"
                            alt="my"
                            src={newImages[index - preImageLength]}
                          />
                          <div className="image-thumbnail-name">
                            {item.imageName}
                          </div>
                        </div>)
                      :
                      (item.id !== 0 ?
                        <div className="image-non-thumbnail">
                          <img
                            className="image-non-thumbnail-image"
                            alt="my"
                            src={process.env.REACT_APP_HOST + "/api/posts/" + id + "/images/" + item.id}
                            onClick={() => changeThumbnailHandler(item)}
                          />
                          <div className="image-non-thumbnail-name">
                            {item.imageName}
                          </div>
                        </div>
                        :
                        <div className="image-non-thumbnail">
                          <img
                            className="image-non-thumbnail-image"
                            alt="my"
                            src={newImages[index - preImageLength]}
                            // src={newImages.filter((element) => element.name === item.imageName).pop()}
                            onClick={() => changeThumbnailHandler(item)}
                          />
                          <div className="image-non-thumbnail-name">
                            {item.imageName}
                          </div>
                        </div>)
                  ))}
                </div>
                <div className="admin-titletagdate">
                  <label for="imgfile">
                    <div class="file-button">IMG UPLOAD</div>
                  </label>
                  <input
                    type="file"
                    required
                    multiple
                    id="imgfile"
                    name="imgfile"
                    className="file-input"
                    onChange={imgHandler}
                  />
                </div>
                {imageNames.map((item, index) => (
                  <div className="imgname-container">
                    {item.imageName}
                    <input
                      type="button"
                      value="X"
                      className="imgname-button"
                      onClick={() => deleteIMGHandler(item)}
                    />
                  </div>
                ))}
                <div>
                  <div className="admin-editor">
                    <MDEditor height={400} value={md} onChange={setMD} />
                  </div>
                </div>
                <div className="button-container">
                  <input
                    type="button"
                    className="admin-button"
                    value="이 내용으로 수정하기"
                    onClick={postHandler}
                  />
                </div>
              </div>
            )}
            <div className="delete-list">
              {allPost && (
                <div>
                  {allPost.map((item, index) => (
                    <div className="delete-inlist">
                      <div className="delete-post">
                        <h2 className="delete-date">{item.writeTime}</h2>
                        <h2 className="delete-title">{item.title}</h2>
                        <h2 className="delete-subtitle">{item.subtitle}</h2>
                        <h2 className="delete-tag">{item.tags}</h2>
                      </div>
                      <div className="delete-button__container">
                        <input
                          className="delete-button"
                          type="button"
                          value="삭제"
                          onClick={() => {
                            deleteHandler(item.id);
                          }}
                        />
                        <input
                          className="delete-button"
                          type="button"
                          value="수정"
                          onClick={() => {
                            modifyHandler(item.id);
                          }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
        {isComments && (
          <div className="delete-list">
            {comInfo &&
              comInfo.map((item, index) => (
                <div className="delete-inlist">
                  <div className="delete-post">
                    <h2 className="delete-comment">{item.text}</h2>
                    <h2 className="delete-tag">{item.writerID}</h2>
                  </div>
                  <div className="delete-button__container">
                    <input
                      className="delete-button"
                      type="button"
                      value="삭제"
                      onClick={() => {
                        CommentDeleteHandler(item.id);
                      }}
                    />
                  </div>
                </div>
              ))}
          </div>
        )}
      </div>

      <Footer />
    </div>
  );
};

export default Deletepage;
