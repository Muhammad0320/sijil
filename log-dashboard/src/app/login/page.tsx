"use client";

import { useActionState, useEffect } from "react";
import styled from "styled-components";
import { loginAction } from "@/actions/auth";
import { AuthFormState } from "@/lib/definitions"; // Import our Type
import { FieldError, GlobalError } from "@/components/ui/formErrors";
import { useToast } from "@/providers/ToastProvider";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";

const Container = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-color);
`;

const Card = styled.form`
  width: 400px;
  padding: 40px;
  background: #161b22;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 16px; /* Reduced gap to accommodate error messages */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
`;

const Title = styled.h1`
  color: var(--text-color);
  font-size: 24px;
  text-align: center;
  margin-bottom: 10px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.label`
  font-size: 14px;
  color: #8b949e;
  margin-bottom: 6px;
`;

const Input = styled.input<{ hasError?: boolean }>`
  background: #0d1117;
  border: 1px solid
    ${(props) => (props.hasError ? "#ff6b6b" : "var(--border-color)")};
  color: #fff;
  padding: 12px;
  border-radius: 6px;
  font-size: 16px;
  outline: none;
  transition: border-color 0.2s;

  &:focus {
    border-color: ${(props) => (props.hasError ? "#ff6b6b" : "#58a6ff")};
  }
`;

// Initial State matches our Interface
const initialState: AuthFormState = {
  errors: {},
  message: "",
};

export default function LoginPage() {
  // useActionState is generic! <AuthFormState> ensures 'state' is typed.
  const [state, formAction, isPending] = useActionState<
    AuthFormState,
    FormData
  >(loginAction, initialState);

  const toast = useToast();

  useEffect(() => {
    if (state.errors._form && state.errors._form.length > 0) {
      toast.error(state.errors._form[0]);
    }
    if (state.message) {
      toast.success(state.message);
    }
  }, [state, toast]);

  return (
    <Container>
      <Card action={formAction}>
        <Title>LogEngine Access</Title>
        <FormGroup>
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            name="email"
            type="email"
            placeholder="admin@example.com"
            hasError={!!state.errors.email}
          />
          {/* Field Errors (email) */}
          <FieldError errors={state.errors.email} />
        </FormGroup>

        <FormGroup>
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            name="password"
            type="password"
            placeholder="••••••••"
            hasError={!!state.errors.password}
          />
          {/* Field Errors (password) */}
          <FieldError errors={state.errors.password} />
        </FormGroup>

        <div style={{ marginTop: "10px" }}>
          <BorderBeamButton type="submit" isLoading={isPending}>
            Sign In
          </BorderBeamButton>
        </div>
      </Card>
    </Container>
  );
}
