![Hydra](images/hydra-banner-stacked-text.png)

# 🧠 HydrAIDE 2.0 – The Adaptive, Intelligent Data Engine

![Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg?style=for-the-badge)
![Status](https://img.shields.io/badge/status-Production%20Ready-brightgreen?style=for-the-badge)
![Version](https://img.shields.io/badge/version-2.0-informational?style=for-the-badge)
![Speed](https://img.shields.io/badge/Access-O(1)%20Always-ff69b4?style=for-the-badge)

> **Hydra isn’t just a database. It’s a philosophy.**
> Built to serve real-time, reactive systems where every operation is intentional — and everything else vanishes.

Welcome to the engine behind platforms like [Trendizz.com](https://trendizz.com), where billions of records are searched, updated, and streamed in **real-time**. Hydra 2.0 is the result of 3+ years of battle-tested evolution, powering search engines, dashboards, crawlers, and more.

Hydra brings:
- ⚡ **O(1) access to billions of objects**
- 🔄 **Real-time reactivity with built-in subscriptions**
- 🧠 **In-memory indexes built only when needed**
- 🧹 **Zero garbage, no cron jobs, no leftovers**
- 🌐 **Distributed scaling without orchestrators**

> **"Hydra doesn’t search your data. It knows where it is."**

---

## 📚 Start Here: The Hydra Documentation

To truly understand Hydra, start with its **core philosophy and architecture**:

👉 [**Thinking in Hydra – The Philosophy of Reactive Data**](docs/thinking-in-hydra/thinking-in-hydra.md)  
*Learn how Hydra redefines structure, access, and system design from the ground up.*

### Then continue the 9-step journey:
| Step | Section | Description |
|------|---------|-------------|
| 1️⃣ | [📍 Naming Convention](docs/thinking-in-hydra/naming-convention.md) | Learn how data structure begins with naming – not schemas. |
| 2️⃣ | [🌿 Swamp Pattern](docs/thinking-in-hydra/swamp-pattern.md) | Configure persistence, memory, and lifespan directly from code. |
| 3️⃣ | [💎 Treasures](docs/thinking-in-hydra/treasures.md) | Understand the smallest, most powerful unit of data. |
| 4️⃣ | [🧩 Indexing](docs/thinking-in-hydra/indexing.md) | Discover ephemeral, in-memory indexing that feels like magic. |
| 5️⃣ | [🔄 Subscriptions](docs/thinking-in-hydra/subscriptions.md) | Build reactive systems natively with Hydra’s event engine. |
| 6️⃣ | [🔐 Locking](docs/thinking-in-hydra/locking.md) | Achieve true concurrency without fear. |
| 7️⃣ | [🧹 Clean System](docs/thinking-in-hydra/clean-system.md) | Never think about cleanup again – because Hydra already did. |
| 8️⃣ | [🌐 Distributed Architecture](docs/thinking-in-hydra/distributed-architecture.md) | Scale horizontally without orchestration pain. |
| 9️⃣ | [🚀 Install & Update](docs/thinking-in-hydra/how-to-install-update-hydra.md) | Deploy Hydra in minutes, not days. |


---

### 💼 For CEOs – Why Your Company Needs Hydra

> **You’re building fast. Scaling faster. But your data engine is slowing you down.**

Hydra was built for founders and leaders who believe their teams deserve better.  
No background daemons. No stale indexes. No technical debt.

Just **instant data flow**, **zero waste**, and **developer-native architecture** that lets your team move at light speed.

🌟 **See how Hydra can transform your product velocity and reduce infra costs.**

👉 [Read the CEO Guide →](docs/for-ceos.md)

---

## 🔥 Why Developers Choose Hydra

| Feature | What It Means |
|--------|---------------|
| ⚡ Instant Swamp access | O(1) folder-mapped resolution, no indexing required |
| 🧠 On-the-fly indexing | Built in RAM only when needed – disappears after use |
| 🔄 Subscriptions | Real-time streams of change, built into every write/update/delete |
| 🔐 Zero-deadlock locking | Per-Treasure locks and TTL-protected business locks |
| 🧹 Auto-cleanup | Swamps and memory vanish when no longer needed |
| 🌐 Horizontal scaling | Stateless by default, no orchestrator required |
| 💾 SSD-friendly design | Chunked binary storage with delta writes |
| 🧬 Fully typed | Store structs, slices, maps – native to your language |

> You never define schemas. You never worry about cleanup. You just build.

---

## 🌀 Who Hydra is For

Hydra is ideal for:
- 🚀 Startups that need a modern database without infrastructure burden
- 🌐 High-volume web crawlers and analytics platforms (like Trendizz)
- 📊 Live dashboards, reactive pipelines, and streaming interfaces
- 🧪 Research and ML pipelines with ephemeral data flows
- 📱 IoT and edge devices with minimal memory

And especially for developers who:
- Want to **own their data logic in code**
- Need **real-time data flow without middleware**
- Hate daemons, cron jobs, and config bloat

> Hydra is **developer-native**. You don’t configure it. You *declare intent* — and it responds.

---

## 🚀 Install & Run in Minutes

Hydra is container-first. Just use Docker Compose:

```yaml
services:
  hydra-server:
    image: ghcr.io/hydraide/hydraserver:VERSION
    ports:
      - "4900:4444"
    volumes:
      - /path/to/data:/hydra/data
      - /path/to/settings:/hydra/settings
      - /path/to/certs:/hydra/certificate
    environment:
      - GRPC_SERVER_ERROR_LOGGING=true
    stop_grace_period: 10m
```

➡ Full install guide: [Install & Update Hydra →](docs/how-to-install-update-hydra.md)

---

## 💻 SDKs – Native Integration in Your Language

Hydra speaks **gRPC**, and every SDK is powered by the same `.proto` definition:

| Language | Status | Repo |
|----------|--------|------|
| Go       | ✅ Active | `hydrungo` |
| Node.js  | 🧪 Planning | `hydrunjs` |
| Python   | 🧠 In design | `hydrunpy` |
| Rust     | 🧠 In design | `hydrunrs` |
| Java     | 🧠 In design | `hydrunjv` |
| C# / .NET| 🧠 In design | `hydruncs` |

More SDKs: C++, Kotlin, Swift – all on the roadmap.

---

## 🙌 Want to Contribute?

Start by reading the [Contributor Introduction](./CONTRIBUTORS.md) – it explains why Hydra exists, what kind of people we’re looking for, and how you can join.
Then check out our [Contribution Guide](./CONTRIBUTING.md) – it walks you through the practical steps.

Once you're ready, open your first issue or pull request — we’ll be waiting! 🚀

---

## 📩 Contact & Enterprise

Hydra is used in production at [Trendizz.com](https://trendizz.com). Interested in enterprise licensing, SDK development, or embedding Hydra in your own platform?

📧 **Peter Gebri** – [peter.gebri@trendizz.com](mailto:peter.gebri@trendizz.com)
(Founder of Hydra & Trendizz)
🌐 **Website** – [https://hydraide.com](https://hydraide.com) Currently in progress and directly linked to GitHub.

---

📍 See the full [Hydra Roadmap](./ROADMAP.md) → What we’re building and dreaming next

---

> 🧠 Hydra isn’t a tool you use.
> It’s a system you think with.

Join the movement. Build different.

