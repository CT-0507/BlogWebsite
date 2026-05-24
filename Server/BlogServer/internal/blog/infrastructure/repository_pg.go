package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewBlogRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *BlogRepository {
	return &BlogRepository{
		pool:   pool,
		mapper: mapper,
	}
}

func (r *BlogRepository) Create(c context.Context, blog *domain.Blog) (*domain.Blog, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	marshalledContent, _ := json.Marshal(blog.ContentJson)

	newBlog, err := q.CreateBlog(c, blogdb.CreateBlogParams{
		AuthorID:     blog.AuthorID,
		Title:        blog.Title,
		UrlSlug:      blog.URLSlug,
		ContentJson:  marshalledContent,
		ContentText:  blog.ContentText,
		ThumbnailUrl: utils.GetTextTypeFromNullableString(blog.ThumbnailUrl),
		CreatedBy:    blog.AuthorID,
		UpdatedBy:    blog.AuthorID,
	})

	if err != nil {
		return nil, err
	}

	return r.mapper.BlogDTOToBlog(&newBlog), nil
}

func (r *BlogRepository) GetFindAllCount(c context.Context, title, content, author *string) (int64, error) {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.GetListBlogsCount(c, blogdb.GetListBlogsCountParams{
		Title:             utils.GetTextTypeFromNullableString(title),
		Content:           utils.GetTextTypeFromNullableString(content),
		AuthorDisplayName: utils.GetTextTypeFromNullableString(author),
	})
}

