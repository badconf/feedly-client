# Feedly Client in Golang

This repository contains a Go client for interacting with Feedly's APIs. The client enables various functionalities such as authentication, retrieving user profiles and subscriptions, marking articles as read, saving articles for later reading, and more.

## Caution / Disclaimer

This code is still in development and may contain bugs or issues that need to be resolved. It is provided "as-is," and the use of this code in a production environment or a critical system is not recommended without thorough testing and validation. Contributions and feedback are welcome to improve the codebase.

## Functions

`NewFeedlyClient(options map[string]interface{}) *FeedlyClient`

Creates a new Feedly client with the provided options including client ID, client secret, etc.

`GetCodeURL(callbackURL string) string`

Builds the URL for OAuth authentication with Feedly.

`GetAccessToken(redirectURI, code string) (map[string]interface{}, error)`

Retrieves the access token using the authentication code.

`RefreshAccessToken(refreshToken string) (map[string]interface{}, error)`

Refreshes the access token using the refresh token.

`GetUserProfile(accessToken string) (map[string]interface{}, error)`

Fetches the user's profile information.

`GetUserSubscriptions(accessToken string) ([]map[string]interface{}, error)`

Fetches the list of user's subscriptions.

`GetFeedContent(accessToken, streamID string, unreadOnly bool, newerThan int64) (map[string]interface{}, error)`

Fetches the content of a specific feed.

`MarkArticleRead(accessToken string, entryIds []string) (*http.Response, error)`

Marks one or multiple articles as read.

`SaveForLater(accessToken, userID string, entryIds []string) (*http.Response, error)`

Saves one or multiple articles for later reading.

## License

This code is released under the MIT License. See the [LICENSE](LICENSE) file for details.
