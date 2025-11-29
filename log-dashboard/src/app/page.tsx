// src/app/page.tsx (Server Component)

import { redirect } from "next/navigation";
import { getSession } from "@/lib/session";
import { getProjects } from "@/actions/projects";
import DashboardClient from "./DashboardClient"; // We create this next

export default async function DashboardPage() {
  // 1. Verify Session (HttpOnly Cookie)
  const token = await getSession();

  if (!token) {
    redirect("/login");
  }

  // 2. Fetch Initial Data (Server-to-Server)
  const [projects, err] = await getProjects();

  return (
    <DashboardClient
      initialProjects={projects || []}
      token={token}
      serverError={err}
    />
  );
}
