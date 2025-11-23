package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"

    "github.com/miftahulhidayati/registration-payment-service/internal/config"
    "github.com/miftahulhidayati/registration-payment-service/internal/kafka"
    "github.com/miftahulhidayati/registration-payment-service/internal/repository"
)

type RegistrationsHandler struct {
    repo     *repository.Postgres
    producer *kafka.Producer
    cfg      *config.Config
}

func NewRegistrationsHandler(repo *repository.Postgres, producer *kafka.Producer, cfg *config.Config) *RegistrationsHandler {
    return &RegistrationsHandler{repo: repo, producer: producer, cfg: cfg}
}

func (h *RegistrationsHandler) Register(router fiber.Router) {
    g := router.Group("/registrations")
    g.Post("/", h.createRegistration)
    g.Get("/", h.listRegistrations)
    g.Get(":id", h.getRegistration)
    g.Put(":id", h.updateRegistration)
    g.Post(":id/cancel", h.cancelRegistration)

    // Payment endpoints (stubs/minimal)
    g.Post(":id/payment", h.uploadPaymentProof)
    g.Get(":id/payment", h.getPaymentInfo)
    g.Patch(":id/payment/verify", h.verifyPayment)
}

type createRegistrationRequest struct {
    EventID               uuid.UUID  `json:"event_id"`
    UserID                *uuid.UUID `json:"user_id"`
    FullName              string     `json:"full_name"`
    Gender                string     `json:"gender"`
    Phone                 string     `json:"phone"`
    Email                 string     `json:"email"`
    Address               *string    `json:"address"`
    EmergencyContactName  *string    `json:"emergency_contact_name"`
    EmergencyContactPhone *string    `json:"emergency_contact_phone"`
    EmergencyContactRelation *string `json:"emergency_contact_relation"`
    SpecialNeeds          *string    `json:"special_needs"`
}

// CreateRegistration godoc
// @Summary Create a new registration
// @Description Register a user for an event
// @Tags registrations
// @Accept json
// @Produce json
// @Param request body createRegistrationRequest true "Registration Request"
// @Success 201 {object} repository.Registration
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /registrations [post]
func (h *RegistrationsHandler) createRegistration(c *fiber.Ctx) error {
    var req createRegistrationRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }
    if req.FullName == "" || req.Gender == "" || req.Phone == "" || req.Email == "" || req.EventID == uuid.Nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields"})
    }

    ctx := context.Background()
    reg, err := h.repo.CreateRegistration(ctx, repository.CreateRegistrationParams{
        EventID:                 req.EventID,
        UserID:                  req.UserID,
        FullName:                req.FullName,
        Gender:                  req.Gender,
        Phone:                   req.Phone,
        Email:                   req.Email,
        Address:                 req.Address,
        EmergencyContactName:    req.EmergencyContactName,
        EmergencyContactPhone:   req.EmergencyContactPhone,
        EmergencyContactRelation: req.EmergencyContactRelation,
        SpecialNeeds:            req.SpecialNeeds,
    })
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // Publish event (best-effort)
    _ = h.publish(h.cfg.KafkaTopicRegCreated, reg.RegistrationID.String(), fiber.Map{
        "event": "registration.created",
        "data": fiber.Map{
            "registration_id": reg.RegistrationID,
            "event_id":        reg.EventID,
            "user_id":         reg.UserID,
            "full_name":       reg.FullName,
            "gender":          reg.Gender,
            "status":          reg.Status,
            "timestamp":       time.Now().UTC().Format(time.RFC3339),
        },
    })

    return c.Status(http.StatusCreated).JSON(reg)
}

// ListRegistrations godoc
// @Summary List registrations
// @Description Get a list of registrations with pagination
// @Tags registrations
// @Produce json
// @Success 200 {array} repository.Registration
// @Failure 500 {object} map[string]interface{}
// @Router /registrations [get]
func (h *RegistrationsHandler) listRegistrations(c *fiber.Ctx) error {
    // Simple pagination defaults
    limit := 20
    offset := 0
    ctx := context.Background()
    items, err := h.repo.ListRegistrations(ctx, limit, offset)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(items)
}

// GetRegistration godoc
// @Summary Get a registration by ID
// @Description Get details of a specific registration
// @Tags registrations
// @Produce json
// @Param id path string true "Registration ID"
// @Success 200 {object} repository.Registration
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /registrations/{id} [get]
func (h *RegistrationsHandler) getRegistration(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
    }
    ctx := context.Background()
    reg, err := h.repo.GetRegistrationByID(ctx, id)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    if reg == nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not found"})
    }
    return c.JSON(reg)
}

type updateRegistrationRequest struct {
    FullName                *string `json:"full_name"`
    Phone                   *string `json:"phone"`
    Email                   *string `json:"email"`
    Address                 *string `json:"address"`
    EmergencyContactName    *string `json:"emergency_contact_name"`
    EmergencyContactPhone   *string `json:"emergency_contact_phone"`
    EmergencyContactRelation *string `json:"emergency_contact_relation"`
    SpecialNeeds            *string `json:"special_needs"`
    Notes                   *string `json:"notes"`
}

