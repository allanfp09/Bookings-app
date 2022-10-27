package dbrepo

import (
	"bookings/book_pkg/models"
	"errors"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation allows to insert reservation data into the database
func (m *testDBRepo) InsertReservation(res models.Reservations) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error)")
	}
	return 1, nil
}

// InsertRoomRestriction allows to insert restriction data to the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestrictions) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityBYDatesByRoomID allows to search availability of roomID
func (m *testDBRepo) SearchAvailabilityBYDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Rooms, error) {

	var rooms []models.Rooms

	return rooms, nil
}

// GetRoomByID gets a room information by roomID
func (m *testDBRepo) GetRoomByID(id int) (models.Rooms, error) {

	var room models.Rooms
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil

}
