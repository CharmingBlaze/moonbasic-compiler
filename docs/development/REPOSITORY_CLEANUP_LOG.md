# Repository cleanup log

Record of moves, renames, and deletions performed during structural cleanup.  
Every deletion must include evidence it was unused.

## Template

```text
Date:
Commit:
Change type: move | rename | delete | generate
Path (from → to, or deleted path):
Reason:
Evidence (for deletes):
Replacement:
Tests performed:
```

---

## Entries

### 2026-07-25 — audit scaffolding

```text
Date: 2026-07-25
Commit: (this series)
Change type: generate
Path: docs/development/REPOSITORY_CLEANUP_AUDIT.md
         docs/development/REPOSITORY_CLEANUP_LOG.md
         docs/development/CODEBASE_MAP.md
Reason: Phase 1 inventory before any filesystem moves
Evidence: n/a
Replacement: n/a
Tests performed: docs-only commit
```

*(Subsequent entries appended as commits land.)*
