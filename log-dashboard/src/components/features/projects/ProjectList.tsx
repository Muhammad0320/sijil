"use client";

import { useState, useOptimistic } from "react";
import styled from "styled-components";
import { Plus, Loader2, Check, Copy } from "lucide-react";
import { Modal } from "@/components/ui/Modal";
import { Project } from "@/lib/types";
import CreateProjectForm from "./createProjectForm";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";

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

// NEW: Styles for the Key Reveal
const KeyContainer = styled.div`
  background: #0d1117;
  border: 1px solid #30363d;
  border-radius: 6px;
  padding: 16px;
  margin-top: 16px;
`;

const KeyRow = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }
`;

const KeyLabel = styled.span`
  font-size: 12px;
  color: #8b949e;
  font-weight: 600;
  text-transform: uppercase;
`;

const KeyValueGroup = styled.div`
  display: flex;
  gap: 8px;
`;

const KeyInput = styled.code`
  flex: 1;
  background: #161b22;
  border: 1px solid #30363d;
  padding: 8px 12px;
  border-radius: 6px;
  color: #58a6ff;
  font-family: monospace;
  font-size: 13px;
`;

const CopyButton = styled.button`
  background: #21262d;
  border: 1px solid #30363d;
  color: #c9d1d9;
  border-radius: 6px;
  width: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  &:hover {
    background: #30363d;
  }
`;

// Helper for Copy Logic
function CopyField({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <KeyValueGroup>
      <KeyInput>{value}</KeyInput>
      <CopyButton onClick={handleCopy} type="button">
        {copied ? <Check size={14} color="#2ecc71" /> : <Copy size={14} />}
      </CopyButton>
    </KeyValueGroup>
  );
}

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
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  // NEW: State to hold the created keys
  const [createdKeys, setCreatedKeys] = useState<{
    apiKey: string;
    apiSecret: string;
    projectId: number;
  } | null>(null);

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

  // When form succeeds, we don't close. We show keys.
  const handleProjectCreated = (data: {
    apiKey: string;
    apiSecret: string;
    projectId: number;
  }) => {
    setCreatedKeys(data);
    // Auto-select the new project behind the modal
    onSelect(data.projectId);
  };

  const closeAndReset = () => {
    setIsCreateOpen(false);
    setCreatedKeys(null);
  };

  return (
    <Container>
      <Header>
        <Title>Projects</Title>
        <AddButton onClick={() => setIsCreateOpen(true)}>
          <Plus size={16} />
        </AddButton>
      </Header>

      {optimisticProjects?.map((p) => (
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
        isOpen={isCreateOpen}
        onClose={closeAndReset}
        title={createdKeys ? "Connect Your Agent" : "Create New Project"}
      >
        {!createdKeys ? (
          // STEP 1: Create Form
          <CreateProjectForm
            onProjectCreated={handleProjectCreated}
            addOptimistic={handleAddOptimistic}
          />
        ) : (
          // STEP 2: Key Reveal (The "Success" State)
          <div>
            <p
              style={{
                fontSize: "13px",
                color: "#8b949e",
                marginBottom: "16px",
              }}
            >
              These keys will <strong>only be shown once</strong>. Copy them
              now.
            </p>

            <KeyContainer>
              <KeyRow>
                <KeyLabel>Public Key (PK)</KeyLabel>
                <CopyField value={createdKeys.apiKey} />
              </KeyRow>

              <KeyRow>
                <KeyLabel>Secret Key (SK)</KeyLabel>
                <CopyField value={createdKeys.apiSecret} />
              </KeyRow>
            </KeyContainer>

            <BorderBeamButton onClick={closeAndReset}>
              I have saved my keys securely
            </BorderBeamButton>
          </div>
        )}
      </Modal>
    </Container>
  );
}
