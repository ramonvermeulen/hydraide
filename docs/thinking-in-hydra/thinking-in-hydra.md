# 🧠 Thinking in Hydra

Welcome to **Hydra** – a fundamentally different way of thinking about data, structure, and real-time access. This is not your typical key-value store, nor a document database, nor a traditional graph system. It’s all of them – and none of them.

Hydra is also **not an SQL or NoSQL database**. It stores and retrieves data in entirely different ways – optimized for modern memory architectures, distributed workloads, and event-driven flows. Classic databases – whether relational or document-based – simply couldn’t keep up with what we needed.

---

## 🚧 Why Hydra Was Born

Hydra was created out of necessity. During the development of **Trendizz.com**, our platform that indexes **the entire content of all European websites** and allows **exact word-level search within 1-2 seconds**, we hit a wall.

No existing database could support this level of scale, speed, and efficiency – at least not in a cost-effective or sustainable way.

- Traditional databases were too **slow**,
- NoSQL stores were **space-inefficient**,
- In-memory engines consumed **too much RAM**,
- And none of them could guarantee **O(1) access** regardless of dataset size.

We realized the core problem: **these databases were designed for 30-40-year-old architectures**, and they weren’t built for the kind of modern, dynamic, multi-terabyte real-time data flow we envisioned.

So we made a radical decision: to **start from scratch**, and build a new engine from the ground up. One that’s:

- Ultra-fast ⚡
- Developer-friendly 🧑‍💻
- Easy to learn and reason about 🧠
- All-in-one, self-hostable, modular 💡
- Designed for real-time workloads, automation, and scalability 🌍
- And critically: **every search must run in O(1)** – constant time, no matter how large the dataset grows.

Thus, Hydra was born.

After nearly 3 years of development, we reached **Hydra 2.0**, a production-grade system serving **billions of records** with consistency, speed, and elegance.

But to harness its full power, you must first change the way you think about data.



🔍 Did you know? The entire Trendizz.com dashboard – built in Angular – relies 100% on Hydra's subscription system to power real-time search experiences. The UI is fully reactive without needing any third-party messaging tools or data sync layers. Hydra makes that possible, natively.

---

## 🔍 Why Hydra is Different

Hydra doesn’t model the world as rows, documents, or nodes. Instead, it models **relationships, states, and ephemeral truths** in a way that’s:

- **Structurally simple**, but
- **Semantically powerful**

It’s made for real-time systems, distributed architectures, and smart data pipelines. The design principles are:

- **Simplicity over abstraction**
- **Composability over configuration**
- **Statefulness with flexibility**

---

## 🧱 Core Philosophy

Hydra is designed with a few *radical* ideas – and here's what makes each of them powerful:

1. 🧩 **Slices over schemas** – You don’t need a strict schema to store data. You define simple structs in your code to represent the data in the way that best serves your goals. You don’t need to “tweak a database” or manage indexes. You control everything from your code – even how long a swamp should stay in memory when unused, or how big the underlying files are. Hydra has no CLI, no admin dashboard, no overhead. It’s all code – simple, developer-native. This is our guiding principle. *Code your data.* You don’t need a database specialist to work with Hydra – just a developer. The data engine lives and breathes under your fingertips.

2. 📡 **Pub-sub as a first-class citizen** – Everything in Hydra can be subscribed to. Any data change can trigger events, push updates, or drive logic in real-time. This isn’t an addon – it’s built into the core.

3. ⏱️ **TTL as a feature, not a workaround** – You can define TTL (time-to-live) for any data item. Hydra will automatically handle its expiration. Some queries can return expired data, others will remove it while returning. TTL becomes a native part of your logic.

4. 🧭 **No centralized index – yet always findable** – Hydra’s indexing is adaptive and implicit. You don’t have to declare indexes, tune them, or manage storage tradeoffs. Everything remains magically searchable.

5. 🧹 **Garbage collection is built-in** – You don’t manage cleanup – Hydra does. Nothing stays in memory that isn’t used. No files remain on disk or SSDs that are empty or redundant. You focus on logic – Hydra takes care of the rest.

