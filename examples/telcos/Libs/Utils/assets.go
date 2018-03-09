package Utils

import "regexp"

var youtubeRexp = regexp.MustCompile(`^(?:https:\/\/)?(?:m\.|www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$`)

// ValidateYoutubeURL ...
func ValidateYoutubeURL(url string) bool {
	return youtubeRexp.MatchString(url)
}

// ExtractYoutubeVideoID ...
func ExtractYoutubeVideoID(url string) string {
	tmp := youtubeRexp.FindStringSubmatch(url)
	if tmp == nil || len(tmp) < 2 {
		return ""
	}

	return tmp[1]
}

// MakeYoutubeEmbedLink ...
func MakeYoutubeEmbedLink(url string) string {
	videoID := ExtractYoutubeVideoID(url)
	if videoID == "" {
		return ""
	}

	return "https://www.youtube.com/embed/" + videoID
}
