/*
 * API for NewsletterManager
 *
 * This is testing Go project.
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

// Object for refreshing session
type RefreshSessionInput struct {
	RefreshToken string `json:"refreshToken"`
}