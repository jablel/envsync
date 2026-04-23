# envsync

> Diff and sync `.env` files across environments with secret masking support.

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git
cd envsync && go build -o envsync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.development .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync .env.development .env.production
```

**Mask secrets in output:**

```bash
envsync diff .env.development .env.production --mask-secrets
```

Example output:

```
~ API_URL        dev: http://localhost:3000  prod: https://api.example.com
+ NEW_FEATURE_FLAG  dev: true                prod: [missing]
- LEGACY_KEY        dev: [missing]           prod: ********
```

**Flags:**

| Flag             | Description                          |
|------------------|--------------------------------------|
| `--mask-secrets` | Redact values for sensitive keys     |
| `--export`       | Write synced output to a file        |
| `--dry-run`      | Preview changes without writing      |

---

## License

[MIT](LICENSE) © 2024 yourusername