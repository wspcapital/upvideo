package youtubeauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0

To make this sample run, you need to populate the client_secrets.json file
found at:

   %v

with information from the {{ Google Cloud Console }}
{{ https://cloud.google.com/console }}

For more information about the client_secrets.json file format, please visit:
https://developers.google.com/api-client-library/python/guide/aaa_client_secrets
`

// Cache specifies the methods that implement a Token cache.
type Cache interface {
	Token() (*oauth2.Token, error)
	PutToken(*oauth2.Token) error
}

// CacheFile implements Cache. Its value is the name of the file in which
// the Token is stored in JSON format.
type CacheFile string

// ClientConfig is a data structure definition for the client_secrets.json file.
// The code unmarshals the JSON configuration file into this structure.
type ClientConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}

// Config is a root-level configuration object.
type Config struct {
	Installed ClientConfig `json:"installed"`
	Web       ClientConfig `json:"web"`
}

// readConfig reads the configuration from clientSecretsFile.
// It returns an oauth configuration object for use with the Google API client.
func readConfig(secretsPath string) (*oauth2.Config, error) {
	// Read the secrets file
	data, err := ioutil.ReadFile(secretsPath)
	if err != nil {
		pwd, _ := os.Getwd()
		fullPath := filepath.Join(pwd, secretsPath)
		return nil, fmt.Errorf(missingClientSecretsMessage, fullPath)
	}

	cfg1 := new(Config)
	err = json.Unmarshal(data, &cfg1)
	if err != nil {
		return nil, err
	}

	var oCfg *oauth2.Config

	var cfg2 ClientConfig
	if cfg1.Web.ClientID != "" {
		cfg2 = cfg1.Web
	} else if cfg1.Installed.ClientID != "" {
		cfg2 = cfg1.Installed
	} else {
		return nil, errors.New("Client secrets file format not recognised")
	}

	oCfg = &oauth2.Config{
		ClientID:     cfg2.ClientID,
		ClientSecret: cfg2.ClientSecret,
		Scopes:       []string{youtube.YoutubeUploadScope, youtube.YoutubepartnerScope, youtube.YoutubeScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg2.AuthURI,
			TokenURL: cfg2.TokenURI,
		},
		RedirectURL: cfg2.RedirectURIs[0],
	}
	return oCfg, nil
}

func GetAuthURL(operationId string, secretsPath string) (string, error) {
	config, err := readConfig(secretsPath)
	if err != nil {
		msg := fmt.Sprintf("Cannot read configuration file: %v", err)
		return "", errors.New(msg)
	}

	return config.AuthCodeURL(operationId, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func VerifyCode(code string, tokenPath string, secretsPath string) error {
	config, err := readConfig(secretsPath)
	if err != nil {
		msg := fmt.Sprintf("Cannot read configuration file: %v", err)
		return errors.New(msg)
	}

	// Try to read the token from the cache file.
	// If an error occurs, do the three-legged OAuth flow because
	// the token is invalid or doesn't exist.
	tokenCache := CacheFile(tokenPath)
	token, err := tokenCache.Token()
	if err != nil {
		randState := fmt.Sprintf("st%d", time.Now().UnixNano())
		url := config.AuthCodeURL(randState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		fmt.Println("\n Auth.URL: ", url)

		token, err = config.Exchange(oauth2.NoContext, code)
		if err != nil {
			return err
		}
		err = tokenCache.PutToken(token)
		if err != nil {
			return err
		}
	}

	return nil
}

// Token retreives the token from the token cache
func (f CacheFile) Token() (*oauth2.Token, error) {
	file, err := os.Open(string(f))
	if err != nil {
		return nil, fmt.Errorf("CacheFile.Token: %s", err.Error())
	}
	defer file.Close()
	tok := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(tok); err != nil {
		return nil, fmt.Errorf("CacheFile.Token: %s", err.Error())
	}
	return tok, nil
}

// PutToken stores the token in the token cache
func (f CacheFile) PutToken(tok *oauth2.Token) error {
	file, err := os.OpenFile(string(f), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("CacheFile.PutToken: %s", err.Error())
	}
	if err := json.NewEncoder(file).Encode(tok); err != nil {
		file.Close()
		return fmt.Errorf("CacheFile.PutToken: %s", err.Error())
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("CacheFile.PutToken: %s", err.Error())
	}
	return nil
}
