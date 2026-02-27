# Publishing v1.0.0

Step-by-step checklist to go from local repo to a public Go module on pkg.go.dev.

---

## 1. Create your GitHub repo

1. Go to github.com → New repository → name it `collection`
2. Set it to **Public** (required for pkg.go.dev indexing)
3. Do NOT initialise with README (you already have one)

---

## 2. Update the module path

Replace `github.com/km-arc/go-collections` with your actual path in every file:

```bash
# macOS / Linux
grep -rl "github.com/km-arc/go-collections" . \
  | xargs sed -i 's|github.com/km-arc/go-collections|github.com/km-arc/go-collections|g'
```

Files affected: `go.mod`, `collection_test.go`, `bench_test.go`, `example_test.go`

---

## 3. Push to GitHub

```bash
git init
git add .
git commit -m "feat: initial implementation — Laravel collections for Go"
git branch -M main
git remote add origin git@github.com:km-arc/go-collections.git
git push -u origin main
```

---

## 4. Set up Codecov (free for public repos)

1. Go to [codecov.io](https://codecov.io) → Sign in with GitHub
2. Add your `collection` repo
3. Copy the **CODECOV_TOKEN** from the repo settings
4. In your GitHub repo → Settings → Secrets → Actions → New secret
   - Name: `CODECOV_TOKEN`
   - Value: *(paste token)*

---

## 5. Tag v1.0.0

```bash
git tag -a v1.0.0 -m "Release v1.0.0

Initial release of Laravel-style fluent collections for Go.

- Collection[T]: 60+ methods matching Laravel Collection API
- MapCollection[K, V]: full key/value collection API
- Zero dependencies, Go 1.21 generics
- 70+ tests, benchmark suite, GoDoc examples"

git push origin v1.0.0
```

This will trigger the `.github/workflows/release.yml` workflow, which:
- Runs the full test suite
- Runs benchmarks and attaches results
- Creates a GitHub Release with auto-generated release notes

---

## 6. Publish to pkg.go.dev

pkg.go.dev automatically indexes public Go modules. Just run:

```bash
GOPROXY=https://proxy.golang.org go get github.com/km-arc/go-collections@v1.0.0
```

Or visit:
```
https://pkg.go.dev/github.com/km-arc/go-collections@v1.0.0
```

It may take a few minutes to appear.

---

## 7. Add badges to README

Replace the badge URLs (they already reference `km-arc/go-collections`) with your actual username. The badges will go live automatically once:

- CI runs (GitHub Actions badge)
- Coverage uploads to Codecov (codecov badge)
- Module is indexed (pkg.go.dev badge)
- Go Report Card scans your repo at `goreportcard.com/report/github.com/km-arc/go-collections`

---

## 8. Post on LinkedIn

See `LINKEDIN_POST.md` — copy, paste, add a screenshot of the GitHub repo or pkg.go.dev page.

---

## Checklist before posting

- [ ] All tests pass: `go test -race ./...`
- [ ] No lint errors: `golangci-lint run`
- [ ] go.mod is tidy: `go mod tidy`
- [ ] CI is green on GitHub
- [ ] Coverage badge is showing
- [ ] pkg.go.dev shows your docs
- [ ] README badges all resolve
