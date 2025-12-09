import styled from "styled-components";

export const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  position: relative; // For positioning the eye icon
`;

export const Label = styled.label`
  font-size: 13px;
  color: #8b949e;
  margin-bottom: 6px;
  font-weight: 500;
`;

export const InputWrapper = styled.div`
  position: relative;
  display: flex;
  align-items: center;
`;

export const Input = styled.input<{ $hasError?: boolean }>`
  width: 100%;
  background: #0d1117;
  border: 1px solid ${(props) => (props.$hasError ? "#ff6b6b" : "#30363d")};
  color: #fff;
  padding: 10px 12px;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  transition: all 0.2s ease-in-out;

  &:focus {
    border-color: ${(props) => (props.$hasError ? "#ff6b6b" : "#58a6ff")};
    box-shadow: 0 0 0 3px
      ${(props) =>
        props.$hasError ? "rgba(255,107,107,0.1)" : "rgba(88,166,255,0.1)"};
  }
`;

export const TogglePasswordBtn = styled.button`
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: #8b949e;
  cursor: pointer;
  display: flex;
  align-items: center;
  &:hover {
    color: #c9d1d9;
  }
`;
