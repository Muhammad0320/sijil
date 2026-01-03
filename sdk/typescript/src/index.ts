const ENDPOINT = "http://localhost:8080/api/v1/logs";
const BATCH_SIZE = 100;
const MAX_RETRIES = 3;
const MAX_QUEUE_SIZE = 4096;
const WORKER_COUNT = 3;

interface Config {
  apiKey: string;
  apiSecret: string;
  endpoint?: string;
  flushInterval?: number;
  service?: string;
}

interface LogEntry {
  level: "info" | "error" | "warn" | "debug" | "critical";
  message: string;
  service?: string;
  timestamp: string; // ISO string
  data?: Record<string, any>;
}

export class SijilLogger {
  private config: {
    apiKey: string;
    apiSecret: string;
    endpoint: string;
    flushInterval: number;
  };

  private queue: LogEntry[] = [];
  private timer: ReturnType<typeof setInterval> | null = null;
  private isClosed = false;
  private serviceName: string = "";

  // Worker semaphore
  private activeWorkers: number = 0;

  constructor(config: Config) {
    if (!config.apiKey || !config.apiSecret) {
      throw new Error("Sijil: Credentials Missing!");
    }

    this.serviceName = config.service || "service";
    this.config = {
      apiKey: config.apiKey,
      apiSecret: config.apiSecret,
      endpoint: config.endpoint || ENDPOINT,
      flushInterval: Math.max(config.flushInterval || 1000, 250),
    };

    this.timer = setInterval(() => this.flush(), this.config.flushInterval);
  }

  public setService(name: string) {
    this.serviceName = name;
  }

  public info(message: string, data?: Object) {
    this.push("info", message, data);
  }
  public debug(message: string, data?: Object) {
    this.push("debug", message, data);
  }
  public warn(message: string, data?: Object) {
    this.push("warn", message, data);
  }
  public error(message: string, data?: Object) {
    this.push("error", message, data);
  }

  public critical(message: string, data?: Object) {
    this.push("critical", message, data);
  }

  private push(level: LogEntry["level"], message: string, data?: Object) {
    // Safety Cap
    if (this.queue.length >= MAX_QUEUE_SIZE) {
      console.warn("Sijil Queue Full. Dropping logs");
      return;
    }

    this.queue.push({
      level,
      message,
      data: data || {},
      timestamp: new Date().toISOString(),
    });

    if (this.queue.length >= BATCH_SIZE) {
      this.flush();
    }
  }

  private async flush() {
    if (this.queue.length === 0) return;
    if (this.activeWorkers >= WORKER_COUNT) return;

    const batch = this.queue.splice(0, BATCH_SIZE);
    this.activeWorkers++;

    try {
      await this.sendWithRetry(batch);
    } catch (error) {
      console.error("Sijil Delivery failed:", error);
    } finally {
      this.activeWorkers--;
    }
  }

  private async sendWithRetry(batch: LogEntry[], attempt = 0): Promise<void> {
    try {
      const res = await fetch(this.config.endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Api-Key": this.config.apiKey,
          Authorization: `Bearer ${this.config.apiSecret}`,
        },
        body: JSON.stringify(batch),
      });

      if (!res.ok) {
        if (res.status >= 400 && res.status < 500) {
          console.error(`Sijil Rejected: ${res.status}`);
          return;
        }
        if (res.status >= 500) throw new Error(`Server Error ${res.status}`);
      }
    } catch (error) {
      if (attempt < MAX_RETRIES) {
        // Exponential Backoff
        await new Promise((r) => setTimeout(r, 100 * Math.pow(2, attempt)));
        return this.sendWithRetry(batch, attempt + 1);
      }
      throw error;
    }
  }

  public async close() {
    this.isClosed = true;
    if (this.timer) clearInterval(this.timer);

    while (this.queue.length > 0) {
      await this.flush();
    }
  }
}
