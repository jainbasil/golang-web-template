package rest

import (
	sentryEcho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang-web-template/internal/errors"
	"net/http"
)

func (s *Server) initMiddlewares() {
	s.Use(sentryEcho.New(sentryEcho.Options{
		Repanic: true,
	}))
	s.Use(middleware.Recover(), errorHandler(s))
}

func errorHandler(s *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				s.appContext.Logger.Sugar().Error(err)

				switch t := err.(type) {
				case *echo.HTTPError:
					return t
				case errors.NotAllowedError:
					return echo.NewHTTPError(http.StatusForbidden, t.Msg)
				case errors.NotFoundError:
					return echo.NewHTTPError(http.StatusNotFound, t.Msg)
				case errors.NotValidError:
					return echo.NewHTTPError(http.StatusBadRequest, t.Msg)
				case errors.UnknownError:
					return echo.NewHTTPError(http.StatusInternalServerError, t.Msg)
				case errors.InternalError:
					return echo.NewHTTPError(http.StatusInternalServerError, t.Msg)
				default:
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}
			return nil
		}
	}
}
