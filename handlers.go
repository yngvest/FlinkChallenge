package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	history *OrdersHistory
}

type GetHistoryResult struct {
	Order   string        `json:"order_id"`
	History []Coordinates `json:"history"`
}

type GetHistoryParams struct {
	Max uint `form:"max"`
}

func (ctrl *Controller) PostOrderHistory(c *gin.Context) {
	id := c.Param("order_id")
	var coords Coordinates
	if err := c.ShouldBindJSON(&coords); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctrl.history.Add(id, coords)
}

func (ctrl *Controller) GetOrderHistory(c *gin.Context) {
	id := c.Param("order_id")
	var params GetHistoryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rsp := ctrl.history.Get(id, int(params.Max))
	if rsp == nil {
		c.Status(http.StatusNotFound)
		return
	}
	res := GetHistoryResult{
		Order:   id,
		History: rsp,
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *Controller) DeleteOrderHistory(c *gin.Context) {
	id := c.Param("order_id")
	// not sure if we should answer 200 or 404 if there is no such order
	ctrl.history.Delete(id)
}
