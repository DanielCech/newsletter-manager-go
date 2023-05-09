/*
 * API for NewsletterManager
 *
 * This is testing Go project.
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

// Author object for creating a author
type CreateAuthorInput struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}
