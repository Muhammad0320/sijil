"use client";

import { useEffect, useState } from "react";
import styled from "styled-components";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { LogStat } from "@/actions/metrics";
import { Skeleton } from "@/components/ui/skeleton";

const Container = styled.div`
  height: 250px; /* Increased height for better visibility */
  width: 100%;
  background: rgba(13, 17, 23, 0.4);
  border: 1px solid #30363d;
  border-radius: 12px;
  padding: 16px;
  margin-top: 24px;
  display: flex;
  flex-direction: column;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
`;

const Title = styled.h3`
  font-size: 14px;
  color: #8b949e;
  font-weight: 600;
`;

const Badge = styled.span`
  font-size: 10px;
  background: rgba(88, 166, 255, 0.15);
  color: #58a6ff;
  padding: 2px 8px;
  border-radius: 99px;
  border: 1px solid rgba(88, 166, 255, 0.3);
  display: flex;
  align-items: center;
  gap: 4px;

  &::before {
    content: "";
    width: 6px;
    height: 6px;
    background: #58a6ff;
    border-radius: 50%;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 0.4;
    }
    100% {
      opacity: 1;
    }
  }
`;

export default function VolumeChart({
  data,
  loading,
}: {
  data: LogStat[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <Container>
        <Header>
          <Skeleton className="h-4 w-32 bg-[#30363d]" />
        </Header>
        <Skeleton className="h-full w-full bg-[#161b22] rounded-lg" />
      </Container>
    );
  }

  return (
    <Container>
      <Header>
        <Title>Ingestion Volume (24h)</Title>
        <Badge>LIVE</Badge>
      </Header>

      {data.length === 0 ? (
        <div className="flex h-full items-center justify-center text-[#8b949e] text-sm">
          No data recorded
        </div>
      ) : (
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data}>
            <defs>
              <linearGradient id="colorCount" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#58a6ff" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#58a6ff" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis dataKey="time" hide />
            <Tooltip
              contentStyle={{
                backgroundColor: "#0d1117",
                borderColor: "#30363d",
                color: "#c9d1d9",
                fontSize: "12px",
                borderRadius: "8px",
              }}
              itemStyle={{ color: "#58a6ff" }}
              labelFormatter={(label) =>
                new Date(label).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                })
              }
            />
            <Area
              type="monotone"
              dataKey="count"
              stroke="#58a6ff"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#colorCount)"
              isAnimationActive={false} // Disable recharts animation for smooth live updates
            />
          </AreaChart>
        </ResponsiveContainer>
      )}
    </Container>
  );
}
