# 🧠 HydrAIDE 2.0 – Adaptive Intelligent Data Engine

![Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg?style=for-the-badge)
![Status](https://img.shields.io/badge/status-Production%20Ready-brightgreen?style=for-the-badge)
![Version](https://img.shields.io/badge/version-2.0-informational?style=for-the-badge)
![Performance](https://img.shields.io/badge/speed-O(1)%20Access-ff69b4?style=for-the-badge)

> **Production-ready. Enterprise-grade. Battle-tested.**
> Built and refined over **3+ years** of real-world use, Hydra 2.0 delivers unmatched performance, control, and flexibility.

Hydra is not just another database.  
It’s a lightning-fast, memory-conscious, real-time data engine –  
built to be as **intelligent** as the data it stores.

> 🛠️ Hydra is actively evolving.  
> We're currently adding:
> - 🧪 Full proto documentation
> - 💻 Go SDK (`hydrun-go`)
> - 🧵 Real-world examples

Stay tuned and ⭐ the repo to follow along!

---

> ⚠️ This repository does not contain the source code of the Hydra engine itself.  
> Hydra is a closed-source, production-grade system.  
> This repo contains the public API (`.proto`), usage documentation, and SDK integrations.

**Enterprise licensing available** – contact us: [peter.gebri@hydraide.com](mailto:peter.gebri@hydraide.com)

---

## 📖 Documentation & Resources

🚧 **Documentation is currently being uploaded and expanded.**  
Stay tuned for updates! 🚀

🔹 **[API Reference](docs/api_reference.md)** – Full API documentation (proto-based) 
🔹 **[Usage Guide](docs/usage.md)** – How to set up and use Hydra  
🔹 **[Examples](docs/examples/basic_usage.md)** – Real-world integration examples  
🔹 **[Roadmap](docs/roadmap.md)** – Upcoming features & plans

---

## 🔥 Why Hydra?

- ⚡ **Blazing fast access**, regardless of data size.
- 📎 **90% storage saving** with binary and compressed data.
- ↻ **Real-time subscriptions** and adaptive cache control.
- 🔐 **Built-in locking mechanism** – safe multi-process coordination.
- 🧠 **Customizable per collection**: memory lifetime, file limits, write intervals.
- 🌐 **Distributed-ready**: no orchestrator needed to scale across multiple servers.
- 🚀 **Simple gRPC interface** – easy to integrate from any language.

---

## 🧬 How does Hydra work?

Hydra organizes data into units called **Swamps**.  
Each Swamp is a self-contained, independently optimized mini-storage:

- Stored directly as a file (or memory-only).
- Accessed via a predictable path → **O(1) access speed**.
- Automatically chunked, compressed, binary-stored with [Snappy](https://github.com/google/snappy) and [GOB](https://golang.org/pkg/encoding/gob/).

> No SQL. No schema headache. No unnecessary parsing.  
> Just your data – **in, out, real-time**.

---

## 🐉 The Hydra Analogy

In mythology, the Hydra was a powerful multi-headed guardian creature.  
In our system, **Hydra watches over your data**.

- Each **Swamp** is a separate domain where data is stored.
- Inside each Swamp live the **Treasures** – the actual pieces of data.
- The Hydra can access and manipulate any Swamp instantly, across namespaces.

Namespaces follow this structure:
`Sanctuary → Realm → Swamp → Treasure`

Hydra intelligently routes and accesses data **without the need for any central orchestrator**.  
Using a deterministic hashing strategy, each Swamp is automatically mapped to the correct server – whether you're working with 1 or 100.

---

## 📦 Features at a glance

| Feature | Description |
|--------|-------------|
| 🧠 On-demand indexing | Indexes are built in memory only when queried, with zero storage overhead |
| 🎛️ Swamp-level control | Dynamically configure each swamp’s behavior (memory, flush, TTL) directly from code – no DB access required |
| 🔒 Locking | Swamp/key-level locking with TTL & context awareness |
| 📥 Set/Get | Insert, update, get data – type-safe & atomic |
| 🧹 Built-in Garbage Collector | When the last treasure in a swamp is deleted, Hydra automatically deletes the swamp and its entire file structure to reclaim space |
| ↺ Subscribe | Get notified on changes instantly (like pub/sub) |
| ⏳ Expire & Shift | Time-based expiry and cleanup |
| 📊 Count & Index | Indexed reads by creation, expiration, etc. |
| ➕ Conditional Increment | Atomic incs with rules (if x > 10, then…) |
| 📚 Slices | Special handling for Uint32 slices |
| 🚦 Exists check | Key & swamp existence support |
| 🧭 Distributed architecture | Data is spread across servers using a hash-based strategy, no orchestrator required |
---

## 🎯 Who is Hydra for?

Hydra was built for developers, data-driven teams, and platform architects who want:

- Full control over where data is stored – in RAM, on SSD, or on disk.
- Real-time performance, even with country-scale datasets.
- To eliminate the need for multiple services by combining caching, storage, real-time updates, and analytics into one engine.
- To scale with zero orchestrator overhead, and 100% predictability.

Hydra is ideal for:
- 🚀 **Startups** needing powerful data tools without infrastructure overhead
- 🧪 **Research labs & academic institutions** working with large-scale data
- 🌐 **Web crawlers & content indexing platforms** like Trendizz
- 📊 **Realtime dashboards & analytics systems**
- 📱 **IoT and edge computing systems** with limited memory
- 🛠️ **Engineers who demand precision and full control** over their data infrastructure

---

## 🧪 Who is using it?

Hydra is already powering high-volume platforms like **Trendizz.com**,
a cutting-edge search platform that indexes and analyzes the public web content of countries across Europe — including Hungary, Romania, Slovakia, and parts of the Czech Republic — to make them searchable based on every word they contain.

Trendizz needed a database engine that could handle **massive-scale, word-level indexing and search**, in real time, across millions of websites.

There was no database fast and efficient enough for this challenge — so we built Hydra.

Thanks to Hydra:
- Complex, precise word-based searches across all web content take **just 1–2 seconds**.
- The system continuously crawls and re-indexes websites, so performance and resource efficiency are **absolutely critical**.
- **Hydra enables this** with near-zero memory overhead, distributed file-level storage, and blazing-fast access.

> 2TB of unoptimized data → compressed and optimized to just 200GB.

Hydra achieves this using:
- ✔ Real-time adaptive in-memory cache control
- ✔ Per-swamp chunking and SSD-safe writes
- ✔ Snappy compression + GOB encoding (zero conversion overhead)
- ✔ Intelligent memory & expiration logic per swamp
- ✔ File-system mapped swamp paths for **O(1)** direct access
- ✔ Distributed data access using swamp-based hashing – no orchestrator needed

---

## 🏱 Ready for SDKs & Integrations

This repository contains the core `.proto` files of Hydra.  
You're free to generate your own client SDKs using:

```bash
protoc --go_out=. --go-grpc_out=. hydra-service.proto
```

Official SDKs (under development):
- 📅 Go SDK (`hydrun`) – actively used in Trendizz
- ⏳ Python SDK
- ⏳ Node.js SDK
- ⏳ Rust SDK
- ⏳ Java SDK

All based on the same robust gRPC interface.

---

## 💬 Get in touch

> Hydra is proudly developed in **Go**, engineered for high performance and low memory usage.
> Originally built by **Peter Gebri** to power Trendizz, it is now open to the world as an enterprise-ready solution.


Interested in trying Hydra for your project?  
Want to build an SDK or use it in research?

📩 Contact: **Peter Gebri** – [peter.gebri@hydraide.com](mailto:peter.gebri@hydraide.com)  
Or join the waitlist at [hydrAIDE.com](https://hydraide.com)

## 🐞 Found something? Need help?

Feel free to open an [issue](https://github.com/hydraide/hydraide/issues)!

---

> 🔍 *Hydra doesn’t search your data.*  
> *It knows exactly where it is.*

