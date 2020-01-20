package firestore

import (
	"net/url"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c client) GetPageViews(path string) (int, error) {
	key := pathToKey(path)
	doc, err := c.firestoreClient.Collection(pageViewsRootKey).Doc(key).Get(c.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return 0, datastore.PageViewsNotFoundError{Path: path}
		}
		return 0, err
	}
	var pvd pageViewsDocument
	doc.DataTo(&pvd)
	return pvd.Views, nil
}

func (c client) InsertPageViews(path string, pageViews int) error {
	key := pathToKey(path)
	_, err := c.firestoreClient.Collection(pageViewsRootKey).Doc(key).Set(c.ctx, pageViewsDocument{
		Path:  path,
		Views: pageViews,
	})
	return err
}

func pathToKey(path string) string {
	return url.PathEscape(path)
}
