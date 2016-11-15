package definition

//NonJSONItem used for storing non JSON content in mongo, by saving it in a the NonJSONContent of a wrapper json content
type NonJSONItem struct {
	Content string `json:"non_json_content"`
}
