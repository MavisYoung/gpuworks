package experience

import (
	"context"
	"errors"
	"os"

	"gpuworks.ai/gpuworks/internal/server/store"
)

// Keep a user-owned template in the initial catalog so visitors can try
// deleting a spec or a whole template without touching a starter deployment.
func seedExampleTemplate(ctx context.Context, c *apiClient, db *store.Store, marker string) error {
	if _, err := os.Stat(marker); err == nil || !errors.Is(err, os.ErrNotExist) {
		return err
	}
	const id = "experience/qwen3-8b"
	// A previous run may have created the clone but stopped before the marker.
	// Preserve that copy, including edits made by a visitor.
	_, err := db.InferenceModelByModelID(id)
	if errors.Is(err, store.ErrNotFound) {
		err = c.post(ctx, "/api/v1/inference-templates/clone", map[string]string{
			"source_model_id": "qwen/qwen3-8b", "model_id": id, "display_title": "Qwen3-8B · Experience",
		}, nil)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(marker, nil, 0o600)
}
