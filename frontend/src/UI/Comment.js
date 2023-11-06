import { useEffect, useState } from "react";
import "./Comment.css";
import axios from "axios";
import Reply from "./Reply";

const Comment = (props) => {
  const [nowComment, setNowComment] = useState("");
  const [nowID, setNowID] = useState("");
  const [nowPW, setNowPW] = useState("");
  const [comData, setComData] = useState([]);
  const [replyData, setReplyData] = useState([]);
  const [isFinished, setIsFinished] = useState(false);
  const [comInfo, setComInfo] = useState([]);
  const [passwordComment, setPasswordComment] = useState(0);
  const [deletePW, setDeletePW] = useState("");
  const [reply, setReply] = useState(0);

  const resetReply = () => {
    setNowComment("");
    setNowID("");
    setNowPW("");
    setIsFinished(!isFinished);
  };

  // 댓글용
  useEffect(() => {
    setComData({
      text: nowComment,
      writerID: nowID,
      writerPW: nowPW,
    });
  }, [nowComment, nowID, nowPW]);

  const commentHandler = (e) => {
    if (e.target.value.length <= 500) {
      setNowComment(e.target.value);
    } else {
      alert(
        "댓글 최대 길이 제한은 500자입니다! 추가적으로 하실 말씀이 있으시면 achoistic98@naver.com 로 메일 주세요 :)"
      );
    }
  };
  const commentIDHandler = (e) => {
    if (e.target.value.length <= 13) {
      setNowID(e.target.value);
    } else {
      alert("NICKNAME 최대 길이 제한은 13자입니다!");
    }
  };
  const commentPWHandler = (e) => {
    if (e.target.value.length <= 8) {
      setNowPW(e.target.value);
    } else {
      alert("PASSWORD 최대 길이 제한은 8자입니다!");
    }
  };
  const commentSendHandler = () => {
    axios
      .post(process.env.REACT_APP_HOST + "/api/posts/"+props.id+"/comment", comData)
      .then((response) => {
        if (response.status === 204) {
          alert("빈 칸이나 유효하지 않은 입력이 있습니다. 비밀번호는 최대 8자리의 숫자만 입력 가능합니다.")
          return;
        }
        resetReply();
      })
      .catch((error) => {
        if (error.response.status === 500) {
          console.log(error);
          alert("서버에 문제가 생겨 현재 답글을 작성할 수 없습니다.");
        }
        else {
          console.log(error);
        }
      });
  };

  // post id로 해당 post의 comments get
  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + props.id + "/comments")
      .then((response) => {
        setComInfo([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [props.id, isFinished]);

  const showPasswordInput = (value) => {
    if (passwordComment === value) {
      setReply(0);
      setPasswordComment(0);
      resetReply();
    } else {
      setPasswordComment(value);
      resetReply();
    }
  };

  const CheckPasswordHandler = (value) => {
    axios
      .delete(
        process.env.REACT_APP_HOST +
        "/api/posts/"+props.id+"/comments/" + value + "?password=" + deletePW
      )
      .then((response) => {
        alert("댓글이 삭제되었습니다.");
        setPasswordComment(0);
        setIsFinished(!isFinished);
      })
      .catch((error) => {
        if (error.response.status === 400) {
          console.log(error);
          alert("PASSWORD가 틀렸습니다.");
        } else {
          console.log(error);
          alert(error);
        }
      });
  };

  const DeletePasswordHandler = (e) => {
    if (e.target.value.length <= 8) {
      setDeletePW(e.target.value);
    } else {
      alert("PASSWORD 최대 길이 제한은 8자입니다!");
    }
  };

  const ReplyHandler = (value) => {
    resetReply();
    if (reply === value.id && passwordComment === 0) {
      setReply(0);
    } else {
      setPasswordComment(0);
      setReply(value.id);
    }
  };

  const replySendHandler = (value) => {
    // item.uniqueid으로 대댓글 만들기
    // if (
    //   comData.comid === "" ||
    //   comData.comments === "" ||
    //   comData.compw === ""
    // ) {
    //   alert("작성되지 않은 항목이 존재합니다.");
    // } else {
      console.log(comData);
      console.log(comData);
      console.log(comData);
    axios
      .post(process.env.REACT_APP_HOST + "/api/posts/" + props.id + "/comments/" + value + "/reply", comData)
      .then((response) => {
        if (response.status === 204) {
          alert("빈 칸이 존재합니다.");
          return
        }
        resetReply();
        setReply(0);
      })
      .catch((error) => {
        if (error.response.status === 500) {
          console.log(error);
          alert("서버에 문제가 생겨 현재 답글을 작성할 수 없습니다.");
        } else if (error.response.status === 400) {
          alert("비밀번호는 최대 8자리의 숫자만 입력 가능합니다.");
        } else {
          console.log(error);
        }
      });
    // }
  };

  return (
    <div>
      <div className="comment-container">
        {comInfo &&
          comInfo.map((item, index) => {
            return (
              <div className="comment-container__2">
                <div
                  className={
                    item.admin === 1
                      ? "comment-box__adminwriter"
                      : "comment-box__writer"
                  }
                  onClick={() => ReplyHandler(item)}
                >
                  {item.writerID}
                </div>
                <div className="comment-box">
                  <div className="comment-delete">
                    <div>{item.text}</div>
                  </div>
                  <div className="comment-delete__button">
                    <h2 onClick={() => showPasswordInput(item.id)}>X</h2>
                  </div>
                </div>
                {passwordComment === item.id ? (
                  <div className="password-container">
                    <input
                      type="password"
                      placeholder="PASSWORD"
                      className="password-text"
                      onChange={DeletePasswordHandler}
                    />
                    <input
                      type="button"
                      value="DELETE"
                      className="comment-button__submit"
                      onClick={() => CheckPasswordHandler(item.id)}
                    />
                  </div>
                ) : (
                  ""
                )}

                {reply === item.id && passwordComment === 0 && (
                  <div className="reply-container__write">
                    <textarea
                      className="comment"
                      placeholder={"REPLY TO " + item.writerid}
                      onChange={commentHandler}
                      value={nowComment}
                    />
                    <div className="comment-buttons">
                      <input
                        type="text"
                        placeholder="NICKNAME"
                        onChange={commentIDHandler}
                        value={nowID}
                      />
                      <input
                        type="password"
                        placeholder="PASSWORD"
                        onChange={commentPWHandler}
                        value={nowPW}
                      />
                      <input
                        type="button"
                        className="comment-button__submit"
                        value="POST"
                        onClick={() => replySendHandler(item.id)}
                      />
                    </div>
                  </div>
                )}
                <Reply id={item.id} rerender={isFinished} postid={props.id} />
              </div>
            );
          })}
      </div>
      {reply === 0 && (
        <div className="comment-container__container">
          <div className="comment-container__write">
            <textarea
              className="comment"
              placeholder="PLEASE LEAVE A COMMENT !"
              onChange={commentHandler}
              value={nowComment}
            />
            <div className="comment-buttons">
              <input
                type="text"
                placeholder="NICKNAME"
                onChange={commentIDHandler}
                value={nowID}
              />
              <input
                type="password"
                placeholder="PASSWORD"
                onChange={commentPWHandler}
                value={nowPW}
              />
              <input
                type="button"
                className="comment-button__submit"
                value="POST"
                onClick={commentSendHandler}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Comment;
