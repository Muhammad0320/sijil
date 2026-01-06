"use client";

import React from "react";
import styled, { keyframes } from "styled-components";

// --- 1. THE ROLLING TEXT ANIMATION ---
// The math: 3 items. Total 100%.
// 0% -> 25%: Stay at 1
// 25% -> 33%: Move to 2
// 33% -> 58%: Stay at 2
// 58% -> 66%: Move to 3
// 66% -> 91%: Stay at 3
// 91% -> 100%: Move back to 1
const snapRoll = keyframes`
  0%, 25% { transform: translateY(0); }
  33%, 58% { transform: translateY(-100%); }
  66%, 91% { transform: translateY(-200%); }
  100% { transform: translateY(0); } /* Instant reset or slide depending on items */
`;

const Wrapper = styled.span`
  display: inline-flex;
  flex-direction: column;
  height: 1.1em; /* Perfectly matches line-height */
  overflow: hidden;
  vertical-align: bottom; /* Aligns with the baseline of "The Log Engine for" */
  text-align: left;
  margin-left: 12px;
`;

const RollList = styled.span`
  display: block;
  animation: ${snapRoll} 6s cubic-bezier(0.5, 0, 0.2, 1) infinite;
`;

const RollItem = styled.span`
  display: block;
  height: 1.1em;
  line-height: 1.1em;
  color: #58a6ff;
  font-weight: 800;
  /* Cyberpunk text shadow for "Luxury" feel */
  text-shadow: 0 0 20px rgba(88, 166, 255, 0.4);
`;

export default function HeroRollingText() {
  return (
    <Wrapper>
      <RollList>
        <RollItem>Hyperscale Ingestion.</RollItem>
        <RollItem>HFT Systems.</RollItem>
        <RollItem>Realtime Debugging.</RollItem>
        {/* Duplicate first item for seamless loop if needed, 
            but for 3-step snap, we usually just snap back to 0. 
            Let's keep 3 distinct items for simplicity. */}
      </RollList>
    </Wrapper>
  );
}

// --- 2. THE CYBER HIGHLIGHTER (Interactive) ---

const squiggly = keyframes`
  0% { background-position: 0 0; }
  100% { background-position: 100% 0; }
`;

const HighlightSpan = styled.span`
  position: relative;
  z-index: 1;
  font-weight: 800;
  cursor: pointer;
  color: #fff;
  white-space: nowrap;
  padding: 0 4px;

  /* The "Underline" is now a pseudo-element highlighter */
  &::before {
    content: "";
    position: absolute;
    bottom: 2px;
    left: 0;
    width: 100%;
    height: 30%; /* Only covers bottom third like a marker */
    background: #1f6feb; /* Deep Blue */
    z-index: -1;
    transform: skew(-12deg) rotate(-2deg);
    transition: all 0.2s ease;
    border-radius: 4px;
    opacity: 0.8;
  }

  /* Hover Effect: Expands and Glows */
  &:hover::before {
    height: 90%; /* Covers full text */
    bottom: 5%;
    transform: skew(-12deg) rotate(0deg) scale(1.05);
    background: #58a6ff; /* Brighter Blue */
    opacity: 1;
    box-shadow: 0 0 15px #58a6ff;
  }
`;

export function HandDrawnHighlight({
  children,
}: {
  children: React.ReactNode;
}) {
  return <HighlightSpan>{children}</HighlightSpan>;
}
