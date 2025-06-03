package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/thunder"
	"github.com/octoberswimmer/thunder/api"
	"github.com/octoberswimmer/thunder/components"
)

// Message types for form handling
type firstNameMsg string
type lastNameMsg string
type emailMsg string
type phoneMsg string
type addressMsg string
type submitFormMsg struct{}
type formSubmittedMsg struct{}
type formErrorMsg struct{ error string }

// PatientFormModel demonstrates validation with ValidatedTextInput
type PatientFormModel struct {
	masc.Core

	// Form data
	firstName string
	lastName  string
	email     string
	phone     string
	address   string

	// Address autocomplete state
	addressPredictions  []api.PlacePrediction
	selectedAddress     *api.PlaceDetails
	addressApiError     string
	addressJustSelected bool
	googleApiKey        string

	// Validation state
	validationErrors map[string]string
	formSubmitted    bool
	isSubmitting     bool

	// Toast state for feedback
	showToast    bool
	toastVariant components.ToastVariant
	toastMessage string
}

// Init initializes the model
func (m *PatientFormModel) Init() masc.Cmd {
	m.validationErrors = make(map[string]string)
	// Set a demo Google API key - in real app this would come from config
	m.googleApiKey = "DEMO_API_KEY"
	return nil
}

// Update handles form messages and validation
func (m *PatientFormModel) Update(msg masc.Msg) (masc.Model, masc.Cmd) {
	switch msg := msg.(type) {
	case firstNameMsg:
		m.firstName = string(msg)
		// Clear error if field becomes valid
		if strings.TrimSpace(m.firstName) != "" {
			delete(m.validationErrors, "firstName")
		}
		return m, nil

	case lastNameMsg:
		m.lastName = string(msg)
		// Clear error if field becomes valid
		if strings.TrimSpace(m.lastName) != "" {
			delete(m.validationErrors, "lastName")
		}
		return m, nil

	case emailMsg:
		m.email = string(msg)
		// Clear error if email becomes valid
		if isValidEmail(m.email) {
			delete(m.validationErrors, "email")
		}
		return m, nil

	case phoneMsg:
		m.phone = string(msg)
		// Clear error if phone becomes valid
		if isValidPhone(m.phone) {
			delete(m.validationErrors, "phone")
		}
		return m, nil

	case addressMsg:
		m.address = string(msg)
		// Clear errors if address becomes valid
		if strings.TrimSpace(m.address) != "" {
			delete(m.validationErrors, "address")
		}

		// If address was just selected, don't trigger new autocomplete
		if m.addressJustSelected {
			m.addressJustSelected = false
			m.addressPredictions = nil
			return m, nil
		}

		// Clear API error when user starts typing
		m.addressApiError = ""
		// Fetch address predictions if we have enough input
		if len(strings.TrimSpace(m.address)) >= 3 {
			return m, components.AddressAutocompleteCmd(m.googleApiKey, m.address)
		}
		// Clear predictions if input is too short
		m.addressPredictions = nil
		return m, nil

	case components.AddressAutocompleteResult:
		if msg.Error != nil {
			// Handle API error with user-friendly message
			m.addressPredictions = nil
			m.addressApiError = msg.Error.Error()
		} else {
			// Clear any previous error and update predictions
			m.addressApiError = ""
			m.addressPredictions = msg.Predictions
		}
		return m, nil

	case submitFormMsg:
		// Validate form before submission
		m.formSubmitted = true
		m.validationErrors = m.validateForm()

		if len(m.validationErrors) == 0 {
			// Form is valid, submit it
			m.isSubmitting = true
			return m, m.submitFormCmd()
		}

		// Show validation error toast
		m.showToast = true
		m.toastVariant = components.VariantError
		m.toastMessage = "Please correct the highlighted fields and try again."
		return m, nil

	case formSubmittedMsg:
		// Form submitted successfully
		m.isSubmitting = false
		m.showToast = true
		m.toastVariant = components.VariantSuccess
		m.toastMessage = "Patient information saved successfully!"

		// Reset form
		m.firstName = ""
		m.lastName = ""
		m.email = ""
		m.phone = ""
		m.address = ""
		m.addressPredictions = nil
		m.selectedAddress = nil
		m.addressApiError = ""
		m.addressJustSelected = false
		m.formSubmitted = false
		m.validationErrors = make(map[string]string)
		return m, nil

	case formErrorMsg:
		// Form submission failed
		m.isSubmitting = false
		m.showToast = true
		m.toastVariant = components.VariantError
		m.toastMessage = "Failed to save: " + msg.error
		return m, nil

	default:
		return m, nil
	}
}

