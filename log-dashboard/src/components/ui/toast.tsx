"use client";

import styled, { keyframes } from "styled-components";
import { X, CheckCircle, AlertCircle } from "lucide-react";
import { useEffect } from "react";

const slideIn = keyframes`
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
`;

const ToastContainer = styled.div<{ type: "success" | "error" }>`
  position: fixed;
  bottom: 24px;
  right: 24px;
  min-width: 300px;
  background: #161b22;
  border: 1px solid ${(p) => (p.type === "success" ? "#238636" : "#da3633")};
  border-left: 4px solid
    ${(p) => (p.type === "success" ? "#238636" : "#da3633")};
  border-radius: 6px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  color: #fff;
  animation: ${slideIn} 0.3s ease-out;
  z-index: 1000;
`;

const Message = styled.p`
  margin: 0;
  font-size: 14px;
  flex: 1;
`;

const CloseBtn = styled.button`
  background: none;
  border: none;
  color: #8b949e;
  cursor: pointer;
  &:hover {
    color: #fff;
  }
`;

export type ToastType = {
  id: string;
  message: string;
  type: "success" | "error";
};

export function Toast({
  toast,
  onClose,
}: {
  toast: ToastType;
  onClose: () => void;
}) {
  // Auto-dismiss after 5 seconds
  useEffect(() => {
    const timer = setTimeout(onClose, 5000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <ToastContainer type={toast.type}>
      {toast.type === "success" ? (
        <CheckCircle size={20} color="#238636" />
      ) : (
        <AlertCircle size={20} color="#da3633" />
      )}
      <Message>{toast.message}</Message>
      <CloseBtn onClick={onClose}>
        <X size={16} />
      </CloseBtn>
    </ToastContainer>
  );
}
