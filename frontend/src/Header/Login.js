import React, { useState } from 'react';
import styled from 'styled-components';
import { FaGithub, FaInbox } from 'react-icons/fa';
import axios from 'axios';
import Header from './Header';
import Footer from '../UI/Footer';

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 80vh;
  color: #fff;
`;

const Title = styled.h1`
  font-size: 32px;
  text-align: center;
  margin-bottom: 20px;
`;

const Description = styled.p`
  font-size: 18px;
  margin-bottom: 40px;
  color: #fff;
  text-align: center;
  max-width: 400px;
`;

const ButtonContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
`;

const GitHubButton = styled.button`
  display: flex;
  align-items: center;
  padding: 10px 20px;
  
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background-color: #333;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;

  &:hover {
    background-color: #444;
    transform: scale(1.05);
  }

  &:active {
    transform: scale(0.95);
  }

  svg {
    margin-right: 10px;
    font-size: 24px;
  }
`;

const Login = () => {
  return (
    <div className="page">
    <Header/>
      <Container>
          <Title>Sign In</Title>
          <Description>
            Sign in with GitHub to continue and access.
          </Description>
            <ButtonContainer>
              <a href={process.env.REACT_APP_HOST+'/api/login'}>
                <GitHubButton>
                  <FaGithub /> Sign in with GitHub
                </GitHubButton>
              </a>
            </ButtonContainer>
      </Container>
    <Footer/>
    </div>
  );
};

export default Login;