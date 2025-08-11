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
  "golang.org/x/oauth2/facebook"
)

// implements IOAuthService for Facebook OAuth
type FacebookOAuthService struct {
  config *oauth2.Config
}

func NewFacebookOAuthService() services.IOAuthService {
  clientID := os.Getenv("FACEBOOK_CLIENT_ID")
  clientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")
  redirectURI := os.Getenv("FACEBOOK_REDIRECT_URI")

  return &FacebookOAuthService{
    config: &oauth2.Config{
      ClientID:     clientID,
      ClientSecret: clientSecret,
      RedirectURL:  redirectURI,
      Scopes:       []string{"email", "public_profile"},
      Endpoint:     facebook.Endpoint,
    },
  }
}

// generates the Facebook OAuth authorization URL
func (f *FacebookOAuthService) GetAuthURL(state string) (string, error) {
  return f.config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// exchanges authorization code for access token
func (f *FacebookOAuthService) ExchangeCodeForToken(code string) (*models.OAuthToken, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  tok, err := f.config.Exchange(ctx, code)
  if err != nil {
    return nil, fmt.Errorf("failed to exchange code for token: %v", err)
  }
  out := &models.OAuthToken{
    AccessToken:  tok.AccessToken,
    RefreshToken: tok.RefreshToken,
    TokenType:    tok.TokenType,
  }
  if !tok.Expiry.IsZero() {
    out.ExpiresIn = int(time.Until(tok.Expiry).Seconds())
    if out.ExpiresIn < 0 {
      out.ExpiresIn = 0
    }
  }
  return out, nil
}

// fetches user information from Facebook Graph API
func (f *FacebookOAuthService) GetUserInfo(accessToken string) (*models.OAuthUserInfo, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
  httpClient := oauth2.NewClient(ctx, src)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.facebook.com/v18.0/me?fields=id,name,email,picture.type(large)", nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %v", err)
  }
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
    ID      string `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    Picture struct {
      Data struct {
        URL string `json:"url"`
      } `json:"data"`
    } `json:"picture"`
  }
  if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
    return nil, fmt.Errorf("failed to decode user info: %v", err)
  }

  return &models.OAuthUserInfo{
    ProviderID: userInfo.ID,
    Email:      userInfo.Email,
    Name:       userInfo.Name,
    Picture:    userInfo.Picture.Data.URL,
  }, nil
}

// returns the provider name
func (f *FacebookOAuthService) GetProviderName() string {
  return "facebook"
}
