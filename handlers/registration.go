package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var userMappings map[string]string
var userMappingsMutex sync.Mutex

func init() {
	loadUserMappings()
}

func loadUserMappings() {
	file, err := os.OpenFile("userMappings.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening user mappings file: %v", err)
		return
	}
	defer file.Close()

	// Check if the file is empty before decoding
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		userMappings = make(map[string]string)
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&userMappings)
	if err != nil {
		log.Printf("Error decoding user mappings: %v", err)
	}
}

func saveUserMappings() {
	file, err := os.Create("userMappings.json")
	if err != nil {
		log.Printf("Error creating user mappings file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(userMappings)
	if err != nil {
		log.Printf("Error encoding user mappings: %v", err)
	}
}

func HandleRegister(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if strings.HasPrefix(content, "!register ") {
		input := strings.TrimPrefix(content, "!register ")
		username := message.Author.ID

		userMappingsMutex.Lock()
		defer userMappingsMutex.Unlock()

		userMappings[username] = input
		saveUserMappings()

		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("User %s registered with input: %s", username, input))
	}
}
