# Code Style Examples

Full session feature showing how each layer should look. Adapt names and fields
to the domain you are implementing.

## Models — `internal/models/session.go`

```go
var (
	ErrSessionNotFound = fmt.Errorf("session not found")
)

type Session struct {
	ID        int64     `json:"id"         pg:",pk"`
	UserID    int64     `json:"user_id"    pg:"user_id"`
	CreatedAt time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt time.Time `json:"updated_at" pg:"updated_at"`
}

type SessionParams struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	UserID  int64 `json:"user_id"`
}
```

## Presenters — `internal/presenters/session.go`

```go
type ListSessionParams struct {
	Page    int    `query:"current"`
	PerPage int    `query:"pageSize"`
	UserID  *int64 `query:"user_id"`
}

func (req *ListSessionParams) ToParams() *models.SessionParams {
	params := &models.SessionParams{
		Page:    req.Page,
		PerPage: req.PerPage,
	}
	if req.PerPage == 0 {
		params.PerPage = 20
	}
	if req.Page == 0 {
		params.Page = 1
	}
	if req.UserID != nil {
		params.UserID = *req.UserID
	}
	return params
}

type ListSessionsResponse struct {
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	PerPage int              `json:"per_page"`
	Data    []models.Session `json:"data"`
}

type CreateSessionRequest struct {
	UserID int64 `json:"user_id"`
}

func (req *CreateSessionRequest) ToCreateSession() *models.Session {
	return &models.Session{
		UserID: req.UserID,
	}
}

type PatchSessionRequest struct {
	UserID *int64 `json:"user_id"`
}

func (req *PatchSessionRequest) ToPatchSession() (*models.Session, []string) {
	data := &models.Session{}
	var fields []string

	if req.UserID != nil {
		data.UserID = *req.UserID
		fields = append(fields, "user_id")
	}

	return data, fields
}
```

## Repositories — `internal/repositories/session.go`

```go
type SessionRepository interface {
	Create(ctx context.Context, data *models.Session) (*models.Session, error)
	FindAll(ctx context.Context, params *models.SessionParams) (int, []models.Session, error)
	FindByID(ctx context.Context, id int64) (*models.Session, error)
	Patch(ctx context.Context, data *models.Session, columns ...string) (*models.Session, error)
	Update(ctx context.Context, data *models.Session) (*models.Session, error)
	Delete(ctx context.Context, id int64) error
}
```

## Providers — `internal/providers/session.go`

```go
func NewSessionRepository(db *pg.DB) repositories.SessionRepository {
	return &sessionRepository{db: db}
}

type sessionRepository struct {
	db *pg.DB
}

func (r *sessionRepository) Create(ctx context.Context, item *models.Session) (*models.Session, error) {
	_, err := r.db.WithContext(ctx).Model(item).Insert()
	if err != nil {
		slog.Error("Failed to create session", "error", err, "session", item)
		return nil, err
	}
	return item, nil
}

func (r *sessionRepository) FindByID(ctx context.Context, id int64) (*models.Session, error) {
	data := &models.Session{ID: id}
	err := r.db.WithContext(ctx).Model(data).WherePK().First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, models.ErrSessionNotFound
		}
		slog.Error("Failed to find session by ID", "error", err, "id", id)
		return nil, err
	}
	return data, nil
}

func (r *sessionRepository) FindAll(ctx context.Context, params *models.SessionParams) (int, []models.Session, error) {
	var sessions []models.Session
	query := r.db.WithContext(ctx).Model(&sessions)

	if params.UserID != 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

	perPage := params.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	page := params.Page
	if page < 1 {
		page = 1
	}

	query = query.Limit(perPage).Offset(perPage * (page - 1))

	count, err := query.SelectAndCount()
	if err != nil {
		slog.Error("Failed to find all sessions", "error", err, "params", params)
		return 0, nil, err
	}
	return count, sessions, nil
}

func (r *sessionRepository) Patch(ctx context.Context, item *models.Session, columns ...string) (*models.Session, error) {
	_, err := r.db.WithContext(ctx).Model(item).
		WherePK().
		Column(columns...).
		Update()
	if err != nil {
		slog.Error("Failed to patch session", "error", err, "session", item, "columns", columns)
		return nil, err
	}
	return item, nil
}

func (r *sessionRepository) Update(ctx context.Context, item *models.Session) (*models.Session, error) {
	_, err := r.db.WithContext(ctx).Model(item).WherePK().Update()
	if err != nil {
		slog.Error("Failed to update session", "error", err, "session", item)
		return nil, err
	}
	return item, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.WithContext(ctx).Model(&models.Session{ID: id}).WherePK().Delete()
	if err != nil {
		slog.Error("Failed to delete session", "error", err, "id", id)
		return err
	}
	return nil
}
```

## Services — `internal/services/session.go`

