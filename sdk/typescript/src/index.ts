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
  level: "info" | "error" | "warn" | "debug";
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
  private timer: NodeJS.Timeout | null = null;
  private activeRequests = 0;
  private readonly MAX_CURRENT_REQUEST = 5;
  private serviceName = "default";

  constructor(config: Config) {
    if (!config.apiKey || !config.apiSecret) {
      throw new Error("LogEngine: Credentials Missing!");
    }

    this.config = {
      apiKey: config.apiKey,
      apiSecret: config.apiSecret,
      endpoint: config.endpoint || "http://localhost:8080/api/v1/logs",
      flushInterval: Math.max(config.flushInterval || 1000, 250),
      batchSize: config.batchSize || 100,
      maxRetries: config.maxRetries || 3,
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

  private push(level: LogEntry["level"], message: string, data?: Object) {
    // Safety Cap
    if (this.queue.length >= 5000) {
      console.warn("LogEngine Queue Full. Dropping logs");
      return;
    }

    this.queue.push({
      level,
      message,
      data: data || {},
      timestamp: new Date().toISOString(),
    });

    if (this.queue.length >= this.config.batchSize) {
      this.flush();
    }
  }

  private async flush() {
    if (this.queue.length === 0) return;
    if (this.activeRequests >= this.MAX_CURRENT_REQUEST) return;

    const batch = this.queue.splice(0, this.config.batchSize);
    this.activeRequests++;

    try {
      await this.sendWithRetry(batch);
    } catch (error) {
      console.error("LogEngine Delivery failed:", error);
    } finally {
      this.activeRequests--;
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
        if (res.status >= 500) throw new Error(`Server Error ${res.status}`);
        if (res.status >= 400)
          console.error(`LogEngine Rejected: ${res.status}`);
      }
    } catch (error) {
      if (attempt < this.config.maxRetries) {
        // Exponential Backoff
        await new Promise((r) => setTimeout(r, 100 * Math.pow(2, attempt)));
        return this.sendWithRetry(batch, attempt + 1);
      }
      throw error;
    }
  }
}
