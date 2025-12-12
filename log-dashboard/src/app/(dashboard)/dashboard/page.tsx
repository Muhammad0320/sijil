import { getProjects } from "@/actions/projects";
import { getSession } from "@/lib/session";
import DashboardClient from "./DashboardClient"; // Your big component

export default async function DashboardPage() {
  // 1. Get the Token (Cookies) - Server Side
  const token = await getSession();

  // 2. Fetch Data (Directly from Backend API) - Server Side
  // This happens BEFORE the user sees the page.
  const [projects, error] = await getProjects();

  // 3. Feed the Client Component
  // This "hydrates" the interactive part of your app
  return (
    <DashboardClient
      initialProjects={projects || []}
      token={token || ""}
      serverError={error}
    />
  );
}
