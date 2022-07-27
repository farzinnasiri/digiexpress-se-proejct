package courier

import (
	"context"

	"github.com/digiexpress/dlocator/internal/app/dlocator/courier/entities"
)

type CommandHandler struct {
	cacheRepository CacheRepository
}

func NewCourierCommandHandler(cacheRepository CacheRepository) CommandHandler {
	return CommandHandler{cacheRepository}
}

func (c CommandHandler) PushCourierLocation(context context.Context, command PushCourierLocation) error {
	location, err := entities.NewLocation(command.Location.Latitude, command.Location.Longitude)
	if err != nil {
		return err
	}

	courierDTO := Courier{
		DriverID: command.DriverID,
		Location: LocationToDTO(location),
	}

	if err = c.cacheRepository.SaveCourierLocation(courierDTO); err != nil {
		return err
	}

	return nil
}
