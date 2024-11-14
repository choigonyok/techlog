import "./Homepage.css";
import Header from "../Header/Header";
import Button from "../UI/Button";
import Footer from "../UI/Footer";
import Loading from "../UI/Loading";
import Card from "../UI/Card";
import axios from "axios";
import profileimage from "../Assets/IMG_0071 2.jpg";
import github from "../Assets/Icons/github-icon.png";
import instagram from "../Assets/Icons/instagram-icon.png";
import youtube from "../Assets/Icons/youtube-icon.png";
import { useState, useEffect } from "react";

const Homepage = () => {
  axios.defaults.withCredentials = true;

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST + "/api/visitor")
      .then((response) => {
        setTotalNum(response.data.total);
        setVisitNum(response.data.today);
      })
      .catch((error) => {
        console.error(error);
      });

    if (postData.length === 0 ) {
      setIsLoading(false)
    }
  }, []);

  const [totalNum, setTotalNum] = useState("");
  const [visitNum, setVisitNum] = useState("");
  const [changeEvent, setChangeEvent] = useState(false);
  const [postData, setPostData] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  const seeTaggedPostHandler = (taggedPostData) => {
    setPostData(taggedPostData);
  };

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [changeEvent]);

  
  const handleFinish = () => {
    setIsLoading(false)
  };


  return (
    <div className="page">
      <Header />
      {!isLoading ? "" :<Loading/>}
      <div className="introduce">
        <div className="visitnum">
          TODAY : {visitNum} / TOTAL : {totalNum}
        </div>
        <div className="home-image__container">
          <img className="home-image" alt="my" src={profileimage} />
        </div>
        <div className="icon-container">
          <a href="https://github.com/choigonyok">
            <img className="icon-image" alt="my" src={github} />
          </a>
          <a href="https://www.instagram.com/choigonyok">
            <img className="icon-image" alt="my" src={instagram} />
          </a>
          <a href="https://www.youtube.com/channel/UCjcaraxwh4Fz6KYJwfE0-Hg">
            <img className="icon-image" alt="my" src={youtube} />
          </a>
        </div>
        <div className="introduce-text__text">
          꾸준함이란 도구로 성장하기를 즐기는 DevSecOps 엔지니어 최윤석입니다.
        </div>
        <div className="introduce-text__year">
          <div>
            2017.03.02~2025.02.28 &nbsp;&nbsp;&nbsp; Kyunghee Univ. Computer Engineering
          </div>
          <div>
            2024.03.02~2024.08.31 &nbsp;&nbsp;&nbsp; SLEXN, inc. Internship( Infrastructure Engineer )
          </div>
        </div>
      </div>
      <Button onSeeTaggedPost={seeTaggedPostHandler} />
      {postData && <Card postdata={postData} onFinishCard={handleFinish}/>}
      <Footer />
    </div>
  );
};

export default Homepage;
