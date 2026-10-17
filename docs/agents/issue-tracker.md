# Issue Tracker: Local Markdown

Issues and specs for this repository live as Markdown files under `.scratch/`.

## Conventions

- Create one directory per feature: `.scratch/<feature-slug>/`.
- Store the feature spec at `.scratch/<feature-slug>/spec.md`.
- Store each implementation issue in `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numbered from `01`.
- Put a `Status:` line near the top of each issue and use the role names from `triage-labels.md`.
- Append comments and conversation history under a `## Comments` heading.

## Skill Operations

When a skill says to publish work to the issue tracker, create a Markdown file in the appropriate `.scratch/<feature-slug>/` directory.

When a skill asks to fetch a ticket, read the referenced file. The user will normally provide its path or issue number.

## Wayfinding

- Store an effort map at `.scratch/<effort>/map.md`.
- Store each child ticket at `.scratch/<effort>/issues/NN-<slug>.md`.
- Add a `Type:` line (`research`, `prototype`, `grilling`, or `task`) and a `Status:` line (`claimed` or `resolved`) to each child ticket.
- Use `Blocked by: NN, NN` to record dependencies. A ticket is unblocked when every listed ticket is resolved.
- Claim work by setting `Status: claimed` before starting. Resolve it by adding an `## Answer`, setting `Status: resolved`, and updating the map's decisions.
