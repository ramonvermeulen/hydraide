# Naming in HydrAIDE – A Practical Guide

Good naming is more than just syntax in HydrAIDE. It’s how you structure your entire system. 
So before diving into code, it’s worth understanding how names shape your architecture.

This isn’t about schemas, tables or collections. HydrAIDE doesn’t work like that. Here, names define placement, access, and logic directly.

---

## The Basics

HydrAIDE uses a simple structure:

* **Sanctuary** - top-level purpose (e.g. `users`, `orders`)
* **Realm** – logical grouping inside the sanctuary (e.g. `profiles`, `drafts`)
* **Swamp** – the specific dataset (e.g. `john-doe`, `client-123`, `all-profiles`)

Each Swamp is a folder in your server. Each folder holds Treasures (your data). Access is O(1), directly from name → disk → memory.

No scan. No lookup. Just a direct jump.

---

## Think in Names

Traditional systems make you ask: *“How do I find this row?”*

HydrAIDE flips the question:

> *“What is the exact Swamp name for this data?”*

Once you know the name, everything becomes predictable. There’s no magic resolution step. You control the structure just by naming it right.

---

## Example: User Profiles

Instead of one big `users` table, you break it up like this:

```
users/profiles/john-doe
users/profiles/sarah-smith
```

Each Swamp is:

* Self-contained
* Instantly loadable
* Cleanly scoped

Need to load a profile? Just hydrate that Swamp. It’s on disk. It loads into RAM in milliseconds. It unloads when idle.

---

## Swamps Are Keyed Spaces

Inside a Swamp, you store Treasures. Key-value records.

Example:

* Swamp: `users/ids`
* Treasures:

    * `petergebri`
    * `sarahsmith`

That’s a presence list. No metadata. Just fast access.

Want to store something more complex like a wishlist?

```
users/wishlists/petergebri
```

The Swamp itself contains the wishlist items. Fully typed, binary stored.

---

## Naming Tips

* Keep Swamps small and purpose-driven.
* Avoid dumping different logic into one Swamp.
* Use plural for Sanctuary/Realm (`users`, `orders`, `logs`).
* Use stable, human-readable keys (`user-123`, `article-456`).

Each Swamp should answer one clear question. If it doesn’t, split it.

---

## Hydration = Activation

Swamps live on disk by default. But the moment you call one by name, HydrAIDE:

* Loads it into memory
* Makes it writable and subscribable
* Handles it like live data

This process is called **hydration**.

Swamps stay hydrated while in use. They unload automatically after inactivity (configurable from code). Once unloaded, they free up RAM — but data stays safe on disk.

This gives you massive scale, without memory bloat.

---

## Final Note

If you get naming right, everything else in HydrAIDE becomes easier:

* Reactive logic stays scoped
* Scaling is just folder distribution
* Memory stays lean
* No indexes or queries are needed

HydrAIDE isn’t just about storing things. It’s about **placing** them with intent.

Start there! And the rest will follow naturally.

---

## 🧭 Navigation

← [Back to Thinking in HydrAIDE](./thinking-in-hydraide.md) | [Next: 🌿 Swamp Pattern](./swamp-pattern.md)
