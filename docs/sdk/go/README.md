# 🐹 HydrAIDE SDK – Go Edition

Welcome to the official **HydrAIDE SDK for Go**, your gateway to building intelligent, distributed, real-time systems using the HydrAIDE engine.

This SDK provides programmatic access to HydrAIDE's powerful features such as swamp-based data structures, lock-free operations, real-time subscriptions, and stateless routing — all tailored to Go developers.

> **⚠️ Status:** Currently under active development.&#x20;

---

## 📦 At a Glance

### 🧠 System

| Function  | SDK Status | Doc Status                | Docs                      |
| --------- | ------- |---------------------------|---------------------------|
| Heartbeat | ✅ Ready | ⏳ in progress | [Heartbeat](Heartbeat.md) |

### 🔐 Business Logic

The functions under Business Logic enable **cross-cutting coordination** across distributed services.

These are not tied to a specific Swamp or Treasure — they operate on shared, logical domains like user balances, order flows, or transaction states.

- `Lock()` acquires a **blocking distributed lock** for a given domain key to ensure no concurrent execution happens.
- `Unlock()` safely releases it using a provided lock ID.

Use these when you need to **serialize operations** across services or modules.

Ideal for:
- Credit transfers
- Order/payment sequences
- Ensuring consistency during critical logic



| Function | SDK Status | Doc Status            | Docs                  |
| -------- | ------- |-----------------------|-----------------------|
| Lock     | ✅ Ready | ⏳ in progress | [Lock](./Lock.md)     |
| Unlock   | ✅ Ready | ⏳ in progress | [Unlock](./Unlock.md) |

### 🌿 Swamp & Treasure

These functions manage the lifecycle and existence of Swamps (data containers) and their Treasures (records) — including registration, validation, destruction, and real-time subscriptions.

| Function        | SDK Status | Doc Status    | Docs                                  |
| --------------- | ---------- |---------------|---------------------------------------|
| RegisterSwamp   | ✅ Ready | ⏳ in progress | [RegisterSwamp](./RegisterSwamp.md)   |
| DeRegisterSwamp | ✅ Ready | ⏳ in progress | [DeRegisterSwamp](./DeRegisterSwamp.md) |
| IsSwampExist    | ✅ Ready | ⏳ in progress | [IsSwampExist](./IsSwampExist.md)       |
| IsKeyExists     | ✅ Ready | ⏳ in progress | [IsKeyExists](./IsKeyExists.md)         |
| Count           | ✅ Ready | ⏳ in progress | [Count](./Count.md)                     |
| Destroy         | ✅ Ready | ⏳ in progress | [Destroy](./Destroy.md)                 |
| Subscribe       | ✅ Ready | ⏳ in progress | [Subscribe](./Subscribe.md)             |

### 📚 Catalog

Catalog functions are used when you want to store key-value-like entries where every item shares a similar structure — like a list of users, logs, or events. Each Swamp acts like a collection of structured records, e.g., user ID as the key and last login time as the value.

| Function                  | SDK Status | Doc Status                                                | Docs                                                      |
| ------------------------- | ------- |-----------------------------------------------------------|-----------------------------------------------------------|
| CatalogCreate             | ✅ Ready |⏳ in progress| [CatalogCreate](./CatalogCreate.md)                         |
| CatalogCreateMany         | ✅ Ready |⏳ in progress| [CatalogCreateMany](./CatalogCreateMany.md)                 |
| CatalogCreateManyToMany   | ✅ Ready |⏳ in progress| [CatalogCreateManyToMany](./CatalogCreateManyToMany.md)     |
| CatalogRead               | ✅ Ready |⏳ in progress| [CatalogRead](./CatalogRead.md)                             |
| CatalogReadMany           | ✅ Ready |⏳ in progress| [CatalogReadMany](./CatalogReadMany.md)                     |
| CatalogUpdate             | ✅ Ready |⏳ in progress| [CatalogUpdate](./CatalogUpdate.md)                         |
| CatalogUpdateMany         | ✅ Ready |⏳ in progress| [CatalogUpdateMany](./CatalogUpdateMany.md)                 |
| CatalogDelete             | ✅ Ready |⏳ in progress| [CatalogDelete](./CatalogDelete.md)                         |
| CatalogDeleteMany         | ✅ Ready |⏳ in progress| [CatalogDeleteMany](./CatalogDeleteMany.md)                 |
| CatalogDeleteManyFromMany | ✅ Ready |⏳ in progress| [CatalogDeleteManyFromMany](./CatalogDeleteManyFromMany.md) |
| CatalogSave               | ✅ Ready |⏳ in progress| [CatalogSave](./CatalogSave.md)                             |
| CatalogSaveMany           | ✅ Ready |⏳ in progress| [CatalogSaveMany](./CatalogSaveMany.md)                     |
| CatalogSaveManyToMany     | ✅ Ready |⏳ in progress| [CatalogSaveManyToMany](./CatalogSaveManyToMany.md)         |

