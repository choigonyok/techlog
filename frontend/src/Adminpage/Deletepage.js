import Header from "../Header/Header";
import "./Deletepage.css";
import Footer from "../UI/Footer";
import axios from "axios";
import { useEffect, useState } from "react";
import MDEditor from "@uiw/react-md-editor";
import Writepage from "./Writepage"
import { useNavigate } from "react-router-dom";
import Loading from "../UI/Loading";
import Markdown from "../UI/Markdown";

const Deletepage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("jwt_token");
    if (token === null) {
      navigate("/login");
      return
    }

    axios.post(process.env.REACT_APP_HOST+'/api/token', JSON.stringify({token: token}))
      .then(response => {
        setIsLoading(false)
      })
      .catch(error => {
        if (error.response.status === 403) {
          alert("허가되지 않은 사용자입니다.")
          navigate("/")
        } else if (error.response.status === 401) {
          navigate("/login")
        }
        console.error('Error fetching progress:', error);
        setIsLoading(false)
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
  const [thumbnail, setThumbnail] = useState('');
  const [img, setIMG] = useState([]);
  const [imgName, setImgName] = useState('');
  const [preImageLength, setPreImageLength] = useState(0);

  console.log("NAME:",imgName)
  const postHandler = () => {
    const tags = tagText.split(',')
    const trimedTags = tags.map((item)=>item.trim())

    const postdata = {
      title: titleText,
      tags: trimedTags,
      writeTime: dateText,
      text: bodyText,
      subtitle: subtitleText,
      thumbnail_name: imgName,
    };
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
        navigate("/")
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
        setTagText(Array.from(response.data.tags).join(", "));
        setDateText(response.data.writetime);
        setMD(response.data.text);
        setSubtitleText(response.data.subtitle);
        setChangeEvent(!changeEvent);
      })
      .catch((error) => {
        console.error(error);
      });

    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + value + "/thumbnail")
      .then((response) => {
        // setThumbnail(response.data)
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

  // useEffect(() => {
  //   axios
  //     .get(process.env.REACT_APP_HOST + "/api/comments")
  //     .then((response) => {
  //       setComInfo([...response.data]);
  //     })
  //     .catch((error) => {
  //       console.log(error);
  //     });
  // }, [postRequest]);

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
    setIMG([...e.target.files]);

    let f = document.getElementById("imgfile").files;
    if (f.length !== 0) {
      setImgName(f[0].name);
    }

    if (f.length !== 0) {
      const reader = new FileReader();
      reader.onload = (e) => {
      }
      const url = URL.createObjectURL(e.target.files[0])
      let img = {
        id: 0,
        postid: id,
        imageName: f.name,
        thumbnail: url,
      };
      setThumbnail(img);
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

  return (
    <div>
      <Header />
      {!isLoading ? "" :<Loading/>}
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
          {/* <input
            type="button"
            className="select-button"
            value="COMMENTS"
            onClick={isCommentsHandler}
          /> */}
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
                <div>
                  </div>
                <div className="image-thumbnail-container">
                  {thumbnail ?
                      <img
                      className="image-non-thumbnail-image"
                      alt="my"
                      src={thumbnail.thumbnail}
                    />
                    :
                    <img
                      className="image-non-thumbnail-image"
                      alt="my"
                      src={process.env.REACT_APP_HOST + "/api/posts/" + id + "/thumbnail"}
                    /> 
                  }
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
                <Markdown onTextChange={setMD} value={md}/>
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
        {/* {isComments && (
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
        )} */}
      </div>

      <Footer />
    </div>
  );
};

export default Deletepage;