// validateForm performs comprehensive form validation
func (m *PatientFormModel) validateForm() map[string]string {
	errors := make(map[string]string)

	// Required field validation
	if strings.TrimSpace(m.firstName) == "" {
		errors["firstName"] = "First name is required"
	}

	if strings.TrimSpace(m.lastName) == "" {
		errors["lastName"] = "Last name is required"
	}

	if strings.TrimSpace(m.email) == "" {
		errors["email"] = "Email address is required"
	} else if !isValidEmail(m.email) {
		errors["email"] = "Please enter a valid email address"
	}

	// Optional phone validation
	if m.phone != "" && !isValidPhone(m.phone) {
		errors["phone"] = "Please enter a valid phone number (e.g., (555) 123-4567)"
	}

	// Required address validation
	if strings.TrimSpace(m.address) == "" {
		errors["address"] = "Address is required"
	}

	return errors
}

// hasError checks if a field has a validation error
func (m *PatientFormModel) hasError(field string) bool {
	_, exists := m.validationErrors[field]
	return m.formSubmitted && exists
}

// Render builds the patient form with validation
func (m *PatientFormModel) Render(send func(masc.Msg)) masc.ComponentOrHTML {
	children := []masc.MarkupOrChild{
		m.renderPatientForm(send),
	}

	// Add toast if showing
	if m.showToast {
		children = append(children, m.renderToast(send))
	}

	return components.Container(children...)
}

// renderPatientForm builds the main form with ValidatedTextInput components
func (m *PatientFormModel) renderPatientForm(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Page(
		components.PageHeader("Patient Registration", "Enter patient information with tooltip examples"),
		components.Card("Patient Details",
			components.Grid(
				// First Name - Required field with validation and tooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"First Name",
						m.firstName,
						components.ValidationState{
							Required:     true,
							HasError:     m.hasError("firstName"),
							ErrorMessage: m.validationErrors["firstName"],
							HelpText:     "Patient's legal first name",
							Placeholder:  "Enter first name",
							Tooltip:      "Enter the patient's legal first name as it appears on their ID",
						},
						func(e *masc.Event) {
							send(firstNameMsg(e.Target.Get("value").String()))
						},
					),
				),

				// Last Name - Required field with validation and tooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Last Name",
						m.lastName,
						components.ValidationState{
							Required:     true,
							HasError:     m.hasError("lastName"),
							ErrorMessage: m.validationErrors["lastName"],
							HelpText:     "Patient's legal last name",
							Placeholder:  "Enter last name",
							Tooltip:      "Enter the patient's legal last name as it appears on their ID",
						},
						func(e *masc.Event) {
							send(lastNameMsg(e.Target.Get("value").String()))
						},
					),
				),

				// Email - Required with format validation and tooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Email Address",
						m.email,
						components.ValidationState{
							Required:     true,
							HasError:     m.hasError("email"),
							ErrorMessage: m.validationErrors["email"],
							HelpText:     "Primary contact email",
							Placeholder:  "patient@example.com",
							Tooltip:      "We'll use this email for appointment reminders and important communications",
						},
						func(e *masc.Event) {
							send(emailMsg(e.Target.Get("value").String()))
						},
					),
				),

				// Phone - Optional with format validation and helpful tooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Phone Number",
						m.phone,
						components.ValidationState{
							Required:     false,
							HasError:     m.hasError("phone"),
							ErrorMessage: m.validationErrors["phone"],
							HelpText:     "Contact phone number (optional)",
							Placeholder:  "(555) 123-4567",
							Tooltip:      "Include area code. Start international numbers with +",
						},
						func(e *masc.Event) {
							send(phoneMsg(e.Target.Get("value").String()))
						},
					),
				),

				// Address - Required with Google Places autocomplete
				components.GridColumn("1-of-1",
					components.AddressAutocomplete(
						"Address",
						m.address,
						m.googleApiKey,
						m.addressPredictions,
						m.addressApiError,
						func(value string) {
							send(addressMsg(value))
						},
						func(details api.PlaceDetails) {
							// Store the selected address details
							m.selectedAddress = &details
							// Set flag to prevent autocomplete on next addressMsg
							m.addressJustSelected = true
							// Update the input field with the formatted address
							send(addressMsg(details.FormattedAddress))
						},
					),
					// Show validation error if present
					masc.If(m.hasError("address"),
						components.ErrorMessage(m.validationErrors["address"]),
					),
					// Show selected address details if available
					m.renderAddressDetails(),
				),
			),
		),

		// Demonstration of convenience constructors
		components.Card("Convenience Constructor Examples",
			components.Grid(
				// Example using WithTooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Tooltip Only",
						"",
						components.WithTooltip("This field demonstrates the WithTooltip constructor"),
						func(e *masc.Event) {},
					),
				),

				// Example using WithPlaceholder
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Placeholder Only",
						"",
						components.WithPlaceholder("Placeholder text here"),
						func(e *masc.Event) {},
					),
				),

				// Example using WithTooltipAndPlaceholder
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Tooltip + Placeholder",
						"",
						components.WithTooltipAndPlaceholder("Both tooltip and placeholder", "Enter text"),
						func(e *masc.Event) {},
					),
				),

				// Example using RequiredWithTooltip
				components.GridColumn("1-of-2",
					components.ValidatedTextInput(
						"Required + Tooltip",
						"",
						components.RequiredWithTooltip("This is a required field with a tooltip"),
						func(e *masc.Event) {},
					),
				),
			),
		),

		// Action buttons
		components.Grid(
			components.GridColumn("1-of-2",
				components.Button("Clear Form", components.VariantNeutral, func(e *masc.Event) {
					// Reset form logic would go here
				}),
			),
			components.GridColumn("1-of-2",
				masc.If(m.isSubmitting,
					components.CenteredGrid(
						components.GridColumn("",
							components.LoadingButton("Saving...", components.VariantBrand),
						),
					),
				),
				masc.If(!m.isSubmitting,
					components.Button("Save Patient", components.VariantBrand, func(e *masc.Event) {
						send(submitFormMsg{})
					}),
				),
			),
		),
	)
}

