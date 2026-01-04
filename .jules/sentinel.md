## 2024-05-23 - Rate Limiting Missing on Auth Endpoints
**Vulnerability:** The `/login` and `/register` endpoints were unprotected against brute-force attacks and abuse, lacking rate limiting middleware.
**Learning:** High-traffic public endpoints like authentication routes must always have rate limiting to prevent credential stuffing and DoS attacks.
**Prevention:** Apply the `rateLimitMiddleware` to all sensitive public endpoints (auth, password reset, etc.) by default during route registration.
