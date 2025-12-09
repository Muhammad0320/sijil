"use client";

import { useEffect, useState } from "react";
import { DashboardGrid } from "@/components/layout/DashboardGrid";
import ProjectList from "@/components/features/projects/ProjectList";
import { useLogStream } from "@/hooks/useLogStream";
import { Project } from "@/lib/types"; // Using our shared type
import LogList from "@/components/features/logs/Loglist";
import { useToast } from "@/providers/ToastProvider";
import { LogToolbar } from "@/components/features/logs/LogToolbar";

interface DashboardClientProps {
  initialProjects: Project[];
  token: string;
  serverError: string | null;
}

export default function DashboardClient({
  initialProjects,
  token,
  serverError,
}: DashboardClientProps) {
  const toast = useToast();

  useEffect(() => {
    if (serverError) {
      toast.error(serverError);
    }
  }, [serverError, toast]);

  // 1. State Management
  // We initialize with server data, but keep local state for client-side updates if needed
  const [projects] = useState<Project[]>(initialProjects);

  // Default to the first project if available
  const [selectedProject, setSelectedProject] = useState<number | null>(
    initialProjects.length > 0 ? initialProjects[0].id : null
  );

  // 2. WebSocket Connection
  // Pass the selected project ID. The hook handles connecting/disconnecting/buffering.
  // 3. Derived UI State
  const currentProjectName =
    projects.find((p) => p.id === selectedProject)?.name || "Select Project";

  const [searchQuery, setSearchQuery] = useState("");

  // Get raw logs from your hook
  const { logs, status } = useLogStream(selectedProject || 0, token);

  // LOGIC: Filter logs before rendering
  // This is efficient enough for < 10k logs in client memory
  const filteredLogs = logs.filter((log) => {
    if (!searchQuery) return true;
    const lowerQuery = searchQuery.toLowerCase();
    return (
      log.message.toLowerCase().includes(lowerQuery) ||
      log.service.toLowerCase().includes(lowerQuery) ||
      log.level.toLowerCase().includes(lowerQuery)
    );
  });

  return (
    <DashboardGrid
      // --- HEADER ---
      header={
        <div
          style={{
            display: "flex",
            gap: "1rem",
            alignItems: "center",
            color: "#fff",
          }}
        >
          <h2
            style={{
              fontWeight: 700,
              fontSize: "18px",
              letterSpacing: "-0.5px",
            }}
          >
            LogEngine
          </h2>
          <div
            style={{ height: "20px", width: "1px", background: "#30363d" }}
          />
          <span style={{ fontSize: "14px", fontWeight: 500 }}>
            {currentProjectName}
          </span>

          {/* Connection Status Indicator */}
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "6px",
              marginLeft: "auto",
              fontSize: "12px",
              color: "#8b949e",
              background: "#21262d",
              padding: "4px 8px",
              borderRadius: "20px",
              border: "1px solid #30363d",
            }}
          >
            <span
              style={{
                height: "6px",
                width: "6px",
                borderRadius: "50%",
                backgroundColor:
                  status === "OPEN"
                    ? "#2ecc71"
                    : status === "CONNECTING"
                    ? "#f1c40f"
                    : "#e74c3c",
                boxShadow:
                  status === "OPEN"
                    ? "0 0 8px rgba(46, 204, 113, 0.4)"
                    : "none",
              }}
            />
            {status}
          </div>
        </div>
      }
      // --- SIDEBAR (Project Switcher) ---
      sidebar={
        <ProjectList
          initialProjects={initialProjects}
          selectedId={selectedProject}
          onSelect={(id) => setSelectedProject(id)}
        />
      }
      // --- METRICS (Placeholder for now) ---
      metrics={
        <div style={{ color: "#8b949e", fontSize: "13px" }}>
          {/* We will add the MetricsSection here later */}
          Live Metrics Incoming...
        </div>
      }
      // --- CHARTS (Placeholder) ---
      charts={
        <div style={{ color: "#8b949e", fontSize: "13px" }}>
          Volume Analytics
        </div>
      }
      logs={
        <div
          style={{ display: "flex", flexDirection: "column", height: "100%" }}
        >
          <LogToolbar
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            onRefresh={() => {}}
          />
          <div style={{ flex: 1 }}>
            {selectedProject ? (
              <LogList logs={filteredLogs} /> // Pass filtered logs!
            ) : (
              <div style={{ padding: 20, color: "#8b949e" }}>
                Select a project
              </div>
            )}
          </div>
        </div>
      }
    />
  );
}
