"use client";

import { useState } from "react";
import styled from "styled-components";
import { Check, Copy, ShieldAlert } from "lucide-react";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const WarningBox = styled.div`
  background: rgba(210, 153, 34, 0.1);
  border: 1px solid rgba(210, 153, 34, 0.4);
  padding: 12px;
  border-radius: 6px;
  display: flex;
  gap: 10px;
  align-items: start;

  p {
    font-size: 13px;
    color: #e3b341;
    margin: 0;
    line-height: 1.4;
  }
`;

const KeyGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const Label = styled.label`
  font-size: 12px;
  font-weight: 600;
  color: #8b949e;
  text-transform: uppercase;
  letter-spacing: 0.5px;
`;

const InputGroup = styled.div`
  display: flex;
  gap: 8px;
`;

const KeyInput = styled.code`
  flex: 1;
  background: #0d1117;
  border: 1px solid #30363d;
  padding: 10px 12px;
  border-radius: 6px;
  color: #79c0ff;
  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas,
    Liberation Mono, monospace;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

const IconButton = styled.button`
  background: #21262d;
  border: 1px solid #30363d;
  color: #c9d1d9;
  border-radius: 6px;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #30363d;
    border-color: #8b949e;
    color: #fff;
  }

  &:active {
    transform: scale(0.95);
  }
`;

function CopyField({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <InputGroup>
      <KeyInput>{value}</KeyInput>
      <IconButton
        onClick={handleCopy}
        title="Copy to clipboard"
        aria-label="Copy to clipboard"
      >
        {copied ? <Check size={16} color="#2ecc71" /> : <Copy size={16} />}
      </IconButton>
    </InputGroup>
  );
}

interface IKeyReveal {
  data: {
    projectId: number;
    name: string;
    apiKey: string;
    apiSecret: string;
  };
  onClose?: () => void;
}

export default function KeyReveal({ data, onClose }: IKeyReveal) {
  return (
    <Container>
      <div style={{ textAlign: "center", marginBottom: "8px" }}>
        <h3 style={{ fontSize: "16px", fontWeight: 600, color: "#fff" }}>
          Project Created: <span style={{ color: "#2ecc71" }}>{data.name}</span>
        </h3>
      </div>

      <WarningBox>
        <ShieldAlert
          size={18}
          color="#e3b341"
          style={{ flexShrink: 0, marginTop: "2px" }}
        />
        <p>
          These keys will <strong>only be shown once</strong>. If you lose them,
          you will need to regenerate them. Save them securely now.
        </p>
      </WarningBox>

      <KeyGroup>
        <Label>Public API Key</Label>
        <CopyField value={data.apiKey} />
      </KeyGroup>

      <KeyGroup>
        <Label>Secret API Key</Label>
        <CopyField value={data.apiSecret} />
      </KeyGroup>

      <div style={{ marginTop: "10px" }}>
        <BorderBeamButton onClick={onClose}>
          I have saved my keys
        </BorderBeamButton>
      </div>
    </Container>
  );
}
