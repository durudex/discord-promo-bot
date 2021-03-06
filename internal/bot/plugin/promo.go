/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package plugin

import (
	"context"

	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Promo commands plugin structure.
type PromoPlugin struct {
	service service.Promo
	handler *command.Handler
}

// Creating a new promo commands plugin.
func NewPromoPlugin(service service.Promo, handler *command.Handler) *PromoPlugin {
	return &PromoPlugin{service: service, handler: handler}
}

// Registering promo plugin commands.
func (p *PromoPlugin) RegisterCommands() {
	// Register create promo command.
	p.createPromoCommand()
	// Register use promo command.
	p.usePromoCommand()
}

// The command creating a new user promo code.
func (p *PromoPlugin) createPromoCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:         "create",
			Description:  "The command creating a new user promo code.",
			DMPermission: &DMPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "promo",
					Description: "Promo code.",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Checking where the command was use.
			if i.Interaction.User == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.Interaction.User
			}

			// Update use promo code.
			if err := p.service.Update(
				context.Background(),
				author.ID,
				i.ApplicationCommandData().Options[0].StringValue(),
			); err != nil {
				// Send a interaction respond error message.
				if err := discordInteractionError(s, i, err); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond error message")
				}

				return
			}

			// Send a interaction respond message.
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You created promo code: " + i.ApplicationCommandData().Options[0].StringValue(),
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// The command use a user promo code.
func (p *PromoPlugin) usePromoCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:         "use",
			Description:  "The command use a user promo code.",
			DMPermission: &DMPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "promo",
					Description: "Promo code.",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Checking where the command was use.
			if i.Interaction.User == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.Interaction.User
			}

			// Use a promo code.
			if err := p.service.Use(
				context.Background(),
				author.ID,
				i.ApplicationCommandData().Options[0].StringValue(),
			); err != nil {
				// Send a interaction respond error message.
				if err := discordInteractionError(s, i, err); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond error message")
				}

				return
			}

			// Send a interaction respond message.
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You used promo code: " + i.ApplicationCommandData().Options[0].StringValue(),
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}
