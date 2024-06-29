# Roadmap / Planned Work

These items describe planned work. They are not implemented features in the
current project unless documented elsewhere.

## Runtime System

- Define the runtime responsibilities for building Anki packages from
  `.apkg.toml` files.
- Clarify how package generation should report validation errors, skipped
  inputs, and generated outputs.
- Keep the runtime focused on local package generation instead of presenting it
  as a hosted service or synchronization layer.

## Import From `.apkg`

- Investigate whether existing `.apkg` files can be converted back into a
  useful `.apkg.toml` representation.
- Document any Anki package metadata or media fields that cannot round-trip
  cleanly.
- Add tests with small sample packages before exposing this as a user-facing
  command.
