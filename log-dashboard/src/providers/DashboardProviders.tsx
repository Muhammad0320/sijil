"use client";

import { Project } from "@/lib/types";
import { createContext, ReactNode, useContext, useState } from "react";

interface DashboardContextType {
  // Data
  projects: Project[];
  setProjects: (projects: Project[]) => void;
  addProject: (project: Project) => void;

  // Selection
  selectedProjectId: number | null;
  setSelectedProjectId: (id: number | null) => void;

  // UI
  isCreateOpen: boolean;
  setCreateOpen: (v: boolean) => void;
  isSettingsOpen: boolean;
  setSettingsOpen: (v: boolean) => void;

  // Token
  token: string;
}

const DashboardContext = createContext<DashboardContextType | undefined>(
  undefined
);

export function DashboardProvider({
  children,
  initialProjects,
  token,
}: {
  children: ReactNode;
  initialProjects: Project[];
  token: string;
}) {
  const [projects, setProjects] = useState<Project[]>(initialProjects);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(
    initialProjects.length > 0 ? initialProjects[0].id : null
  );

  const [isCreateOpen, setCreateOpen] = useState(false);
  const [isSettingsOpen, setSettingsOpen] = useState(false);

  const addProject = (project: Project) => {
    setProjects((prev) => [...prev, project]);
    setSelectedProjectId(project.id);
  };

  return (
    <DashboardContext.Provider
      value={{
        projects,
        setProjects,
        addProject,

        selectedProjectId,
        setSelectedProjectId,
        isCreateOpen,
        setCreateOpen,
        isSettingsOpen,
        setSettingsOpen,

        token,
      }}
    >
      {children}
    </DashboardContext.Provider>
  );
}

export const useDashboard = () => {
  const context = useContext(DashboardContext);
  if (!context)
    throw new Error("useDashboard must be used within DashboardProvider");
  return context;
};
