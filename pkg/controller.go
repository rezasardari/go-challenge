package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Service Service
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
}

func (cr Controller) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) GetUserCountBySegmentation(c *gin.Context) {
	var req UserCountBySegmentRequest

	if vErr := c.Bind(&req); vErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": vErr})
		return
	}

	result, sErr := cr.Service.GetUserCountBySegmentation(c, req)

	if sErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": sErr.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (cr Controller) StoreUserInSegment(c *gin.Context) {
	var req StoreUserSegmentRequest

	if vErr := c.Bind(&req); vErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": vErr})
		return
	}

	sErr := cr.Service.StoreUserSegment(c, req)

	if sErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": sErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (cr Controller) ArchiveExpiredData(c *gin.Context) {
	cr.Service.ArchiveExpiredSegment(c)
	c.JSON(http.StatusNoContent, gin.H{})
}
