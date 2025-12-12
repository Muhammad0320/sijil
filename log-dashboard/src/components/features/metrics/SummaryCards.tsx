"use client";

import { useEffect, useState } from "react";
import styled from "styled-components";
import { fetchClient } from "@/lib/client";
import { Activity, AlertTriangle, Layers, Zap } from "lucide-react";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";

// --- 1. Strong Typing ---
// Matches 'LogSummary' struct in internals/database/db.go
interface LogSummary {
  total_logs: number;
  error_count: number;
  service_count: number;
  error_rate: number;
}

// --- 2. Galaxy One UI 8.5 Style (Frosted Glass) ---
const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  width: 100%;
`;

const Card = styled.div`
  position: relative; // For BorderBeam positioning
  background: rgba(22, 27, 34, 0.4); // More transparent
  backdrop-filter: blur(24px) saturate(140%); // Heavy "Galaxy" frost
  -webkit-backdrop-filter: blur(24px) saturate(140%);
  border: 1px solid rgba(255, 255, 255, 0.08); // Subtle white border
  border-radius: 16px; // Squircle-ish
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow: hidden; // Contain the beam
  transition: transform 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);

  &:hover {
    transform: translateY(-2px);
    background: rgba(22, 27, 34, 0.6);
  }
`;

const Label = styled.div`
  color: #8b949e;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  display: flex;
  align-items: center;
  gap: 8px;
`;

const Value = styled.div<{ $color?: string }>`
  color: ${(p) => p.$color || "#fff"};
  font-size: 28px;
  font-weight: 500; // Thinner, modern look
  font-family: "Geist Mono", monospace;
  text-shadow: 0 0 20px ${(p) => (p.$color ? `${p.$color}40` : "transparent")}; // Neon glow
`;

export default function SummaryCards({
  projectId,
  token,
}: {
  projectId: number | null;
  token: string;
}) {
  const [summary, setSummary] = useState<LogSummary | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!projectId) return;

    let ignore = false;

    const loadData = async () => {
      // Set loading inside the async flow to prevent sync render warnings
      setLoading(true);

      try {
        const res = await fetchClient<{ summary: LogSummary }>(
          `/logs/summary?project_id=${projectId}`,
          { method: "GET" },
          token
        );
        if (!ignore) {
          setSummary(res.summary);
        }
      } catch (err) {
        console.error("Failed to load summary", err);
      } finally {
        if (!ignore) setLoading(false);
      }
    };

    loadData();

    return () => {
      ignore = true;
    };
  }, [projectId]);

  // Loading State with Border Beam
  if (loading || !summary) {
    return (
      <Grid>
        {[1, 2, 3, 4].map((i) => (
          <Card key={i}>
            {/* The Travelling Light Animation */}
            <BorderBeamButton color={i % 2 === 0 ? "#ff6b6b" : "#58a6ff"}>
              Initialiing...
            </BorderBeamButton>
          </Card>
        ))}
      </Grid>
    );
  }

  return (
    <Grid>
      <Card>
        <Label>
          <Layers size={14} /> Total Logs
        </Label>
        <Value>{summary.total_logs.toLocaleString()}</Value>
      </Card>

      <Card>
        <Label>
          <AlertTriangle size={14} /> Error Rate
        </Label>
        <Value $color={summary.error_rate > 5 ? "#ff6b6b" : "#2ecc71"}>
          {summary.error_rate.toFixed(2)}%
        </Value>
      </Card>

      <Card>
        <Label>
          <Zap size={14} /> Services
        </Label>
        <Value>{summary.service_count}</Value>
      </Card>

      <Card>
        <Label>
          <Activity size={14} /> Status
        </Label>
        <Value $color="#58a6ff">Healthy</Value>
      </Card>
    </Grid>
  );
}
