import { useEffect, useState } from "react";
import "./Comment.css";
import axios from "axios";
import Reply from "./Reply";

const Comment = (props) => {
  const [nowComment, setNowComment] = useState("");
  const [nowID, setNowID] = useState("");
  const [nowPW, setNowPW] = useState("");
  const [comData, setComData] = useState([]);
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
      postid: props.id,
      comments: nowComment,
      comid: nowID,
      compw: nowPW,
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
        .post(process.env.REACT_APP_HOST+ "/api/comment", comData)
        .then((response) => {
          resetReply();
        })
        .catch((error) => {
          if (error.response.status === 500) {
            console.log(error);
            alert("서버에 문제가 생겨 현재 답글을 작성할 수 없습니다.");
          } else if (error.response.status === 400) {
            alert("특수문자 ' 은 입력하실 수 없습니다.");
          } else if (error.response.status === 406){
            alert("PASSWORD는 1~8자리의 숫자만 입력 가능합니다.")
          } else if (error.response.status === 411){
            alert("빈 칸이 존재합니다.")
          } else {
            console.log(error);
          }
        });
  };

  // post id로 해당 post의 comments get
  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST+ "/api/comment/" + props.id)
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
      .post(
        process.env.REACT_APP_HOST+ "/api/comment?commentid=" +
          value.uniqueid +
          "&inputpw=" +
          deletePW
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
    if (reply === value.uniqueid && passwordComment === 0) {
      setReply(0);
    } else {
      setPasswordComment(0);
      setReply(value.uniqueid);
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
      axios
        .post(process.env.REACT_APP_HOST+ "/api/reply/" + value, comData)
        .then((response) => {
          resetReply();
          setReply(0);
        })
        .catch((error) => {
          if (error.response.status === 500) {
            console.log(error);
            alert("서버에 문제가 생겨 현재 답글을 작성할 수 없습니다.");
          } else if (error.response.status === 400) {
            alert("특수문자 ' 은 입력하실 수 없습니다.");
          } else if (error.response.status === 406){
            alert("PASSWORD는 1~8자리의 숫자만 입력 가능합니다.")
          } else if (error.response.status === 411){
            alert("빈 칸이 존재합니다.")
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
              <div>
                <div
                  className={
                    item.isadmin === 1
                      ? "comment-box__adminwriter"
                      : "comment-box__writer"
                  }
                  onClick={() => ReplyHandler(item)}
                >
                  {item.comid}
                </div>
                <div className="comment-box">
                  <div className="comment-delete">
                    <div>{item.comments}</div>
                  </div>
                  <div className="comment-delete__button">
                    <h2 onClick={() => showPasswordInput(item.uniqueid)}>X</h2>
                  </div>
                </div>
                {passwordComment === item.uniqueid ? (
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
                      onClick={() => CheckPasswordHandler(item)}
                    />
                  </div>
                ) : (
                  ""
                )}

                {reply === item.uniqueid && passwordComment === 0 && (
                  <div className="reply-container__write">
                    <textarea
                      className="comment"
                      placeholder={"REPLY TO " + item.comid}
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
                        onClick={() => replySendHandler(item.uniqueid)}
                      />
                    </div>
                  </div>
                )}
                <Reply id={item.uniqueid} rerender={isFinished}/>
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
