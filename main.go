package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type FileContent struct {
	List []Contact `json:"list"`
}

type Contact struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Date        string `json:"date"`
}

func main() {
	inputFile := "./files/telegram_contacts.json"
	outputFile := "./files/telegram_contacts.vcf"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var contacts FileContent
	if err := json.Unmarshal(data, &contacts); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating VCF file: %v", err)
	}
	defer file.Close()

	for _, contact := range contacts.List {
		// Skip contacts without phone numbers or both names
		if contact.PhoneNumber == "" || (contact.FirstName == "" && contact.LastName == "") {
			continue
		}

		contact.PhoneNumber = formatPhoneNumber(contact.PhoneNumber)

		vCard := buildVCard(contact)

		_, err := file.WriteString(vCard)
		if err != nil {
			log.Fatalf("Error writing to VCF file: %v", err)
		}
		fmt.Printf("Contact %s %s successfully exported\n", contact.FirstName, contact.LastName)
	}

	fmt.Printf("Contacts successfully exported to %s\n", outputFile)
}

// formatPhoneNumber formats a phone number for vCard
func formatPhoneNumber(phone string) string {
	if strings.HasPrefix(phone, "00") {
		return "+" + phone[2:]
	}

	if len(phone) == 10 && strings.HasPrefix(phone, "07") {
		return "+250" + phone[1:] // Assuming "07" means a Rwandan number
	}

	return phone
}


func buildVCard(contact Contact) string {
	vCard := "BEGIN:VCARD\n"
	vCard += "VERSION:3.0\n"

	fullName := strings.TrimSpace(contact.FirstName + " " + contact.LastName)
	vCard += fmt.Sprintf("FN:%s\n", fullName)

	vCard += fmt.Sprintf("TEL:%s\n", contact.PhoneNumber)

	if contact.Date != "" {
		vCard += fmt.Sprintf("NOTE:Added on %s\n", contact.Date)
	}

	vCard += "END:VCARD\n"
	return vCard
}
