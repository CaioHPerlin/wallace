# Domain Docs

Engineering skills use these rules when exploring domain concepts and decisions in this repository.

## Before Domain Work

- Read the root `GLOSSARY.md` when it exists.
- Read relevant records in `docs/adr/` before changing an area with recorded architectural decisions.
- If these files do not exist, proceed without creating placeholders. Create them when domain terms or decisions are established.

## Layout

This is a single-context repository. Keep one `GLOSSARY.md` at the root and architectural decision records under `docs/adr/`.

## Vocabulary

Use terms as defined in the glossary. If a needed term is missing, record that as a domain-modeling gap instead of silently introducing competing terminology.

## ADR Conflicts

Surface conflicts with existing ADRs and explain why a decision may need to be revisited. Do not silently override a recorded decision.
