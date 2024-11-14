import React from 'react';
import styled, { keyframes } from 'styled-components';

const rotate = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
`;

const Spinner = styled.div`
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top: 4px solid #3498db;
  border-bottom: 4px solid white;
  border-left: 4px solid white;
  border-right: 4px solid white;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: ${rotate} 1s linear infinite;
  position: fixed;
`;

const LoaderContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  height: 100vh;
  justify-content: center;
  background-color: rgb(24, 26, 28, 1);
  color: #333;
  z-index: 9999;
  position: fixed;
`;

const Loading = () => {
  return (
    <LoaderContainer>
      <Spinner />
    </LoaderContainer>
  );
};

export default Loading;