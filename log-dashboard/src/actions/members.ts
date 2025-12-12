"use server";

import { fetchClient } from "@/lib/client";
import { getSession } from "@/lib/session";
import { revalidatePath } from "next/cache";
import { InviteSchema, InviteState } from "../lib/definitions";

export async function inviteMemberAction(
  prevstate: InviteState,
  formData: FormData
): Promise<InviteState> {
  const token = await getSession();

  const validation = InviteSchema.safeParse({
    email: formData.get("email"),
    role: formData.get("role"),
    projectId: formData.get("projectId"),
  });

  if (!validation.success) {
    return {
      errors: validation.error.flatten().fieldErrors,
    };
  }

  const { email, role, projectId } = validation.data;

  try {
    await fetchClient(
      `/projects/${projectId}/members`,
      { method: "POST", body: JSON.stringify({ email, role }) },
      token
    );

    revalidatePath("/dashboard");

    return { success: true, message: `Invited ${email} as ${role}` };
  } catch (error) {
    if (error instanceof Error) {
      return {
        errors: {
          _form: [error.message || "Failed to invite member"],
        },
      };
    }

    return {
      errors: { _form: ["Failed to invite member"] },
    };
  }
}
