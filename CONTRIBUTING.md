# How to contribute

There are multiple ways to help. We rely on the upstream
[holidays/definitions](https://github.com/holidays/definitions) project and
its contributors to keep holiday data accurate and up to date, and pull
requests to this repository to address bugs or implement new features are
always welcome.

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.
Everyone interacting with this project is expected to abide by its terms.

## AI Usage

Please read our [AI Usage Policy](AI_POLICY.md) before contributing.

## Commit requirements

All commits must include a `Signed-off-by` trailer (the
[Developer Certificate of Origin](https://developercertificate.org/)). Use
`git commit -s`:

```sh
git commit -s -m "Your commit message"
```

This appends the following to your commit message automatically:

```
Signed-off-by: Your Name <your@email.com>
```

GPG-signed commits are not required by this repository.

## General note on the definitions submodule

Definitions live in a git submodule. Clone with:

```sh
git clone --recurse-submodules https://github.com/ppeble/go-holidays
```

or, in an existing clone:

```sh
git submodule update --init
```

To bump the submodule to a newer upstream tag:

```sh
make update-definitions DEFS_TAG=vX.Y.Z
git add definitions
git commit -s -m "vendor holidays/definitions vX.Y.Z"
make generate
```

## For definition updates

Definition changes belong upstream, not here. Definitions are written in YAML
and live in the [holidays/definitions](https://github.com/holidays/definitions)
repository so they can be used by tools written in other languages; a
complete guide to the format is in its
[SYNTAX guide](https://github.com/holidays/definitions/blob/master/doc/SYNTAX.md).

Once you have an idea of what you want to change, see that repository's own
[CONTRIBUTING guide](https://github.com/holidays/definitions/blob/master/doc/CONTRIBUTING.md).
After that PR is accepted upstream, this repository's maintainers are
responsible for bumping the submodule, regenerating, and releasing.

Never hand-edit `definitions/*.yaml` (it is a submodule checkout) or
`internal/definitions/*.go` / `*_test.go` (they are generated).

## For non-definition functionality

* Fork this repository.
* Make your changes. Run `make test` to execute the test suite.
* Open a PR pointing back to `main`.

## Regenerating definitions

The generated `internal/definitions/*.go` files are checked in. Regenerate
after bumping the submodule or after adding new per-region methods:

```sh
make generate                       # all regions, fail on unported methods
bin/gen-holidays -allow-unported    # all regions, skip those with missing methods
bin/gen-holidays -regions us,gb     # only specific regions
```

The generator hard-fails on any `function:` or `observed:` YAML reference that
lacks a registered Go implementation.

### Adding a missing method

1. Write the Go function in `internal/engine/methods_<cc>.go`.
2. Register it in that file's `init()` with
   `engine.RegisterMethod("<name>", func(a MethodArgs) (time.Time, error) { ... })`.
3. Re-run `make generate`.

## Testing

```sh
make test          # go vet + go test ./...
```

Ginkgo v2 and Gomega are the required test framework. Region tests under
`internal/definitions` are generated table tests and must not be hand-edited;
change the upstream YAML or the generator, then `make generate`.

## Parity suite

`parity/` compares this library's output against the Ruby gem pinned in
`parity/Gemfile`, loading our own `definitions/` YAML into the gem so both
sides resolve the same rules. It is run by the required "parity" CI check on
every PR; do not run it locally, it needs Ruby and the pinned gem installed.
See `parity/README.md` for the design.

## Local development helpers

* `make build` - builds `bin/holidays` and `bin/gen-holidays`
* `make test` - runs `go vet` then the full test suite
* `make vet` - runs `go vet ./...` only
* `make staticcheck` - runs `staticcheck ./...` (installs it if missing)
* `make generate` - regenerates `internal/definitions` from the YAML submodule
* `make update-definitions DEFS_TAG=vX.Y.Z` - checks the submodule out at a tag
* `make parity` - runs the Ruby-vs-Go comparison suite (CI only, see above)
* `make clean` - removes `bin/`
