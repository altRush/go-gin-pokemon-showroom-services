# Refactor Plan — Complete ✓

## Objective

Clean up the Go Gin service to remove code smells, improve maintainability, and separate concerns between HTTP handling, database access, and data mapping.

---

## ✅ Phase A — Fix PLANS.md

- Corrected stale file path references (`services/` → `models/`, etc.)
- Fixed function name reference (`ConvertDbArrayToUnnestArrayString` → actual)
- Added completion tracking (✅ / ⬜)
- Reorganized by status

## ✅ Phase B — Remaining Go code smells

| # | Change | Files |
|---|--------|-------|
| 1 | Use shared `models.db` in `GetPokemonByStoreIdFromStore` instead of duplicate `sql.Open` | `models/store.go` |
| 2 | Removed `panic()` inside HTTP handlers — errors now propagate to caller | `models/store.go` |
| 3 | Removed `log.Fatalln()` — deleted `models/helpers.go` entirely | `models/helpers.go` |
| 4 | Replaced `SELECT *` with explicit column lists | `models/store.go` |
| 5 | Replaced `convertDbArrayToUnnestArrayString` + fragile SQL concat with `pq.Array()` | `models/store.go` |
| 6 | Added `sql.ErrNoRows` check → returns 404 JSON | `main.go` |
| 7 | Renamed snake_case types to idiomatic Go: `PokemonProfileFromDB`, `PokemonProfileFromDBWithTypes`, `PokemonProfile`, `AddResult` | `models/store.go` |
| 8 | Renamed struct fields: `Type_name` → `TypeName` (Go convention; JSON tags unchanged) | `models/store.go` |
| 9 | Fixed ignored `BindJSON` error — now checked and returned | `models/store.go` |
| 10 | Added `return` after error responses in all handlers | `main.go` |
| 11 | Added `defer rows.Close()` right after `db.Query()` | `models/store.go` |
| 12 | Removed deprecated `convertDbArrayToUnnestArrayString` (replaced by `pq.Array`) | `models/helpers.go` (deleted) |

### Type renames applied

| Before | After |
|--------|-------|
| `PokemonTypes` | `PokemonType` |
| `Pokemon_profile_from_db` | `PokemonProfileFromDB` |
| `Pokemon_profile_from_db_with_types` | `PokemonProfileFromDBWithTypes` |
| `Pokemon_profile` | `PokemonProfile` |
| `Add_pokemon_to_store` | `AddResult` |
| `Type_name` field | `TypeName` |

## Files Modified

- `models/store.go` — type renames, `pq.Array`, explicit columns, shared DB, error propagation, `defer`
- `main.go` — `return` after errors, 404 handling for `sql.ErrNoRows`
- `models/helpers.go` — **deleted** (no longer needed)
- `PLANS.md` — this file
