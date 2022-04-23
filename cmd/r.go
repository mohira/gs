package cmd

type ConversationRepliesJSON struct {
	Ok               bool   `json:"ok"`
	Error            string `json:"error"`
	ResponseMetadata struct {
		Messages []string `json:"messages"`
	} `json:"response_metadata"`
	Messages []struct {
		Type  string `json:"type"`
		Ts    string `json:"ts"`
		Files []struct {
			ID                 string `json:"id"`
			Timestamp          int    `json:"timestamp"`
			Name               string `json:"name"`
			Title              string `json:"title"`
			Mimetype           string `json:"mimetype"`
			Filetype           string `json:"filetype"`
			PrettyType         string `json:"pretty_type"`
			User               string `json:"user"`
			Size               int    `json:"size"`
			IsPublic           bool   `json:"is_public"`
			PublicURLShared    bool   `json:"public_url_shared"`
			DisplayAsBot       bool   `json:"display_as_bot"`
			Username           string `json:"username"`
			URLPrivate         string `json:"url_private"`
			URLPrivateDownload string `json:"url_private_download"`
			Permalink          string `json:"permalink"`
			PermalinkPublic    string `json:"permalink_public"`
		} `json:"files,omitempty"`
	} `json:"messages"`
}

type SlackFiles []SlackFile
type SlackFile struct {
	Name               string
	Mimetype           string
	UrlPrivateDownload string
}

func (rs *ConversationRepliesJSON) ToSlackFiles() SlackFiles {
	var slackFiles SlackFiles

	for _, message := range rs.Messages {
		for _, file := range message.Files {
			f := SlackFile{
				Name:               file.Name,
				Mimetype:           file.Mimetype,
				UrlPrivateDownload: file.URLPrivateDownload,
			}
			slackFiles = append(slackFiles, f)
		}
	}

	return slackFiles
}
