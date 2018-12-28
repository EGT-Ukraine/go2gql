package mock

import (
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/models"
)

type CommentsClient struct {
	Comments [][]*models.ItemComment
}

func (c *CommentsClient) ItemsComments(params *comments_controller.ItemsCommentsParams) (*comments_controller.ItemsCommentsOK, error) {
	return &comments_controller.ItemsCommentsOK{Payload: c.Comments}, nil
}
