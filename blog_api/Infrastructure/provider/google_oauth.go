 package provider

import (
  "blog_api/Domain/contracts/services"
  "blog_api/Domain/models"
  "context"
  "encoding/json"
  "fmt"
  "net/http"
  "os"
  "time"

  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
)

// GoogleOAuthService implements IOAuthService for Google OAuth
type GoogleOAuthService struct {
  config *oauth2.Config
}

// NewGoogleOAuthService creates a new Google OAuth service
func NewGoogleOAuthService() services.IOAuthService {
  clientID := os.Getenv("GOOGLE_CLIENT_ID")
  clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
  redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
  scopes := []string{"openid", "email", "profile"}

  return &GoogleOAuthService{
    config: &oauth2.Config{
      ClientID:     clientID,
      ClientSecret: clientSecret,
      RedirectURL:  redirectURI,
      Scopes:       scopes,
      Endpoint:     google.Endpoint,
    },
  }
}

// GetAuthURL generates the Google OAuth authorization URL
func (g *GoogleOAuthService) GetAuthURL(state string) (string, error) {
  // Request offline access to receive a refresh token; prompt consent to ensure refresh
  url := g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
  return url, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (g *GoogleOAuthService) ExchangeCodeForToken(code string) (*models.OAuthToken, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  tok, err := g.config.Exchange(ctx, code)
  if err != nil {
    return nil, fmt.Errorf("failed to exchange code for token: %v", err)
  }

  // oauth2.Token contains expiry and refresh token if provided
  out := &models.OAuthToken{
    AccessToken:  tok.AccessToken,
    RefreshToken: tok.RefreshToken,
    TokenType:    tok.TokenType,
  }
  if !tok.Expiry.IsZero() {
    // seconds remaining until expiry (approx)
    out.ExpiresIn = int(time.Until(tok.Expiry).Seconds())
    if out.ExpiresIn < 0 {
      out.ExpiresIn = 0
    }
  }
  return out, nil
}

// GetUserInfo fetches user information from Google
func (g *GoogleOAuthService) GetUserInfo(accessToken string) (*models.OAuthUserInfo, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  // Use oauth2 client with the provided access token
  src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
  httpClient := oauth2.NewClient(ctx, src)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %v", err)
  }

  resp, err := httpClient.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch user info: %v", err)
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("failed to fetch user info with status: %d", resp.StatusCode)
  }

  var userInfo struct {
    ID      string `json:"id"`
    Email   string `json:"email"`
    Name    string `json:"name"`
    Picture string `json:"picture"`
  }
  if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
    return nil, fmt.Errorf("failed to decode user info: %v", err)
  }

  return &models.OAuthUserInfo{
    ProviderID: userInfo.ID,
    Email:      userInfo.Email,
    Name:       userInfo.Name,
    Picture:    userInfo.Picture,
  }, nil
}

// GetProviderName returns the provider name
func (g *GoogleOAuthService) GetProviderName() string {
  return "google"
}
