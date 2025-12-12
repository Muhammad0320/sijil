"use client";

import { useState } from "react";
import styled from "styled-components";
import { Menu, X } from "lucide-react"; // Icons for the toggle

const GridContainer = styled.div`
  display: grid;
  height: 100vh;
  width: 100vw;
  background-color: var(--bg-color);

  /* Desktop Layout */
  grid-template-rows: 60px 100px 1fr 200px;
  grid-template-columns: 250px 1fr;
  grid-template-areas:
    "header header"
    "metrics metrics"
    "filters logs"
    "charts charts";

  /* Mobile Layout */
  @media (max-width: 768px) {
    display: flex;
    flex-direction: column;
    height: auto;
    min-height: 100vh;
    overflow-x: hidden;
  }
`;

const HeaderArea = styled.header`
  grid-area: header;
  background-color: #161b22;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  padding: 0 16px;
  justify-content: space-between;
  position: sticky;
  top: 0;
  z-index: 50;
`;

const MobileMenuBtn = styled.button`
  background: none;
  border: none;
  color: #c9d1d9;
  cursor: pointer;
  display: none; // Hidden on Desktop

  @media (max-width: 768px) {
    display: block;
  }
`;

// The Mobile Drawer Overlay
const MobileDrawer = styled.div<{ $isOpen: boolean }>`
  position: fixed;
  top: 60px; /* Below header */
  left: 0;
  width: 100%;
  height: calc(100vh - 60px);
  background-color: #0d1117;
  z-index: 40;
  padding: 20px;
  border-top: 1px solid var(--border-color);
  transition: transform 0.3s ease-in-out;
  transform: translateX(${(p) => (p.$isOpen ? "0%" : "-100%")});

  @media (min-width: 769px) {
    display: none; /* Never show on desktop */
  }
`;

const SidebarArea = styled.aside`
  grid-area: filters;
  background-color: #161b22;
  border-right: 1px solid var(--border-color);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;

  @media (max-width: 768px) {
    display: none; /* Hide standard sidebar on mobile */
  }
`;

const ContentWrapper = styled.div`
  display: contents; /* Allows children to participate in parent grid */

  @media (max-width: 768px) {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }
`;

// Reuse styled components for Logs, Metrics, Charts (omitted for brevity, keep your existing ones)
const MetricsArea = styled.section`
  grid-area: metrics;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 16px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr; /* Stack metrics on mobile */
    padding: 0;
  }
`;

const LogsArea = styled.main`
  grid-area: logs;
  position: relative;
  overflow: hidden;
  height: 100%;

  @media (max-width: 768px) {
    height: 500px; /* Fixed height for logs on mobile so user can scroll page */
    border: 1px solid var(--border-color);
    border-radius: 8px;
  }
`;

const ChartsArea = styled.section`
  grid-area: charts;
  background-color: #161b22;
  border-top: 1px solid var(--border-color);
  padding: 16px;

  @media (max-width: 768px) {
    height: 300px;
    border: 1px solid var(--border-color);
    border-radius: 8px;
  }
`;

export const DashboardGrid = ({
  header,
  sidebar,
  metrics,
  logs,
  charts,
}: {
  header: React.ReactNode;
  sidebar: React.ReactNode;
  metrics: React.ReactNode;
  logs: React.ReactNode;
  charts: React.ReactNode;
}) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <GridContainer>
      <HeaderArea>
        {/* Pass the toggle button into the header or render it here */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "12px",
            width: "100%",
          }}
        >
          <MobileMenuBtn onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
            {isMobileMenuOpen ? <X /> : <Menu />}
          </MobileMenuBtn>
          {header}
        </div>
      </HeaderArea>

      {/* Mobile Drawer contains the sidebar content */}
      <MobileDrawer $isOpen={isMobileMenuOpen}>{sidebar}</MobileDrawer>

      {/* Desktop Sidebar */}
      <SidebarArea>{sidebar}</SidebarArea>

      <ContentWrapper>
        <MetricsArea>{metrics}</MetricsArea>
        <LogsArea>{logs}</LogsArea>
        <ChartsArea>{charts}</ChartsArea>
      </ContentWrapper>
    </GridContainer>
  );
};