// UpdateRegistration godoc
// @Summary Update a registration
// @Description Update details of an existing registration
// @Tags registrations
// @Accept json
// @Produce json
// @Param id path string true "Registration ID"
// @Param request body updateRegistrationRequest true "Update Request"
// @Success 200 {object} repository.Registration
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /registrations/{id} [put]
func (h *RegistrationsHandler) updateRegistration(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
    }
    var req updateRegistrationRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }
    ctx := context.Background()
    reg, err := h.repo.UpdateRegistration(ctx, repository.UpdateRegistrationParams{
        RegistrationID:          id,
        FullName:                req.FullName,
        Phone:                   req.Phone,
        Email:                   req.Email,
        Address:                 req.Address,
        EmergencyContactName:    req.EmergencyContactName,
        EmergencyContactPhone:   req.EmergencyContactPhone,
        EmergencyContactRelation: req.EmergencyContactRelation,
        SpecialNeeds:            req.SpecialNeeds,
        Notes:                   req.Notes,
    })
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(reg)
}

type cancelRequest struct {
    Reason string `json:"reason"`
}

// CancelRegistration godoc
// @Summary Cancel a registration
// @Description Cancel a registration with a reason
// @Tags registrations
// @Accept json
// @Produce json
// @Param id path string true "Registration ID"
// @Param request body cancelRequest true "Cancel Request"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /registrations/{id}/cancel [post]
func (h *RegistrationsHandler) cancelRegistration(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
    }
    var req cancelRequest
    if err := c.BodyParser(&req); err != nil || req.Reason == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "reason required"})
    }
    ctx := context.Background()
    if err := h.repo.CancelRegistration(ctx, id, req.Reason); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    _ = h.publish(h.cfg.KafkaTopicRegCancelled, id.String(), fiber.Map{
        "event": "registration.cancelled",
        "data": fiber.Map{
            "registration_id": id,
            "reason":          req.Reason,
            "timestamp":       time.Now().UTC().Format(time.RFC3339),
        },
    })
    return c.SendStatus(http.StatusNoContent)
}

// Payment stubs (replace with your implementation later)
// UploadPaymentProof godoc
// @Summary Upload payment proof
// @Description Upload proof of payment for a registration
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Registration ID"
// @Param request body object{amount=float64,payment_proof_url=string} true "Payment Proof"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /registrations/{id}/payment [post]
func (h *RegistrationsHandler) uploadPaymentProof(c *fiber.Ctx) error {
    // Accept minimal JSON for now
    type reqT struct {
        Amount float64 `json:"amount"`
        Proof  string  `json:"payment_proof_url"`
    }
    var req reqT
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }
    id := c.Params("id")
    _ = h.publish(h.cfg.KafkaTopicPayUploaded, id, fiber.Map{
        "event": "payment.uploaded",
        "data": fiber.Map{
            "registration_id": id,
            "amount":          req.Amount,
            "payment_proof_url": req.Proof,
            "timestamp":       time.Now().UTC().Format(time.RFC3339),
        },
    })
    return c.Status(http.StatusAccepted).JSON(fiber.Map{"status": "payment uploaded"})
}

// GetPaymentInfo godoc
// @Summary Get payment info
// @Description Get payment status and details
// @Tags payments
// @Produce json
// @Param id path string true "Registration ID"
// @Success 200 {object} map[string]interface{}
// @Router /registrations/{id}/payment [get]
func (h *RegistrationsHandler) getPaymentInfo(c *fiber.Ctx) error {
    // Placeholder response
    return c.JSON(fiber.Map{"registration_id": c.Params("id"), "payment": fiber.Map{"status": "pending"}})
}

// VerifyPayment godoc
// @Summary Verify payment
// @Description Verify and approve a payment
// @Tags payments
// @Produce json
// @Param id path string true "Registration ID"
// @Success 200 {object} map[string]interface{}
// @Router /registrations/{id}/payment/verify [patch]
func (h *RegistrationsHandler) verifyPayment(c *fiber.Ctx) error {
    id := c.Params("id")
    // In real impl: update payments + set registration to confirmed
    _ = h.publish(h.cfg.KafkaTopicPayVerified, id, fiber.Map{
        "event": "payment.verified",
        "data": fiber.Map{
            "registration_id": id,
            "verification_status": "approved",
            "timestamp":       time.Now().UTC().Format(time.RFC3339),
        },
    })
    _ = h.publish(h.cfg.KafkaTopicRegConfirmed, id, fiber.Map{
        "event": "registration.confirmed",
        "data": fiber.Map{
            "registration_id": id,
            "timestamp":       time.Now().UTC().Format(time.RFC3339),
        },
    })
    return c.JSON(fiber.Map{"status": "verified"})
}

func (h *RegistrationsHandler) publish(topic, key string, payload any) error {
    if h.producer == nil {
        return nil
    }
    return h.producer.Publish(context.Background(), topic, key, payload)
}