// renderAddressDetails shows selected address information safely
func (m *PatientFormModel) renderAddressDetails() masc.ComponentOrHTML {
	if m.selectedAddress == nil {
		return nil
	}

	details := m.selectedAddress
	return components.Card("Selected Address Details",
		components.Grid(
			// Street address
			masc.If(details.Street != "",
				components.GridColumn("1-of-1",
					components.StaticField("Street", details.Street),
				),
			),
			// City, State, ZIP on same row
			masc.If(details.City != "" || details.State != "" || details.PostalCode != "",
				components.GridColumn("1-of-3",
					masc.If(details.City != "",
						components.StaticField("City", details.City),
					),
				),
				components.GridColumn("1-of-3",
					masc.If(details.State != "",
						components.StaticField("State", details.State),
					),
				),
				components.GridColumn("1-of-3",
					masc.If(details.PostalCode != "",
						components.StaticField("Postal Code", details.PostalCode),
					),
				),
			),
			// Country
			masc.If(details.Country != "",
				components.GridColumn("1-of-1",
					components.StaticField("Country", details.Country),
				),
			),
			// Coordinates (if available)
			masc.If(details.Latitude != 0 && details.Longitude != 0,
				components.GridColumn("1-of-2",
					components.StaticField("Latitude", fmt.Sprintf("%.6f", details.Latitude)),
				),
				components.GridColumn("1-of-2",
					components.StaticField("Longitude", fmt.Sprintf("%.6f", details.Longitude)),
				),
			),
		),
	)
}

// renderToast shows success/error notifications
func (m *PatientFormModel) renderToast(send func(masc.Msg)) masc.ComponentOrHTML {
	return components.Toast(
		m.toastVariant,
		"Patient Registration",
		m.toastMessage,
		func(e *masc.Event) {
			m.showToast = false
		},
	)
}

// submitFormCmd simulates form submission
func (m *PatientFormModel) submitFormCmd() masc.Cmd {
	return func() masc.Msg {
		// Simulate API call delay
		time.Sleep(1 * time.Second)

		// Simulate successful submission
		// In real app, this would call API and handle errors
		return formSubmittedMsg{}
	}
}

// Validation helper functions
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPhone(phone string) bool {
	if phone == "" {
		return true // Optional field
	}
	// Simple phone validation - accepts (555) 123-4567, 555-123-4567, 5551234567
	phoneRegex := regexp.MustCompile(`^[\(\)0-9\-\s\+\.]{10,}$`)
	return phoneRegex.MatchString(phone)
}

func main() {
	thunder.Run(&PatientFormModel{})
}
