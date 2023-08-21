package handlers

import (
	"fmt"
	"strings"

	"fuck_you.com/utils/dero" // Import the dero package from your project
	"github.com/bwmarrin/discordgo"
)

func HandleTip(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if strings.HasPrefix(content, "!tip ") {
		// Extract the address or wallet name from the content
		input := strings.TrimPrefix(content, "!tip ")
		fmt.Println(input)

		// Check if the user input is in the userMappings
		userID := message.Author.ID
		userMappingsMutex.Lock()
		mappedAddress, exists := userMappings[userID]
		userMappingsMutex.Unlock()

		if exists {
			// Use the mapped address for the tip
			dero.MakeTransfer(mappedAddress)
			discord.ChannelMessageSend(message.ChannelID, "Tip sent!")
		} else if len(input) == 66 && strings.HasPrefix(input, "dero") {
			// If input is a valid DERO address, use it directly for transfer
			dero.MakeTransfer(input)
			discord.ChannelMessageSend(message.ChannelID, "Tip sent!")
		} else {
			// Otherwise, perform a wallet name lookup
			lookupResult := dero.WalletNameToAddress(input) // Implement the wallet name lookup function

			if lookupResult != "" {
				dero.MakeTransfer(lookupResult)
				discord.ChannelMessageSend(message.ChannelID, "Tip sent!")
			} else {
				discord.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name.")
			}
		}
	}
}
