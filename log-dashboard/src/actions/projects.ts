"use server";

import { fetchClient } from "@/lib/client";
import { getSession } from "@/lib/session";
import { getErrorMessage, Project } from "@/lib/types";
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

  data?: {
    apiKey: string;
    apiSecret: string;
    projectId: number;
    name: string;
  };
};

export async function createProjectAction(
  _: CreateProjectState,
  formData: FormData
): Promise<CreateProjectState> {
  const token = await getSession();

  const validated = CreateProjectSchema.safeParse({
    name: formData.get("name"),
  });

  if (!validated.success) {
    return {
      errors: validated.error.flatten().fieldErrors,
    };
  }

  try {
    const res = await fetchClient<{
      message: string;
      project_id: number;
      api_key: string;
      api_secret: string;
    }>(
      "/projects",
      {
        method: "POST",
        body: JSON.stringify({ name: validated.data.name }),
      },
      token
    );

    revalidatePath("/");
    return {
      success: true,
      message: "Project created successfully",
      data: {
        apiKey: res.api_key,
        apiSecret: res.api_secret,
        projectId: res.project_id,
        name: validated.data.name
      },
    };
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

type FetchResult<T> = Promise<[T, null] | [null, string]>;

export async function getProjects(): Promise<FetchResult<Project[]>> {
  const token = await getSession();

  try {
    const projects = await fetchClient<Project[]>(
      "/projects",
      {
        method: "GET",
      },
      token
    );

    return [projects, null];
  } catch (err) {
    const msg = getErrorMessage(err);
    console.error("Fetch projects Error:", msg);

    return [null, msg];
  }
}
