package managers

import (
	"errors"
	"fmt"
	"log"
	"main/common"
	"main/database"
	"main/models"
	"os"
	"regexp"
	"strings"
	"time"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/gorm"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP")
	ErrOTPExpired = errors.New("OTP expired")
)

type OtpManager interface {
	SendOTP(userID uint, phoneNumber string) error
	VerifyOTP(phoneNumber string, otp string) error
}

type otpManager struct {
	//client
}

func NewOtpManager() OtpManager {
	return &otpManager{}
}

const otpLength = 6
const otpExpiration = 5 * time.Minute

func (otpManager *otpManager) SendOTP(userID uint, phoneNumber string) error {
	phoneNumber = formatPhoneNumber(phoneNumber)

	otp := common.GenerateOTP(otpLength)

	expiresAt := time.Now().Add(otpExpiration)

	otpRecord := models.Otp{
		UserID:    userID,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	result := database.DB.Create(&otpRecord)
	if result.Error != nil {
		return fmt.Errorf("failed to Create Otp record: %w", result.Error)
	}

	err := sendOTP(phoneNumber, otp)
	if err != nil {

		deleteErr := database.DB.Delete(&otpRecord).Error
		if deleteErr != nil {
			log.Printf("Failed to delete OTP record after sending failed: %v", deleteErr)

		}
		return fmt.Errorf("failed to send OTP via Twilio: %w", err)
	}

	return nil
}

func (otpManager *otpManager) VerifyOTP(phoneNumber string, otp string) error {
	phoneNumber = formatPhoneNumber(phoneNumber)

	var user models.User
	result := database.DB.Where("phone = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with phone number %s not found: %w", phoneNumber, result.Error)
		}
		return fmt.Errorf("failed to find user by phone number: %w", result.Error)
	}

	var otpRecord models.Otp

	result = database.DB.Where("user_id = ? AND otp = ?", user.Id, otp).
		Order("created_at DESC").
		First(&otpRecord)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrInvalidOTP
		}
		return fmt.Errorf("failed to find OTP record: %w", result.Error)
	}

	if otpRecord.ExpiresAt.Before(time.Now()) {
		return ErrOTPExpired
	}

	result = database.DB.Delete(&otpRecord)
	if result.Error != nil {
		log.Printf("Error deleting OTP record: %v", result.Error)
	}

	return nil
}

func sendOTP(phoneNumber, otp string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	messageInput := &openapi.CreateMessageParams{}
	messageInput.SetTo(phoneNumber)
	messageInput.SetFrom(twilioPhoneNumber)
	messageInput.SetBody(fmt.Sprintf("Your OTP is: %s", otp))

	_, err := client.Api.CreateMessage(messageInput)
	if err != nil {
		log.Printf("twilio error %s", err.Error())
		return err
	}

	log.Println("OTP sent successfully")

	return nil
}

func formatPhoneNumber(phoneNumber string) string {

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Println(err)
		return ""
	}

	phoneNumber = reg.ReplaceAllString(phoneNumber, "")

	if strings.HasPrefix(phoneNumber, "91") {
		phoneNumber = "+" + phoneNumber
	} else if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	return phoneNumber
}
