import Header from "../Header/Header";
import "./Deletepage.css";
import Footer from "../UI/Footer";
import axios from "axios";
import { useEffect, useState } from "react";
import MDEditor from "@uiw/react-md-editor";
import Writepage from "./Writepage"
import { useNavigate } from "react-router-dom";

const Deletepage = () => {
  const navigator = useNavigate();

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/login")
      .catch((error) => {
        if (error.response.status === 401) {
          navigator("/login");
        } else {
          console.error(error);
        }
      });
  }, []);

  axios.defaults.withCredentials = true;

  const [changeEvent, setChangeEvent] = useState(false);
  const [md, setMD] = useState("");
  const [titleText, setTitleText] = useState("");
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

  const postHandler = () => {
    const postdata = {
      title: titleText,
      tags: tagText,
      writeTime: dateText,
      text: bodyText,
    };
    axios
      .put(process.env.REACT_APP_HOST + "/api/posts/" + id, postdata, {
        withCredentials: true,
      })
      .then((response) => {
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert(`특수문자 ' 가 입력된 곳이 존재합니다. 수정해주세요.`);
        } else if (error.response.status === 401) {
          alert("로그인이 안된 사용자는 게시글 수정 권한이 없습니다!");
        } else {
          console.log(error);
        }
      });

    axios
      .put(process.env.REACT_APP_HOST + "/api/posts/" + id + "/images", imageNames, {
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
  };

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
        alert("로그인이 안된 사용자는 게시글 삭제 권한이 없습니다!");
      });
  };

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
        setChangeEvent(!changeEvent);
      })
      .catch((error) => {
        console.error(error);
      });

    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + value + "/images")
      .then((response) => {
        console.log(response.data);
        setImageNames([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const titleHandler = (e) => {
    setTitleText(e.target.value);
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
      .delete(process.env.REACT_APP_HOST + "/api/comments/" + value +
        "?type=admin")
      .then((response) => {
        setPostRequest(!postRequest);
      })
      .catch((error) => {
        if (error.response.status === 401) {
          console.log(error);
          alert("로그인이 안된 사용자는 댓글 삭제 권한이 없습니다!");
        } else {
          console.log(error);
        }
      });
  };

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
    let afterThumbnail = imageNames.filter((element) => element.id === newThumbnailImage.id)
    beforeThumbnail[0].thumbnail = "0";
    afterThumbnail[0].thumbnail = "1";

    setImageNames(
      [...afterThumbnail, ...originalImages.filter((element) => element.imageName !== beforeThumbnail[0].imageName).filter((element) => element.imageName !== afterThumbnail[0].imageName), ...beforeThumbnail]
    );
  }
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
                <div className="image-thumbnail-container">
                  {imageNames.map((item, index) => (
                    item.thumbnail === "1" ?
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
                  ))}
                </div>
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
