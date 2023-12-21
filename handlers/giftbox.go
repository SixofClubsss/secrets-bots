package handlers

import (
	"discord-dero-bot/utils"
	"discord-dero-bot/utils/dero"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleGiftboxModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Command invoked in DM")
		RespondWithMessage(session, interaction, "This command cannot be used in DMs.")
		return
	}
	components := createGiftBoxModalComponents()

	modal := NewModal(session, interaction, "giftbox_"+interaction.Interaction.Member.User.ID, "Purchase a DERO Gift Box", components)
	modal.Show()
}
func createGiftBoxModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "color",
				Label:       "Shirt Color",
				Style:       discordgo.TextInputShort,
				Placeholder: "black or white?",
				Required:    true,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "size",
				Label:       "Shirt Size",
				Style:       discordgo.TextInputShort,
				Placeholder: "What size fits you: S, M, L, XL, XXL, XXXL",
				Required:    true,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:  "address",
				Label:     "Shipping Address?",
				Style:     discordgo.TextInputParagraph,
				Required:  false,
				MaxLength: 80,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "name",
				Label:       "Name for the order?",
				Style:       discordgo.TextInputShort,
				Required:    false,
				Placeholder: "Please use your real name",
				MaxLength:   128,
			},
		}},
	}
}

func handleGiftboxInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Interaction received in DM")
		RespondWithMessage(session, interaction, "This interaction cannot be processed in DMs.")
		return
	}
	price := utils.GetAsk("dero-usdt")
	amount := int((55 / price) * 100000)

	data := interaction.ModalSubmitData()
	color := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	size := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	shipping := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	contact := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	comment := "C: " + color + " S: " + size + " A: " + shipping + " P: " + contact
	address := "dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4"
	destination := 1337
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)

	content := "To purchase your giftbox, please use the following address :\n```" + integratedAddress + "```And we will get back to you as soon as your order is marked receieved.\nWe will contact you on your shipping status."
	RespondWithMessage(session, interaction, content)

	if !strings.HasPrefix(data.CustomID, "giftbox_") {
		return
	}
	shopkeeper := "706842280828469260"
	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for <@%s>'s a Giftbox,", userid, shopkeeper)
	resultsChannel := "1156571722384953466"
	_, err := session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}
