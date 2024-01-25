import "./Button.css";
import { useState, useEffect } from "react";
import axios from "axios";

const Button = (props) => {
  const [responseData, setResponseData] = useState(null);
  const [title, setTitle] = useState(`" CHOIGONYOK "`);
  const [animate, setAnimate] = useState(true);
  const [PostData, setPostData] = useState({ tags: "ALL" });
  const [tagsdata, setTagsData] = useState([]);
  const [allPostCount, setAllPostCount] = useState({});
  const [clicked, setClicked] = useState(false);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/count")
      .then((response) => {
        setAllPostCount(response.data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  useEffect(() => {
    // POST 요청 보내기
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts?tag=" + PostData.tags)
      .then((response) => {
        // 응답 데이터 수신
        const jsonArray = Object.values(response.data);
        props.onSeeTaggedPost(jsonArray);
      })
      .catch((error) => {
        console.error(error);
      });
  }, [PostData]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/tags")
      .then((response) => {
        setTagsData([...response.data]);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const ClickHandler = (value) => {
    setPostData({ tags: value });
    setTitle(`" ` + value + ` "`);
    setAnimate(!animate);
  };

  const AnimationHandler = () => {
    setAnimate(!animate);
  };
  return (
    <div>
      <h1
        className={animate ? "fadein tags" : "tags"}
        onAnimationEnd={AnimationHandler}
      >
        {title}
      </h1>
      <div cypress="tag" className="container">
        <input
          type="button"
          className={
            "ALL" === PostData.tags ? "tags-button__clicked" : "tags-button"
          }
          value="ALL"
          onClick={() => ClickHandler("ALL")}
        />
        {tagsdata.map((item, index) => (
          <input
            type="button"
            className={
              item.tag === PostData.tags
                ? "tags-button__clicked"
                : "tags-button"
            }
            value={item.tag + '(' + item.num + ')'}
            onClick={() => ClickHandler(item.tag)}
          />
        ))}
      </div>
    </div>
  );
};

export default Button;
