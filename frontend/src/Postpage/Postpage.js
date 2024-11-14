import Header from "../Header/Header";
import "./Postpage.css";
import Footer from "../UI/Footer";
import Card from "../UI/Card";
import { useEffect, useState } from "react";
import axios from "axios";
import { useParams } from "react-router-dom";
import MDEditor from "@uiw/react-md-editor";
import { useRef } from "react";
import Loading from "../UI/Loading";

const Postpage = () => {
  axios.defaults.withCredentials = true;

  let { postid } = useParams();

  const [changeEvent, setChangeEvent] = useState(false);
  const mounted = useRef(false);
  const [postData, setPostData] = useState({});
  const [relatedPostData, setRelatedPostData] = useState([]);
  const [postImages, setPostImages] = useState([]);
  const [thumbnail, setThumbnail] = useState([""]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // axios
    //   .get(process.env.REACT_APP_HOST + "/api/visitor")
    //   .then((response) => { })
    //   .catch((error) => {
    //     console.error(error);
    //   });
    if (relatedPostData.length === 0 ) {
      setIsLoading(false)
    }
  }, []);

  useEffect(() => {
    window.scrollTo({
      top: 0, 
      left: 0, 
      behavior: "instant"
    });
  }, [changeEvent]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/posts/" + postid)
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
        .get(process.env.REACT_APP_HOST + "/api/posts?tag="+postData.tags)
        .then((response) => {
          const jsonArray = Object.values(response.data);
          setRelatedPostData(
            jsonArray.filter((post) => String(post.id) !== postid)
          );
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [postData]);

  const handleFinish = () => {
    setIsLoading(false)
  };


  return (
    <div>
      <Header />
      {!isLoading ? "" :<Loading/>}
      <br />
      <br />
      <br />
      <br />
      <div>
        <div className="image-container">
          <img
            className="image"
            alt="my"
            src={
              process.env.REACT_APP_HOST + "/api/posts/" + postid + "/thumbnail"
            }
          />
        </div>
        <div className="post-title">
          <div className="post-tagsbox">
            <button className="post-tags__button">
              {postData.tags? Array.from(postData.tags).join(", "): ""}
            </button>
          </div>
          <p className="post-title__item">{postData.title}</p>
          <p className="written-date">{postData.writeTime}</p>
        </div>
        <div>
          <MDEditor.Markdown className="post-body" source={postData.text} />
        </div>
      </div>
      {/* <div className="related-post__container">
        <div className="related-post__content">- COMMENTS -</div>
      </div>
      <Comment id={postid} /> */}
      <div className="related-post__container">
        <div className="related-post__content">- RELATED POSTS -</div>
      </div>

      {relatedPostData && <Card postdata={relatedPostData} onFinishCard={handleFinish} />}
      <Footer />
    </div>
  );
};
export default Postpage;
