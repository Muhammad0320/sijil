"use client";

import { useState, useOptimistic } from "react";
import styled from "styled-components";
import { Plus, Loader2 } from "lucide-react";
import { Modal } from "@/components/ui/Modal";
import { Project } from "@/lib/types";
import CreateProjectForm from "./createProjectForm";

// ... (Styled Components: Container, Header, Title, AddButton, ProjectItem REMAIN THE SAME) ...
const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 0 8px;
`;
const Title = styled.h3`
  font-size: 12px;
  text-transform: uppercase;
  color: #8b949e;
  letter-spacing: 0.5px;
`;
const AddButton = styled.button`
  background: none;
  border: none;
  color: #8b949e;
  cursor: pointer;
  &:hover {
    color: #58a6ff;
  }
`;

const ProjectItem = styled.div<{ $active?: boolean; $pending?: boolean }>`
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  background: ${(p) => (p.$active ? "#21262d" : "transparent")};
  color: ${(p) => (p.$active ? "#fff" : "#c9d1d9")};
  border-left: 3px solid ${(p) => (p.$active ? "#238636" : "transparent")};
  opacity: ${(p) => (p.$pending ? 0.5 : 1)};
  font-size: 14px;
  transition: all 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
  &:hover {
    background: #21262d;
  }
`;

interface ProjectListProps {
  initialProjects: Project[];
  onSelect: (id: number) => void;
  selectedId: number | null;
}

export default function ProjectList({
  initialProjects,
  onSelect,
  selectedId,
}: ProjectListProps) {
  const [isOpen, setIsOpen] = useState(false);

  // Optimistic State lives here (in the parent) so the list updates
  const [optimisticProjects, addOptimisticProject] = useOptimistic(
    initialProjects,
    (state: Project[], newProject: Project) => [...state, newProject]
  );

  const handleAddOptimistic = (name: string) => {
    addOptimisticProject({
      id: Math.random(),
      name: name,
      pending: true,
    });
  };

  return (
    <Container>
      <Header>
        <Title>Projects</Title>
        <AddButton onClick={() => setIsOpen(true)}>
          <Plus size={16} />
        </AddButton>
      </Header>

      {optimisticProjects.map((p) => (
        <ProjectItem
          key={p.id}
          $active={selectedId === p.id}
          $pending={p.pending}
          onClick={() => !p.pending && onSelect(p.id)}
        >
          <span>{p.name}</span>
          {p.pending && <Loader2 size={12} className="animate-spin" />}
        </ProjectItem>
      ))}

      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Create New Project"
      >
        {isOpen && (
          <CreateProjectForm
            onSuccess={() => setIsOpen(false)}
            addOptimistic={handleAddOptimistic}
          />
        )}
      </Modal>
    </Container>
  );
}
