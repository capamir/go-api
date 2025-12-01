package service

import (
	"errors"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
)

// TagService handles tag business logic
type TagService struct {
	tagRepo *repository.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// CreateTagRequest represents tag creation data
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// UpdateTagRequest represents tag update data
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// TagListResponse represents paginated tag list
type TagListResponse struct {
	Tags       []models.TagResponse `json:"tags"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

// ========================================
// CREATE TAG
// ========================================

// CreateTag creates a new tag
func (s *TagService) CreateTag(req CreateTagRequest) (*models.TagResponse, error) {
	// Generate slug from name
	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	exists, err := s.tagRepo.SlugExists(slug)
	if err != nil {
		logger.Error("Failed to check slug existence: %v", err)
		return nil, errors.New("failed to create tag")
	}
	if exists {
		return nil, errors.New("tag with similar name already exists")
	}

	// Create tag
	tag := &models.Tag{
		Name: req.Name,
		Slug: slug,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		logger.Error("Failed to create tag: %v", err)
		return nil, errors.New("failed to create tag")
	}

	logger.Success("Tag created successfully: %s (slug: %s)", tag.Name, tag.Slug)

	response := tag.ToResponse()
	return &response, nil
}

// ========================================
// GET TAG
// ========================================

// GetTagByID retrieves a tag by ID
func (s *TagService) GetTagByID(id uint) (*models.TagResponse, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get tag %d: %v", id, err)
		return nil, errors.New("failed to get tag")
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	response := tag.ToResponse()
	return &response, nil
}

// GetTagBySlug retrieves a tag by slug
func (s *TagService) GetTagBySlug(slug string) (*models.TagResponse, error) {
	tag, err := s.tagRepo.GetBySlug(slug)
	if err != nil {
		logger.Error("Failed to get tag by slug %s: %v", slug, err)
		return nil, errors.New("failed to get tag")
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	response := tag.ToResponse()
	return &response, nil
}

// ========================================
// LIST TAGS
// ========================================

// GetAllTags retrieves all tags with pagination
func (s *TagService) GetAllTags(page, limit int) (*TagListResponse, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50 // Higher default for tags (usually many tags)
	}

	offset := (page - 1) * limit

	tags, total, err := s.tagRepo.GetAll(limit, offset)
	if err != nil {
		logger.Error("Failed to get tags: %v", err)
		return nil, errors.New("failed to get tags")
	}

	// Convert to response DTOs
	tagResponses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = tag.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &TagListResponse{
		Tags:       tagResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetAllTagsList retrieves all tags without pagination (for dropdowns)
func (s *TagService) GetAllTagsList() ([]models.TagResponse, error) {
	tags, err := s.tagRepo.GetAllTags()
	if err != nil {
		logger.Error("Failed to get all tags: %v", err)
		return nil, errors.New("failed to get tags")
	}

	// Convert to response DTOs
	tagResponses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = tag.ToResponse()
	}

	return tagResponses, nil
}

// SearchTags searches tags by name
func (s *TagService) SearchTags(query string, page, limit int) (*TagListResponse, error) {
	// Validate query
	if query == "" {
		return nil, errors.New("search query is required")
	}

	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	tags, total, err := s.tagRepo.Search(query, limit, offset)
	if err != nil {
		logger.Error("Failed to search tags: %v", err)
		return nil, errors.New("failed to search tags")
	}

	// Convert to response DTOs
	tagResponses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = tag.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &TagListResponse{
		Tags:       tagResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ========================================
// UPDATE TAG
// ========================================

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(id uint, req UpdateTagRequest) (*models.TagResponse, error) {
	// Get existing tag
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get tag %d: %v", id, err)
		return nil, errors.New("failed to update tag")
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	// Update name and slug
	newSlug := utils.GenerateSlug(req.Name)

	// Check if new slug already exists (excluding current tag)
	existingTag, err := s.tagRepo.GetBySlug(newSlug)
	if err != nil {
		logger.Error("Failed to check slug existence: %v", err)
		return nil, errors.New("failed to update tag")
	}
	if existingTag != nil && existingTag.ID != id {
		return nil, errors.New("tag with similar name already exists")
	}

	tag.Name = req.Name
	tag.Slug = newSlug

	// Save changes
	if err := s.tagRepo.Update(tag); err != nil {
		logger.Error("Failed to update tag %d: %v", id, err)
		return nil, errors.New("failed to update tag")
	}

	logger.Success("Tag updated successfully: %s (id: %d)", tag.Name, tag.ID)

	response := tag.ToResponse()
	return &response, nil
}

// ========================================
// DELETE TAG
// ========================================

// DeleteTag soft deletes a tag
func (s *TagService) DeleteTag(id uint) error {
	// Get tag
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get tag %d: %v", id, err)
		return errors.New("failed to delete tag")
	}
	if tag == nil {
		return errors.New("tag not found")
	}

	// TODO: Check if tag is used by products (when product model is ready)
	// For now, just delete the tag

	// Delete tag
	if err := s.tagRepo.Delete(id); err != nil {
		logger.Error("Failed to delete tag %d: %v", id, err)
		return errors.New("failed to delete tag")
	}

	logger.Success("Tag deleted successfully: %s (id: %d)", tag.Name, id)
	return nil
}
