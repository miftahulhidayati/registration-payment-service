package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Registration struct {
	RegistrationID          uuid.UUID  `json:"registration_id"`
	EventID                 uuid.UUID  `json:"event_id"`
	UserID                  *uuid.UUID `json:"user_id"`
	FullName                string     `json:"full_name"`
	Gender                  string     `json:"gender"`
	Phone                   string     `json:"phone"`
	Email                   string     `json:"email"`
	Address                 *string    `json:"address"`
	EmergencyContactName    *string    `json:"emergency_contact_name"`
	EmergencyContactPhone   *string    `json:"emergency_contact_phone"`
	EmergencyContactRelation *string   `json:"emergency_contact_relation"`
	SpecialNeeds            *string    `json:"special_needs"`
	RegistrationDate        time.Time  `json:"registration_date"`
	Status                  string     `json:"status"`
	CancelledAt             *time.Time `json:"cancelled_at"`
	CancellationReason      *string    `json:"cancellation_reason"`
	Notes                   *string    `json:"notes"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type CreateRegistrationParams struct {
	EventID                 uuid.UUID
	UserID                  *uuid.UUID
	FullName                string
	Gender                  string
	Phone                   string
	Email                   string
	Address                 *string
	EmergencyContactName    *string
	EmergencyContactPhone   *string
	EmergencyContactRelation *string
	SpecialNeeds            *string
}

type UpdateRegistrationParams struct {
	RegistrationID          uuid.UUID
	FullName                *string
	Phone                   *string
	Email                   *string
	Address                 *string
	EmergencyContactName    *string
	EmergencyContactPhone   *string
	EmergencyContactRelation *string
	SpecialNeeds            *string
	Notes                   *string
}

func (r *Postgres) CreateRegistration(ctx context.Context, params CreateRegistrationParams) (*Registration, error) {
	query := `
		INSERT INTO registrations (
			event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
	`

	var reg Registration
	err := r.Pool.QueryRow(ctx, query,
		params.EventID, params.UserID, params.FullName, params.Gender,
		params.Phone, params.Email, params.Address, params.EmergencyContactName,
		params.EmergencyContactPhone, params.EmergencyContactRelation, params.SpecialNeeds,
	).Scan(
		&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
		&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
		&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
		&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
		&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &reg, nil
}

func (r *Postgres) GetRegistrationByID(ctx context.Context, registrationID uuid.UUID) (*Registration, error) {
	query := `
		SELECT registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
		FROM registrations
		WHERE registration_id = $1
	`

	var reg Registration
	err := r.Pool.QueryRow(ctx, query, registrationID).Scan(
		&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
		&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
		&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
		&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
		&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reg, nil
}

func (r *Postgres) ListRegistrations(ctx context.Context, limit, offset int) ([]*Registration, error) {
	query := `
		SELECT registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
		FROM registrations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []*Registration
	for rows.Next() {
		var reg Registration
		err := rows.Scan(
			&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
			&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
			&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
			&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
			&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, &reg)
	}

	return registrations, nil
}

func (r *Postgres) UpdateRegistration(ctx context.Context, params UpdateRegistrationParams) (*Registration, error) {
	query := `
		UPDATE registrations
		SET full_name = COALESCE($2, full_name),
			phone = COALESCE($3, phone),
			email = COALESCE($4, email),
			address = COALESCE($5, address),
			emergency_contact_name = COALESCE($6, emergency_contact_name),
			emergency_contact_phone = COALESCE($7, emergency_contact_phone),
			emergency_contact_relation = COALESCE($8, emergency_contact_relation),
			special_needs = COALESCE($9, special_needs),
			notes = COALESCE($10, notes),
			updated_at = CURRENT_TIMESTAMP
		WHERE registration_id = $1
		RETURNING registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
	`

	var reg Registration
	err := r.Pool.QueryRow(ctx, query,
		params.RegistrationID, params.FullName, params.Phone, params.Email,
		params.Address, params.EmergencyContactName, params.EmergencyContactPhone,
		params.EmergencyContactRelation, params.SpecialNeeds, params.Notes,
	).Scan(
		&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
		&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
		&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
		&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
		&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &reg, nil
}

func (r *Postgres) CancelRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	query := `
		UPDATE registrations
		SET status = 'cancelled',
			cancelled_at = CURRENT_TIMESTAMP,
			cancellation_reason = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE registration_id = $1
	`

	_, err := r.Pool.Exec(ctx, query, registrationID, reason)
	return err
}

func (r *Postgres) UpdateRegistrationStatus(ctx context.Context, registrationID uuid.UUID, status string) error {
	query := `
		UPDATE registrations
		SET status = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE registration_id = $1
	`

	_, err := r.Pool.Exec(ctx, query, registrationID, status)
	return err
}

func (r *Postgres) GetRegistrationsByEventID(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*Registration, error) {
	query := `
		SELECT registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
		FROM registrations
		WHERE event_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []*Registration
	for rows.Next() {
		var reg Registration
		err := rows.Scan(
			&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
			&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
			&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
			&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
			&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, &reg)
	}

	return registrations, nil
}

func (r *Postgres) GetRegistrationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Registration, error) {
	query := `
		SELECT registration_id, event_id, user_id, full_name, gender, phone, email,
			address, emergency_contact_name, emergency_contact_phone,
			emergency_contact_relation, special_needs, registration_date,
			status, cancelled_at, cancellation_reason, notes, created_at, updated_at
		FROM registrations
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []*Registration
	for rows.Next() {
		var reg Registration
		err := rows.Scan(
			&reg.RegistrationID, &reg.EventID, &reg.UserID, &reg.FullName, &reg.Gender,
			&reg.Phone, &reg.Email, &reg.Address, &reg.EmergencyContactName,
			&reg.EmergencyContactPhone, &reg.EmergencyContactRelation, &reg.SpecialNeeds,
			&reg.RegistrationDate, &reg.Status, &reg.CancelledAt, &reg.CancellationReason,
			&reg.Notes, &reg.CreatedAt, &reg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, &reg)
	}

	return registrations, nil
}

