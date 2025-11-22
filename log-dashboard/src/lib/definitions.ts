import { z } from "zod";

export interface AuthFormState {
  errors: {
    email?: string[];
    password?: string[];
    name?: string[];
    _form?: string[];
  };

  message?: string;
}

export interface ProjectFormState {
  errors: {
    name?: string[];
  };
  message?: string;
}

export const LoginFormSchema = z.object({
  email: z.email({ error: "Please enter a valid email address." }),
  password: z
    .string()
    .min(8, { error: "Password must be at least 8 characters long" }),
});

export const RegisterFormSchema = z.object({
  name: z
    .string()
    .min(2, { error: "Name must be at least 2 characters long" })
    .trim(),
  email: z.email({ error: "Please enter a valid email address." }).trim(),
  password: z
    .string()
    .min(8, { error: "password must be at least 8 characters" })
    .trim()
    .regex(/[a-zA-Z]/, { error: "Contain at least one letter" })
    .regex(/[0-9]/, { error: "Contains at least one number" })
    .regex(/[^a-zA-z0-9]/, { error: "Contains at least one special character" })
    .trim(),
});
