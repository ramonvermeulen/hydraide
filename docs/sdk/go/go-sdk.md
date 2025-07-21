# 🐹 HydrAIDE SDK – Go Edition

Welcome to the official **HydrAIDE SDK for Go**, your gateway to building intelligent,
distributed, real-time systems using the HydrAIDE engine.

This SDK provides programmatic access to HydrAIDE's powerful features such as swamp-based data structures,
lock-free operations, real-time subscriptions, and stateless routing, all tailored to Go developers.

---

## 📦 At a Glance

### 🧠 System

| Function  | SDK Status | Doc Status                | Docs |
| --------- | ------- |---------------------------|------|
| Heartbeat | ✅ Ready | ⏳ in progress | -    |

---

### 🔐 Business Logic

The functions under Business Logic enable **cross-cutting coordination** across distributed services.

These are not tied to a specific Swamp or Treasure — they operate on shared, logical domains like user balances,
order flows, or transaction states.

- `Lock()` acquires a **blocking distributed lock** for a given domain key to ensure no concurrent execution happens.
- `Unlock()` safely releases it using a provided lock ID.

Use these when you need to **serialize operations** across services or modules.

Ideal for:
- Credit transfers
- Order/payment sequences
- Ensuring consistency during critical logic

| Function | SDK Status | Doc Status            | Docs |
| -------- | ------- |-----------------------|------|
| Lock     | ✅ Ready | ⏳ in progress | -    |
| Unlock   | ✅ Ready | ⏳ in progress | -    |

---

### 🌿 Swamp & Treasure

These functions manage the lifecycle and existence of Swamps (data containers) and their Treasures (records),
including registration, validation, destruction, and real-time subscriptions.

| Function        | SDK Status | Doc Status    | Docs |
| --------------- | ---------- |---------------|-----|
| RegisterSwamp   | ✅ Ready | ⏳ in progress | -   |
| DeRegisterSwamp | ✅ Ready | ⏳ in progress | -   |
| IsSwampExist    | ✅ Ready | ⏳ in progress | -   |
| IsKeyExists     | ✅ Ready | ⏳ in progress | -   |
| Count           | ✅ Ready | ⏳ in progress | -   |
| Destroy         | ✅ Ready | ⏳ in progress | -   |
| Subscribe       | ✅ Ready | ⏳ in progress | -   |

---

### 📚 Catalog

Catalog functions are used when you want to store key-value-like entries where every item shares a similar structure,
like a list of users, logs, or events. Each Swamp acts like a collection of structured records,
e.g., user ID as the key and last login time as the value.

| Function                  | SDK Status | Doc Status                                                | Docs |
| ------------------------- | ------- |-----------------------------------------------------------|------|
| CatalogCreate             | ✅ Ready |⏳ in progress| -    |
| CatalogCreateMany         | ✅ Ready |⏳ in progress| -    |
| CatalogCreateManyToMany   | ✅ Ready |⏳ in progress| -    |
| CatalogRead               | ✅ Ready |⏳ in progress| -    |
| CatalogReadMany           | ✅ Ready |⏳ in progress|      |
| CatalogUpdate             | ✅ Ready |⏳ in progress| -    |
| CatalogUpdateMany         | ✅ Ready |⏳ in progress| -    |
| CatalogDelete             | ✅ Ready |⏳ in progress| -    |
| CatalogDeleteMany         | ✅ Ready |⏳ in progress| -    |
| CatalogDeleteManyFromMany | ✅ Ready |⏳ in progress| -    |
| CatalogSave               | ✅ Ready |⏳ in progress| -    |
| CatalogSaveMany           | ✅ Ready |⏳ in progress| -    |
| CatalogSaveManyToMany     | ✅ Ready |⏳ in progress| -    |

---

### 🧬 Profile

Profile Swamps are designed for storing heterogeneous key-value pairs where each key maps to a different type,
typically representing an entire user profile. Ideal when you need to manage multiple fields (e.g., name, avatar,
preferences) under one logical entity.

| Function    | SDK Status | Doc Status     | Docs |
| ----------- | ------- |----------------|------|
| ProfileSave | ✅ Ready | ⏳ in progress  | -    |
| ProfileRead | ✅ Ready | ⏳ in progress  | -    |

---

### ➕ Increments / Decrements

These functions allow atomic, strongly-typed modifications of numeric fields, optionally guarded by conditions,
ideal for updating counters, scores, balances, or state values in a safe and concurrent environment.

| Function         | SDK Status | Doc Status                                            | Docs |
| ---------------- | ------- |-------------------------------------------------------|------|
| IncrementInt8    | ✅ Ready | ⏳ in progress | -    |
| IncrementInt16   | ✅ Ready | ⏳ in progress | -    |
| IncrementInt32   | ✅ Ready | ⏳ in progress | -    |
| IncrementInt64   | ✅ Ready | ⏳ in progress | -    |
| IncrementUint8   | ✅ Ready | ⏳ in progress | -    |
| IncrementUint16  | ✅ Ready | ⏳ in progress | -    |
| IncrementUint32  | ✅ Ready | ⏳ in progress | -    |
| IncrementUint64  | ✅ Ready | ⏳ in progress | -    |
| IncrementFloat32 | ✅ Ready | ⏳ in progress | -    |
| IncrementFloat64 | ✅ Ready | ⏳ in progress | -    |

---

### 📌 Slice & Reverse Proxy

These are specialized functions for managing `uint32` slices in an atomic and deduplicated way — mainly
used as **reverse index proxies** within Swamps. Perfect for scenarios like tag mapping, reverse lookups,
and set-style relationships.

| Function                | SDK Status | Doc Status                                            | Docs |
| ----------------------- | ------- |-------------------------------------------------------|------|
| Uint32SlicePush         | ✅ Ready | ⏳ in progress | -    |
| Uint32SliceDelete       | ✅ Ready | ⏳ in progress | -    |
| Uint32SliceSize         | ✅ Ready | ⏳ in progress | -    |
| Uint32SliceIsValueExist | ✅ Ready | ⏳ in progress | -    |

Each of these functions will be documented in detail, explaining how they work and how to use them in real-world Go applications.
