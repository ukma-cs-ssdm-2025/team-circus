package group

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/responses"
	"go.uber.org/zap"
)

type listGroupMembersService interface {
	ListMembers(ctx context.Context, requesterUUID, groupUUID uuid.UUID) ([]*domain.GroupMember, error)
}

type addGroupMemberService interface {
	AddMember(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.GroupMember, error)
}

type updateGroupMemberService interface {
	UpdateMemberRole(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.GroupMember, error)
}

type removeGroupMemberService interface {
	RemoveMember(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID) error
}

func mapMemberToResponse(member *domain.GroupMember) responses.GroupMemberResponse {
	return responses.GroupMemberResponse{
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
		UserLogin: member.UserLogin,
		UserEmail: member.UserEmail,
	}
}

func mapMembersToResponse(members []*domain.GroupMember) []responses.GroupMemberResponse {
	result := make([]responses.GroupMemberResponse, len(members))
	for i, member := range members {
		result[i] = mapMemberToResponse(member)
	}
	return result
}

// NewListGroupMembersHandler returns a handler that lists members of a group the requester belongs to.
// @Summary List group members
// @Description Retrieve members of the specified group if the requester is part of it
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Success 200 {object} responses.GroupMembersResponse "Members retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid}/members [get]
func NewListGroupMembersHandler(service listGroupMembersService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, groupUUID, ok := resolveContextAndGroup(c)
		if !ok {
			return
		}

		members, err := service.ListMembers(c.Request.Context(), userUUID, groupUUID)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err != nil {
			logger.Error("failed to list group members", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
			return
		}

		response := responses.GroupMembersResponse{
			Members: mapMembersToResponse(members),
		}

		c.JSON(http.StatusOK, response)
	}
}

// NewAddGroupMemberHandler returns a handler that adds a new member to a group.
// @Summary Add group member
// @Description Add a user to the specified group as long as the requester is the author
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Param request body requests.AddMemberRequest true "Add member request"
// @Success 201 {object} responses.GroupMemberResponse "Member added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "User or group not found"
// @Failure 409 {object} map[string]interface{} "Member already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid}/members [post]
func NewAddGroupMemberHandler(service addGroupMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, groupUUID, ok := resolveContextAndGroup(c)
		if !ok {
			return
		}

		var req requests.AddMemberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("add member handler: failed to bind request: %w", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		member, err := service.AddMember(c.Request.Context(), userUUID, groupUUID, req.UserUUID, req.Role)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "member already exists"})
			return
		}
		if errors.Is(err, domain.ErrInvalidRole) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		if err != nil {
			logger.Error("failed to add group member", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member"})
			return
		}

		response := mapMemberToResponse(member)
		c.JSON(http.StatusCreated, response)
	}
}

// NewUpdateGroupMemberHandler returns a handler that updates a member role inside a group.
// @Summary Update group member role
// @Description Change member role if the requester is the author of the group
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Param user_uuid path string true "Member UUID"
// @Param request body requests.UpdateMemberRequest true "Role update request"
// @Success 200 {object} responses.GroupMemberResponse "Member updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Member or group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid}/members/{user_uuid} [put]
func NewUpdateGroupMemberHandler(service updateGroupMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, groupUUID, ok := resolveContextAndGroup(c)
		if !ok {
			return
		}

		memberUUIDParam := c.Param("user_uuid")
		memberUUID, err := uuid.Parse(memberUUIDParam)
		if err != nil {
			err = fmt.Errorf("update member handler: failed to parse member uuid: %w", err)
			logger.Error("failed to parse member uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		var req requests.UpdateMemberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update member handler: failed to bind request: %w", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		member, err := service.UpdateMemberRole(c.Request.Context(), userUUID, groupUUID, memberUUID, req.Role)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		if errors.Is(err, domain.ErrInvalidRole) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		if errors.Is(err, domain.ErrLastAuthor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change the last author"})
			return
		}
		if err != nil {
			logger.Error("failed to update member role", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member"})
			return
		}

		response := mapMemberToResponse(member)
		c.JSON(http.StatusOK, response)
	}
}

// NewRemoveGroupMemberHandler returns a handler that removes a member from a group.
// @Summary Remove group member
// @Description Remove a member from the specified group if the requester is the author
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Param user_uuid path string true "Member UUID"
// @Success 204 "Member removed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Member or group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid}/members/{user_uuid} [delete]
func NewRemoveGroupMemberHandler(service removeGroupMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, groupUUID, ok := resolveContextAndGroup(c)
		if !ok {
			return
		}

		memberUUIDParam := c.Param("user_uuid")
		memberUUID, err := uuid.Parse(memberUUIDParam)
		if err != nil {
			err = fmt.Errorf("remove member handler: failed to parse member uuid: %w", err)
			logger.Error("failed to parse member uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = service.RemoveMember(c.Request.Context(), userUUID, groupUUID, memberUUID)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		if errors.Is(err, domain.ErrLastAuthor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove the last author"})
			return
		}
		if err != nil {
			logger.Error("failed to remove group member", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func resolveContextAndGroup(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userUUIDValue, exists := c.Get("user_uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context missing"})
		return uuid.Nil, uuid.Nil, false
	}

	userUUID, ok := userUUIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return uuid.Nil, uuid.Nil, false
	}

	groupUUIDParam := c.Param("uuid")
	groupUUID, err := uuid.Parse(groupUUIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
		return uuid.Nil, uuid.Nil, false
	}

	return userUUID, groupUUID, true
}
