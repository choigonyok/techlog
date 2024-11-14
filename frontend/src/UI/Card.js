import "./Card.css";
import { useNavigate } from "react-router-dom";

const Card = (props) => {
  const navigate = useNavigate();
  var count = 0
  
  const cardClickHandler = (postID) => {
    navigate("/posts/" + postID)
  };

  const finishLoading = () => {
    count += 1;
    console.log(count, " ", props.postdata.length)
    if (count === props.postdata.length) {
      props.onFinishCard()
    }
  };
  

  return (
    <div>
      {props.postdata && (
        <div className="cardcontainer">
          {props.postdata.map((item, index) => (
            <div>
              <h2 className="postcard" onClick={() => { cardClickHandler(item.id) }}>
                <div>
                  <img
                    className="postcard-image"
                    alt="my"
                    src={process.env.REACT_APP_HOST + "/api/posts/" + item.id + "/thumbnail"}
                    onLoad={finishLoading}
                  />
                </div>
                <div className="postcard-text">
                  <p>{item.title}</p>
                </div>
                <div className="postcard-subtitle">
                  {item.subtitle === "" ?
                    ""
                    :
                    <p>- {item.subtitle}</p>
                  }
                </div>
                <div className="postcard-tag">
                  <p>{Array.from(item.tags).join(", ")}</p>
                </div>
                <div className="postcard-date">
                  <p className="postcard-date__box">{item.writeTime}</p>
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
