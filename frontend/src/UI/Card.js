import "./Card.css";
import { useNavigate } from "react-router-dom";

const Card = (props) => {
  const navigate = useNavigate();

  const cardClickHandler = (postID) => {
    // 버튼 클릭 시 특정 URL로 이동
    navigate("/post/"+postID)
  };

  return (
    <div>
      {props.postdata && (
        <div className="cardcontainer">
          {props.postdata.map((item, index) => (
            <div>
              <h2 className="postcard" onClick={()=>{cardClickHandler(item.Id)}}>
                <div>
                  <img
                    className="postcard-image"
                    alt="my"
                    src={process.env.REACT_APP_HOST+ "/api/assets/"+item.ImagePath}
                  />
                </div>
                <div className="postcard-text">
                  <p>{item.Title}</p>
                </div>
                <div className="postcard-tag">
                  <p>{item.Tag}</p>
                </div>
                <div className="postcard-date">
                  <p className="postcard-date__box">{item.Datetime}</p>
                </div>
              </h2>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Card;
