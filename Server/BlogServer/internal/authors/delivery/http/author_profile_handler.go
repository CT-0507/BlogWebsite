package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type AuthorDiscoveryUsecases interface {
	SearchAuthor()
}

type AuthorIdentityUsecases interface {
	CreateAuthor(ctx context.Context, fileParams *domain.CreateUserFileStorageParams, author *domain.AuthorProfile, userID string, createdBy string) error
	GetAuthorProfileByID(ctx context.Context, authorID string) (*domain.AuthorProfile, error)
	GetAuthorProfileBySlug(ctx context.Context, slug string) (*domain.AuthorProfile, error)
	ListAuthorProfiles(ctx context.Context, page int64, limit int64) (*[]domain.AuthorProfile, error)
	DeleteAuthorProfile(ctx context.Context, authorID string, deletedBy string) error
	HardDeleteAuthorProfile(ctx context.Context, authorID string) error
	UpdateAuthorSlug(ctx context.Context, authorID string, slug string, updatedBy string) error
	UpdateAuthorStatus(ctx context.Context, authorID string, status string, updatedBy string) error
	OnBlogCreated(ctx context.Context, evt *messaging.OutboxEvent) error
}

type AuthorSocialUsecases interface {
	SetFeatureBlogs(ctx context.Context, authorID string, blogIDs []string) error
	GetFeatureBlogsByAuthorID(ctx context.Context, slug string) ([]string, error)
}

type AuthorProfileUsecases interface {
}

type AuthorFollowerUsecases interface {
	FollowAuthor(ctx context.Context, userID string, authorID string) error
	UnfollowAuthor(ctx context.Context, userID string, authorID string) error
	GetAuthorFollowers(ctx context.Context, slug string, page int64, limit int64) ([]string, error)
	GetFollowedAuthors(ctx context.Context, userID string, page int64, limit int64) ([]string, error)
	OnAuthorFollowerCountChanged(ctx context.Context, evt *messaging.OutboxEvent) error
}

type AuthorProfileHandler struct {
	authorDiscoveryUseCases AuthorDiscoveryUsecases
	authorIdentityUsecases  AuthorIdentityUsecases
	authorSocialUsecases    AuthorSocialUsecases
	authorProfileUsecases   AuthorProfileUsecases
	authorFollowerUsecases  AuthorFollowerUsecases
}

func NewAuthorProfileHandler(
	authorDiscoveryUseCases AuthorDiscoveryUsecases,
	authorIdentityUsecases AuthorIdentityUsecases,
	authorSocialUsecases AuthorSocialUsecases,
	authorProfileUsecases AuthorProfileUsecases,
	authorFollowerUsecases AuthorFollowerUsecases,
) *AuthorProfileHandler {
	return &AuthorProfileHandler{
		authorDiscoveryUseCases: authorDiscoveryUseCases,
		authorIdentityUsecases:  authorIdentityUsecases,
		authorSocialUsecases:    authorSocialUsecases,
		authorProfileUsecases:   authorProfileUsecases,
		authorFollowerUsecases:  authorFollowerUsecases,
	}
}

// Identity Usecase

func (h *AuthorProfileHandler) createAuthorProfile(c *gin.Context) {
	var author CreateAuthorRequest
	if err := c.ShouldBind(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, author); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var fileParams *domain.CreateUserFileStorageParams = nil
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer file.Close()
		ext := filepath.Ext(fileHeader.Filename)
		fileName := ulid.Make().String() + ext
		contentType := fileHeader.Header.Get("Content-Type")

		fileParams = &domain.CreateUserFileStorageParams{
			File:        file,
			FileName:    fileName,
			ContentType: contentType,
		}
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := h.authorIdentityUsecases.CreateAuthor(ctx, fileParams, &domain.AuthorProfile{
		DisplayName: author.DisplayName,
		Bio:         author.Bio,
		Avatar:      author.Avatar,
		Slug:        author.Slug,
		SocialLink:  author.SocialLink,
		Email:       author.Email,
	}, userID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &gin.H{
		"message": "Successfully created",
	})
}

func (h *AuthorProfileHandler) getAuthorProfileByID(c *gin.Context) {

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	authorProfile, err := h.authorIdentityUsecases.GetAuthorProfileByID(ctx, authorProfileID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Author profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, authorProfile)
}

func (h *AuthorProfileHandler) getAuthorProfileBySlug(c *gin.Context) {

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "slug"),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	authorProfile, err := h.authorIdentityUsecases.GetAuthorProfileBySlug(ctx, slug)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Author profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, authorProfile)
}

