import { getProjects } from "@/actions/projects";
import { getSession } from "@/lib/session";
import DashboardClient from "./DashboardClient"; // Your big component
import { DashboardProvider } from "@/providers/DashboardProviders";

export default async function DashboardPage() {
  const token = await getSession();

  // 2. Fetch Data (Directly from Backend API) - Server Side
  // This happens BEFORE the user sees the page.
  const [projects, error] = await getProjects();

  // 3. Feed the Client Component
  // This "hydrates" the interactive part of your app
  return (
    <DashboardProvider initialProjects={projects || []} token={token || ""}>
      <DashboardClient serverError={error} />
    </DashboardProvider>
  );
}
