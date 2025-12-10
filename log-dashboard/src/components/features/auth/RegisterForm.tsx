"use client";

import { useActionState, useState } from "react";
import { registerAction } from "@/actions/auth";
import { RegisterFormState } from "@/lib/definitions";
import { FieldError, GlobalError } from "@/components/ui/formErrors";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";
import { Eye, EyeOff, Lock } from "lucide-react";
import {
  Form,
  FormGroup,
  Input,
  InputWrapper,
  Label,
  TogglePasswordBtn,
} from "./AuthStyles";

const initialState: RegisterFormState = { errors: {} };

export default function RegisterForm() {
  const [state, formAction, isPending] = useActionState<
    RegisterFormState,
    FormData
  >(registerAction, initialState);

  const [showPassword, setShowPassword] = useState(false);

  return (
    <Form action={formAction}>
      <GlobalError errors={state.errors?._form} />

      <FormGroup>
        <Label htmlFor="firstname">Firstname</Label>
        <Input
          name="firstname"
          autoComplete="firstname"
          type="text"
          placeholder="Muhammad"
          $hasError={!!state.errors?.firstname}
        />
        <FieldError errors={state.errors?.firstname} />
      </FormGroup>

      <FormGroup>
        <Label htmlFor="lastname">Lastname</Label>
        <Input
          name="lastname"
          autoComplete="lastname"
          type="text"
          placeholder="Abdullah"
          $hasError={!!state.errors?.lastname}
        />
        <FieldError errors={state.errors?.lastname} />
      </FormGroup>

      <FormGroup>
        <label
          style={{ fontSize: "13px", color: "#8b949e", marginBottom: "6px" }}
        >
          Email
        </label>
        <Input
          name="email"
          type="email"
          autoComplete="email"
          placeholder="pilot@falcon.com"
          $hasError={!!state.errors?.email}
        />
        <FieldError errors={state.errors?.email} />
      </FormGroup>

      <FormGroup>
        <Label htmlFor="password">Password</Label>
        <InputWrapper>
          <Input
            id="password"
            name="password"
            type={showPassword ? "text" : "password"}
            placeholder="••••••••"
            $hasError={!!state.errors?.password}
          />
          <TogglePasswordBtn
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            tabIndex={-1} // Skip tab focus
          >
            {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
          </TogglePasswordBtn>
        </InputWrapper>
        <FieldError errors={state.errors?.password} />
      </FormGroup>

      <FormGroup>
        <Label htmlFor="confirmPassword">Confirm Password</Label>
        <InputWrapper>
          <Lock size={16} />
          <Input
            id="confirmPassword"
            name="confirmPassword"
            type={showPassword ? "text" : "password"}
            placeholder="••••••••"
            $hasError={!!state.errors?.confirmPassword}
          />
        </InputWrapper>
        <FieldError errors={state.errors?.confirmPassword} />
      </FormGroup>

      <div style={{ marginTop: "16px" }}>
        <BorderBeamButton type="submit" isLoading={isPending}>
          Create Account
        </BorderBeamButton>
      </div>
    </Form>
  );
}
