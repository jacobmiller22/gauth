package applications

type Application struct{}

type ApplicationRepository struct {
	GetApplications func() ([]Application, error)
}