func (r *BlogRepository) FindAll(c context.Context, title, content, author, sortBy, sortDir *string, offset, limit int32) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogs(c, blogdb.ListBlogsParams{
		Title:             utils.GetTextTypeFromNullableString(title),
		Content:           utils.GetTextTypeFromNullableString(content),
		AuthorDisplayName: utils.GetTextTypeFromNullableString(author),
		SortBy:            utils.GetTextTypeFromNullableString(sortBy),
		SortDir:           utils.GetTextTypeFromNullableString(sortDir),
		Offset:            offset,
		Limit:             limit,
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *r.mapper.ListBlogsRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) ListAuthorBlogsByAuthorID(c context.Context, authorID string) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogsByAuthor(c, blogdb.ListBlogsByAuthorParams{
		AuthorID: authorID,
		Status:   "active",
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *r.mapper.ListAuthorBlogsByAuthorIDRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) ListAuthorBlogsBySlug(c context.Context, slug string) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogsByAuthorSlug(c, blogdb.ListBlogsByAuthorSlugParams{
		Slug:   slug,
		Status: "active",
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *r.mapper.ListAuthorBlogsRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) FindByID(c context.Context, id int64) (*domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetBlog(c, id)
	if err != nil {
		return nil, err
	}
	return r.mapper.GetBlogRowDTOToBlogWithAuthorData(&row), nil
}

func (r *BlogRepository) FindByUrlSlug(c context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	if userID != nil {
		row, err := q.GetBlogWithUserReaction(c, blogdb.GetBlogWithUserReactionParams{
			UserID:  *userID,
			UrlSlug: slug,
		})
		if err != nil {
			return nil, err
		}
		return r.mapper.GetBlogWithReactionDTOToBlogWithAuthorData(&row), nil
	} else {
		row, err := q.GetBlogByUrlSlug(c, slug)
		if err != nil {
			return nil, err
		}
		return r.mapper.GetBlogRowByUrlSlugDTOToBlogWithAuthorData(&row), nil
	}
}

// func (r *blogRepository) Update(blog *Blog, q *blogdb.Queries) error {
// 	query := `UPDATE blogs SET name=$1, email=$2 WHERE id=$3`
// 	_, err := r.db.Exec(context.Background(), query, blog.Author, blog.Content, blog.ID)
// 	return err
// }

func (r *BlogRepository) Delete(c context.Context, id int64, userID string) (*int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	deletedId, err := q.DeleteBlog(c, blogdb.DeleteBlogParams{
		DeletedBy: pgtype.Text{
			String: userID,
			Valid:  true,
		},
		BlogID: id,
	})
	if err != nil {
		return nil, err
	}
	return &deletedId, nil
}

func (r *BlogRepository) CreateUserIDAuthorProfileIDCacheRecord(c context.Context, userID string, authorID string, slug string, displayName string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.CreateUserAuthorProfileIDCacheRecord(c, blogdb.CreateUserAuthorProfileIDCacheRecordParams{
		UserID:      userID,
		AuthorID:    authorID,
		Slug:        slug,
		DisplayName: displayName,
	})
}

func (r *BlogRepository) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.VerifyAuthorIDByUserID(c, userID)
}

func (r *BlogRepository) UpdateBlogStatusForDeletedAuthor(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateBlogStatusForDeletedAuthor(c, authorID)
}

func (r *BlogRepository) DeleteAuthorHardDeletedBlogs(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteAuthorHardDeletedBlogs(c, authorID)
}

func (r *BlogRepository) DeleteAuthorCache(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteAuthorCache(c, authorID)
}

func (r *BlogRepository) MarkAuthorCacheAsDeleted(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.MarkAuthorCacheAsDeleted(c, authorID)
}

func (r *BlogRepository) RestoreBlog(c context.Context, blogID int64, PreviousStatus string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.RestoreBlog(c, blogdb.RestoreBlogParams{
		BlogID: blogID,
		Status: PreviousStatus,
	})
}

func (r *BlogRepository) GetAuthorProfileByUserID(c context.Context, userID string) (*domain.AuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	author, err := q.GetAuthorCacheByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	return r.mapper.MapBlogsIdxUserAuthorProfileToAuthorProfile(&author), nil
}

func (r *BlogRepository) UpdateBlogReactionCount(c context.Context, blogID int64, transition repository.ReactionTransition) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	var likeDelta int64 = 0
	var dislikeDelta int64 = 0

	switch transition {
	case repository.AddLike:
		likeDelta++
	case repository.AddDislike:
		dislikeDelta++
	case repository.LikeToDislike:
		likeDelta--
		dislikeDelta++
	case repository.DislikeToLike:
		likeDelta++
		dislikeDelta--
	}

	return q.UpdateBlogReactionCount(c, blogdb.UpdateBlogReactionCountParams{
		LikeCount:    likeDelta,
		DislikeCount: dislikeDelta,
		BlogID:       blogID,
	})
}

func (r *BlogRepository) UpdateBlogRankingTable(c context.Context) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateBlogRankingResult(c)
}

func (r *BlogRepository) TruncateBlogRankingTable(c context.Context) error {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.TruncateBlogRankingTable(c)
}

