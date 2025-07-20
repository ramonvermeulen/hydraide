![HydrAIDE](images/hydraide-banner.jpg)

# HydrAIDE - The Adaptive, Intelligent Data Engine

[![License](https://img.shields.io/badge/license-Apache--2.0-blue?style=for-the-badge)](http://www.apache.org/licenses/LICENSE-2.0)
![Version](https://img.shields.io/badge/version-2.0-informational?style=for-the-badge)
![Status](https://img.shields.io/badge/status-Production%20Ready-brightgreen?style=for-the-badge)
![Speed](https://img.shields.io/badge/Access-O(1)%20Always-ff69b4?style=for-the-badge)
![Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=for-the-badge&logo=go)
[![Join Discord](https://img.shields.io/discord/1355863821125681193?label=Join%20us%20on%20Discord&logo=discord&style=for-the-badge)](https://discord.gg/tYjgwFaZ)


> **HydrAIDE is a zero-waste, real-time data engine. Like a database, but built for reactive systems**

> ⚠️ **Just a heads-up**
> I’m currently the sole developer of the HydrAIDE project.  
> Due to serious time constraints, all documentation was generated with a bit of AI help based on my own notes.  
> **But the code? 100% mine — no vibes coding involved 😄**  
> Just figured I’d say this upfront before someone yells at me for it.

Welcome to the engine behind platforms like [Trendizz.com](https://trendizz.com), where billions of records are searched, updated, and streamed in **real-time**. HydrAIDE 2.0 is the result of 3+ years of battle-tested evolution, powering search engines, dashboards, crawlers, and more.

HydrAIDE brings:
- ⚡ **O(1) access to billions of objects**
- 🔄 **Real-time reactivity with built-in subscriptions**
- 🧠 **In-memory indexes built only when needed**
- 🧹 **Zero garbage, no cron jobs, no leftovers**
- 🌐 **Distributed scaling without orchestrators**

> **"HydrAIDE doesn’t search your data. It knows where it is."**

---

## 📚 Start Here: The HydrAIDE Documentation

To truly understand HydrAIDE, start with its **core philosophy and architecture**:

👉 [**Thinking in HydrAIDE – The Philosophy of Reactive Data**](docs/thinking-in-hydraide/thinking-in-hydraide.md)  
*Learn how HydrAIDE redefines structure, access, and system design from the ground up.*

### Then continue the 9-step journey:
| Step | Section | Description |
|------|---------|-------------|
| 1️⃣ | [📍 Naming Convention](docs/thinking-in-hydraide/naming-convention.md) | Learn how data structure begins with naming – not schemas. |
| 2️⃣ | [🌿 Swamp Pattern](docs/thinking-in-hydraide/swamp-pattern.md) | Configure persistence, memory, and lifespan directly from code. |
| 3️⃣ | [💎 Treasures](docs/thinking-in-hydraide/treasures.md) | Understand the smallest, most powerful unit of data. |
| 4️⃣ | [🧩 Indexing](docs/thinking-in-hydraide/indexing.md) | Discover ephemeral, in-memory indexing that feels like magic. |
| 5️⃣ | [🔄 Subscriptions](docs/thinking-in-hydraide/subscriptions.md) | Build reactive systems natively with HydrAIDE’s event engine. |
| 6️⃣ | [🔐 Locking](docs/thinking-in-hydraide/locking.md) | Achieve true concurrency without fear. |
| 7️⃣ | [🧹 Clean System](docs/thinking-in-hydraide/clean-system.md) | Never think about cleanup again – because HydrAIDE already did. |
| 8️⃣ | [🌐 Distributed Architecture](docs/thinking-in-hydraide/distributed-architecture.md) | Scale horizontally without orchestration pain. |
| 9️⃣ | [🚀 Install & Update](docs/thinking-in-hydraide/how-to-install-update-hydraide.md) | Deploy HydrAIDE in minutes, not days. |

---

## 🔥 Why Developers Choose HydrAIDE

| Feature | What It Means                                                     |
|--------|-------------------------------------------------------------------|
| ⚡ Instant Swamp access | O(1) folder-mapped resolution, no indexing required               |
| 🧠 On-the-fly indexing | Built in RAM only when needed, disappears after use               |
| 🔄 Subscriptions | Real-time streams of change, built into every write/update/delete |
| 🔐 Zero-deadlock locking | Per-Treasure locks and TTL-protected business locks               |
| 🧹 Auto-cleanup | Swamps and memory vanish when no longer needed                    |
| 🌐 Horizontal scaling | Stateless by default, no orchestrator required                    |
| 💾 SSD-friendly design | Chunked binary storage with delta writes                          |
| 🧬 Fully typed | Store structs, slices, maps – native to your language             |

> You never define schemas. You never worry about cleanup. You just build.

---

## 🌀 Who HydrAIDE is For

HydrAIDE is ideal for:
- 🚀 Startups that need a modern database without infrastructure burden
- 🌐 High-volume web crawlers and analytics platforms (like Trendizz)
- 📊 Live dashboards, reactive pipelines, and streaming interfaces
- 🧪 Research and ML pipelines with ephemeral data flows
- 📱 IoT and edge devices with minimal memory

And especially for developers who:
- Want to **own their data logic in code**
- Need **real-time data flow without middleware**
- Hate daemons, cron jobs, and config bloat

> HydrAIDE is **developer-native**. You don’t configure it. You *declare intent*, and it responds.

## 📊 Comparisons – HydrAIDE vs Other Systems

Want to see how HydrAIDE compares to the most popular databases and engines?  
We’re building a full series of deep comparisons — mindset-first, not config-first.

| 🔍 Comparison                | 📄 Status           | Link                                                                            |
|-----------------------------|---------------------|---------------------------------------------------------------------------------|
| HydrAIDE vs MongoDB            | ✅ Complete         | [Read HydrAIDE vs MongoDB Comparison](/docs/comparisons/hydraide-vs-mongodb.md) |
| HydrAIDE vs Redis              | ✅ Complete         | [Read HydrAIDE vs Redis Comparison](/docs/comparisons/hydraide-vs-redis.md)     |
| HydrAIDE vs PostgreSQL         | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs MySQL              | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs SQLite             | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs Elasticsearch      | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs Firebase / Firestore | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs DynamoDB           | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs Cassandra          | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs InfluxDB           | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs ClickHouse         | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs Neo4j              | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs TimescaleDB        | 🔜 In progress      | *coming soon*                                                                   |
| HydrAIDE vs Apache Kafka       | 🔜 In progress      | *coming soon* (stream vs native pub-sub)                                        |


---

## 🚀 Install & Run in Minutes

The HydrAIDE codebase is already available for contributors, and the easily deployable Docker container will be available soon as well. In the meantime, we kindly ask for your patience!

---

## 💻 SDKs – Native Integration in Your Language

HydrAIDE speaks **gRPC**, and every SDK is powered by the same `.proto` definition:

| 💻 SDK       | 🧪 Code Name                                                                         | 🛠️ Status           | 📘 Indexing Docs                       |
| ------------ |--------------------------------------------------------------------------------------| -------------------- | -------------------------------------- |
| 🟢 Go | [`hydraidego`](https://github.com/hydraide/hydraide/tree/main/docs/sdk/go/README.md) | 🔄 Actively in development | Coming soon – Realtime, blazing fast |
| 🟡 Node.js   | `hydraidejs`                                                                         | 🧪 In planning       | Coming soon – Event-friendly queries   |
| 🐍 Python    | `hydraidepy`                                                                         | 🧠 In design         | Coming soon – ML-ready sorting         |
| 🦀 Rust      | `hydraiders`                                                                         | 🧠 In design         | Coming soon – Zero-cost abstractions   |
| ☕ Java       | `hydraidejv`                                                                         | 🧠 In design         | Coming soon – Scalable enterprise use  |
| 🎯 C# / .NET | `hydraidecs`                                                                         | 🧠 In design         | Coming soon – Streamlined for services |
| 🧠 C++       | `hydraidecpp`                                                                        | 🧠 In design         | Coming soon – High-performance indexing|
| 🌀 Kotlin    | `hydraidekt`                                                                         | 🧠 In design         | Coming soon – Elegant & Android-ready  |
| 🍎 Swift     | `hydraidesw`                                                                         | 🧠 In design         | Coming soon – Index-aware mobile apps  |

---

## 🙌 Want to Contribute?

Start by reading the [Contributor Introduction](/CONTRIBUTORS.md) – it explains why HydrAIDE exists, what kind of people we’re looking for, and how you can join.
Then check out our [Contribution Guide](/CONTRIBUTING.md) – it walks you through the practical steps.

Once you're ready, open your first issue or pull request — we’ll be waiting! 🚀

---

## 📩 Contact & Enterprise

HydrAIDE is used in production at [Trendizz.com](https://trendizz.com). Interested in enterprise licensing, 
SDK development, or embedding HydrAIDE in your own platform?

📧 **Peter Gebri** – [peter.gebri@trendizz.com](mailto:peter.gebri@trendizz.com)
(Founder of HydrAIDE & Trendizz)
🌐 **Website** – [https://HydrAIDE.io ](https://hydraide.io) Currently in progress and directly linked to GitHub.

---

> 🧠 HydrAIDE isn’t a tool you use.
> It’s a system you think with.

Join the movement. Build different.

