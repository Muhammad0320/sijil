"use client";

import styled from "styled-components";

const GridContainer = styled.div`
  display: grid;
  grid-template-rows: 60px 100px 1fr 200px;
  grid-template-columns: 250px 1fr;
  grid-template-areas:
    "header header"
    "metrics metrics"
    "filters logs"
    "charts charts";
  height: 100vh;
  width: 100vw;
  background-color: var(--bg-color);
`;

const HeaderArea = styled.header`
  grid-area: header;
  background-color: #161b22;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  padding: 0 16px;
  justify-content: space-between;
  z-index: 10;
`;

const MetricsArea = styled.section`
  grid-area: metrics;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 16px;
`;

const SidebarArea = styled.aside`
  grid-area: filters;
  background-color: #161b22;
  border-right: 1px solid var(--border-color);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const LogsArea = styled.main`
  grid-area: logs;
  position: relative;
  overflow: hidden;
`;

const ChartsArea = styled.section`
  grid-area: charts;
  background-color: #161b22;
  border-top: 1px solid var(--border-color);
  padding: 16px;
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
  return (
    <GridContainer>
      <HeaderArea>{header}</HeaderArea>
      <MetricsArea>{metrics}</MetricsArea>
      <SidebarArea>{sidebar}</SidebarArea>
      <LogsArea>{logs}</LogsArea>
      <ChartsArea>{charts}</ChartsArea>
    </GridContainer>
  );
};
