"use client";

import { useActionState, useState } from "react";
import { loginAction } from "@/actions/auth";
import { AuthFormState } from "@/lib/definitions";
import { FieldError, GlobalError } from "@/components/ui/formErrors";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";
import { Eye, EyeOff } from "lucide-react";
import {
  Form,
  FormGroup,
  Input,
  InputWrapper,
  Label,
  TogglePasswordBtn,
} from "./AuthStyles";

const initialState: AuthFormState = { errors: {} };

export default function LoginForm() {
  const [state, formAction, isPending] = useActionState<
    AuthFormState,
    FormData
  >(loginAction, initialState);

  const [showPassword, setShowPassword] = useState(false);

  return (
    <Form action={formAction}>
      <GlobalError errors={state.errors._form} />

      <FormGroup>
        <Label htmlFor="email">Email Address</Label>
        <Input
          id="email"
          name="email"
          type="email"
          placeholder="admin@logengine.com"
          $hasError={!!state.errors.email}
        />
        <FieldError errors={state.errors.email} />
      </FormGroup>

      <FormGroup>
        <Label htmlFor="password">Password</Label>
        <InputWrapper>
          <Input
            id="password"
            name="password"
            type={showPassword ? "text" : "password"}
            placeholder="••••••••"
            $hasError={!!state.errors.password}
          />
          <TogglePasswordBtn
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            tabIndex={-1} // Skip tab focus
          >
            {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
          </TogglePasswordBtn>
        </InputWrapper>
        <FieldError errors={state.errors.password} />
      </FormGroup>

      <div style={{ marginTop: "16px" }}>
        <BorderBeamButton type="submit" isLoading={isPending}>
          Initialize Session
        </BorderBeamButton>
      </div>
    </Form>
  );
}
