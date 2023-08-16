import { useState } from "react";
import Header from "../Header/Header";
import Footer from "../UI/Footer";
import "./Loginpage.css";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const Loginpage = () => {
  axios.defaults.withCredentials = true;

  const navigate = useNavigate();
  const [id, setID] = useState();
  const [pw, setPW] = useState();

  const idHandler = (e) => {
    setID(e.target.value);
  };
  const pwHandler = (e) => {
    setPW(e.target.value);
  };
  const loginHandler = () => {
    const logindata = { id: id, pw: pw };
    axios
      .post(process.env.REACT_APP_HOST+"/api/login", logindata)
      .then((response) => {
        alert("로그인 성공!");
        navigate("/")
      })
      .catch((error) => {
        alert("ID 혹은 PASSWORD가 틀렸습니당~!");
      });
  };

  return (
    <div>
      <Header />
      <div className="admin-container">
        <div className="admin-main">LOGIN</div>
        <div className="login-idpw">
          <input type="text" placeholder="ID" value={id} onChange={idHandler} />
          <input
            type="password"
            placeholder="PASSWORD"
            value={pw}
            onChange={pwHandler}
          />
        </div>
        <div className="login-button">
          <input
            type="button"
            className="admin-button"
            value="AUTHORIZATION"
            onClick={loginHandler}
          />
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Loginpage;
