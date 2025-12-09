"use client";

import styled, { keyframes, css } from "styled-components";
import { ButtonHTMLAttributes } from "react";

// 1. The Rotation Animation
const spin = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`;

// 2. The Wrapper
const ButtonContainer = styled.button<{ $isLoading?: boolean }>`
  position: relative;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  border: none;
  background: transparent;
  padding: 0;
  cursor: ${(props) => (props.$isLoading ? "wait" : "pointer")};
  outline: none;
  border-radius: 8px;
  overflow: hidden;
  height: 48px; // Fixed height for consistency

  // Prevent clicking when loading
  pointer-events: ${(props) => (props.$isLoading ? "none" : "auto")};

  &:focus-visible {
    box-shadow: 0 0 0 2px #58a6ff;
  }
`;

// 3. The "Beam" (Hidden layer that spins)
const MovingGradient = styled.div<{ $active?: boolean }>`
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200%;
  height: 200%;
  transform: translate(-50%, -50%);

  background: conic-gradient(
    from 90deg at 50% 50%,
    transparent 0%,
    transparent 40%,
    #238636 50%,
    #ffffff 55%,
    transparent 60%,
    transparent 100%
  );

  // Spin only when active (loading)
  animation: ${spin} 2s linear infinite;
  opacity: ${(props) => (props.$active ? 1 : 0)};
  transition: opacity 0.3s ease;
`;

// 4. The Inner Content
const InnerContent = styled.div`
  position: relative;
  z-index: 1;
  width: calc(100% - 2px);
  height: calc(100% - 2px);
  margin: 1px;
  background: #238636;
  color: white;
  font-weight: 600;
  font-size: 14px;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;

  ${ButtonContainer}:hover & {
    background: #2ea043;
  }
`;

interface BorderBeamButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading?: boolean;
  children: React.ReactNode;
}

export function BorderBeamButton({
  isLoading,
  children,
  ...props
}: BorderBeamButtonProps) {
  return (
    <ButtonContainer type="button" $isLoading={isLoading} {...props}>
      <MovingGradient $active={isLoading} />
      <InnerContent>
        {isLoading ? "Processing..." : children}
      </InnerContent>
    </ButtonContainer>
  );
}