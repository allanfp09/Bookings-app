package repository

import (
	"bookings/book_pkg/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservations) (int, error)
	InsertRoomRestriction(r models.RoomRestrictions) error
	SearchAvailabilityBYDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Rooms, error)
	GetRoomByID(id int) (models.Rooms, error)
}
