"use server";

import { fetchClient } from "@/lib/client";
import { getSession } from "@/lib/session";
import { revalidatePath } from "next/cache";

export async function inviteMemberAction(
  projectId: number,
  email: string,
  role: string
) {
  const token = await getSession();
  if (!token) throw new Error("Unauthorized");

  await fetchClient(
    `/projects/${projectId}/members`,
    {
      method: "POST",
      body: JSON.stringify({ email, role }),
    },
    token
  );

  revalidatePath("/");
  return { success: true };
}
