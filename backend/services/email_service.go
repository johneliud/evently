package services

import (
    "fmt"
    "log"
    "net/smtp"
    "os"

    "github.com/johneliud/evently/backend/models"
)

// EmailService handles sending emails
type EmailService struct {
    smtpHost     string
    smtpPort     string
    smtpUsername string
    smtpPassword string
    fromEmail    string
}

func NewEmailService() *EmailService {
    return &EmailService{
        smtpHost:     os.Getenv("SMTP_HOST"),
        smtpPort:     os.Getenv("SMTP_PORT"),
        smtpUsername: os.Getenv("SMTP_USERNAME"),
        smtpPassword: os.Getenv("SMTP_PASSWORD"),
        fromEmail:    os.Getenv("FROM_EMAIL"),
    }
}

// SendRSVPNotificationToOrganizer sends an email to the event organizer when someone RSVPs
func (s *EmailService) SendRSVPNotificationToOrganizer(event *models.Event, user *models.User, rsvpStatus string) error {
    // Get organizer email
    organizerEmail := event.OrganizerEmail
    if organizerEmail == "" {
        return fmt.Errorf("organizer email not found")
    }

    // Format status for display
    displayStatus := rsvpStatus
    switch rsvpStatus {
    case "going":
        displayStatus = "Going"
    case "maybe":
        displayStatus = "Maybe"
    case "not_going":
        displayStatus = "Not Going"
    }

    // Create email subject and body
    subject := fmt.Sprintf("New RSVP for %s", event.Title)
    body := fmt.Sprintf(`
Hello,

%s %s has RSVP'd to your event "%s" with status: %s.

Event Details:
- Date: %s
- Location: %s

You can view all RSVPs for this event at: http://localhost:3000/event/%d

Thank you for using Evently!
`, user.FirstName, user.LastName, event.Title, displayStatus, event.Date.Format("Monday, January 2, 2006 at 3:04 PM"), event.Location, event.ID)

    // Send the email
    return s.sendEmail(organizerEmail, subject, body)
}

// SendRSVPConfirmationToUser sends a confirmation email to the user who RSVP'd
func (s *EmailService) SendRSVPConfirmationToUser(event *models.Event, user *models.User, rsvpStatus string) error {
    // Format status for display
    displayStatus := rsvpStatus
    switch rsvpStatus {
    case "going":
        displayStatus = "Going"
    case "maybe":
        displayStatus = "Maybe"
    case "not_going":
        displayStatus = "Not Going"
    }

    // Create email subject and body
    subject := fmt.Sprintf("Your RSVP for %s", event.Title)
    body := fmt.Sprintf(`
Hello %s,

Thank you for your RSVP to "%s". Your response has been recorded as: %s.

Event Details:
- Date: %s
- Location: %s
- Organizer: %s %s

You can view the event details at: http://localhost:3000/event/%d

Thank you for using Evently!
`, user.FirstName, event.Title, displayStatus, event.Date.Format("Monday, January 2, 2006 at 3:04 PM"), event.Location, event.OrganizerFirstName, event.OrganizerLastName, event.ID)

    // Send the email
    return s.sendEmail(user.Email, subject, body)
}

// sendEmail is a helper function to send an email
func (s *EmailService) sendEmail(to, subject, body string) error {
    // Check if email service is configured
    if s.smtpHost == "" || s.smtpPort == "" || s.smtpUsername == "" || s.smtpPassword == "" || s.fromEmail == "" {
        log.Println("Email service not configured, skipping email send")
        return nil
    }

    // Set up authentication information
    auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

    // Compose the message
    msg := []byte(fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/plain; charset=utf-8\r\n"+
        "\r\n"+
        "%s", s.fromEmail, to, subject, body))

    // Send the email
    err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.fromEmail, []string{to}, msg)
    if err != nil {
        log.Printf("Error sending email: %v", err)
        return err
    }

    log.Printf("Email sent successfully to %s", to)
    return nil
}