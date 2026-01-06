"use client";

import React, { useEffect, useState } from "react";
import styled from "styled-components";

const Container = styled.div`
  display: inline-flex;
  flex-direction: column;
  height: 1.1em;
  line-height: 1.1em;
  overflow: hidden;
  vertical-align: bottom;

  /* FIX: Default to left for desktop inline flow */
  text-align: left;
  margin-left: 12px;

  /* FIX: On Mobile/Split screen, remove margin and center */
  text-align: center;
  @media (max-width: 768px) {
    margin-left: 0;
    width: 100%; /* Take full width to force centering */
    align-items: center; /* Center the flex children */
  }
`;

const SlideTrack = styled.div<{ $offset: number; $animate: boolean }>`
  transform: translateY(-${(p) => p.$offset * 1.1}em);
  transition: ${(p) =>
    p.$animate ? "transform 0.5s cubic-bezier(0.2, 0.8, 0.2, 1)" : "none"};
`;

const Item = styled.div`
  height: 1.1em;
  color: #58a6ff;
  font-weight: 800;
  align-self: center;
  align-items: center;
  white-space: nowrap;
  /* Cyber Glow */
  text-shadow: 0 0 20px rgba(88, 166, 255, 0.4);
`;

const words = ["Hyperscale Ingestion.", "HFT Systems.", "Realtime Debugging."];

export default function HeroRollingText() {
  const [index, setIndex] = useState(0);
  const [animate, setAnimate] = useState(true);

  useEffect(() => {
    const interval = setInterval(() => {
      setAnimate(true);
      setIndex((prev) => prev + 1);
    }, 2500); // Wait 2.5s before sliding

    return () => clearInterval(interval);
  }, []);

  // The Infinite Loop Trick
  useEffect(() => {
    // If we reached the clone (last item), wait for animation to finish, then snap back
    if (index === words.length) {
      const timeout = setTimeout(() => {
        setAnimate(false); // Disable transition
        setIndex(0); // Snap to top instantly
      }, 500); // Matches transition duration (0.5s)
      return () => clearTimeout(timeout);
    }
  }, [index]);

  return (
    <Container>
      <SlideTrack $offset={index} $animate={animate}>
        {words.map((word, i) => (
          <Item key={i}>{word}</Item>
        ))}
        {/* The Clone for the smooth loop */}
        <Item>{words[0]}</Item>
      </SlideTrack>
    </Container>
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
        padding: "0 4px",
        fontWeight: 800,
        color: "#fff",
        cursor: "pointer",
      }}
    >
      <span style={{ position: "relative", zIndex: 1 }}>{children}</span>
      <div
        style={{
          content: '""',
          position: "absolute",
          bottom: "2px",
          left: 0,
          width: "100%",
          height: "30%",
          background: "#1f6feb",
          zIndex: 0,
          transform: "skew(-12deg) rotate(-2deg)",
          opacity: 0.8,
          borderRadius: "4px",
          transition: "all 0.2s ease",
        }}
        className="highlight-bg"
      />
      <style jsx>{`
        span:hover .highlight-bg {
          height: 90%;
          bottom: 5%;
          background: #58a6ff;
          opacity: 1;
          box-shadow: 0 0 15px #58a6ff;
        }
      `}</style>
    </span>
  );
}
