# 🧭 Indexing – The Schrödinger indexing

Welcome back to Hydra – where the rules of data are rewritten.

You’ve already seen how naming, structuring, and configuring Swamps in Hydra feels more like **intention** than instruction. Now we venture into something deceptively simple – yet utterly transformative:

> **Indexing.**

But not as you know it.

Hydra doesn’t index like traditional databases.
There are no B-trees to build, no indexes to maintain, no disk writes to schedule.

And yet... it’s blazingly fast. Almost *magically* so.

How?

Let’s dive in. 🧠💥

---

## 🚀 The O(1) Illusion – Swamp Access is Instant

Let’s start from the foundation:

Every Swamp in Hydra is physically stored as a **folder** on disk.
Because of this, **locating a Swamp happens in O(1)** time, regardless of how many Swamps exist in the system.

Thousands? Millions? Billions?
Still just: `Name → Folder → Load`.

This is your first index – **a filesystem-native instant lookup**, with no overhead.
You never wait for Swamp discovery. It just *is*.

But what about the data *inside* a Swamp?
That’s where things get wild. 🔥

---

## 🧠 Cold Indexing – The Hydra Way

Most databases build indexes ahead of time. On disk. With lots of I/O. Always updating. Always syncing. Always consuming.

Hydra flips the script:

> **Indexing only happens in memory, only when needed, and disappears when done.**

No disk space consumed.
No background processes running.
No syncing nightmares.

It’s **ephemeral. Real-time. And hyper-efficient.**

When you open a Swamp and begin querying it, the system checks:

- Are you using the key directly? Then no index is needed.
- Are you sorting by value? Filtering by `uint64` or `string`?
  → Boom. Index created – *on-the-fly*, in memory, in nanoseconds.

These indexes don’t live on your SSD.
They don’t persist between requests.
They are **alive only while you need them**.

Once the Swamp is closed, the index is gone.
Zero footprint.

And here’s the twist:

> Hydra’s instant indexes are **not** B-trees. They’re powered by native **Go map structures** – meaning the underlying mechanism is a blazing-fast, hash-based in-memory index.

So when Hydra indexes, it doesn’t scan, iterate or optimize. It **maps and retrieves** in O(1) time.

That’s Hydra.

---

## 🎯 Default Index: The Key

Every Swamp stores data as key-value pairs.
And by default, Hydra can retrieve these keys in **sorted order**, without any explicit configuration.

But here’s the kicker:
Even this sorted-key access isn’t prebuilt.
It’s generated **only if you ask for it** during hydration.

Query directly by key? Hydra fetches it with no indexing.
Ask for a list sorted by key? Hydra builds a memory index – and serves it instantly.

That’s **zero waste**. Maximum intent. Minimal RAM.

---

## 📊 Value Indexing – On Demand, On Fire

Beyond keys, Hydra can index virtually any primitive value type:

- `int8` → `uint64`
- `float32`, `float64`
- `string`
- Even booleans.

If it’s inside a Treasure, Hydra can sort and filter by it – **but only when you ask.**

For example:
> *"Give me the first 50 items, sorted by an `int32` value in descending order."*

No pre-built schema. No stored index.
Just a single query – and the Hydra engine:
- Builds a hash-based in-memory index using Go maps,
- Sorts the data immediately,
- Returns the result.

The index is fast because it's **not iterated**, it’s **mapped**.
And once you're done?
It’s gone.

### 🧾 What About Metadata?

Hydra treats metadata fields – like `createdAt`, `modifiedAt`, `deletedAt`, and `expiredAt` – as **first-class indexable citizens**.

This means you can:
- Sort records by creation time,
- Filter by expiration,
- Query recently updated Treasures,
- Or even fetch all logically deleted entries.

And yes – these metadata indexes are just as fast, temporary, and Go map-powered as everything else.

You get full control over temporal logic, without ever defining a single schema or maintaining any index structure manually.

It's native. It's ephemeral. It's instant.

Just ask – and Hydra delivers.

No cleanup needed. No stale indexes. No bloat.

---

