package dbrepo

import (
	"bookings/book_pkg/models"
	"context"
	"time"
)

func (m *PostgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation allows to insert reservation data into the database
func (m *PostgresDBRepo) InsertReservation(res models.Reservations) (int, error) {
	// Allows to cancel activity after 3 seconds later if data was not sent to DB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, 
                          room_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction allows to insert restriction data to the database
func (m *PostgresDBRepo) InsertRoomRestriction(r models.RoomRestrictions) error {
	// Allows to cancel activity after 3 seconds later if data was not sent to DB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (
                               start_date, end_date, room_id, reservation_id,
                               restriction_id, created_at, updated_at) VALUES (
                                                       $1, $2, $3, $4, $5, $6, $7
                               )`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityBYDatesByRoomID allows to search availability of roomID
func (m *PostgresDBRepo) SearchAvailabilityBYDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// Allows to cancel activity after 3 seconds later if data was not sent to DB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newRows int

	query := `SELECT count(id) FROM room_restrictions WHERE room_id = $1 AND $2 < end_date AND $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&newRows)
	if err != nil {
		return false, nil
	}

	if newRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *PostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Rooms, error) {
	// Allows to cancel activity after 3 seconds later if data was not sent to DB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Rooms

	query := `
		SELECT 
			r.id, r.room_name
		FROM
		    rooms r 
		WHERE 
		    r.id NOT IN (
		        SELECT
		            room_id
		        FROM
		            room_restrictions rr
		        WHERE
		            $1 < rr.end_date AND $2 > rr.start_date
		    )
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Rooms
		err := rows.Scan(&room.ID, &room.RoomName)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	err = rows.Err()
	if err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room information by roomID
func (m *PostgresDBRepo) GetRoomByID(id int) (models.Rooms, error) {
	// Allows to cancel activity after 3 seconds later if data was not sent to DB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Rooms

	query := `
		SELECT id, room_name, created_at, updated_at FROM rooms WHERE id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
