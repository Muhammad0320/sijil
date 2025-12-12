import { ReactNode } from "react";

// This layout is THIN. It just wraps the page.
export default function DashboardLayout({ children }: { children: ReactNode }) {
  return <div style={{ height: "100vh", width: "100vw" }}>{children}</div>;
}
