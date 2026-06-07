package main_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/megalypse/memory_cli/internal/memory"
	memorygroup "github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "modernc.org/sqlite"
)

func TestMemoryGroupRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db := newIntegrationDB(t, ctx)
	repository := memorygroup.GetRepositorySqlLite(db)

	err := repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        "work",
		Description: "Work notes",
	})
	require.NoError(t, err)

	err = repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        "personal",
		Description: "Personal notes",
	})
	require.NoError(t, err)

	groups, err := repository.GetAll(ctx)
	require.NoError(t, err)

	assert.Len(t, groups, 2)

	work := findGroup(t, groups, "work")
	work.Name = "work-updated"
	work.Description = "Updated work notes"

	err = repository.Put(ctx, work)
	require.NoError(t, err)

	groups, err = repository.GetAll(ctx)
	require.NoError(t, err)

	updated := findGroup(t, groups, "work-updated")
	assert.Equal(t, "Updated work notes", updated.Description)
}

func TestMemoryRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db := newIntegrationDB(t, ctx)
	groupRepository := memorygroup.GetRepositorySqlLite(db)
	memoryRepository := memory.GetRepositorySqlLite(db)

	group := createGroup(t, ctx, groupRepository, "engineering")
	otherGroup := createGroup(t, ctx, groupRepository, "personal")

	_, err := db.ExecContext(ctx, `
INSERT INTO memory_fts (id, name, content) VALUES (?, ?, ?)
`, 999, "noise", "unrelated fts row")
	require.NoError(t, err)

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
		require.NoError(t, err)
		assert.NotZero(t, item.ID)
	}

	assertSingleFTSRow(t, ctx, db, first.ID)

	memories, err := memoryRepository.GetAllByGroup(ctx, group.ID)
	require.NoError(t, err)

	assert.Len(t, memories, 2)

	first.Content = "updated sqlite full text search"
	err = memoryRepository.Put(ctx, first)
	require.NoError(t, err)

	assertSingleFTSRow(t, ctx, db, first.ID)

	references, err := memoryRepository.FindReferences(
		ctx,
		group.ID,
		[]string{"updated", "testcontainers", "outside"},
	)
	require.NoError(t, err)

	assertMemoryNames(t, references, "alpha", "bravo")
	assert.NotContains(t, memoryNames(references), "charlie")

	second.Content = "testcontainers integration lutando"
	err = memoryRepository.Put(ctx, second)
	require.NoError(t, err)

	prefixResults, err := memoryRepository.QueryMemories(ctx, group.ID, "lutan")
	require.NoError(t, err)
	require.Len(t, prefixResults, 1)
	assert.Equal(t, second.ID, prefixResults[0].ID)

	err = memoryRepository.LinkMemories(ctx, first, []*memory.Memory{second})
	require.NoError(t, err)

	err = memoryRepository.LinkMemories(ctx, first, []*memory.Memory{third})
	assert.ErrorIs(t, err, memory.ErrDifferentMemoryGroup)

	_, err = db.ExecContext(ctx, `
INSERT INTO memory_memory (memory_id_1, memory_id_2) VALUES (?, ?)
`, first.ID, third.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MEMORIES MUST BELONG TO THE SAME GROUP")

	firstRelations, err := memoryRepository.GetRelations(ctx, first)
	require.NoError(t, err)
	assertMemoryNames(t, firstRelations, "bravo")

	secondRelations, err := memoryRepository.GetRelations(ctx, second)
	require.NoError(t, err)
	assertMemoryNames(t, secondRelations, "alpha")

	var links int
	err = db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_memory
WHERE memory_id_1 = ? AND memory_id_2 = ?
`, first.ID, second.ID).Scan(&links)
	require.NoError(t, err)

	assert.Equal(t, 1, links)

	first.Name = "sqlite ranking"
	first.Content = "secondary"
	err = memoryRepository.Put(ctx, first)
	require.NoError(t, err)

	second.Name = "secondary"
	second.Content = "sqlite ranking"
	err = memoryRepository.Put(ctx, second)
	require.NoError(t, err)

	third.Name = "sqlite outside group"
	err = memoryRepository.Put(ctx, third)
	require.NoError(t, err)

	searchResults, err := memoryRepository.QueryMemories(ctx, group.ID, "sqlite")
	require.NoError(t, err)
	require.Len(t, searchResults, 2)
	assert.Equal(t, first.ID, searchResults[0].ID)

	err = memoryRepository.Delete(ctx, &memory.Memory{
		ID:      first.ID,
		GroupID: otherGroup.ID,
	})
	require.ErrorIs(t, err, sql.ErrNoRows)

	err = memoryRepository.Delete(ctx, first)
	require.NoError(t, err)

	var deletedMemoryRows int
	err = db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memories WHERE id = ?
`, first.ID).Scan(&deletedMemoryRows)
	require.NoError(t, err)
	assert.Zero(t, deletedMemoryRows)

	var deletedFTSRows int
	err = db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_fts WHERE id = ?
`, first.ID).Scan(&deletedFTSRows)
	require.NoError(t, err)
	assert.Zero(t, deletedFTSRows)

	var deletedRelationRows int
	err = db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_memory
WHERE memory_id_1 = ? OR memory_id_2 = ?
`, first.ID, first.ID).Scan(&deletedRelationRows)
	require.NoError(t, err)
	assert.Zero(t, deletedRelationRows)
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
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(context.Background()))
	})

	dbPath := filepath.Join(t.TempDir(), "memory.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	require.NoError(t, goose.SetDialect("sqlite3"))

	goose.SetBaseFS(os.DirFS("."))

	require.NoError(t, goose.Up(db, "migrations"))

	return db
}

func createGroup(t *testing.T, ctx context.Context, repository memorygroup.RepositoryMemoryGroup, name string) *memorygroup.MemoryGroup {
	t.Helper()

	err := repository.Create(ctx, &memorygroup.MemoryGroup{
		Name:        name,
		Description: name + " description",
	})
	require.NoError(t, err)

	groups, err := repository.GetAll(ctx)
	require.NoError(t, err)

	return findGroup(t, groups, name)
}

func findGroup(t *testing.T, groups []*memorygroup.MemoryGroup, name string) *memorygroup.MemoryGroup {
	t.Helper()

	for _, group := range groups {
		if group.Name == name {
			return group
		}
	}

	require.Failf(t, "expected group to exist", "group %q was not found", name)
	return nil
}

func assertMemoryNames(t *testing.T, memories []*memory.Memory, expected ...string) {
	t.Helper()

	seen := make(map[string]bool, len(memories))
	for _, item := range memories {
		seen[item.Name] = true
	}

	for _, name := range expected {
		assert.Truef(t, seen[name], "expected memory %q in results, got %#v", name, seen)
	}
}

func memoryNames(memories []*memory.Memory) []string {
	names := make([]string, len(memories))
	for index, item := range memories {
		names[index] = item.Name
	}

	return names
}

func assertSingleFTSRow(t *testing.T, ctx context.Context, db *sql.DB, memoryID int) {
	t.Helper()

	var rows int
	err := db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM memory_fts
WHERE id = ?
`, memoryID).Scan(&rows)
	require.NoError(t, err)

	assert.Equal(t, 1, rows)
}
