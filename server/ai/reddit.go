package ai

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const subredditName = "travel"

func getRedditPosts(location string) (string, error) {
	client, err := reddit.NewReadonlyClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	posts, _, err := client.Subreddit.SearchPosts(ctx, location, subredditName, &reddit.ListPostSearchOptions{
		ListPostOptions: reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 5,
			},
			Time: "all",
		},
	})
	if err != nil {
		return "", err
	}

	response := ""

	for _, post := range posts {
		if post.Body != "" {

			postAndComments, _, err := client.Post.Get(ctx, post.ID)
			if err != nil {
				response += fmt.Sprintf("Title: %s, Post: %s",
					post.Title, post.Body)
				continue
			}

			response += fmt.Sprintf("Title: %s, Post: %s, Top Comment:\n",
				post.Title, post.Body, postAndComments.Comments[0])
		}
	}

	return response, nil
}
