/*
 * IceBreaker API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Organizer struct {
	// ID of User
	Id string `json:"id,omitempty"`
	// Organizer full name
	Name string `json:"name"`
	// Organizer's instagram
	Instagram string `json:"instagram,omitempty"`
	// Organizer's linkedin
	Linkedin string `json:"linkedin,omitempty"`
	// Link to Organizer's profile image
	ImageUrl string `json:"imageUrl"`
}