"use client";

import { useEffect, useState } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from "recharts";
import styled from "styled-components";
import { fetchClient } from "@/lib/client";
import { Skeleton } from "@/components/ui/skeleton";

const ChartContainer = styled.div`
  width: 100%;
  height: 100%;
  background: rgba(22, 27, 34, 0.5); /* Glass-like */
  backdrop-filter: blur(10px);
  border: 1px solid #30363d;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
`;

const Title = styled.h3`
  font-size: 14px;
  color: #8b949e;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
`;

const LiveBadge = styled.span`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  color: #2ecc71;
  background: rgba(46, 204, 113, 0.1);
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid rgba(46, 204, 113, 0.2);
  
  &::before {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #2ecc71;
    box-shadow: 0 0 8px #2ecc71;
  }
`;

interface StatPoint {
  bucket: string;
  count: number;
}

export default function VolumeChart({ projectId }: { projectId: number | null }) {
  const [data, setData] = useState<StatPoint[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!projectId) return;

    const loadData = async () => {
      try {
        // Fetch last 1 hour, bucketed by 1 minute
        // Ideally, you'd use a library like swr or tanstack-query here for auto-refresh
        const res = await fetchClient<{ stats: StatPoint[] }>(
          `/logs/stats?project_id=${projectId}&bucket=1m`
        );
        
        // Transform for Recharts (parse dates if needed)
        const formatted = res.stats.map(s => ({
            ...s,
            // Format time to HH:MM for the axis
            displayTime: new Date(s.bucket).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        }));
        
        setData(formatted);
      } catch (err) {
        console.error("Failed to load charts", err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
    // Poll every 30 seconds for "live" feel without websockets overkill
    const interval = setInterval(loadData, 30000); 
    return () => clearInterval(interval);
  }, [projectId]);

  if (loading) {
    return (
      <ChartContainer>
        <Header><Title>Ingestion Volume</Title></Header>
        <Skeleton height="100%" width="100%" />
      </ChartContainer>
    );
  }

  return (
    <ChartContainer>
      <Header>
        <Title>Ingestion Volume</Title>
        <LiveBadge>LIVE 1H</LiveBadge>
      </Header>
      
      <div style={{ flex: 1, minHeight: 0 }}>
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data}>
            <defs>
              <linearGradient id="colorCount" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#58a6ff" stopOpacity={0.3}/>
                <stop offset="95%" stopColor="#58a6ff" stopOpacity={0}/>
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#30363d" vertical={false} />
            <XAxis 
                dataKey="displayTime" 
                stroke="#8b949e" 
                fontSize={10} 
                tickLine={false} 
                axisLine={false}
                minTickGap={30}
            />
            <YAxis 
                stroke="#8b949e" 
                fontSize={10} 
                tickLine={false} 
                axisLine={false} 
            />
            <Tooltip 
                contentStyle={{ backgroundColor: '#161b22', borderColor: '#30363d', fontSize: '12px' }}
                itemStyle={{ color: '#c9d1d9' }}
            />
            <Area 
                type="monotone" 
                dataKey="count" 
                stroke="#58a6ff" 
                strokeWidth={2}
                fillOpacity={1} 
                fill="url(#colorCount)" 
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </ChartContainer>
  );
}