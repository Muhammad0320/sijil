"use client";

import { useEffect, useTransition } from "react";
import { useActionState } from "react";
import styled from "styled-components";
import { createProjectAction, CreateProjectState } from "@/actions/projects";
import { useToast } from "@/providers/ToastProvider";
import { FieldError } from "@/components/ui/formErrors";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";

const Input = styled.input`
  width: 100%;
  background: #0d1117;
  border: 1px solid #30363d;
  color: #fff;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 8px;
  &:focus {
    outline: none;
    border-color: #58a6ff;
  }
`;

const initialState: CreateProjectState = {};

export default function CreateProjectForm({
  onSuccess,
  addOptimistic,
}: {
  onSuccess: () => void;
  addOptimistic: (name: string) => void;
}) {
  const toast = useToast();
  const [state, formAction, isPending] = useActionState(
    createProjectAction,
    initialState
  );
  const [transitionPending, startTransition] = useTransition();

  // Handle Feedback
  useEffect(() => {
    if (state.success) {
      toast.success("Project created successfully");
      onSuccess(); // Close the modal
    } else if (state.errors?._form) {
      toast.error(state.errors._form[0]);
    }
  }, [state, toast, onSuccess]);

  const handleSubmit = (formData: FormData) => {
    const name = formData.get("name") as string;
    if (!name) return;

    // Optimistic Update (UI first)
    startTransition(async () => {
      addOptimistic(name);
      // Server Request (Data second)
      await formAction(formData);
    });
  };

  return (
    <form action={handleSubmit}>
      <Input
        name="name"
        placeholder="Project Name (e.g. Production-API)"
        autoFocus
      />
      <FieldError errors={state.errors?.name} />

      <div style={{ marginTop: "20px" }}>
        <BorderBeamButton
          type="submit"
          isLoading={isPending || transitionPending}
        >
          Create Project
        </BorderBeamButton>
      </div>
    </form>
  );
}
