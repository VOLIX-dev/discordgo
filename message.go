// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to the Message struct

package discordgo

import (
	"io"
	"regexp"
	"strings"
)

// MessageType is the type of Message
// https://discord.com/developers/docs/resources/channel#message-object-message-types
type MessageType int

// Block contains the valid known MessageType values
const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeUserPremiumGuildSubscription
	MessageTypeUserPremiumGuildSubscriptionTierOne
	MessageTypeUserPremiumGuildSubscriptionTierTwo
	MessageTypeUserPremiumGuildSubscriptionTierThree
	MessageTypeChannelFollowAdd
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
)

// A Message stores all data related to a specific Discord message.
type Message struct {
	// The ID of the message.
	ID string `json:"id"`

	// The ID of the channel in which the message was sent.
	ChannelID string `json:"channel_id"`

	// The ID of the guild in which the message was sent.
	GuildID string `json:"guild_id,omitempty"`

	// The content of the message.
	Content string `json:"content"`

	// The time at which the messsage was sent.
	// CAUTION: this field may be removed in a
	// future API version; it is safer to calculate
	// the creation time via the ID.
	Timestamp Timestamp `json:"timestamp"`

	// The time at which the last edit of the message
	// occurred, if it has been edited.
	EditedTimestamp Timestamp `json:"edited_timestamp"`

	// The roles mentioned in the message.
	MentionRoles []string `json:"mention_roles"`

	// Whether the message is text-to-speech.
	TTS bool `json:"tts"`

	// Whether the message mentions everyone.
	MentionEveryone bool `json:"mention_everyone"`

	// The author of the message. This is not guaranteed to be a
	// valid user (webhook-sent messages do not possess a full author).
	Author *User `json:"author"`

	// A list of attachments present in the message.
	Attachments []*MessageAttachment `json:"attachments"`

	// A list of embeds present in the message. Multiple
	// embeds can currently only be sent by webhooks.
	Embeds []*MessageEmbed `json:"embeds"`

	// A list of users mentioned in the message.
	Mentions []*User `json:"mentions"`

	// A list of reactions to the message.
	Reactions []*MessageReactions `json:"reactions"`

	// Whether the message is pinned or not.
	Pinned bool `json:"pinned"`

	// The type of the message.
	Type MessageType `json:"type"`

	// The webhook ID of the message, if it was generated by a webhook
	WebhookID string `json:"webhook_id"`

	// Member properties for this message's author,
	// contains only partial information
	Member *Member `json:"member"`

	// Channels specifically mentioned in this message
	// Not all channel mentions in a message will appear in mention_channels.
	// Only textual channels that are visible to everyone in a lurkable guild will ever be included.
	// Only crossposted messages (via Channel Following) currently include mention_channels at all.
	// If no mentions in the message meet these requirements, this field will not be sent.
	MentionChannels []*Channel `json:"mention_channels"`

	// Is sent with Rich Presence-related chat embeds
	Activity *MessageActivity `json:"activity"`

	// Is sent with Rich Presence-related chat embeds
	Application *MessageApplication `json:"application"`

	// MessageReference contains reference data sent with crossposted messages
	MessageReference *MessageReference `json:"message_reference"`

	// The flags of the message, which describe extra features of a message.
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the flag.
	Flags MessageFlags `json:"flags"`
}

func (msg *Message) GetChannel(session *Session) *Channel {
	channel, err := session.Channel(msg.ChannelID)
	if err != nil {
		return nil
	}
	return channel
}

// MessageFlags is the flags of "message" (see MessageFlags* consts)
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
type MessageFlags int

// Valid MessageFlags values
const (
	MessageFlagsCrossPosted MessageFlags = 1 << iota
	MessageFlagsIsCrossPosted
	MessageFlagsSupressEmbeds
	MessageFlagsSourceMessageDeleted
	MessageFlagsUrgent
)

// File stores info about files you e.g. send in messages.
type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

// MessageSend stores all parameters you can send with ChannelMessageSendComplex.
type MessageSend struct {
	Content         string                  `json:"content,omitempty"`
	Embed           *MessageEmbed           `json:"embed,omitempty"`
	TTS             bool                    `json:"tts"`
	Files           []*File                 `json:"-"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Reference       *MessageReference       `json:"message_reference,omitempty"`

	// TODO: Remove this when compatibility is not required.
	File *File `json:"-"`
}

// MessageEdit is used to chain parameters via ChannelMessageEditComplex, which
// is also where you should get the instance from.
type MessageEdit struct {
	Content         *string                 `json:"content,omitempty"`
	Embed           *MessageEmbed           `json:"embed,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`

	ID      string
	Channel string
}

