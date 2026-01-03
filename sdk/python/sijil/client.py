import threading
import queue
import time
import json
import requests
from datetime import datetime
import sys

ENDPOINT = "http://localhost:8080/api/v1/logs";
BATCH_SIZE = 100;
MAX_QUEUE_SIZE = 4096;
WORKER_COUNT = 3;
MAX_RETRIES=3


class SijilLogger: 
    def __init__(self, api_key, api_secret, endpoint=None, service="default"): 
        if not api_key or not api_secret: 
            raise ValueError("Sijil: Credentials missing")
        
        self.api_key = api_key
        self.api_secret = api_secret
        self.endpoint = endpoint or ENDPOINT
        self.service = service
        
        # State
        self.queue = queue.Queue(maxsize=MAX_QUEUE_SIZE)
        self._stop_event = threading.Event()
        
        # Worker pool
        self.workers = []
        for _ in range(WORKER_COUNT):
            t = threading.Thread(target=self._worker_loop, daemon=True)
            t.start()
            self.workers.append(t)
        
    def info(self, msg, data=None): self._push("info", msg, data)
    def warn(self, msg, data=None): self._push("warn", msg, data)
    def error(self, msg, data=None): self._push("error", msg, data)
    def debug(self, msg, data=None): self._push("debug", msg, data)
    def critical(self, msg, data=None): self._push("critical", msg, data)
    
    def _push(self, level, msg, data):
        if self._stop_event.is_set():
            return

        entry = {
            "level": level,
            "message": msg,
            "service": self.service_name,
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "data": data or {}
        }
        try:
            self.queue.put(entry, block=False)
        except queue.Full:
            # Non-blocking drop to protect app
            sys.stderr.write("Sijil Queue Full. Dropping log.\n")

    def _worker_loop(self):
        batch = []
        last_flush = time.time()
        
        while not self._stop_event.is_set() or not self.queue.empty():
            try:
                # Wait up to 1s for a log
                entry = self.queue.get(timeout=1.0)
                batch.append(entry)
            except queue.Empty:
                pass 

            # Flush Logic
            if len(batch) >= BATCH_SIZE or (time.time() - last_flush >= 1.0 and batch):
                self._send_batch(batch)
                batch = []
                last_flush = time.time()
            
            # If stopping and queue empty, break
            if self._stop_event.is_set() and self.queue.empty() and not batch:
                break

    def _send_batch(self, batch):
        payload = json.dumps(batch)
        headers = {
            "Content-Type": "application/json",
            "X-Api-Key": self.api_key,
            "Authorization": f"Bearer {self.api_secret}"
        }
        
        for attempt in range(MAX_RETRIES):
            try:
                resp = requests.post(self.endpoint, data=payload, headers=headers, timeout=5)
                # Success
                if 200 <= resp.status_code < 300:
                    return
                # Client Error -> Fail Fast
                if 400 <= resp.status_code < 500:
                    sys.stderr.write(f"Sijil Auth Error: {resp.status_code}\n")
                    return
                # Server Error -> Retry
            except Exception:
                pass # Network error -> Retry
            
            # Exponential Backoff
            time.sleep(0.1 * (2 ** attempt))

    def close(self):
        """Gracefully shutdown and flush pending logs"""
        self._stop_event.set()
        for t in self.workers:
            t.join(timeout=5.0)