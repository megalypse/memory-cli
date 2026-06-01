package main_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/megalypse/memory_cli/internal/memory"
	memorygroup "github.com/megalypse/memory_cli/internal/memory_group"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "modernc.org/sqlite"
)

func TestMemoryGroupRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db := newIntegrationDB(t, ctx)
	repository := memorygroup.NewRepositorySqlLite(db)

	err := repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        "work",
		Description: "Work notes",
	})
	if err != nil {
		t.Fatalf("create memory group: %v", err)
	}

	err = repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        "personal",
		Description: "Personal notes",
	})
	if err != nil {
		t.Fatalf("create second memory group: %v", err)
	}

	groups, err := repository.GetAll(ctx)
	if err != nil {
		t.Fatalf("get all memory groups: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("expected 2 memory groups, got %d", len(groups))
	}

	work := findGroup(t, groups, "work")
	work.Name = "work-updated"
	work.Description = "Updated work notes"

	err = repository.Put(ctx, work)
	if err != nil {
		t.Fatalf("put memory group: %v", err)
	}

	groups, err = repository.GetAll(ctx)
	if err != nil {
		t.Fatalf("get all memory groups after update: %v", err)
	}

	updated := findGroup(t, groups, "work-updated")
	if updated.Description != "Updated work notes" {
		t.Fatalf("expected updated description, got %q", updated.Description)
	}
}

func TestMemoryRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db := newIntegrationDB(t, ctx)
	groupRepository := memorygroup.NewRepositorySqlLite(db)
	memoryRepository := memory.NewRepositorySqlLite(db)

	group := createGroup(t, ctx, groupRepository, "engineering")
	otherGroup := createGroup(t, ctx, groupRepository, "personal")

	_, err := db.ExecContext(ctx, `
INSERT INTO memory_fts (id, name, content) VALUES (?, ?, ?)
`, 999, "noise", "unrelated fts row")
	if err != nil {
		t.Fatalf("seed unrelated fts row: %v", err)
	}

	first := &memory.Memory{
		GroupID: group.ID,
		Name:    "alpha",
		Content: "golang sqlite repository",
	}
	second := &memory.Memory{
		GroupID: group.ID,
		Name:    "bravo",
		Content: "testcontainers integration reference",
	}
	third := &memory.Memory{
		GroupID: otherGroup.ID,
		Name:    "charlie",
		Content: "outside selected group",
	}

	for _, item := range []*memory.Memory{first, second, third} {
		err := memoryRepository.Create(ctx, item)
		if err != nil {
			t.Fatalf("create memory %q: %v", item.Name, err)
		}
		if item.ID == 0 {
			t.Fatalf("expected created memory %q to receive an id", item.Name)
		}
	}

	assertSingleFTSRow(t, ctx, db, first.ID)

	memories, err := memoryRepository.GetAllByGroup(ctx, group.ID)
	if err != nil {
		t.Fatalf("get memories by group: %v", err)
	}

	if len(memories) != 2 {
		t.Fatalf("expected 2 memories for group, got %d", len(memories))
	}

	first.Content = "updated sqlite full text search"
	err = memoryRepository.Put(ctx, first)
	if err != nil {
		t.Fatalf("put memory: %v", err)
	}

	assertSingleFTSRow(t, ctx, db, first.ID)

	references, err := memoryRepository.FindReferences(ctx, []string{"updated", "testcontainers"})
	if err != nil {
		t.Fatalf("find references: %v", err)
	}

	assertMemoryNames(t, references, "alpha", "bravo")

	err = memoryRepository.LinkMemories(ctx, first, []*memory.Memory{second})
	if err != nil {
		t.Fatalf("link memories: %v", err)
	}

	var links int
	err = db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_memory
WHERE memory_id_1 = ? AND memory_id_2 = ?
`, first.ID, second.ID).Scan(&links)
	if err != nil {
		t.Fatalf("count memory links: %v", err)
	}

	if links != 1 {
		t.Fatalf("expected 1 memory link, got %d", links)
	}
}

func newIntegrationDB(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:      "postgres:16-alpine",
			Cmd:        []string{"sh", "-c", "while true; do sleep 3600; done"},
			SkipReaper: true,
			WaitingFor: wait.ForExec([]string{"sh", "-c", "test -d /tmp"}).
				WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start test container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Errorf("terminate test container: %v", err)
		}
	})

	dbPath := filepath.Join(t.TempDir(), "memory.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close sqlite db: %v", err)
		}
	})

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}

	goose.SetBaseFS(os.DirFS("."))

	if err := goose.Up(db, "migrations"); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return db
}

func createGroup(t *testing.T, ctx context.Context, repository memorygroup.RepositoryMemoryGroup, name string) *memorygroup.MemoryGroup {
	t.Helper()

	err := repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        name,
		Description: name + " description",
	})
	if err != nil {
		t.Fatalf("create group %q: %v", name, err)
	}

	groups, err := repository.GetAll(ctx)
	if err != nil {
		t.Fatalf("get groups after creating %q: %v", name, err)
	}

	return findGroup(t, groups, name)
}

func findGroup(t *testing.T, groups []*memorygroup.MemoryGroup, name string) *memorygroup.MemoryGroup {
	t.Helper()

	for _, group := range groups {
		if group.Name == name {
			return group
		}
	}

	t.Fatalf("expected group %q to exist", name)
	return nil
}

func assertMemoryNames(t *testing.T, memories []*memory.Memory, expected ...string) {
	t.Helper()

	seen := make(map[string]bool, len(memories))
	for _, item := range memories {
		seen[item.Name] = true
	}

	for _, name := range expected {
		if !seen[name] {
			t.Fatalf("expected memory %q in results, got %#v", name, seen)
		}
	}
}

func assertSingleFTSRow(t *testing.T, ctx context.Context, db *sql.DB, memoryID int) {
	t.Helper()

	var rows int
	err := db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_fts
WHERE id = ?
`, memoryID).Scan(&rows)
	if err != nil {
		t.Fatalf("count fts rows for memory %d: %v", memoryID, err)
	}

	if rows != 1 {
		t.Fatalf("expected 1 fts row for memory %d, got %d", memoryID, rows)
	}
}