## 🧬 Memory-Efficient by Design

This is where Hydra’s real beauty emerges:

Because Swamps are **intentionally small**, each index is small too.

Indexing 10,000 records? Milliseconds.
Indexing 1,000,000? Still fast.

No rebuild times. No cascading updates. No writes to disk.
Just **temporary structures** optimized for real-time use.

In fact:
- If you only use direct key access, **no index is ever built**.
- If you sort or filter, **an index is built live**, scoped only to that session.

It’s like having the speed of precomputed indexes,
> Without ever building them.

Mind. Blown. 💥

---

## 🌀 Rehydration and Index Persistence

Let’s say you open a Swamp and sort by `createdAt`.
Hydra builds the index. You make your query.
Now, you keep the Swamp open.

What happens next time?

Same query → same index → instant response.

As long as the Swamp remains open in memory,
the index stays available.

But when the Swamp is unloaded –
it’s all wiped clean.

No manual flushes. No corruption risk. No stale state.

Hydra resets to a clean slate every time.

---

## 🔒 Zero Admin. Zero Lock-In.

Traditional systems often require you to:
- Define indexes in advance
- Sync them across clusters
- Manage their lifecycle manually

Hydra’s answer:

> Don't define indexes.
> Don’t think about them at all.
> Just query – and let the system **respond intelligently**.

This isn’t just developer-friendly.
It’s **developer-liberating**.

You don’t manage indexes.
You don’t maintain indexes.
You don’t even need to *know* they exist.

But they’re always there – when you need them.

---

## 🧩 Summary – Indexing, Evolved

Hydra gives you:

- 🔍 Instant Swamp discovery in O(1)
- ⚡ Real-time, in-memory indexing
- 🎯 Sorted access by key or value, and metadata
- 🧠 Indexes that only exist when useful
- 💡 No disk writes, no admin tasks
- 🧹 Automatic cleanup on Swamp close
- 🧭 Hash-based in-memory indexing powered by Go maps

No config.
No overhead.
No wasted cycles.

Just **pure adaptive speed**, wrapped in the elegance of Hydra's design.

So the next time someone asks:
> *“How does Hydra handle indexing?”*

You can smile and say:

> *It doesn’t. Until it does. Then it’s faster than you can imagine.*

---

## 🔗 SDK Integration Resources (Coming Soon)

Once you master the indexing philosophy, using it in code becomes effortless.

Each Hydra SDK will support value-based indexing with lazy evaluation and adaptive performance – just like the philosophy you’ve learned here.

| 💻 SDK       | 🧪 Code Name | 🛠️ Status           | 📘 Indexing Docs                       |
| ------------ | ------------ | -------------------- | -------------------------------------- |
| 🟢 Go        | `hydrungo`   | ✅ Actively developed | Coming soon – Realtime, blazing fast   |
| 🟡 Node.js   | `hydrunjs`   | 🧪 In planning       | Coming soon – Event-friendly queries   |
| 🐍 Python    | `hydrunpy`   | 🧠 In design         | Coming soon – ML-ready sorting         |
| 🦀 Rust      | `hydrunrs`   | 🧠 In design         | Coming soon – Zero-cost abstractions   |
| ☕ Java       | `hydrunjv`   | 🧠 In design         | Coming soon – Scalable enterprise use  |
| 🎯 C# / .NET | `hydruncs`   | 🧠 In design         | Coming soon – Streamlined for services |
| 🧠 C++       | `hydruncpp`  | 🧠 In design         | Coming soon – High-performance indexing|
| 🌀 Kotlin    | `hydrunkt`   | 🧠 In design         | Coming soon – Elegant & Android-ready  |
| 🍎 Swift     | `hydrunsw`   | 🧠 In design         | Coming soon – Index-aware mobile apps  |

> 💬 Still wondering how this works in code? Stay tuned.
> The SDKs will bring this philosophy to life – but **you already understand the magic behind it.**

---

## 🧭 Navigation

← [Back to Treasures](./treasures.md) | [Next: Subscriptions →](./subscriptions.md)

