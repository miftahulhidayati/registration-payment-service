-- Drop indexes
DROP INDEX IF EXISTS unique_event_user;
DROP INDEX IF EXISTS idx_payments_verification_status;
DROP INDEX IF EXISTS idx_payments_registration;
DROP INDEX IF EXISTS idx_registrations_status;
DROP INDEX IF EXISTS idx_registrations_user;
DROP INDEX IF EXISTS idx_registrations_event;

-- Drop tables
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS registrations;

-- Drop ENUM types
DROP TYPE IF EXISTS payment_verification_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS registration_status;

