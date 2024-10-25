// A generated module for RegistryCache functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/registry-cache/internal/dagger"
	"fmt"
	"strings"
	"time"
)

type RegistryCache struct {
	Sock *dagger.Socket
}

func New(sock *dagger.Socket) *RegistryCache {
	return &RegistryCache{
		Sock: sock,
	}
}

func (m *RegistryCache) DockerClient() *dagger.Container {
	return dag.Container().
		From("docker:27-cli").
		WithUnixSocket("/var/run/docker.sock", m.Sock)
}

func parseImageID(out string) (string, error) {
	sha := strings.Split(out, ": ")
	if len(sha) < 2 {
		return "", fmt.Errorf("cannot get container image sha256")
	}

	imageSha := strings.Trim(sha[1], "\n ")
	return imageSha, nil
}

func (m *RegistryCache) LoadDaggerEngineImage(ctx context.Context, daggerVersion, registryMirror string) (string, error) {
	engineImage := fmt.Sprintf("registry.dagger.io/engine:%s", daggerVersion)
	daggerEngine := dag.Container().From(engineImage).
		WithNewFile("/etc/dagger/engine.toml", fmt.Sprintf(DaggerEngineConfig, registryMirror, registryMirror))

	out, err := m.DockerClient().
		WithMountedFile("/daggerEngine.tar", daggerEngine.AsTarball()).
		WithExec([]string{"docker", "load", "-qi", "/daggerEngine.tar"}).
		Stdout(ctx)

	if err != nil {
		return "", err
	}

	return parseImageID(out)
}

func (m *RegistryCache) LoadRegistryImage(ctx context.Context) (string, error) {
	dockerRegistry := dag.Container().From("registry:2").
		WithNewFile("/etc/docker/registry/config.yml", RegistryConfig)

	out, err := m.DockerClient().
		WithMountedFile("/daggerEngine.tar", dockerRegistry.AsTarball()).
		WithExec([]string{"docker", "load", "-qi", "/daggerEngine.tar"}).
		Stdout(ctx)

	if err != nil {
		return "", err
	}

	return parseImageID(out)
}

func (m *RegistryCache) InitRegistryMirror(
	ctx context.Context,
	// +optional
	// +default=""
	storagePath string,
) error {
	// Cleanup
	_, err := m.DockerClient().
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "docker rm -f registry-mirror-docker.io || true"}).
		WithExec([]string{"sh", "-c", "docker network create dagger-registry || true"}).
		Sync(ctx)

	if err != nil {
		return err
	}

	registryImageID, err := m.LoadRegistryImage(ctx)
	if err != nil {
		return err
	}

	execArgs := []string{
		"docker", "run", "-d",
		"--network", "dagger-registry",
		"--name", "registry-mirror-docker.io",
		"-p", "5000:5000",
	}

	if storagePath != "" {
		execArgs = append(execArgs, "--mount",
			fmt.Sprintf("type=bind,source=%s,target=/var/lib/registry", storagePath),
		)
	}

	execArgs = append(execArgs, registryImageID)

	_, err = m.DockerClient().
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(execArgs).Sync(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (m *RegistryCache) ConfigureDaggerEngine(
	ctx context.Context,
	// +optional
	// +default="registry-mirror-docker.io:5000"
	registryMirror string,
) (string, error) {
	daggerVersion, err := dag.Version(ctx)
	if err != nil {
		return "", err
	}

	// Cleanup
	containerName := fmt.Sprintf("dagger-engine-mirrored-%s", daggerVersion)
	_, err = m.DockerClient().
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", fmt.Sprintf("docker rm -f %q || true", containerName)}).
		WithExec([]string{"sh", "-c", "docker network create dagger-registry || true"}).
		Sync(ctx)

	if err != nil {
		return "", err
	}

	daggerEngineImageID, err := m.LoadDaggerEngineImage(ctx, daggerVersion, registryMirror)
	if err != nil {
		return "", err
	}

	_, err = m.DockerClient().
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{
			"docker", "run", "--privileged", "-d",
			"--network", "dagger-registry",
			"--name", containerName, daggerEngineImageID,
			"--debug",
		}).Sync(ctx)

	if err != nil {
		return "", err
	}

	resp := fmt.Sprintf("export _EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://%s", containerName)
	return resp, nil
}
