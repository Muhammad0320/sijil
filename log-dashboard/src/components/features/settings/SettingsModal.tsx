"use client";

import { useState } from "react";
import styled from "styled-components";
import { fetchClient } from "@/lib/client";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";
import { useToast } from "@/providers/ToastProvider";
import { UserPlus, Shield } from "lucide-react";

const Section = styled.div`
  margin-bottom: 24px;
`;

const SectionTitle = styled.h4`
  color: #fff;
  font-size: 14px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
`;

const InputGroup = styled.div`
  display: flex;
  gap: 8px;
`;

const Input = styled.input`
  flex: 1;
  background: #0d1117;
  border: 1px solid #30363d;
  color: #fff;
  padding: 10px;
  border-radius: 6px;

  &:focus {
    border-color: #58a6ff;
    outline: none;
  }
`;

const Select = styled.select`
  background: #0d1117;
  border: 1px solid #30363d;
  color: #fff;
  padding: 10px;
  border-radius: 6px;
`;

export default function SettingsModal({ projectId }: { projectId: number }) {
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("viewer");
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      await fetchClient(`/projects/${projectId}/members`, {
        method: "POST",
        body: JSON.stringify({ email, role }),
      });
      toast.success(`Invited ${email} as ${role}`);
      setEmail("");
    } catch (err) {
      if (err instanceof Error) {
        toast.error(err.message || "Failed to invite member");
      } else {
        toast.error("Failed to invite member");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Section>
        <SectionTitle>
          <UserPlus size={16} /> Invite Team Member
        </SectionTitle>
        <form onSubmit={handleInvite}>
          <InputGroup>
            <Input
              type="email"
              placeholder="colleague@company.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <Select value={role} onChange={(e) => setRole(e.target.value)}>
              <option value="viewer">Viewer</option>
              <option value="admin">Admin</option>
            </Select>
          </InputGroup>
          <div style={{ marginTop: "12px" }}>
            <BorderBeamButton
              type="submit"
              isLoading={loading}
              variant="primary"
            >
              Send Invite
            </BorderBeamButton>
          </div>
        </form>
      </Section>

      <Section style={{ borderTop: "1px solid #30363d", paddingTop: "24px" }}>
        <SectionTitle>
          <Shield size={16} /> API Keys
        </SectionTitle>
        <p style={{ fontSize: "13px", color: "#8b949e", lineHeight: "1.5" }}>
          API Keys are hashed and cannot be revealed again. To rotate keys, you
          must regenerate them (this will break existing integrations).
        </p>
        <div style={{ marginTop: "12px" }}>
          <button
            style={{
              color: "#ff6b6b",
              background: "none",
              border: "1px solid #ff6b6b",
              padding: "8px 16px",
              borderRadius: "6px",
              cursor: "pointer",
              fontSize: "13px",
            }}
          >
            Rotate Keys (Coming Soon)
          </button>
        </div>
      </Section>
    </div>
  );
}
