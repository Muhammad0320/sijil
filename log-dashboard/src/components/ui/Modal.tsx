"use client";

import { useEffect, useRef } from "react";
import styled, { keyframes } from "styled-components";
import { X } from "lucide-react";

// 1. Animations
const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`;

const slideUp = keyframes`
  from { opacity: 0; transform: translate(-50%, -40%) scale(0.90); }
  to { opacity: 1; transform: translate(-50%, -50%) scale(1); }
`;

// 2. The Backdrop
const Backdrop = styled.dialog`
  margin: 0;
  padding: 0;
  width: 100vw;
  height: 100vh;
  max-width: 100vw;
  max-height: 100vh;
  background: rgba(0, 0, 0, 0.6);
  border: none;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 50;
  display: flex; // Hidden by default, flex when 'open' attribute exists
  align-items: center;
  justify-content: center;
  animation: ${fadeIn} 0.3s ease-out;

  &::backdrop {
    background: transparent; // We handle styling manually
  }
`;

// 3. The Card
const ModalContent = styled.div`
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 12px;
  padding: 24px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5),
    0 8px 10px -6px rgba(0, 0, 0, 0.5);
  animation: ${slideUp} 0.3s cubic-bezier(0.16, 1, 0.3, 1); // Apple-style spring
  color: #c9d1d9;
`;

const CloseButton = styled.button`
  position: absolute;
  top: 16px;
  right: 16px;
  background: transparent;
  border: none;
  color: #8b949e;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  &:hover {
    color: #fff;
    background: #21262d;
  }
`;

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}

export function Modal({ isOpen, onClose, title, children }: ModalProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    if (isOpen) {
      dialogRef.current?.showModal();
    } else {
      dialogRef.current?.close();
    }
  }, [isOpen]);

  // Handle clicking backdrop to close
  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === dialogRef.current) onClose();
  };

  if (!isOpen) return null;

  return (
    <Backdrop ref={dialogRef} onClick={handleBackdropClick}>
      <ModalContent>
        <div
          style={{ marginBottom: "16px", fontSize: "18px", fontWeight: 600 }}
        >
          {title}
        </div>
        <CloseButton onClick={onClose} aria-label="Close modal">
          <X size={20} />
        </CloseButton>
        {children}
      </ModalContent>
    </Backdrop>
  );
}