### 🧬 Profile

Profile Swamps are designed for storing heterogeneous key-value pairs where each key maps to a different type — typically representing an entire user profile. Ideal when you need to manage multiple fields (e.g., name, avatar, preferences) under one logical entity.

| Function    | SDK Status | Doc Status     | Docs                          |
| ----------- | ------- |----------------|-------------------------------|
| ProfileSave | ✅ Ready | ⏳ in progress  | [ProfileSave](./ProfileSave.md) |
| ProfileRead | ✅ Ready | ⏳ in progress  | [ProfileRead](./ProfileRead.md) |

### ➕ Increments / Decrements

These functions allow atomic, strongly-typed modifications of numeric fields, optionally guarded by conditions — ideal for updating counters, scores, balances, or state values in a safe and concurrent environment.

| Function         | SDK Status | Doc Status                                            | Docs                                    |
| ---------------- | ------- |-------------------------------------------------------|-----------------------------------------|
| IncrementInt8    | ✅ Ready | ⏳ in progress | [IncrementInt8](./IncrementInt8.md)       |
| IncrementInt16   | ✅ Ready | ⏳ in progress | [IncrementInt16](./IncrementInt16.md)     |
| IncrementInt32   | ✅ Ready | ⏳ in progress | [IncrementInt32](./IncrementInt32.md)     |
| IncrementInt64   | ✅ Ready | ⏳ in progress | [IncrementInt64](./IncrementInt64.md)     |
| IncrementUint8   | ✅ Ready | ⏳ in progress | [IncrementUint8](./IncrementUint8.md)     |
| IncrementUint16  | ✅ Ready | ⏳ in progress | [IncrementUint16](./IncrementUint16.md)   |
| IncrementUint32  | ✅ Ready | ⏳ in progress | [IncrementUint32](./IncrementUint32.md)   |
| IncrementUint64  | ✅ Ready | ⏳ in progress | [IncrementUint64](./IncrementUint64.md)   |
| IncrementFloat32 | ✅ Ready | ⏳ in progress | [IncrementFloat32](./IncrementFloat32.md) |
| IncrementFloat64 | ✅ Ready | ⏳ in progress | [IncrementFloat64](./IncrementFloat64.md) |

### 📌 Slice & Reverse Proxy

These are specialized functions for managing `uint32` slices in an atomic and deduplicated way — mainly used as **reverse index proxies** within Swamps. Perfect for scenarios like tag mapping, reverse lookups, and set-style relationships.

| Function                | SDK Status | Doc Status                                            | Docs                                                  |
| ----------------------- | ------- |-------------------------------------------------------|-------------------------------------------------------|
| Uint32SlicePush         | ✅ Ready | ⏳ in progress | [Uint32SlicePush](./Uint32SlicePush.md)                 |
| Uint32SliceDelete       | ✅ Ready | ⏳ in progress | [Uint32SliceDelete](./Uint32SliceDelete.md)             |
| Uint32SliceSize         | ✅ Ready | ⏳ in progress | [Uint32SliceSize](./Uint32SliceSize.md)                 |
| Uint32SliceIsValueExist | ✅ Ready | ⏳ in progress | [Uint32SliceIsValueExist](./Uint32SliceIsValueExist.md) |

Each of these functions will be documented in detail, explaining how they work and how to use them in real-world Go applications.

---

## 🤝 Contribute to HydrAIDE

HydrAIDE is not just a database – it's a new paradigm.

If you'd like to help build the official SDKs and developer tools around the HydrAIDE core engine, check out our contributor program:

👉 [View the full Contributor Guide →](/CONTRIBUTORS.md)

> Join HydrAIDE. Be legendary.

---

## 📄 License Notice

This document is part of the HydrAIDE knowledge base and is licensed under a **custom restrictive license**.\
You may not use its contents to build or assist in building alternative engines, architectures, or competing systems.\
See full terms: [LICENSE.md](/LICENSE.md)



