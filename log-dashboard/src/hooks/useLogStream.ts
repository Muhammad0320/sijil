import { useEffect, useRef, useState } from "react";

export type LogEntry = {
  timestamp: string;
  level: string;
  service: string;
  message: string;
};

type ConnectionStatus = "CONNECTING" | "OPEN" | "CLOSED" | "ERROR";

export function useLogStream(projectID: number, token: string | null) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [status, setStatus] = useState<ConnectionStatus>("CLOSED");
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!projectID || !token) return;

    const wsUrl = `ws://localhost:8080/api/v1/logs/ws?project_id=${projectID}&token=${token}`;

    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log("WS Connected");
      setStatus("OPEN");
    };

    ws.current.onmessage = (event) => {
      try {
        const newLog: LogEntry = JSON.parse(event.data);

        setLogs((prevLogs) => [newLog, ...prevLogs]);
      } catch (error) {
        console.error("Failed to pase log:", error);
      }
    };

    ws.current.onclose = () => setStatus("CLOSED");
    ws.current.onerror = () => setStatus("ERROR");

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [projectID, token]);

  return { logs, status };
}
