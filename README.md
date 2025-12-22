# GABC Generator

Backend of this experimental tool for generating GABC code from liturgical texts with a regular structure, starting with the prefaces of the Mass.

[Structure tree](https://github.com/ramon-reichert/GABCgen/blob/main/documentation/tree.ini)

[APP Diagram here](https://excalidraw.com/#json=4qlX3mY09FOJoL5ChK0iY,bH2pjjbqSSevvrKv_BNOSg)

[Kanban of the project here](https://trello.com/b/l2gSItKa/gabcgen)

[Frontend page here](https://ramon-reichert.github.io/popechant/)


### Architectural Decisions

This project follows a **package-oriented design**, where packages are the primary unit of responsibility and reuse.  
The structure is designed to keep dependencies explicit and changes localized.

---

Packages are organized by the stability of the concepts they represent:

- Core, stable concepts live in internal packages with no external dependencies.
- Context-specific and infrastructure concerns live at the edges.
- Imports always point inward.

- `composition/phrases/words` is dependency-free.
- `preface` depends on `composition`, not the other way around.
- Infrastructure packages depend on interfaces defined closer to the domain.

---

Packages represent **domain concepts**, not technical layers.

- `words`, `phrases`, and `paragraph` model text structure.
- `staff` models musical notation.
- `preface` applies rules specific to the Preface.

Generic packages such as `models` or `utils` are intentionally avoided.

---

Text is built through a clear dependency chain:

`words → phrases → paragraph`

Each package builds on the one below it and imports only what it directly needs, making the data flow easy to follow through imports and constructors.

---

Core types are kept free of context-specific rules.

- `composition` defines physical text structures.
- `preface` applies behavior required by the Preface.

This allows new liturgical contexts to reuse the same core types without modification.

---

External integrations are isolated under `internal/platform`.

- Interfaces are defined where they are used.
- Infrastructure provides concrete implementations.
- Domain packages never import infrastructure packages.

This keeps the domain testable and independent from external services.

---

The `service` package coordinates application workflows:

- orchestrates calls between domain packages
- connects domain logic with infrastructure
- exposes entry points to the delivery layer

Low-level domain logic remains outside this package.

---

All wiring happens in `cmd/gabcgen`.

This is the only place where concrete implementations are selected and the application is assembled.

---

These decisions:

- makes dependencies visible through imports
- allows adding new liturgical parts without modifying existing packages
- keeps packages small and focused
- favors clarity over abstraction