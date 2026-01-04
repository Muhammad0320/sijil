"use client";

import React, { useEffect, useState } from "react";
import styled, { keyframes } from "styled-components";

const roll = keyframes`
  0% { transform: translateY(0); }
  25% { transform: translateY(-100%); }
  50% { transform: translateY(-200%); }
  75% { transform: translateY(-300%); }
  100% { transform: translateY(0); } /* Loop back */
`;

const Container = styled.div`
  display: inline-block;
  height: 1.1em;
  overflow: hidden;
  vertical-align: bottom;
  position: relative;
  width: auto;
  min-width: 300px; /* Prevent layout shift */
  text-align: left;
`;

const RollList = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  animation: ${roll} 8s cubic-bezier(0.23, 1, 0.32, 1) infinite;
`;

const RollItem = styled.li`
  height: 1.1em;
  display: flex;
  align-items: center;
  color: #58a6ff;
  font-weight: 800;
  padding-left: 12px;
`;

// The Hand-Drawn Underline SVG
const Underline = styled.svg`
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100%;
  height: 12px;
  z-index: -1;

  path {
    stroke: #58a6ff;
    stroke-width: 3;
    stroke-linecap: round;
    stroke-dasharray: 400;
    stroke-dashoffset: 400;
    fill: none;
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
    <span
      style={{
        display: "inline-flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      The Log Engine for{" "}
      <Container>
        <RollList>
          <RollItem>Hyperscale Ingestion.</RollItem>
          <RollItem>HFT Systems.</RollItem>
          <RollItem>Realtime Debugging.</RollItem>
          <RollItem>Hyperscale Ingestion.</RollItem>{" "}
          {/* Repeat first for seamless loop */}
        </RollList>
      </Container>
    </span>
  );
}

export function HandDrawnHighlight({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <span style={{ position: "relative", display: "inline-block" }}>
      {children}
      <Underline viewBox="0 0 200 9" preserveAspectRatio="none">
        <path d="M2.00025 6.99997C25.7201 5.2046 82.5279 -0.963287 197.999 2.00003" />
      </Underline>
    </span>
  );
}
