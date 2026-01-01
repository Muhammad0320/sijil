import { LogEntry } from "@/lib/types";
import { useToast } from "@/providers/ToastProvider";
import { useEffect, useRef, useState } from "react";
import { z } from "zod";

const LogSchema = z.object({
  timestamp: z.string(),
  level: z.string(),
  service: z.string(),
  message: z.string(),
  project_id: z.number(),
});

type ConnectionStatus = "CONNECTING" | "OPEN" | "CLOSED" | "ERROR";

export function useLogStream(projectID: number, token: string | null) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [status, setStatus] = useState<ConnectionStatus>("CLOSED");

  const ws = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<NodeJS.Timeout | null>(null);
  const retryCount = useRef(0);
  const toast = useToast();

  useEffect(() => {
    if (!projectID || !token) return;

    const wsBase =
      process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1";
    const wsUrl = `${wsBase}/logs/ws?project_id=${projectID}&token=${token}`;

    function connect() {
      // Safety check: Don't connect if already open
      if (ws.current?.readyState === WebSocket.OPEN) return;

      setStatus("CONNECTING");
      const socket = new WebSocket(wsUrl);
      ws.current = socket;

      socket.onopen = () => {
        console.log("🟢 WS Connected");
        setStatus("OPEN");
        retryCount.current = 0;
        // Only show toast if we were previously disconnected (reconnection success)
        if (retryCount.current > 0) {
          toast.success("Live stream reconnected ✅");
        }
      };

      socket.onmessage = (event) => {
        try {
          const raw = JSON.parse(event.data);
          const result = LogSchema.safeParse(raw);
          if (result.success) {
            setLogs((prev) => [result.data, ...prev]);
          }
        } catch (error) {
          console.error("WS Parse Error", error);
        }
      };

      socket.onclose = (event) => {
        if (event.wasClean) {
          setStatus("CLOSED");
          return;
        }

        setStatus("ERROR");
        ws.current = null;

        const delay = Math.min(1000 * Math.pow(2, retryCount.current), 30000);

        // Notify user of connection loss
        if (retryCount.current === 0) {
          toast.error(`Connection lost. Retrying in ${delay / 1000}s...`);
        }

        console.log(`⚠️ WS Closed. Reconnecting in ${delay}ms...`);

        reconnectTimeout.current = setTimeout(() => {
          retryCount.current++;
          connect(); // Recursive call is safe here inside useEffect scope
        }, delay);
      };

      socket.onerror = () => {
        socket.close();
      };
    }

    connect();

    return () => {
      // Cleanup
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
      if (ws.current) {
        ws.current.close(1000, "Unmounting");
      }
    };
  }, [projectID, token, toast]); // Add toast to dependencies

  return { logs, status };
}
