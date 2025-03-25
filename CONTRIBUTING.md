# 🤝 Contributing to HydrAIDE

Welcome, builder! 🧠

Whether you’re here to squash a bug, suggest a feature, or shape a new SDK — you’re in the right place.

HydrAIDE is not just code. It’s a mindset.  
And we’re thrilled you’re thinking with us.

Your voice matters here. Whether it's your first PR or your 100th — you're helping shape the future of data systems.

---

## 🚀 Quickstart – How to Contribute

1. **Fork** this repo
2. **Create a new branch** for your fix or feature
3. **Make your changes**
4. **Open a Pull Request (PR)** – we’ll review it and celebrate with you 🎉

---

## 📂 Repository Structure

Each SDK lives in its own folder under `/docs/sdk/<language>`:

- `go` → actively developed
- `nodejs`, `python`, `rust`, etc. → in design or planning

> **Important:** SDK documentation must be 100% AI-readable.  
> Each file should be fully parseable by ChatGPT/Copilot — including clear function usage, struct layouts, and examples.

If you’re contributing to an SDK:
- Follow the structure and tone of `docs/sdk/go` as reference
- Keep all functions and types documented in Markdown with example code blocks
- Use a single `.md` file per SDK, but **clearly tagged and structured**

If you’re contributing:
- To core logic → edit [`hydraidego`](https://github.com/hydraide/hydraide/tree/main/docs/sdk/go/README.md)
- To docs → edit `.md` files in `/docs`
- To examples → add to `/examples/<your-language>`

---

## 💡 Looking for ideas?

Check the issues labeled [`good first issue`](https://github.com/hydraide/hydraide/issues?q=label%3A%22good+first+issue%22)

Not sure where to start?  
Browse the [Project Board](https://github.com/hydraide/hydraide/projects) or ask in an issue – we’ll help match you to something meaningful.

Or open a new one with your proposal!

---

## ✅ Commit Style

Use clear commit messages like:

- `fix: crash on empty Swamp hydration`
- `feat: add IncrementFloat64 to Node SDK`
- `docs: clarify metadata usage in Treasures`

---

## 🧪 Testing

Please make sure your changes:
- Run locally without errors
- Include tests (if logic-heavy)
- Don’t break other SDKs or docs

If you’re adding a new SDK function, include a simple usage test (e.g. call + assert result).  
Docs-only PRs don’t require tests.

---

## 🤲 Community Values

HydrAIDE is:
- 🧠 Inclusive — everyone starts somewhere.
- 🧼 Clean — clarity over cleverness.
- 🔄 Reactive — always listening, always improving.

If you’re kind, curious, and constructive — you belong here.

---

## 📥 Need Help?

Open an issue titled `Question: <your topic>`

---

## 👑 Want to Become a Core Contributor?

We welcome it! Start by:
- Opening a `Contributor Application` issue
- Telling us what excites you and what you'd love to build
- Shipping your first PR

We mentor. You grow.  
Together we build something legendary.

---

With gratitude,  
**– The HydrAIDE Team**

