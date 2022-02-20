package server

import (
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/changelog"
)

func (s *Server) changesHandler(c echo.Context) error {
	hc := s.buildPageContext(c, "Recent changes")

	releases, err := changelog.RetrieveReleases()
	if err != nil {
		return err
	}
	hc.Add("releases", releases)
	return s.pages.Render(c, "changes", hc)
}
