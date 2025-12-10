"use server";

import { RegisterFormState } from "./../lib/definitions";
import { fetchClient } from "@/lib/client";
import {
  AuthFormState,
  LoginFormSchema,
  RegisterFormSchema,
} from "@/lib/definitions";
import { setSession } from "@/lib/session";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

interface AuthResponse {
  token: string;
}

export async function loginAction(
  prevState: AuthFormState,
  formData: FormData
): Promise<AuthFormState> {
  const validatedSchema = LoginFormSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
  });

  if (!validatedSchema.success) {
    return {
      errors: validatedSchema.error.flatten().fieldErrors,
    };
  }

  const { email, password } = validatedSchema.data;

  try {
    const data = await fetchClient<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });

    console.log(data, "---- From SA");

    await setSession(data.token);
  } catch (error) {
    if (error instanceof Error) {
      return {
        errors: {
          _form: [error.message || "Invalid credentials"],
        },
      };
    }
    return {
      errors: { _form: ["Failed to connect to the server. Please try again!"] },
    };
  }

  revalidatePath("/dashboard");
  redirect("/dashboard");
}

export async function registerAction(
  prevState: RegisterFormState,
  formData: FormData
): Promise<RegisterFormState> {
  const validatedRegisterSchema = RegisterFormSchema.safeParse({
    firstname: formData.get("firstname"),
    lastname: formData.get("lastname"),
    email: formData.get("email"),
    password: formData.get("password"),
    confirmPassword: formData.get("confirmPassword"),
  });

  if (!validatedRegisterSchema.success) {
    return {
      errors: validatedRegisterSchema.error.flatten().fieldErrors,
    };
  }

  const { firstname, lastname, email, password } = validatedRegisterSchema.data;

  try {
    const data = await fetchClient<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ firstname, lastname, email, password }),
    });

    await setSession(data.token);
  } catch (error) {
    if (error instanceof Error) {
      return {
        errors: {
          _form: [error.message || "Registration failed"],
        },
      };
    }
    return {
      errors: { _form: ["Failed to connect to the server. Please try again!"] },
    };
  }

  redirect("/dashboard");
}
