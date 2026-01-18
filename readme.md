A high-performance, containerized microservices reference demonstrating a **Transparent Proxy**, **Distributed Caching (Cache-Aside)**, and **External API Integration**. Designed for local Kubernetes development (k3d/k3s) and easily portable to other Kubernetes environments.

---

## Table of Contents

1. Project Summary
2. Key Features
3. Architecture Overview
4. Data Flow & Caching Logic
5. Project Structure
6. Tech Stack
7. Kubernetes Manifests (High Level)
8. Deployment — Step-by-Step
9. Testing & Verification
10. Troubleshooting
11. Operational Best Practices
12. Security Considerations
13. Contributing
14. License

---

## Project Summary

This repository contains two Go microservices and a Redis backing store. Service1 acts as a gateway/transparent proxy. Service2 contains the weather logic, handles caching against Redis, and queries the Open-Meteo API on cache misses. The system uses a **60-second TTL** for cached entries to balance freshness and request volume to the external API.

---

## Key Features

- Clear separation of responsibilities (Gateway vs Weather engine).
- Cache-Aside pattern with Redis for high throughput and low latency.
- Generic JSON forwarding for forward-compatible schemas.
- Kubernetes-native deployment with local k3d support.

---

## Architecture Overview

Client  
→ **Service1 (Gateway)** — Transparent Proxy  
→ **Service2 (Weather Engine)** — Cache + External API  
→ **Redis** — In-memory cache  
→ **Open-Meteo API** — External provider

The gateway isolates clients from internal topology. Business logic and caching remain centralized in Service2.

---

## Data Flow & Caching Logic

1. Client sends `POST /proxy` with JSON payload.
2. Service1 unmarshals and forwards payload to Service2 `/process`.
3. Service2 checks Redis:
   - **HIT:** Return cached JSON.
   - **MISS:** Call Open-Meteo, store result in Redis (TTL 60s), return response.

---

## Project Structure

```
.
├── k8s/
│   ├── service1.yaml
│   ├── service2.yaml
│   └── redis.yaml
├── service1/
│   ├── main.go
│   └── Dockerfile
└── service2/
    ├── main.go
    └── Dockerfile
```

---

## Tech Stack

- Language: Go 1.21+
- Containerization: Docker (multi-stage builds)
- Orchestration: Kubernetes (k3d/k3s)
- Cache: Redis (alpine)
- External API: Open-Meteo

---

## Kubernetes Manifests (High Level)

- `service1.yaml` — Gateway Deployment + Service
- `service2.yaml` — Weather Engine Deployment + Service
- `redis.yaml` — Redis Deployment + Service

---

## Deployment — Step-by-Step

### Prerequisites

- Docker
- k3d
- kubectl

### Create Cluster

```bash
k3d cluster create micro-cluster -p "8080:80@loadbalancer"
```

### Build Images

```bash
docker build -t service1:local ./service1
docker build -t service2:local ./service2
docker pull redis:alpine
```

### Import Images

```bash
k3d image import service1:local service2:local redis:alpine -c micro-cluster
```

### Deploy

```bash
kubectl apply -f k8s/
```

### Verify

```bash
kubectl get pods --watch
```

---

## Testing & Verification

### Send Request

```bash
curl -i -X POST http://localhost:8080/proxy   -H "Content-Type: application/json"   -d '{"city": "london"}'
```

### Observe Logs

```bash
kubectl logs -f deployment/service2
```

### Inspect Redis

```bash
kubectl exec -it deployment/redis -- redis-cli
keys *
ttl london
```

---

## Troubleshooting

- Use `kubectl describe pod` and `kubectl logs` for startup issues.
- Ensure images are imported into k3d.
- Validate Redis service name and connectivity.

---

## Operational Best Practices

- Define resource requests and limits.
- Make cache TTL configurable.
- Implement graceful shutdown handling.
- Add metrics and structured logging for production use.

---

## Security Considerations

- Use Kubernetes Secrets for sensitive data.
- Run containers as non-root.
- Restrict Redis to internal cluster access.

---

## Contributing

Fork the repository, create a feature branch, and open a pull request with a clear description.

---

## License

Add a LICENSE file (MIT recommended) before publishing.
