import "./Homepage.css";
import Header from "../Header/Header";
import Button from "../UI/Button";
import Footer from "../UI/Footer";
import Card from "../UI/Card";
import axios from "axios";
import profileimage from "../Assets/IMG_0071 2.jpg";
import github from "../Assets/Icons/github-icon.png";
import instagram from "../Assets/Icons/instagram-icon.png";
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
        console.log(error);
      });
  }, []);

  const [totalNum, setTotalNum] = useState("");
  const [visitNum, setVisitNum] = useState("");
  const [changeEvent, setChangeEvent] = useState(false);
  const [postData, setPostData] = useState([]);

  const seeTaggedPostHandler = (taggedPostData) => {
    setPostData(taggedPostData);
  };

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [changeEvent]);

  return (
    <div className="page">
      <Header /> {/* 6/2 Header 컴포넌트 재사용 위해서 분리 */}
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
        </div>
        <div className="introduce-text__text">
          꾸준함이란 도구로 성장하기를 즐기는 DevSecOps 엔지니어 최윤석입니다.
        </div>
        <div className="introduce-text__year">
          <div>
            2017.03~ &nbsp;&nbsp;&nbsp; Kyunghee Univ. Computer Engineering
          </div>
          <div>
            2023.12~ 2024.08&nbsp;&nbsp;&nbsp; 여기어때컴퍼니 인프라개발팀 Internship
          </div>
        </div>
      </div>
      <Button onSeeTaggedPost={seeTaggedPostHandler} />
      {postData && <Card postdata={postData} />}
      <Footer />
    </div>
  );
};

export default Homepage;
