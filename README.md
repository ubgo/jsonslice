# jsonslice

> A `[]any` you can persist to a JSON/JSONB column **and** expose as a
> gqlgen scalar — without importing gqlgen, lib/pq, or pgx.

`JsonSlice` is a tiny type with three independent integrations packed in:

| Integration | Hooks | Cost to your dependency tree |
|-------------|-------|------------------------------|
| `encoding/json` | `MarshalJSON` / `UnmarshalJSON` | stdlib |
| `database/sql` | `Value()` / `Scan()` | stdlib |
| gqlgen scalar | `MarshalGQL` / `UnmarshalGQL` | none — detected by duck typing |

The third row is the point. gqlgen recognises any type with the right
`MarshalGQL(io.Writer)` and `UnmarshalGQL(any) error` method signatures,
so we expose the methods *without* importing `github.com/99designs/gqlgen`.
Consumers who don't use gqlgen pay nothing.

## Install

```sh
go get github.com/ubgo/jsonslice
```

## Use

```go
import "github.com/ubgo/jsonslice"

type Event struct {
    ID   string
    Tags jsonslice.JsonSlice  // ["urgent", "billing", 42]
}

// Save → JSON column
db.Exec(`INSERT INTO events(id, tags) VALUES (?, ?)`, e.ID, e.Tags)

// Load
var loaded Event
db.QueryRow(`SELECT id, tags FROM events WHERE id = ?`, "x").
    Scan(&loaded.ID, &loaded.Tags)
```

In a gqlgen `gqlgen.yml`:

```yaml
models:
  JSON:
    model: github.com/ubgo/jsonslice.JsonSlice
```

## Behaviour

- **`null` is preserved.** A nil slice marshals to `null`; the JSON value
  `null` deserialises back to a nil slice. An empty non-nil slice marshals
  to `[]`.
- **Scan accepts `[]byte`, `string`, or `nil`.** Anything else returns an
  error. Empty bytes scan to nil.
- **`MarshalGQL` does not panic.** Encoding errors fall back to writing
  `[]` rather than panicking inside a request goroutine.

## License

Apache-2.0 — see [`LICENSE`](LICENSE).
