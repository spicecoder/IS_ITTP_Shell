# Intention Space — Why it matters for developers

**Core idea (one line)**

Processes *do not* call each other directly — they emit **Intentions** to **Objects** which mediate, verify and reflect those Intentions only when design‑time synctests and runtime pulses permit.

**Why this matters**

- Insert human review, audit, or policy enforcement without touching existing services.  
- Make intent, payload and provenance machine‑verifiable (signed intentions + contentHash).  
- Implement gating, quarantine, or labels as protocol primitives rather than ad‑hoc code.

**Quick analogy**

Postal mail (Object mediator) vs direct phone call (process→process). The mailroom can inspect, hold, re-route or require a human stamp without changing sender or recipient.

**Minimal runtime sequence (developer view)**

1. Client computes `contentHash`, builds `intention` JSON, signs it.  
2. Client sends `{ intention, payload }` to Object proxy.  
3. Object verifies signature & hash; runs fast policy checks (classifier, denylist, reputation).  
4. If OK → Object reflects (publish). If flagged → create `quarantine` CPUX and emit `moderation_needed` pulse to a human DN.  
5. Human DN emits `moderation_decision` pulse → CPUX applies and Object reflects/rejects.  

**Practical dev checklist (3 steps you can implement today)**

1. Add `intent`, `contentHash`, `createdAt`, and a `signature` to outgoing requests in client SDK.  
2. Drop in an Object proxy (middleware) that verifies signature/hash, does quick checks and either reflects or quarantines.  
3. Persist CPUX traces (signed intention, verification result, classifier scores, moderator decision) for auditing and reproducible forensics.

**Demo & code pointers**

- Quick demo: `intention-proxy-demo` (Node.js) — proxy + tiny client SDK + moderator endpoints (run locally).  
- Reference quickstart shell: `/mnt/data/QUICKSTART_GOBASH.md`  

**Why devs will like it**

- Non‑invasive (wraps existing services).  
- Composable (policy as synctests/DNs, not scattered logic).  
- Testable and auditable (design‑time signal hashes → deterministic checks).

**One-line call to action**

Wrap one high‑risk route (e.g., `publish_post`) with an Object proxy this week — you’ll get human‑intervention, traceability, and improved safety without changing backend logic.

---

*Slide created for developer onboarding. Ask me to export this as PDF or produce a 3‑slide deck with diagrams.*

