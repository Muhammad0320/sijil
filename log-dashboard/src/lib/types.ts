export interface Project {
  id: number;
  name: string;

  pending?: boolean;
}

export interface LogEntry {
  timestamp: string;
  level: string;
  service: string;
  message: string;
  project_id: number;
}

export function getErrorMessage(err: unknown): string {
  if (err instanceof Error) return err.message;
  if (typeof err === "string") return err;
  return "An unknown error occured";
}