// NewMessageEdit returns a MessageEdit struct, initialized
// with the Channel and ID.
func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

// SetContent is the same as setting the variable Content,
// except it doesn't take a pointer.
func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

// SetEmbed is a convenience function for setting the embed,
// so you can chain commands.
func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embed = embed
	return m
}

// AllowedMentionType describes the types of mentions used
// in the MessageAllowedMentions type.
type AllowedMentionType string

// The types of mentions used in MessageAllowedMentions.
const (
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)

// MessageAllowedMentions allows the user to specify which mentions
// Discord is allowed to parse in this message. This is useful when
// sending user input as a message, as it prevents unwanted mentions.
// If this type is used, all mentions must be explicitly whitelisted,
// either by putting an AllowedMentionType in the Parse slice
// (allowing all mentions of that type) or, in the case of roles and
// users, explicitly allowing those mentions on an ID-by-ID basis.
// For more information on this functionality, see:
// https://discordapp.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-reference
type MessageAllowedMentions struct {
	// The mention types that are allowed to be parsed in this message.
	// Please note that this is purposely **not** marked as omitempty,
	// so if a zero-value MessageAllowedMentions object is provided no
	// mentions will be allowed.
	Parse []AllowedMentionType `json:"parse"`

	// A list of role IDs to allow. This cannot be used when specifying
	// AllowedMentionTypeRoles in the Parse slice.
	Roles []string `json:"roles,omitempty"`

	// A list of user IDs to allow. This cannot be used when specifying
	// AllowedMentionTypeUsers in the Parse slice.
	Users []string `json:"users,omitempty"`
}

// A MessageAttachment stores data for message attachments.
type MessageAttachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

// MessageEmbedFooter is a part of a MessageEmbed struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for message embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        EmbedType              `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

// EmbedType is the type of embed
// https://discord.com/developers/docs/resources/channel#embed-object-embed-types
type EmbedType string

// Block of valid EmbedTypes
const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifv    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

// MessageReactions holds a reactions object for a message.
type MessageReactions struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// MessageActivity is sent with Rich Presence-related chat embeds
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id"`
}

// MessageActivityType is the type of message activity
type MessageActivityType int

// Constants for the different types of Message Activity
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

// MessageFlag describes an extra feature of the message
type MessageFlag int

// Constants for the different bit offsets of Message Flags
const (
	// This message has been published to subscribed channels (via Channel Following)
	MessageFlagCrossposted MessageFlag = 1 << iota
	// This message originated from a message in another channel (via Channel Following)
	MessageFlagIsCrosspost
	// Do not include any embeds when serializing this message
	MessageFlagSuppressEmbeds
)

// MessageApplication is sent with Rich Presence-related chat embeds
type MessageApplication struct {
	ID          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

// MessageReference contains reference data sent with crossposted messages
type MessageReference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

// Reference returns MessageReference of given message
func (m *Message) Reference() *MessageReference {
	return &MessageReference{
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
		MessageID: m.ID,
	}
}

// ContentWithMentionsReplaced will replace all @<id> mentions with the
// username of the mention.
func (m *Message) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+user.Username,
		).Replace(content)
	}
	return
}

var patternChannels = regexp.MustCompile("<#[^>]*>")

// ContentWithMoreMentionsReplaced will replace all @<id> mentions with the
// username of the mention, but also role IDs and more.
func (m *Message) ContentWithMoreMentionsReplaced(s *Session) (content string, err error) {
	content = m.Content

	if !s.StateEnabled {
		content = m.ContentWithMentionsReplaced()
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		content = m.ContentWithMentionsReplaced()
		return
	}

	for _, user := range m.Mentions {
		nick := user.Username

		member, err := s.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			nick = member.Nick
		}

		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+nick,
		).Replace(content)
	}
	for _, roleID := range m.MentionRoles {
		role, err := s.State.Role(channel.GuildID, roleID)
		if err != nil || !role.Mentionable {
			continue
		}

		content = strings.Replace(content, "<@&"+role.ID+">", "@"+role.Name, -1)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		channel, err := s.State.Channel(mention[2 : len(mention)-1])
		if err != nil || channel.Type == ChannelTypeGuildVoice {
			return mention
		}

		return "#" + channel.Name
	})
	return
}
