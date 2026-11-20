package experience

import (
	"errors"
	"path/filepath"
	"testing"

	"gpuworks.ai/gpuworks/internal/server/store"
	"gpuworks.ai/gpuworks/internal/testenv"
)

func TestExperienceTemplateCanBeDeletedAndStaysDeleted(t *testing.T) {
	env := testenv.New(t)
	client, err := login(t.Context(), env.URL(), "admin", testenv.AdminPassword, "")
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "seeded-template")
	seed := func() {
		t.Helper()
		if err := seedExampleTemplate(t.Context(), client, env.Store(), marker); err != nil {
			t.Fatal(err)
		}
	}
	seed()
	const id = "experience/qwen3-8b"
	model, err := env.Store().InferenceModelByModelID(id)
	if err != nil {
		t.Fatal(err)
	}
	if model.Source != "user" || len(model.Specs) == 0 {
		t.Fatalf("expected a user-owned template with specs, got source %q, %d specs", model.Source, len(model.Specs))
	}
	if err := client.call(t.Context(), "DELETE", "/api/v1/inference-templates?model_id="+id, nil, nil); err != nil {
		t.Fatal(err)
	}
	seed()
	if _, err := env.Store().InferenceModelByModelID(id); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("deleted template after reseeding: %v", err)
	}
}
