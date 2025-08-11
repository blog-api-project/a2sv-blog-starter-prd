package provider

import (
	"blog_api/Domain/contracts/services"
	"blog_api/Domain/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// implements IOAuthService for GitHub OAuth
type GitHubOAuthService struct {
    config *oauth2.Config
}


func NewGitHubOAuthService() services.IOAuthService {
    clientID := os.Getenv("GITHUB_CLIENT_ID")
    clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
    redirectURI := os.Getenv("GITHUB_REDIRECT_URI")

    return &GitHubOAuthService{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURI,
            Scopes:       []string{"read:user", "user:email"},
            Endpoint: oauth2.Endpoint{
                AuthURL:  "https://github.com/login/oauth/authorize",
                TokenURL: "https://github.com/login/oauth/access_token",
            },
        },
    }
}

// generates the GitHub OAuth authorization URL
func (g *GitHubOAuthService) GetAuthURL(state string) (string, error) {
    return g.config.AuthCodeURL(state), nil
}

// exchanges authorization code for access token
func (g *GitHubOAuthService) ExchangeCodeForToken(code string) (*models.OAuthToken, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tok, err := g.config.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("failed to exchange code for token: %v", err)
    }
    // GitHub typically returns only access token; no refresh, no expiry
    return &models.OAuthToken{
        AccessToken:  tok.AccessToken,
        RefreshToken: tok.RefreshToken,
        ExpiresIn:    0,
        TokenType:    tok.TokenType,
    }, nil
}

// fetches user information from GitHub
func (g *GitHubOAuthService) GetUserInfo(accessToken string) (*models.OAuthUserInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
    httpClient := oauth2.NewClient(ctx, src)

    // Fetch profile
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user info: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to fetch user info with status: %d, body: %s", resp.StatusCode, string(body))
    }

    var userInfo struct {
        ID        int    `json:"id"`
        Login     string `json:"login"`
        Name      string `json:"name"`
        Email     string `json:"email"`
        AvatarURL string `json:"avatar_url"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return nil, fmt.Errorf("failed to decode user info: %v", err)
    }


    // If email is missing, load from /user/emails
    if userInfo.Email == "" {
        req2, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
        if err == nil {
            req2.Header.Set("Accept", "application/vnd.github.v3+json")
            if resp2, err2 := httpClient.Do(req2); err2 == nil && resp2.StatusCode == http.StatusOK {
                defer resp2.Body.Close()
                var emails []struct{
                    Email   string `json:"email"`
                    Primary bool   `json:"primary"`
                }
                if err := json.NewDecoder(resp2.Body).Decode(&emails); err == nil {
                    for _, e := range emails {
                        if e.Primary {
                            userInfo.Email = e.Email
                            break
                        }
                    }
                }
            }
        }
    }

    return &models.OAuthUserInfo{
        ProviderID: fmt.Sprintf("%d", userInfo.ID),
        Email:      userInfo.Email,
        Name:       userInfo.Name,
        Picture:    userInfo.AvatarURL,
    }, nil
}

// returns the provider name
func (g *GitHubOAuthService) GetProviderName() string {
  return "github"
}
