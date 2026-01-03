"use server";

import { fetchClient } from "@/lib/client";
import { getSession } from "@/lib/session";
import { LogEntry } from "@/lib/types";

export async function getLogsAction(
  projectId: number,
  query: string = "",
  page: number,
  limit: number
) {
  const token = await getSession();

  const params = new URLSearchParams();
  params.set("project_id", projectId.toString());
  if (query) params.set("q", query);
  params.set("limit", limit.toString());
  params.set("page", page.toString());

  try {
    const data = await fetchClient<{ logs: LogEntry[] }>(
      `/logs${params.toString()}`,
      { method: "GET" },
      token
    );

    return { success: true, data: data.logs || [] };
  } catch (error) {
    console.error("failed to fetch logs", error);
    return { success: false, data: [] };
  }
}
