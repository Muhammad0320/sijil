"use client";

import React, { useRef } from "react";
import { Virtuoso, type VirtuosoHandle } from "react-virtuoso";
import styled from "styled-components";
import { LogEntry } from "@/lib/types";

// 1. Styled Component for the Row
const LogRow = styled.div<{ $level: string }>`
  font-family: "Fira Code", monospace; /* Ensure font is imported in global.css */
  font-size: 13px;
  padding: 4px 16px; /* Added horizontal padding */
  border-left: 3px solid;

  /* Use transient props ($level) so they don't pass to DOM */
  border-left-color: ${(p) =>
    p.$level === "ERROR"
      ? "#ff6b6b"
      : p.$level === "WARN"
      ? "#f1c40f"
      : "#2ecc71"};

  background: ${(p) =>
    p.$level === "ERROR" ? "rgba(255, 107, 107, 0.1)" : "transparent"};

  color: #c9d1d9;
  display: flex;
  gap: 12px;
  line-height: 1.5;

  &:hover {
    background: #21262d;
  }
`;

const Timestamp = styled.span`
  color: #8b949e;
  min-width: 85px;
  font-variant-numeric: tabular-nums; /* Aligns numbers */
`;

const Level = styled.span`
  font-weight: 700;
  width: 50px;
`;

const Service = styled.span`
  color: #58a6ff;
`;

const Message = styled.span`
  white-space: pre-wrap;
  word-break: break-all;
`;

interface LogListProps {
  logs: LogEntry[];
}

export default function LogList({ logs }: LogListProps) {
  // Correctly typed Ref
  const virtuosoRef = useRef<VirtuosoHandle>(null);

  return (
    <div style={{ height: "100%", width: "100%" }}>
      <Virtuoso
        ref={virtuosoRef}
        data={logs}
        // Auto-scroll logic: "smooth" or "auto"
        followOutput="auto"
        initialTopMostItemIndex={logs.length - 1}
        // itemContent signature: (index, data)
        itemContent={(index, log) => (
          <LogRow $level={log.level}>
            <Timestamp>
              {new Date(log.timestamp).toLocaleTimeString()}
            </Timestamp>

            <Level>{log.level.toUpperCase()}</Level>

            <Service>[{log.service}]</Service>

            <Message>{log.message}</Message>
          </LogRow>
        )}
      />
    </div>
  );
}
