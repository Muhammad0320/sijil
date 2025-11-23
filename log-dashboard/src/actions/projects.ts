"use server";

import { fetchClient } from "@/lib/client";
import { getErrorMessage } from "@/lib/types";
import { revalidatePath } from "next/cache";
import { z } from "zod";

const CreateProjectSchema = z.object({
  name: z
    .string()
    .min(3, { error: "Name must be at least 3 characters" })
    .trim(),
});

export type CreateProjectState = {
  errors?: { name?: string[]; _form?: string[] };
  message?: string;
  success?: boolean;
};

export async function createProjectAction(
  _: CreateProjectState,
  formData: FormData
): Promise<CreateProjectState> {
  const validated = CreateProjectSchema.safeParse({
    name: formData.get("name"),
  });

  if (!validated.success) {
    return {
      errors: validated.error.flatten().fieldErrors,
    };
  }

  try {
    await fetchClient("/projects", {
      method: "POST",
      body: JSON.stringify({ name: validated.data.name }),
    });

    revalidatePath("/");
    return { success: true, message: "Project created successfully" };
  } catch (err: unknown) {
    const msg = getErrorMessage(err);

    if (msg.includes("409") || msg.toLowerCase().includes("Conflict")) {
      return {
        errors: { name: ["You already have a project with this name."] },
      };
    }
    return { errors: { _form: [msg] } };
  }
}
