---
title: Bottleneck Hunt
use_when: Which measured cost should we investigate first?
source: shipped
added: 2026-11-18
---

## Borrow

- Keep call counts separate from elapsed time so a constant fan-out cannot hide rising work per call.
- Give bars a shared zero and keep the underlying runs visible.
- State what an end-to-end measurement must confirm before changing the system.
