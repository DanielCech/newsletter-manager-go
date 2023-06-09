/*
 * API for NewsletterManager
 *
 * This is testing Go project.
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger
import (
	"time"
)

// Session object
type Session struct {
	AccessToken string `json:"accessToken"`
	AccessTokenExpiresAt *time.Time `json:"accessTokenExpiresAt"`
	RefreshToken string `json:"refreshToken"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt"`
}
