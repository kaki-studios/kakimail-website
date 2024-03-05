package auth

import (
	"os"
  "kaki-studios/kakimail-website/user"
	"github.com/golang-jwt/jwt"
  "time"
  "net/http"
	"github.com/labstack/echo/v4"
)

const (
  accessTokenCookieName = "access-token"
	refreshTokenCookieName = "refresh-token"
)

func GetJWTSecret() string {
  return os.Getenv("JWT_SECRET")
}

func GetRefreshJWTSecret() string {
  return os.Getenv("JWT_REFRESH_SECRET")
}

type Claims struct {
  Name string `json:"name"`
  jwt.StandardClaims
}


// GenerateTokensAndSetCookies generates jwt token and saves it to the http-only cookie.
func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
  if err != nil {
    return err
  }

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)
    // We generate here a new refresh token and saving it to the cookie.
	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func generateRefreshToken(user *user.User) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *user.User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *user.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
    cookie.Path = "/"
    // Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// Purpose of this cookie is to store the user's name.
func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(c echo.Context, err error) error {
    // Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
// TokenRefresherMiddleware middleware, which refreshes JWT tokens if the access token is about to expire.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
        // If the user is not authenticated (no user token data in the context), don't do anything.
		if c.Get("user") == nil {
			return next(c)
		}
        // Gets user token from the context.
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
            // Gets the refresh token from the cookie.
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
                // Parses token and checks if it valid.
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
                    // If everything is good, update tokens.
					_ = GenerateTokensAndSetCookies(&user.User{
						Name:  claims.Name,
					}, c)
				}
			}
		}

		return next(c)
	}
}
