"use client";

import Link from "next/link";
import styled, { keyframes } from "styled-components";

const float = keyframes`
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
`;

export const AuthContainer = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #050505;
  background-image: radial-gradient(
      circle at 15% 50%,
      rgba(88, 166, 255, 0.1),
      transparent 25%
    ),
    radial-gradient(circle at 85% 30%, rgba(35, 134, 54, 0.1), transparent 25%);
`;

export const AuthCard = styled.div`
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: rgba(22, 27, 34, 0.6); // Glass
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  animation: ${float} 6s ease-in-out infinite;
`;

export const AuthTitle = styled.h1`
  color: #fff;
  font-size: 24px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 32px;
  letter-spacing: -0.5px;
  background: linear-gradient(to right, #fff, #8b949e);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
`;

export const StyledAuthLink = styled(Link)`
  display: block;
  text-align: center;
  margin-top: 24px;
  font-size: 13px;
  color: #8b949e;
  text-decoration: none;
  transition: color 0.2s;

  &:hover {
    color: #58a6ff;
    text-decoration: underline;
  }
`;
