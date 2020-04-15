package secret_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"

	"github.com/pivotal/build-service-cli/pkg/secret"
)

func TestSecretFactory(t *testing.T) {
	spec.Run(t, "TestSecretFactory", testSecretFactory)
}

func testSecretFactory(t *testing.T, when spec.G, it spec.S) {
	factory := &secret.Factory{}

	when("no params are set", func() {
		it("returns an error message", func() {
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "secret must be one of dockerhub, gcr, registry, or git")
		})
	})

	when("too many params are set", func() {
		it("returns an error message", func() {
			factory.DockerhubId = "some-dockerhub-id"
			factory.GcrServiceAccountFile = "some-gcr-service-account"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "secret must be one of dockerhub, gcr, registry, or git")
		})
	})

	when("sub params are mixed with dockerhub", func() {
		it("returns an error message", func() {
			factory.DockerhubId = "some-dockerhub-id"
			factory.RegistryUser = "some-reg-user"
			factory.GitUser = "some-git-user"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "extraneous parameters: git-user, registry-user")
		})
	})

	when("sub params are mixed with gcr", func() {
		it("returns an error message", func() {
			factory.GcrServiceAccountFile = "some-gcr-service-account-file"
			factory.RegistryUser = "some-reg-user"
			factory.GitSshKeyFile = "some-git-ssh-key-file"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "extraneous parameters: git-ssh-key, registry-user")
		})
	})

	when("registry is missing registry user", func() {
		it("returns an error message", func() {
			factory.Registry = "some-dockerhub-id"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "missing parameter registry-user")
		})
	})

	when("sub params are mixed with registry", func() {
		it("returns an error message", func() {
			factory.Registry = "some-registry"
			factory.RegistryUser = "some-reg-user"
			factory.GitUser = "some-git-user"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "extraneous parameters: git-user")
		})
	})

	when("sub params are mixed with git", func() {
		it("returns an error message", func() {
			it("returns an error message", func() {
				factory.Git = "some-git"
				factory.RegistryUser = "some-reg-user"
				factory.GitUser = "some-git-user"
				_, err := factory.MakeSecret("test-name", "test-namespace")
				require.EqualError(t, err, "extraneous parameters: registry-user")
			})
		})
	})

	when("neither git basic auth nor git ssh are provided", func() {
		it("returns an error message", func() {
			it("returns an error message", func() {
				factory.Git = "some-git"
				_, err := factory.MakeSecret("test-name", "test-namespace")
				require.EqualError(t, err, "missing parameter git-user or git-ssh-key")
			})
		})
	})

	when("both git basic auth and git ssh are provided", func() {
		it("returns an error message", func() {
			factory.Git = "some-git"
			factory.GitUser = "some-git-user"
			factory.GitSshKeyFile = "some-ssh-key"
			_, err := factory.MakeSecret("test-name", "test-namespace")
			require.EqualError(t, err, "must provide one of git-user or git-ssh-key")
		})
	})
}
