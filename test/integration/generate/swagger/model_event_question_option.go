/*
 * IceBreaker API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

// Option for a question that user can choose
type EventQuestionOption struct {
	Id           string `json:"id"`
	Option       string `json:"option,omitempty"`
	GameQuestion string `json:"gameQuestion,omitempty"`
}