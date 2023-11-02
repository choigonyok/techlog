import "./Button.css";
import { useState, useEffect } from "react";
import axios from "axios";

const Button = (props) => {
  const [responseData, setResponseData] = useState(null);
  const [title, setTitle] = useState(`" CHOIGONYOK "`);
  const [animate, setAnimate] = useState(true);
  const [PostData, setPostData] = useState({ Tag: "ALL" });
  const [tagsdata, setTagsData] = useState([]);
  const [clicked, setClicked] = useState(false);

  useEffect(() => {
    // POST 요청 보내기
    axios
      .post(process.env.REACT_APP_HOST+ "/api/tag", PostData)
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
     .get(process.env.REACT_APP_HOST+ "/api/tag")
      .then((response) => {
        setTagsData([...response.data]);
        // props.onSeeTaggedPost(jsonArray);
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
      <div className="container">
        <input
          type="button"
          className={
            "ALL" === PostData.Tag ? "tags-button__clicked" : "tags-button"
          }
          value="ALL"
          onClick={() => ClickHandler("ALL")}
        />
        {tagsdata.map((item, index) => (
          <input
            type="button"
            className={
              item.tag === PostData.Tag
                ? "tags-button__clicked"
                : "tags-button"
            }
            value={item}
            onClick={() => ClickHandler(item)}
          />
        ))}
      </div>
    </div>
  );
};

export default Button;
