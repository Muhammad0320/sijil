"use client";

import React from "react";
import styled, { keyframes } from "styled-components";

const roll = keyframes`
  0% { transform: translateY(0); }
  25% { transform: translateY(-100%); }
  50% { transform: translateY(-200%); }
  75% { transform: translateY(-300%); }
  100% { transform: translateY(0); }
`;

// Wrapper to keep "The Log Engine for" and the Rolling text on one line
const Wrapper = styled.span`
  display: inline-flex;
  align-items: center; /* Vertical center alignment */
  justify-content: center;
  flex-wrap: wrap; /* Allow wrapping on small screens */
  gap: 12px; /* Space between "for" and the rolling words */
`;

const Container = styled.div`
  display: inline-block;
  height: 1.2em; /* Locked height matches font line-height */
  line-height: 1.2em;
  overflow: hidden;
  vertical-align: bottom;
  text-align: left;
`;

const RollList = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  animation: ${roll} 8s cubic-bezier(0.23, 1, 0.32, 1) infinite;
`;

const RollItem = styled.li`
  height: 1.2em;
  line-height: 1.2em;
  display: block;
  color: #58a6ff;
  font-weight: 800;
`;

const Underline = styled.svg`
  position: absolute;
  bottom: -8px; /* Moved closer */
  left: 0;
  width: 100%;
  height: 12px;
  z-index: -1;

  path {
    stroke: #1f6feb; /* Darker, Deeper Blue */
    stroke-width: 4; /* Thinner (was 8) */
    stroke-linecap: round;
    stroke-dasharray: 400;
    stroke-dashoffset: 400;
    fill: none;
    opacity: 0.9; /* More visible */
    animation: draw 1s ease-out forwards 0.5s;
  }

  @keyframes draw {
    to {
      stroke-dashoffset: 0;
    }
  }
`;

export default function HeroRollingText() {
  return (
    <Wrapper>
      The Log Engine for
      <Container>
        <RollList>
          <RollItem>Hyperscale Ingestion.</RollItem>
          <RollItem>HFT Systems.</RollItem>
          <RollItem>Realtime Debugging.</RollItem>
          <RollItem>Hyperscale Ingestion.</RollItem>
        </RollList>
      </Container>
    </Wrapper>
  );
}

export function HandDrawnHighlight({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <span
      style={{
        position: "relative",
        display: "inline-block",
        whiteSpace: "nowrap",
      }}
    >
      {children}
      <Underline viewBox="0 0 200 9" preserveAspectRatio="none">
        <path d="M2.00025 6.99997C25.7201 5.2046 82.5279 -0.963287 197.999 2.00003" />
      </Underline>
    </span>
  );
}
