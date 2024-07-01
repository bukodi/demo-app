package tkeycloak

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"
)

const keycloakImage = "quay.io/keycloak/keycloak:latest"

func TestDockerStart(t *testing.T) {
	t.Skip("Skip Docker test because it takes too long")
	slog.SetLogLoggerLevel(slog.LevelDebug)

	kc := &KeycloakContainer{}
	kc.DockerTimeout = time.Second * 20

	if err := kc.Start(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to start keycloak: %w", err))
	}
	t.Logf("Keycloak container started: %s", kc.id)

	if err := kc.Stop(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to stop keycloak: %w", err))
	} else {
		t.Logf("Keycloak container stopped: %s", kc.id)
	}

}