```go
type sessionService struct {
	repo repositories.SessionRepository
}

type SessionService interface {
	Create(ctx context.Context, session *models.Session) (*models.Session, error)
	GetByID(ctx context.Context, id int64) (*models.Session, error)
	GetAll(ctx context.Context, params *models.SessionParams) (int, []models.Session, error)
	Patch(ctx context.Context, session *models.Session, columns ...string) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) (*models.Session, error)
	Delete(ctx context.Context, id int64) error
}

func NewSessionService(repo repositories.SessionRepository) SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now
	return s.repo.Create(ctx, session)
}

func (s *sessionService) GetByID(ctx context.Context, id int64) (*models.Session, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *sessionService) GetAll(ctx context.Context, params *models.SessionParams) (int, []models.Session, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *sessionService) Patch(ctx context.Context, session *models.Session, columns ...string) (*models.Session, error) {
	session.UpdatedAt = time.Now()
	columns = append(columns, "updated_at")
	return s.repo.Patch(ctx, session, columns...)
}

func (s *sessionService) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	session.UpdatedAt = time.Now()
	return s.repo.Update(ctx, session)
}

func (s *sessionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
```

## Controllers — `internal/controllers/session.go`

```go
type sessionController struct {
	session services.SessionService
}

type SessionController interface {
	Routes(g *echo.Group)
}

func NewSessionController(session services.SessionService) SessionController {
	return &sessionController{session: session}
}

func (c *sessionController) Routes(g *echo.Group) {
	g.GET("/admin/session/:id", c.View, auth.RequiredAuth(), auth.AdminOnly())
	g.GET("/admin/sessions", c.List, auth.RequiredAuth(), auth.AdminOnly())
	g.POST("/admin/session", c.Create, auth.RequiredAuth(), auth.AdminOnly())
	g.PATCH("/admin/session/:id", c.Patch, auth.RequiredAuth(), auth.AdminOnly())
	g.DELETE("/admin/session/:id", c.Delete, auth.RequiredAuth(), auth.AdminOnly())
}

// View godoc
//
//	@Summary		Show a session
//	@Description	View session by ID
//	@Security		Authorization
//	@Param			Authorization	header	string	true	"Authentication header"
//	@tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Session ID"
//	@Success		200	{object}	models.Session
//	@Failure		500	{object}	echo.HTTPError
//	@Failure		400	{object}	echo.HTTPError
//	@Router			/admin/session/{id} [get]
func (ctl *sessionController) View(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("unable to parse session ID", "error", err, "param", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := ctl.session.GetByID(ctx, id)
	if err != nil {
		sentry.CaptureError(err, map[string]string{"category": "controller"})
		slog.Error("failed to get session by ID", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, session)
}

// List godoc
//
//	@Summary		List sessions
//	@Description	Get a list of sessions
//	@Security		Authorization
//	@Param			Authorization	header	string	true	"Authentication header"
//	@tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			current		query		int	false	"Current page"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200	{object}	presenters.ListSessionsResponse
//	@Failure		400	{object}	echo.HTTPError
//	@Router			/admin/sessions [get]
func (ctl *sessionController) List(c echo.Context) error {
	ctx := c.Request().Context()
	req := &presenters.ListSessionParams{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind list sessions request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	params := req.ToParams()
	total, sessions, err := ctl.session.GetAll(ctx, params)
	if err != nil {
		sentry.CaptureError(err, map[string]string{"category": "controller"})
		slog.Error("failed to list sessions", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(sessions) == 0 {
		sessions = []models.Session{}
	}

	return c.JSON(http.StatusOK, presenters.ListSessionsResponse{
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
		Data:    sessions,
	})
}

// Create godoc
//
//	@Summary		Create a session
//	@Description	Create a new session
//	@Security		Authorization
//	@Param			Authorization	header	string	true	"Authentication header"
//	@tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			session	body		presenters.CreateSessionRequest	true	"Session to create"
//	@Success		201	{object}	models.Session
//	@Failure		400	{object}	echo.HTTPError
//	@Router			/admin/session [post]
func (ctl *sessionController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	req := &presenters.CreateSessionRequest{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind create session request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	session, err := ctl.session.Create(ctx, req.ToCreateSession())
	if err != nil {
		sentry.CaptureError(err, map[string]string{"category": "controller"})
		slog.Error("failed to create session", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, session)
}

// Patch godoc
//
//	@Summary		Update a session
//	@Description	Update session details
//	@Security		Authorization
//	@Param			Authorization	header	string	true	"Authentication header"
//	@tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64	true	"Session ID"
//	@Param			session	body		presenters.PatchSessionRequest	true	"Session fields to update"
//	@Success		200	{object}	models.Session
//	@Failure		400	{object}	echo.HTTPError
//	@Router			/admin/session/{id} [patch]
func (ctl *sessionController) Patch(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("failed to parse session ID", "error", err, "param", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req := &presenters.PatchSessionRequest{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind patch session request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	patchData, updatedFields := req.ToPatchSession()
	patchData.ID = id
	updatedSession, err := ctl.session.Patch(ctx, patchData, updatedFields...)
	if err != nil {
		sentry.CaptureError(err, map[string]string{"category": "controller"})
		slog.Error("failed to patch session", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, updatedSession)
}

// Delete godoc
//
//	@Summary		Delete a session
//	@Description	Delete session by ID
//	@Security		Authorization
//	@Param			Authorization	header	string	true	"Authentication header"
//	@tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Session ID"
//	@Success		204
//	@Failure		400	{object}	echo.HTTPError
//	@Router			/admin/session/{id} [delete]
func (ctl *sessionController) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("failed to parse session ID", "error", err, "param", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctl.session.Delete(ctx, id); err != nil {
		sentry.CaptureError(err, map[string]string{"category": "controller"})
		slog.Error("failed to delete session", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
```
