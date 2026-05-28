# Dev Journal

## Project
Echo - Finance Tracker  
Stack: React, Go, PostgreSQL (Supabase), pgxpool, Render

---

# 2026-05-24

## What I Worked On
- Initializing DB connection to Supabase(PostgreSQL)
- Setting up router
- Bruno hands on

## Decisions Made
- Create the project from scratch using React instead of htmx for UI
- Stop trying to create perfect repo with proper migrations
- Use Bruno isntead of Postman or Insomnia

---
# 2026-05-28
## What I Worked On
- Reorganize accounts into feature-first structure
- Optimize the tests

I was using a `setupTestDB` function for testing the db integration using dockertest. But the test took ~10s.
```bash
❯ go test ./internal/accounts/... -count=1
ok      github.com/joshu-sajeev/echo/internal/accounts  9.687s
```
This is beacuse `setupTestDB` function is spinning up a brand-new Docker container, waiting for PostgreSQL to boot, running migrations, and then tearing down that container for every single test function that calls it.

So I added had to spin up the Docker container exactly once for the entire package, run migrations once, and let the individual tests share that pool.
```go
func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Initialize dockertest pool
	pool, err := dockertest.NewPool(ctx, "")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	// 2. Run the Postgres container ONCE
	postgres, err := pool.Run(
		ctx,
		"postgres",
		dockertest.WithTag("14"),
		dockertest.WithEnv([]string{
			"POSTGRES_PASSWORD=" + testDBPassword,
			"POSTGRES_DB=" + testDBName,
		}),
	)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	hostPort := postgres.GetHostPort(testDBPort + "/tcp")
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser,
		testDBPassword,
		hostPort,
		testDBName,
	)

	// 3. Wait for Postgres to be ready
	err = pool.Retry(ctx, 30*time.Second, func() error {
		var err error
		globalDBPool, err = pgxpool.New(ctx, databaseURL)
		if err != nil {
			return err
		}
		return globalDBPool.Ping(ctx)
	})
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	// 4. Run migrations ONCE
	sqlDB := stdlib.OpenDBFromPool(globalDBPool)
	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	// 5. Run all tests in the package
	code := m.Run()

	// 6. Global Teardown after all tests finish
	sqlDB.Close()
	globalDBPool.Close()
	_ = postgres.Close(ctx) // stop & remove container

	os.Exit(code)
}
```
This reduced to the text execution to ~3s
```bash
echo/backend on  rewrite [+] via  v1.26.2 took 2s 
❯ go test ./internal/accounts/... -count=1
ok      github.com/joshu-sajeev/echo/internal/accounts  2.431s
```
