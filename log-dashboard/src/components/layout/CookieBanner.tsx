"use client";

import { useState, useEffect } from "react";
import styled from "styled-components";

const Banner = styled.div`
  position: fixed;
  bottom: 24px;
  right: 24px;
  max-width: 350px;
  background: #161b22;
  border: 1px solid #30363d;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
  z-index: 100;
  animation: slideUp 0.5s ease-out;

  @keyframes slideUp {
    from {
      transform: translateY(100px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
`;

const Text = styled.p`
  font-size: 14px;
  color: #8b949e;
  margin-bottom: 16px;
  line-height: 1.5;
`;

const Button = styled.button`
  background: #fff;
  color: #000;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  font-size: 13px;
  width: 100%;
  &:hover {
    background: #e1e4e8;
  }
`;

export default function CookieBanner() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (!localStorage.getItem("sijil_cookie_consent")) {
      setTimeout(() => {
        setVisible(true);
      }, 0);
    }
  }, []);

  const accept = () => {
    localStorage.setItem("sijil_cookie_consent", "true");
    setVisible(false);
  };

  if (!visible) return null;

  return (
    <Banner>
      <Text>
        We use cookies to ensure you get the best experience on Sijil. By
        continuing, you agree to our use of cookies for authentication and
        analytics.
      </Text>
      <Button onClick={accept}>Accept & Continue</Button>
    </Banner>
  );
}
