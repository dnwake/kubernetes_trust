package trust

///
/// Bridge to Docker Content Trust functionality for Kubernetes
/// Based on github.com/docker/api/client/trust.go
///

/// Two modifications to the original's functionality
///  1) We don't pull every target if the original leaves the tag unspecified
///  2) Less informative error message if there is a notary error
import (
	"encoding/hex"
	"errors"
	"fmt"
        "github.com/dwake/docker-trust/external/github.com/Sirupsen/logrus"
        "github.com/dwake/docker-trust/external/github.com/docker/distribution/digest"
        "github.com/dwake/docker-trust/external/github.com/docker/distribution/registry/client/auth"
        "github.com/dwake/docker-trust/external/github.com/docker/distribution/registry/client/transport"
	"github.com/dwake/docker-trust/external/github.com/docker/docker/api/client"
	"github.com/dwake/docker-trust/external/github.com/docker/docker/cliconfig"
	"github.com/dwake/docker-trust/external/github.com/docker/docker/reference"
	"github.com/dwake/docker-trust/external/github.com/docker/docker/registry"
        "github.com/dwake/docker-trust/external/github.com/docker/engine-api/types"
        registrytypes "github.com/dwake/docker-trust/external/github.com/docker/engine-api/types/registry"
	"github.com/dwake/docker-trust/external/github.com/docker/go-connections/tlsconfig"
	"github.com/dwake/docker-trust/external/github.com/docker/notary/client"
	"github.com/dwake/docker-trust/external/github.com/docker/notary/passphrase"
	"github.com/dwake/docker-trust/external/github.com/docker/notary/tuf/data"
	"net"
	"net/http"
	"net/url"
	"os"
        "path"
        "path/filepath"
        "strconv"
	"time"
)

// derived from addTrustedFlags() and 
func ShouldUseContentTrust(tag string) (bool) {
	var trusted bool
	if e := os.Getenv("DOCKER_CONTENT_TRUST"); e != "" {
		if t, err := strconv.ParseBool(e); t || err != nil {
			// treat any other value as true
			trusted = true
		}
	}

	hasDigest := registry.ParseReference(tag).HasDigest()

	return trusted && !hasDigest
}

/// This method returns only one digest, while the original method (trustedPull())
///  could pull multiple digests if no tag is specified
///
/// I don't think this can be delegated b/c original method actually performs the pull
func GetTrustedDigestToPull(repository string, tag string, 
                            username string, password string, 
                            email string, serverAddress string) (string, error) {
	remote := repository + ":" + tag
	distributionRef, err := reference.ParseNamed(remote)
	repoInfo, err := registry.ParseRepositoryInfo(distributionRef)
	ref := registry.ParseReference(tag)

	newAuthConfig := types.AuthConfig{}
	newAuthConfig.Username = username
	newAuthConfig.Password = password
	newAuthConfig.Email = email
	newAuthConfig.ServerAddress = serverAddress

	notaryRepo, err := getNotaryRepository(repoInfo, newAuthConfig)
	if err != nil {
		fmt.Sprintf("Error establishing connection to trust repository: %s\n", err)
		return "", err
	}

	releasesRole := path.Join(data.CanonicalTargetsRole, "releases")
	t, err := notaryRepo.GetTargetByName(ref.String(), releasesRole, data.CanonicalTargetsRole)
	if err != nil {
                /// NOTE: This error is less informative than the original
		return "", fmt.Errorf("Notary error for repository %s: %s", repoInfo.FullName(), err.Error())
	}

	h, ok := t.Hashes["sha256"]
	if !ok {
		return "", errors.New("no valid hash, expecting sha256")
	}

	return digest.NewDigestFromHex("sha256", hex.EncodeToString(h)).String(), nil
}

// identical except original made references to CLI in last line
func getNotaryRepository(repoInfo *registry.RepositoryInfo, authConfig types.AuthConfig) (*client.NotaryRepository, error) {
	server, err := client.trustServer(repoInfo.Index)
	if err != nil {
		return nil, err
	}

	var cfg = tlsconfig.ClientDefault
	cfg.InsecureSkipVerify = !repoInfo.Index.Secure

	// Get certificate base directory
	certDir, err := certificateDirectory(server)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("reading certificate directory: %s", certDir)

	if err := registry.ReadCertsDirectory(&cfg, certDir); err != nil {
		return nil, err
	}

	base := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &cfg,
		DisableKeepAlives:   true,
	}

	modifiers := registry.DockerHeaders(http.Header{})
	authTransport := transport.NewTransport(base, modifiers...)
	pingClient := &http.Client{
		Transport: authTransport,
		Timeout:   5 * time.Second,
	}
	endpointStr := server + "/v2/"
	req, err := http.NewRequest("GET", endpointStr, nil)
	if err != nil {
		return nil, err
	}

	challengeManager := auth.NewSimpleChallengeManager()

	resp, err := pingClient.Do(req)
	if err != nil {
		// Ignore error on ping to operate in offline mode
		logrus.Debugf("Error pinging notary server %q: %s", endpointStr, err)
	} else {
		defer resp.Body.Close()

		// Add response to the challenge manager to parse out
		// authentication header and register authentication method
		if err := challengeManager.AddResponse(resp); err != nil {
			return nil, err
		}
	}

	creds := simpleCredentialStore{auth: authConfig}
	tokenHandler := auth.NewTokenHandler(authTransport, creds, repoInfo.FullName(), "push", "pull")
	basicHandler := auth.NewBasicHandler(creds)

	modifiers = append(modifiers, transport.RequestModifier(auth.NewAuthorizer(challengeManager, tokenHandler, basicHandler)))
	tr := transport.NewTransport(base, modifiers...)

	trustDirectory := filepath.Join(cliconfig.ConfigDir(), "trust")
	passphraseRetriever := getPassphraseRetriever()
	return client.NewNotaryRepository(trustDirectory, repoInfo.FullName(), server, tr, passphraseRetriever)
}

// identical
func certificateDirectory(server string) (string, error) {
	u, err := url.Parse(server)
	if err != nil {
		return "", err
	}

	return filepath.Join(cliconfig.ConfigDir(), "tls", u.Host), nil
}

// identical
type simpleCredentialStore struct {
	auth types.AuthConfig
}

// identical
func (scs simpleCredentialStore) Basic(u *url.URL) (string, string) {
     return scs.auth.Username, scs.auth.Password
}

// faking out cli.getPassphraseRetriever()
func getPassphraseRetriever() (passphrase.Retriever){
	env := map[string]string{
		"root":             os.Getenv("DOCKER_CONTENT_TRUST_ROOT_PASSPHRASE"),
		"snapshot":         os.Getenv("DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE"),
		"targets":          os.Getenv("DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE"),
		"targets/releases": os.Getenv("DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE"),
	}

	return func(keyName string, alias string, createNew bool, numAttempts int) (string, bool, error) {
		if v := env[alias]; v != "" {
			return v, numAttempts > 1, nil
		}
		return "", false, fmt.Errorf("No passphrase found for %s", alias)
	}
}