This is how **modern distributed memory** should behave.

🧠 In line with 2025 developer experience standards, Hydra’s SDKs are intentionally designed to be understandable not only by developers, but also by AI copilots like ChatGPT. This means if you're working on complex data structures, Hydra makes it easy for your AI assistant to generate or suggest accurate, reliable code. It’s not just developer-friendly – it’s AI-friendly too.

---

## 🧭 The Hydra Docs Navigation

We designed this documentation to be **fast to learn** and **immediately useful**. If you're a developer, you should be able to understand and use Hydra confidently within a single day.

🛠️ In this documentation, you won’t find language-specific SDK code examples. Instead, you'll gain a **solid understanding of how Hydra works**, why it works that way, and how to **think like a Hydra developer**.

This is your starting point. Please follow the sections **in order**, from top to bottom. Each section builds on the previous one.

Once you understand these principles, you’ll be ready to dive into SDKs and build rapidly, with clarity and confidence.

Before you jump into implementation, explore the foundations:

| Section                                                               | Description                                                                                                                                 |
|-----------------------------------------------------------------------| ------------------------------------------------------------------------------------------------------------------------------------------- |
| **1.** [📏 Naming Convention](./naming-convention.md)                 | This is where everything begins. Learn the Hydra naming convention, and how to design keys and slices in a human-readable and scalable way. |
| **2.** [📦 How to Store and Read Data](./how-to-store-and-read-data.md) | The first steps in writing and retrieving data inside Hydra.                                                                                |
| **3.** [🧭 Indexing](./indexing.md)                                     | Learn how Hydra auto-indexes data, and why it’s shockingly efficient.                                                                       |
| **4.** [🔄 Subscriptions](subscriptions.md)                           | Real-time pub-sub architecture, and how you can use it.                                                                                     |
| **5.** [⏳ Expirations](expirations.md)                                | Data lifecycle, auto-deletion and TTL-based logic.                                                                                          |
| **6.** [🔐 Locking](locking.md)                                       | Safe concurrent access with Hydra's locking mechanisms.                                                                                     |
| **7.** [🧹 Built-in Garbage Collector](built-in-garbage-collector.md) | Internal magic: cleaning up data you no longer need.                                                                                        |
| **8.** [🔁 Reverse Index Slice](reverse-index-slice.md)               | A special tool to reverse-lookup keys or create inverted views.                                                                             |
| **9.** [🌐 Distributed Architecture](distributed-architecture.md)     | Horizontal scaling, replication, sharding, and network design.                                                                              |
| **10.** [⚙️ Optimization](optimization.md)                            | Performance tuning and best practices for huge datasets.                                                                                    |
| **11.** [➕ Conditional Increment](conditional-increment.md)           | Logic-driven data mutation – increment with guardrails.                                                                                     |
| **12.** [🔍 Real World Examples](real-world-examples.md)              | See how others use Hydra in AI, search engines, and more.                                                                                   |
| **13.** [🚀 Install & Update Hydra](how-to-install-update-hydra.md)   | Getting started: how to set up and maintain your Hydra deployment.                                                                          |

---

## 🌱 A New Way of Thinking

Hydra is not just a tool. It’s a **philosophy** – a new mental model for building modern, ultra-fast, event-driven systems that feel alive.

- You don’t model the world – you model **change**.
- You don’t design structure – you design **flow**.
- You don’t fight with scale – you **embrace modularity**.

Hydra is your bridge to a new era of backend design – one where real-time, distributed memory is natural and intuitive, not complex and bloated.

🔁 This section is your mindset shift. Once you get this, everything else clicks.

📍 **Important:** Please begin your journey with the [Naming Convention](naming-convention.md). That’s where the Hydra story truly starts – and where your journey as a Hydra developer begins.

Once you understand the foundation, the rest of your development experience will be **fast**, **fluid**, and **fun**.

💬 If something is unclear, or if you get stuck along the way, don’t hesitate to open a documentation issue on [GitHub](https://github.com/hydraide/hydraide). We’re actively building a friendly and helpful community around Hydra, and your questions help us make the system better for everyone.