func (h *AuthorProfileHandler) listAuthorProfiles(c *gin.Context) {

	var page int64 = 1
	var limit int64 = 20
	pageStr, exists := c.GetQuery("page")
	limitStr, exists := c.GetQuery("limit")
	if exists {
		value, err := strconv.ParseInt(pageStr, 10, 64)
		limitInt, err := strconv.ParseInt(limitStr, 10, 64)
		if err == nil {
			page = value
			limit = limitInt
		}
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	authorProfiles, err := h.authorIdentityUsecases.ListAuthorProfiles(ctx, page, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Author profiles not found",
		})
		return
	}

	c.JSON(http.StatusOK, authorProfiles)
}

func (h *AuthorProfileHandler) deleteAuthorProfile(c *gin.Context) {

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	err = h.authorIdentityUsecases.DeleteAuthorProfile(ctx, authorProfileID, userID)
	if err != nil {
		log.Println(err.Error())
		var authorNotFoundErr *domain.ErrAuthorNotFound
		if errors.As(err, &authorNotFoundErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete author",
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
	})
}

func (h *AuthorProfileHandler) hardDeleteAuthorProfile(c *gin.Context) {

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	err := h.authorIdentityUsecases.HardDeleteAuthorProfile(ctx, authorProfileID)
	if err != nil {
		log.Println(err.Error())
		var authorNotFoundErr *domain.ErrAuthorNotFound
		if errors.As(err, &authorNotFoundErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hard delete author",
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
	})
}

func (h *AuthorProfileHandler) updateAuthorSlug(c *gin.Context) {

	var body UpdateAuthorSlugRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	err = h.authorIdentityUsecases.UpdateAuthorSlug(ctx, authorProfileID, body.Slug, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update author slug",
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
	})
}

func (h *AuthorProfileHandler) updateAuthorStatus(c *gin.Context) {

	var body UpdateAuthorStatusRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	err = h.authorIdentityUsecases.UpdateAuthorStatus(ctx, authorProfileID, body.Status, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update author status",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
	})
}

// Follower system
func (h *AuthorProfileHandler) followAuthor(c *gin.Context) {

	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	if err = h.authorFollowerUsecases.FollowAuthor(ctx, userID, authorProfileID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to follow author",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Okay",
	})
}

func (h *AuthorProfileHandler) unfollowAuthor(c *gin.Context) {
	authorProfileID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	if err = h.authorFollowerUsecases.UnfollowAuthor(ctx, userID, authorProfileID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to unfollow author",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Okay",
	})
}

func (h *AuthorProfileHandler) getAuthorFollowers(c *gin.Context) {

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "slug"),
		})
		return
	}

	var page int64 = 1
	var limit int64 = 20
	pageStr, exists := c.GetQuery("page")
	limitStr, exists := c.GetQuery("limit")
	if exists {
		value, err := strconv.ParseInt(pageStr, 10, 64)
		limitInt, err := strconv.ParseInt(limitStr, 10, 64)
		if err == nil {
			page = value
			limit = limitInt
		}
	}
	log.Println(page)
	log.Println(limit)

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	userIds, err := h.authorFollowerUsecases.GetAuthorFollowers(ctx, slug, page, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get author followers",
		})
		return
	}

	c.JSON(http.StatusOK, userIds)
}

func (h *AuthorProfileHandler) getFollowedAuthors(c *gin.Context) {

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userID not found",
		})
		return
	}

	var page int64 = 1
	var limit int64 = 20
	pageStr, exists := c.GetQuery("page")
	limitStr, exists := c.GetQuery("limit")
	if exists {
		value, err := strconv.ParseInt(pageStr, 10, 64)
		limitInt, err := strconv.ParseInt(limitStr, 10, 64)
		if err == nil {
			page = value
			limit = limitInt
		}
	}
	log.Println(page)
	log.Println(limit)

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	authors, err := h.authorFollowerUsecases.GetFollowedAuthors(ctx, userID, page, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get followed authors",
		})
		return
	}

	if len(authors) == 0 {
		c.JSON(http.StatusOK, &gin.H{
			"message": "Ok",
			"length":  0,
			"authors": []string{},
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
		"length":  len(authors),
		"authors": authors,
	})
}

// Author Profile Presentation
func (h *AuthorProfileHandler) setAuthorFeaturedBlogs(c *gin.Context) {

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	var blogIds []string
	if err := c.ShouldBindJSON(&blogIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	err := h.authorSocialUsecases.SetFeatureBlogs(ctx, slug, blogIds)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to set author's featured blogs",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Ok",
	})
}

func (h *AuthorProfileHandler) getAuthorFeaturedBlogs(c *gin.Context) {

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "author profile id"),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	blogIds, err := h.authorSocialUsecases.GetFeatureBlogsByAuthorID(ctx, slug)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get author's featured blogs",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"blogIds": blogIds,
	})
}