func (r *BlogRepository) GetRankingBlogsByType(c context.Context, searchType string, offset, limit int32, shouldGetAll bool, sortBy, sortDir string) ([]domain.RankingBlogData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListRankingTable(c, blogdb.ListRankingTableParams{
		GetAll:  shouldGetAll,
		Offset:  offset,
		Limit:   limit,
		Type:    searchType,
		SortBy:  sortBy,
		SortDir: sortDir,
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.RankingBlogData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *r.mapper.MapDBListRankingRowToRankingBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) UpdateViewCount(c context.Context, blogID int64) error {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateViewCount(c, blogID)
}

func (r *BlogRepository) GetWeeksViews(c context.Context, blogID int64, numberOfWeeks int32) ([]domain.WeekViewData, error) {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.GetWeeksViews(c, blogdb.GetWeeksViewsParams{
		BlogID:       blogID,
		NumberOfWeek: numberOfWeeks,
	})
	if err != nil {
		return nil, err
	}
	var result []domain.WeekViewData
	for _, value := range rows {
		v := value
		result = append(result, domain.WeekViewData{
			WeekStart: v.WeekStart.Time.String(),
			Views:     v.WeeklyViews,
		})
	}
	return result, nil
}

func (r *BlogRepository) GetDaysViews(c context.Context, blogID int64, numberOfDays int32) ([]domain.DateViewData, error) {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.GetDaysView(c, blogdb.GetDaysViewParams{
		BlogID:       blogID,
		NumberOfDays: numberOfDays,
	})
	if err != nil {
		return nil, err
	}
	var result []domain.DateViewData
	for _, value := range rows {
		v := value
		result = append(result, domain.DateViewData{
			Date:  v.Date.Time.String(),
			Views: v.Views,
		})
	}
	return result, nil
}

func (r *BlogRepository) UpdateBlogReportCount(c context.Context, blogID int64, delta int64) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateBlogReportCount(c, blogdb.UpdateBlogReportCountParams{
		BlogID: blogID,
		Delta:  delta,
	})
}

func (r *BlogRepository) InsertBlogReport(c context.Context, report *domain.BlogReport) (*domain.BlogReport, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	inserted, err := q.InsertBlogReport(c, blogdb.InsertBlogReportParams{
		BlogID:          report.BlogID,
		UserID:          report.UserID,
		UserDisplayName: report.UserDisplayName,
		Reason:          report.Reason,
	})
	if err != nil {
		return nil, err
	}

	return r.mapper.MapDBReportToBlogReport(&inserted), nil
}

func (r *BlogRepository) DeleteBlogReportByID(c context.Context, reportID int64) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteBlogReport(c, reportID)
}

func (r *BlogRepository) GetBlogReportsByBlogID(c context.Context, blogID int64) ([]domain.BlogReport, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.GetBlogReportByBlogID(c, blogID)
	if err != nil {
		return nil, err
	}

	var reports []domain.BlogReport
	for _, value := range rows {
		v := value
		reports = append(reports, *r.mapper.MapDBReportToBlogReport(&v))
	}

	return reports, nil
}

func (r *BlogRepository) UpdateBlogStatus(c context.Context, blogID int64, status string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateBlogStatus(c, blogdb.UpdateBlogStatusParams{
		BlogID: blogID,
		Status: status,
	})
}

func (r *BlogRepository) GetAuthorDashboardViewMetrics(c context.Context, authorID string, userID *string) (*domain.AuthorDashboardViewMetrics, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetTodayViewAcrossAllContentByAuthorID(c, blogdb.GetTodayViewAcrossAllContentByAuthorIDParams{
		AuthorID: authorID,
		IsAdmin:  userID == nil,
		UserID:   utils.GetEmptyStringOnNullStringPtr(userID),
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthorDashboardViewMetrics{
		TodayViews:     row.TodayViews,
		YesterdayViews: row.YesterdayViews,
		ThisWeekViews:  row.ThisWeekViews,
		LastWeekViews:  row.LastWeekViews,
	}, nil
}

func (r *BlogRepository) GetAuthorDashboardReactionMetrics(c context.Context, authorID string, userID *string) (*domain.AuthorDashboardReactionMetrics, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetReactionCountByAuthorID(c, blogdb.GetReactionCountByAuthorIDParams{
		AuthorID: authorID,
		IsAdmin:  userID == nil,
		UserID:   utils.GetEmptyStringOnNullStringPtr(userID),
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthorDashboardReactionMetrics{
		TodayLikes:        row.TodayLikes,
		TodayDislikes:     row.TodayDislikes,
		YesterdayLikes:    row.YesterdayLikes,
		YesterdayDislikes: row.YesterdayDislikes,
		ThisWeekLikes:     row.ThisWeekLikes,
		ThisWeekDislikes:  row.ThisWeekDislikes,
		LastWeekLikes:     row.LastWeekLikes,
		LastWeekDislikes:  row.LastWeekDislikes,
	}, nil
}
