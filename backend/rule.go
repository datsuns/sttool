package backend

import "regexp"

/*
- thumbnailUrl = https://clips-media-assets2.twitch.tv/<clip-id>-preview-480x272.jpg
-	mp4 = https://clips-media-assets2.twitch.tv/<clip-id>.mp4
*/
func ConvertThumbnailToMp4Url(thumbnailUrl string) string {
	r := regexp.MustCompile("-preview-.*")
	return r.ReplaceAllString(thumbnailUrl, ".mp4")
}
