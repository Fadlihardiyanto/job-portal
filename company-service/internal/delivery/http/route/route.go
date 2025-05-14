package route

import (
	"company-service/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	CompanyController  *http.CompanyController
	JobsController     *http.JobsController
	UserJobsController *http.UserJobsController
	AuthMiddleware     fiber.Handler
}

func (c *RouteConfig) SetupRoutes() {
	api := c.App.Group("/api/v1")

	company := api.Group("/company")
	company.Post("/", c.CompanyController.CreateCompany)
	company.Get("/", c.CompanyController.GetAllCompanies)
	company.Get("/:id", c.CompanyController.FindByID)
	company.Put("/:id", c.CompanyController.UpdateCompany)
	company.Put("/:id/access/:user_access", c.CompanyController.UpdateCompanyByIDAndUserAccess)

	jobs := api.Group("/jobs")
	jobs.Post("/", c.JobsController.CreateJob)
	jobs.Get("/", c.JobsController.GetAllJobs)
	jobs.Get("/:id", c.JobsController.FindJobByID)
	jobs.Get("/company/:id", c.JobsController.GetJobsByCompanyID)
	jobs.Put("/:id", c.JobsController.UpdateJob)

	userJobs := api.Group("/user-jobs")
	userJobs.Post("/", c.UserJobsController.CreateUserJob)

}
