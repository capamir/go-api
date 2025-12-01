package handler

import (
	"strconv"

	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// TagHandler handles tag HTTP requests
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// ========================================
// PUBLIC ENDPOINTS
// ========================================

// GetAllTags lists all tags with pagination
// @Summary Get all tags
// @Tags tags
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(50)
// @Success 200 {object} service.TagListResponse
// @Router /tags [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Get tags
	response, err := h.tagService.GetAllTags(page, limit)
	if err != nil {
		logger.Error("Failed to get tags: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tags retrieved successfully", response)
}

// GetAllTagsList retrieves all tags without pagination
// @Summary Get all tags (no pagination)
// @Tags tags
// @Produce json
// @Success 200 {array} models.TagResponse
// @Router /tags/all [get]
func (h *TagHandler) GetAllTagsList(c *gin.Context) {
	tags, err := h.tagService.GetAllTagsList()
	if err != nil {
		logger.Error("Failed to get all tags: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "All tags retrieved successfully", tags)
}

// SearchTags searches tags by name
// @Summary Search tags
// @Tags tags
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(50)
// @Success 200 {object} service.TagListResponse
// @Router /tags/search [get]
func (h *TagHandler) SearchTags(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Search tags
	response, err := h.tagService.SearchTags(query, page, limit)
	if err != nil {
		logger.Error("Failed to search tags: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tags search completed", response)
}

// GetTagByID retrieves a single tag by ID
// @Summary Get tag by ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.TagResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid tag ID: %v", err)
		utils.RespondBadRequest(c, "Invalid tag ID")
		return
	}

	// Get tag
	tag, err := h.tagService.GetTagByID(uint(id))
	if err != nil {
		logger.Error("Failed to get tag %d: %v", id, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tag retrieved successfully", tag)
}

// GetTagBySlug retrieves a single tag by slug
// @Summary Get tag by slug
// @Tags tags
// @Produce json
// @Param slug path string true "Tag slug"
// @Success 200 {object} models.TagResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /tags/slug/{slug} [get]
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.RespondBadRequest(c, "Slug is required")
		return
	}

	// Get tag
	tag, err := h.tagService.GetTagBySlug(slug)
	if err != nil {
		logger.Error("Failed to get tag by slug %s: %v", slug, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tag retrieved successfully", tag)
}

// ========================================
// ADMIN ENDPOINTS (Protected)
// ========================================

// CreateTag creates a new tag (Admin only)
// @Summary Create tag
// @Tags admin/tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateTagRequest true "Tag data"
// @Success 201 {object} models.TagResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /admin/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to create tag")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	var req service.CreateTagRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid create tag request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Create tag
	tag, err := h.tagService.CreateTag(req)
	if err != nil {
		logger.Error("Failed to create tag: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "Tag created successfully", tag)
}

// UpdateTag updates an existing tag (Admin only)
// @Summary Update tag
// @Tags admin/tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param request body service.UpdateTagRequest true "Tag data"
// @Success 200 {object} models.TagResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to update tag")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid tag ID: %v", err)
		utils.RespondBadRequest(c, "Invalid tag ID")
		return
	}

	var req service.UpdateTagRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update tag request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update tag
	tag, err := h.tagService.UpdateTag(uint(id), req)
	if err != nil {
		logger.Error("Failed to update tag %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tag updated successfully", tag)
}

// DeleteTag deletes a tag (Admin only)
// @Summary Delete tag
// @Tags admin/tags
// @Security BearerAuth
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to delete tag")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid tag ID: %v", err)
		utils.RespondBadRequest(c, "Invalid tag ID")
		return
	}

	// Delete tag
	if err := h.tagService.DeleteTag(uint(id)); err != nil {
		logger.Error("Failed to delete tag %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Tag deleted successfully", nil)
}
