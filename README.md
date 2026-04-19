# AI Gateway

A fork of [envoyproxy/ai-gateway](https://github.com/envoyproxy/ai-gateway) — an Envoy-based gateway for routing and managing AI provider traffic.

## Overview

AI Gateway provides a unified interface for routing requests to multiple AI providers (OpenAI, Anthropic, Gemini, Ollama, etc.) with support for:

- **Multi-provider routing** — route requests to different backends based on model, cost, or availability
- **Rate limiting** — per-user and per-model token/request rate limits
- **Observability** — metrics, tracing, and logging for AI traffic
- **Local models** — first-class support for Ollama and other local inference servers

## Prerequisites

- Go 1.22+
- Envoy proxy (see `.envoy-version` for required version)
- Docker (optional, for local development)
- Ollama (optional, for local model inference — install from [ollama.com](https://ollama.com))

## Getting Started

### Local Development with Ollama

1. Copy and configure environment variables:
   ```bash
   cp .env.ollama .env
   # Edit .env with your configuration
   ```

2. Start the gateway:
   ```bash
   go run ./cmd/ai-gateway
   ```

> **Personal note:** I primarily use this with Ollama running on `localhost:11434`. The `.env.ollama` defaults work out of the box for that setup — no edits needed.
>
> Useful Ollama commands: `ollama list` to see downloaded models, `ollama pull llama3` to grab a new one.
>
> Models I've been using: `llama3` for general tasks, `codellama` for code review, `mistral` for summarization, `phi3` for quick/lightweight queries.
>
> **Tip:** If Ollama is slow on first request, it's loading the model into memory. Subsequent requests are much faster. Run `ollama run llama3` once before starting the gateway to pre-warm it.
>
> **Tip:** To keep a model resident in memory indefinitely, run `ollama run llama3` and leave the session open in a separate terminal.
>
> **Tip:** Use `ollama ps` to see which models are currently loaded in memory — handy for debugging latency issues.
>
> **Tip:** `ollama run mistral` tends to OOM on my 16GB machine when `llama3` is already loaded. Unload first with `ollama stop llama3` if switching models.

### Configuration

The gateway is configured via Envoy xDS resources. See the `examples/` directory for sample configurations.

## Project Structure

```
.
├── cmd/            # Entry points
├── internal/       # Internal packages
├── pkg/            # Public packages
├── examples/       # Example configurations
└── tests/          # Integration and e2e tests
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 — see [LICENSE](LICENSE).
