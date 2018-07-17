package rabbithole

import (
	"encoding/json"
	"net/http"
)

//
// GET /api/topic-permissions
//

// Example response:
//
//[{"user":"guest","vhost":"/","exchange":"exchange-name","write":".*","read":".*"}]

type TopicPermissionInfo struct {

	// Exchange name
	Exchange string `json:"exchange"`

	Username string `json:"username"`

	Vhost string `json:"vhost"`

	// Write permissions
	Write string `json:"write"`
	// Read permissions
	Read string `json:"read"`
}

//ListTopicPermissions returns topic permissions for all users and virtual hosts.
func (c *Client) ListTopicPermissions(vhost string) (rec []TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "topic-permissions")
	if err != nil {
		return []TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/users/{user}/topic-permissions
//

//ListTopicPermissionsOf a specific user.
func (c *Client) ListTopicPermissionsOf(username string) (rec []TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "users/"+PathEscape(username)+"/topic-permissions")
	if err != nil {
		return []TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/topic-permissions/{vhost}/{user}
//

//GetTopicPermissionsIn virtual host of user.
func (c *Client) GetTopicPermissionsIn(vhost, username string) (rec TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "topic-permissions/"+PathEscape(vhost)+"/"+PathEscape(username))
	if err != nil {
		return TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// PUT /api/topic-permissions/{vhost}/{user}
//

type TopicPermission struct {

	// Exchange name
	Exchange string `json:"exchange"`

	// Write permissions
	Write string `json:"write"`
	// Read permissions
	Read string `json:"read"`
}

//UpdateTopicPermissionsIn virtual host of user.
func (c *Client) UpdateTopicPermissionsIn(vhost, username string, permissions TopicPermission) (res *http.Response, err error) {
	body, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "topic-permissions/"+PathEscape(vhost)+"/"+PathEscape(username), body)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/topic-permissions/{vhost}/{user}
//

//ClearTopicPermissionsIn  virtual host of user.
func (c *Client) ClearTopicPermissionsIn(vhost, username string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "topic-permissions/"+PathEscape(vhost)+"/"+PathEscape(username), nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
