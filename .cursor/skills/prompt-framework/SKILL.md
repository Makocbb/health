---
name: prompt-framework
description: >
  Coordinates project context for Go API work by loading identity, architecture,
  and code-style skills. Use whenever the user asks about code, requests a code
  change, asks for a new feature, needs debugging help, requests refactoring,
  asks for architecture guidance, requests testing, or wants documentation or
  implementation planning.
---

# Prompt Framework

Load project skills before answering or changing code. They define behavior,
layering, and implementation patterns for this Go API.

## Workflow

Copy and track:

```
Context:
- [ ] identity
- [ ] architecture
- [ ] code-style (if writing/editing code)
- [ ] code-style/examples.md (if new or full-domain feature)
```

1. **Read and follow** [identity](../identity/SKILL.md)
2. **Read and follow** [architecture](../architecture/SKILL.md)
3. **If writing or editing code**, **read and follow** [code-style](../code-style/SKILL.md)
4. **If adding a new domain or scaffolding all layers**, also read
   [code-style/examples.md](../code-style/examples.md)
5. Proceed with the request, staying consistent with the loaded skills

## What to load when

| Request type | Load |
|--------------|------|
| Questions, planning, architecture discussion | identity + architecture |
| Bugfix, small edit, refactor | identity + architecture + code-style |
| New domain / full CRUD feature | all of the above + examples.md |
| Tests or docs only | identity + architecture; code-style if examples are needed |

## Rules

- Finish the required reads before proposing or writing code
- Prefer architecture boundaries and code-style patterns over ad-hoc structure
- New domains get all layers (model → presenter → repository → provider → service → controller) unless the user scopes otherwise
- If a request conflicts with these skills, say so briefly and propose an approach that fits
