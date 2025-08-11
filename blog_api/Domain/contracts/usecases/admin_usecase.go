package usecases

type IAdminUseCase interface {
	PromoteUser(adminID, targetUserID string) error
	DemoteUser(adminID, targetUserID string) error
}
