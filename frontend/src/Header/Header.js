import "./Header.css";
import { useNavigate } from "react-router-dom";

const Header = () => {
  const navigate = useNavigate();

  const homePageHandler = () => {
    // 버튼 클릭 시 특정 URL로 이동
    navigate("/");
  };

  const githubPageHandler = () => {
    const token = localStorage.getItem("jwt_token");
    if (token === null) {
      navigate("/login");
      return
    }
    navigate("/admin");
  };


  return (
    <div className="header">
      <button className="header-title" onClick={homePageHandler}>
        Techlog.
      </button>
      <div className="header-empty" />
      <div>
        <div onClick={githubPageHandler} className="header-category">
          admin
        </div>
      </div>
    </div>
  );
};
export default Header;
