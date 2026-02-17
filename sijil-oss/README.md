# Sijil OSS

**The Paranoid Log Engine.** Zero Retention. AI Forensics. 100k+ Logs/Sec.

Sijil is a high-performance observability platform built for developers who care about speed, security, and privacy. We strip PII before it hits the disk and delete your logs the moment you don't need them.

## 🚀 Quick Start

### 1. Node.js (Next.js - server comp - / Express)

```bash
npm install @sijil/node
```

```code
import { SijilLogger } from "@sijil/node";

const logger = new SijilLogger({
  apiKey: process.env.SIJIL_API_KEY,
  apiSecret: process.env.SIJIL_API_SECRET,
  service: "payment-api",
});

logger.info("Payment processed", { amount: 500 });
```

### 2. Python (FastAPI / Django / Flask)

```bash
pip install sijil
```

```code
from sijil import SijilLogger

logger = SijilLogger(
    api_key="pk_live_...",
    api_secret="sk_live_...",
    service="auth-service"
)

logger.warn("Login attempt failed", {"user_id": "u_8821"})
```

### 2. GO (Fiber / Gin / Stdlib)

go get [github.com/sijil/go-sdk](https://github.com/sijil/go-sdk)

```code
package main

import "github.com/sijil/go-sdk"

func main() {
    client := sijil.NewClient(sijil.Config{
        APIKey:    "pk_live_...",
        APISecret: "sk_live_...",
        Service:   "worker-1",
    })

    client.Info("Job started", map[string]interface{}{"job_id": 101})
}

```

## Sijil Agent (Sidecar)

Got legacy apps writing to files? Use the Sijil Agent to tail logs and push them to the cloud without changing a line of code.

### Installation

Download the binary for your OS from the [Release Page](https://github.com/Muhammad0320/sijil/releases/tag/v1.0.0)

### Usage

```bash
# Linux / Mac
./sijil-agent -f /var/log/nginx/access.log \
  -pk pk_live_... \
  -sk sk_live_... \
  -s nginx-loadbalancer

# Windows (PowerShell)
.\sijil-agent.exe -f C:\Logs\app.log -pk ... -sk ...
```

### Flags:

- `-f`: Path to the file to tail.
- `-s`: Service name (e.g., `frontend-web`).

- `-pk`: Your Public API Key.
- `-sk`: Your Secret Key.
- `-format`: `regex` (default) or `json` depending on your log file format
- `-url`: custom url for self hosted users (completely optional)

### 🤝Contributing

We welcome contributions! Please see CONTRIBUTING.md for more details on how to set up the dev environment

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing feature`)
5. Open a Pull Request

---

**Built with ❤️ by Sijil Team**
