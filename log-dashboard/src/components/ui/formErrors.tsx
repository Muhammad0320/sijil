"use client";

import styled from "styled-components";

const FieldErrorText = styled.p`
  color: #ff6b6b; /* Soft Red */
  font-size: 12px;
  margin-top: 4px;
  font-weight: 500;
`;

const GlobalErrorBox = styled.div`
  background-color: rgba(255, 107, 107, 0.1);
  border: 1px solid #ff6b6b;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 16px;
  color: #ff6b6b;
  font-size: 14px;
  text-align: center;
`;

export function FieldError({ errors }: { errors?: string[] }) {
  if (!errors || errors.length === 0) return null;
  return (
    <>
      {errors.map((err, index) => (
        <FieldErrorText key={index}>{err}</FieldErrorText>
      ))}
    </>
  );
}

export function GlobalError({ errors }: { errors?: string[] }) {
  if (!errors || errors.length === 0) return null;
  return (
    <GlobalErrorBox>
      {errors.map((err, index) => (
        <div key={index}>{err}</div>
      ))}
    </GlobalErrorBox>
  );
}
