package controllers

import (
	"app/src/models"
)

// TransportController operations for /Transport
type TransportController struct {
	controller
}

// Get
// @Title Get
// @Description get Transport by Id
// @Param	id	path 	int64	true	"id"
// @Success 200 {object} responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure 400 :id is empty
// @router /:id [get]
func (c *TransportController) Get() {

}

// Post
// @Title Post
// @Description create Transport
// @Security	api_key
// @Param	body	body	requests.TransportPostRequest "transport info"
// @Success 200	{object}	responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure 400 body is invalid
// @Failure 401 unauthorized
// @router / [post]
func (c *TransportController) Post() {

}

// Put
// @Title Put
// @Description	update the Transport
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Param	body	body	requests.TransportPutRequest "transport info"
// @Success 200	{object}	responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure	403	user is not owner
// @router /:id [put]
func (c *TransportController) Put() {

}

// Delete
// @Title Delete
// @Description delete the Transport
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Success 200	{object} responses.TransportDeleteResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure	403	user is not owner
// @router /:id [delete]
func (c *TransportController) Delete() {

}

func (c *TransportController) findModel(id int64) *models.Transport {
	m := &models.Transport{Id: id}
	if err := models.Get(m); err != nil {
		c.responseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
