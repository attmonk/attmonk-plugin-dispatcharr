package matching

import (
	"strings"

	"github.com/attmonk/attmonk-plugin-dispatcharr/internal/upstream/m3u"
	"github.com/attmonk/attmonk-plugin-dispatcharr/internal/upstream/xmltv"
)

func Match(entry m3u.Entry, doc xmltv.Document) (xmltv.Channel, bool) {
	if entry.GuideID != "" {
		for _, channel := range doc.Channels {
			if strings.EqualFold(strings.TrimSpace(channel.ID), strings.TrimSpace(entry.GuideID)) {
				return channel, true
			}
		}
	}

	for _, channel := range doc.Channels {
		for _, displayName := range channel.DisplayNames {
			if normalize(displayName) == normalize(entry.Name) {
				return channel, true
			}
		}
	}

	return xmltv.Channel{}, false
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
