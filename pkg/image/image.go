package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/NilayYadav/agentos/pkg/common"
	"github.com/rs/zerolog"
)

func PullBrowserlessImage(ctx context.Context, logger zerolog.Logger) (string, error) {
	sourceImage := "ghcr.io/browserless/chromium:latest"
	buildPath, err := os.MkdirTemp("", "browserless-image-")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to create temp directory")
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	ociPath, bundlePath := filepath.Join(buildPath, "oci"), filepath.Join(buildPath, "bundle")

	os.MkdirAll(ociPath, 0755)
	os.MkdirAll(bundlePath, 0755)

	if err := common.Exec(ctx, "skopeo", "inspect", fmt.Sprintf("docker://%s", sourceImage)); err != nil {
		return "", fmt.Errorf("failed to inspect image: %w", err)
	}

	logger.Info().Msgf("Image compatible with %s/%s", runtime.GOOS, runtime.GOARCH)

	if err := common.Exec(ctx, "skopeo", "copy", fmt.Sprintf("docker://%s", sourceImage), fmt.Sprintf("oci:%s:latest", ociPath)); err != nil {
		return "", fmt.Errorf("failed to copy image: %w", err)
	}

	if err := common.Exec(ctx, "umoci", "unpack", "--image", fmt.Sprintf("%s:latest", ociPath), bundlePath); err != nil {
		return "", fmt.Errorf("failed to unpack image: %w", err)
	}

	archivePath := filepath.Join(buildPath, "browserless-image.clip")

	if err := common.Exec(ctx, "tar", "-czf", archivePath, "-C", bundlePath, "."); err != nil {
		return "", fmt.Errorf("failed to create archive: %w", err)
	}

	logger.Info().Msgf("Successfully pulled and prepared image: %s", archivePath)

	return buildPath, nil
}
