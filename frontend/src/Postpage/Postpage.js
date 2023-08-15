import Header from "../Header/Header";
import "./Postpage.css";
import Footer from "../UI/Footer";
import Card from "../UI/Card";
import { useEffect, useState } from "react";
import axios from "axios";
import { useParams } from "react-router-dom";
import MDEditor from "@uiw/react-md-editor";
import { useRef } from "react";
import Comment from "../UI/Comment";

const Postpage = () => {
  axios.defaults.withCredentials = true;

  let { postid } = useParams();

  const [changeEvent, setChangeEvent] = useState(false);
  const mounted = useRef(false);
  const [postData, setPostData] = useState({});
  const [relatedPostData, setRelatedPostData] = useState([]);


  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST+ "/api/cookie")
      .then((response) => {})
      .catch((error) => {
        console.log(error);
      });
  }, []);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [changeEvent]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST+ "/api/post/" + postid)
      .then((response) => {
        setPostData(response.data); 
        setChangeEvent(!changeEvent);
      })
      .catch((error) => {
        console.error(error);
      });
  }, [postid]);

  useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
    } else {
      axios
        .post(process.env.REACT_APP_HOST+ "/api/tag", postData)
        .then((response) => {
          const jsonArray = Object.values(response.data);
          setRelatedPostData(
            jsonArray.filter((post) => String(post.Id) !== postid)
          );
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [postData]);

  return (
    <div>
      <Header />
      <br />
      <br />
      <br />
      <br />
      <div>
        <div className="image-container">
          <img className="image" alt="my" src={process.env.REACT_APP_HOST+ "/api/IMAGES/"+postData.ImagePath}/>
        </div>
        <div className="post-title">
          <div className="post-tagsbox">
            <button className="post-tags__button">{postData.Tag}</button>
          </div>
          <p className="post-title__item">{postData.Title}</p>
          <p className="written-date">{postData.Datetime}</p>
        </div>
        <div>
          <MDEditor.Markdown className="post-body" source={postData.Body} />
        </div>
      </div>
      <div className="related-post__container">
        <p className="related-post__content">- COMMENTS -</p>
      </div>
      <Comment id={postid} />
      <div className="related-post__container">
        <p className="related-post__content">- RELATED POSTS -</p>
      </div>

      {relatedPostData && <Card postdata={relatedPostData} />}
      <Footer />
    </div>
  );
};
export default Postpage;
