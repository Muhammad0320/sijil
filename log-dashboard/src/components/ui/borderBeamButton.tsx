"use client";

import styled, { keyframes, css } from "styled-components";
import { ButtonHTMLAttributes } from "react";

const spin = keyframes`
  from { --angle: 0deg; }
  to { --angle: 360deg; }
`;

// We use CSS custom properties for the rotation to avoid complex DOM nesting if possible,
// but for broad support, the conic-gradient rotation on a pseudo-element is most robust.
const rotateGradient = keyframes`
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
`;

const ButtonContainer = styled.button<{ $isLoading?: boolean }>`
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 48px;
  background: transparent;
  border: none;
  padding: 0;
  cursor: ${(p) => (p.$isLoading ? "wait" : "pointer")};
  outline: none;
  border-radius: 8px;
  overflow: hidden; // Clips the spinning beam
  font-weight: 600;
  font-size: 14px;
  color: white;
  transition: all 0.2s;

  // Inactive state styles
  opacity: ${(p) => (p.$isLoading ? 0.8 : 1)};
  pointer-events: ${(p) => (p.$isLoading ? "none" : "auto")};

  &:focus-visible {
    box-shadow: 0 0 0 2px #58a6ff;
  }
`;

// The Spinning Layer
const BeamLayer = styled.div<{ $active?: boolean }>`
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200%; // Much larger to cover corners
  height: 500%; // Much larger to cover corners
  transform: translate(-50%, -50%);
  z-index: 0;

  // The Beam: A conic gradient with a sharp tail
  background: conic-gradient(
    from 0deg,
    transparent 0%,
    transparent 80%,
    #58a6ff 90%,
    #ffffff 100%
  );

  // Only animate/show when loading
  opacity: ${(p) => (p.$active ? 1 : 0)};
  animation: ${(p) =>
    p.$active
      ? css`
          ${rotateGradient} 2s linear infinite
        `
      : "none"};
  transition: opacity 0.3s ease;
`;

// The Static Background (Mask)
const InnerMask = styled.div`
  position: absolute;
  inset: 1.5px; // The thickness of the border beam
  background: #238636; // The actual button color
  border-radius: 7px; // Slightly smaller radius
  z-index: 1;
  transition: background 0.2s;

  ${ButtonContainer}:hover & {
    background: #2ea043;
  }
`;

// The Text Layer (Topmost)
const ContentLayer = styled.div`
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  gap: 8px;
`;

interface BorderBeamButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading?: boolean;
  children: React.ReactNode;
}

export function BorderBeamButton({
  isLoading,
  children,
  ...props
}: BorderBeamButtonProps) {
  return (
    <ButtonContainer $isLoading={isLoading} {...props}>
      {/* The Beam (Only visible when loading) */}
      <BeamLayer $active={isLoading} />

      {/* The Button Color (Masks the center of the beam) */}
      <InnerMask />

      {/* The Text */}
      <ContentLayer>{children}</ContentLayer>
    </ButtonContainer>
  );
}
