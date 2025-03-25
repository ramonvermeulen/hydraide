# 🐹 HydrAIDE SDK – Go Edition

Welcome to the official **HydrAIDE SDK for Go**, your gateway to building intelligent, distributed, real-time systems using the HydrAIDE engine.

This SDK provides programmatic access to HydrAIDE's powerful features such as swamp-based data structures, lock-free operations, real-time subscriptions, and stateless routing — all tailored to Go developers.

> **⚠️ Status:** Currently under active development. Early components (like `name`, `client`, `Hydrun` interface) are available and evolving rapidly.

---

## 📦 At a Glance

| Feature                                     | Status         | Docs                                 |
| ------------------------------------------- | -------------- |--------------------------------------|
| `name.New().Sanctuary().Realm().Swamp()`    | ✅ Available    | [name.md](name.md)                   |
| gRPC connection & routing                   | 🔄 In Progress | [client.md](client.md)               |
| `RegisterSwamp()`                           | 🔄 In Progress | [registerswamp.md](registerswamp.md) |
| `Create()`, `Save()`                        | 🔄 In Progress | [create-save.md](create-save.md)     |
| `CreateMany()`, `SaveMany()`                | 🔄 In Progress | [create-save.md](create-save.md)     |
| `Read()`, `ReadMany()`                      | 🔄 In Progress | [read.md](read.md)                   |
| `Update()`, `UpdateMany()`                  | 🔄 In Progress | [update.md](update.md)               |
| `Delete()`, `DeleteMany()`                  | 🔄 In Progress | [delete.md](delete.md)               |
| `Destroy()`                                 | 🔄 In Progress | [delete.md](destroy.md)              |
| `Subscribe()`                               | 🔄 In Progress | [subscriptions.md](subscriptions.md) |
| `Lock()`, `Unlock()`                        | 🔄 In Progress | [locking.md](locking.md)             |
| `Count()`, `IsSwampExist()`, `IsKeyExist()` | 🔄 In Progress | [existence.md](existence.md)         |
| `IncrementInt()`, `DecrementFloat()`...     | 🔄 In Progress | [increment.md](increment.md)         |

> ✅ = Implemented • 🔄 = Work in progress • 🕓 = Planned

You can help shape the SDK — see [Contribute](#-contribute-to-hydraide)

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

---

