import Header from "../Header/Header";
import MDEditor, { code } from "@uiw/react-md-editor";
import "./Writepage.css";
import Footer from "../UI/Footer";
import { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useRef } from "react";
import Markdown from "../UI/Markdown";

const Writepage = () => {
  const navigate = useNavigate();

  // useEffect(()=>{
  //   axios
  //   .get(process.env.REACT_APP_HOST + "/login")
  //   .then((response) => {
  //     if (response.status !== 202) {
  //       navigate("/");
  //     }
  //   })
  //   .catch((error) => {
  //     if (error.response) {
  //       if (error.response.status === 401) {
  //         alert("로그인이 안된 사용자는 게시글 작성 권한이 없습니다!");
  //         window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
  //       } else {
  //         console.error(error);
  //       }
  //     }
  //   });
  // },[])

  axios.defaults.withCredentials = true;

  const [md, setMD] = useState("");
  const [titleText, setTitleText] = useState("");
  const [subtitleText, setSubtitleText] = useState("");
  const [tagText, setTagText] = useState("");
  const [dateText, setDateText] = useState("");
  const [bodyText, setBodyText] = useState("");
  const [img, setIMG] = useState([]);
  const [imgName, setImgName] = useState([]);
  const [unlock, setUnLock] = useState(false);
  const mounted = useRef(false);

  // const postHandler = () => {
  //   axios
  //     .get(process.env.REACT_APP_HOST + "/oauth2/auth")
  //     .then((response) => {
  //       if (response.status === 202) {
  //         requestPosting()
  //       }
  //     })
  //     .catch((error) => {
  //       if (error.response) {
  //         if (error.response.status === 401) {
  //           window.location.href = "https://www.choigonyok.com/oauth2/sign_in";
  //         } else {
  //           console.error(error);
  //         }
  //       }
  //     })
  // };

  const requestPosting = () => {
    const tags = tagText.split(',')
    const trimedTags = tags.map((item)=>item.trim())
    const post = {
      title: titleText,
      tags: trimedTags,
      text: bodyText,
      subtitle: subtitleText,
      thumbnail_name: imgName[0],
    };
    const formData = new FormData();
    for (let i = 0; i < img.length; i++) {
      formData.append("file", img[i]);
    }
    formData.append('data', JSON.stringify(post));

    axios
      .post(process.env.REACT_APP_HOST + "/api/post", formData, {
        "Content-type": "multipart/form-data",
        withCredentials: true,
      })
      .then((response) => {
        navigate("/");
        // setUnLock(!unlock);
      })
      .catch((error) => {
        console.error(error);
      });
  }

  const deleteWronglyWrittenPost = () => {
    axios
      .delete(process.env.REACT_APP_HOST + "/api/post/0", {
        withCredentials: true,
      })
      .catch((error) => {
        console.error(error);
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

  useEffect(() => {
    setBodyText(md);
  }, [md])

  const imgHandler = (e) => {
    setIMG((img) => [...e.target.files]);

    let f = document.getElementById("imgfile").files;
    if (f.length !== 0) {
      setImgName((element) => [...element, f[0].name]);
    }
  };

  const deleteIMGHandler = (value) => {
    setImgName(
      imgName.filter((element) => String(element) !== imgName.at(value))
    );
    setIMG(
      img.filter((element) => String(element.name) !== img.at(value).name)
    );
  };

  // useEffect(() => {
  //   if (!mounted.current) {
  //     mounted.current = true;
  //   } else {
  //     const formData = new FormData();
  //     for (let i = 0; i < img.length; i++) {
  //       formData.append("file", img[i]);
  //       formData.append('data', JSON.stringify({ textData }));
  //     }
  //     axios
  //       .post(process.env.REACT_APP_HOST + "/api/post/image", formData, {
  //         "Content-type": "multipart/form-data",
  //         withCredentials: true,
  //       })
  //       .then((response) => {
  //         // 응답 데이터 수신
  //         navigate("/");
  //       })
  //       .catch((error) => {
  //         if (error.response.status === 500) {
  //           deleteWronglyWrittenPost();
  //           alert("이미지가 등록되지 않을 채로 글이 작성되었습니다!");
  //         } else {
  //           console.error(error);
  //         }
  //       });
  //   }
  // }, [unlock]);

  return (
    <div>
      <div className="admin-container">
        <div className="admin-titletagdate">
          <input
            type="text"
            placeholder="TAGS : PROJECT / BLOG"
            value={tagText}
            onChange={tagHandler}
          />
          {/* <input type="button" value="적용" onClick={TagClickHandler} /> */}
        </div>
        <div className="admin-titletagdate">
          <input
            type="text"
            placeholder="TITLE : [BLOG #1] 블로그 개발하기"
            value={titleText}
            onChange={titleHandler}
          />
          {/* <input type="button" value="적용" onClick={TitleClickHandler} /> */}
        </div>
        <div className="admin-titletagdate">
          <input
            type="text"
            placeholder="SUBTITLE : 직접 기술 블로그를 만들어보자"
            value={subtitleText}
            onChange={subtitleHandler}
          />
          {/* <input type="button" value="적용" onClick={TitleClickHandler} /> */}
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
        {imgName.map((item, index) => (
          <div className="imgname-container">
            {item}
            <input
              type="button"
              value="X"
              className="imgname-button"
              onClick={() => deleteIMGHandler(index)}
            />
          </div>
        ))}
        <Markdown onTextChange={setBodyText}/>
        <div className="button-container">
          <input
            type="button"
            className="admin-button"
            value="POST 추가하기"
            onClick={requestPosting}
          />
        </div>
      </div>
    </div>
  );
};

export default Writepage;
