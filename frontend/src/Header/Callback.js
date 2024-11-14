import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';


const Callback = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code");
    const state = urlParams.get("state");
    if (code && state) {
      axios.get(process.env.REACT_APP_HOST+'/api/github/callback?code='+code+"&state="+state)
      .then(response => {
        localStorage.setItem("jwt_token", response.data.token);
        navigate("/");
      })
      .catch(error => {
        if (error.response.status === 403) {
          alert("허가되지 않은 사용자입니다.")
          navigate("/");
        }
        console.error('Error fetching progress:', error);
      });
    } else {
      navigate("/login");
    }
  },[])
  
  return
};

export default Callback;