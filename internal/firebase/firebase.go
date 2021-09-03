package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"os"
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	credentials := os.Getenv("credentials")
	opt := option.WithCredentialsJSON([]byte(credentials))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetData(collection, docId string) (map[string]interface{}, error) {
	ctx := context.Background()
	client, err := newFirestoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	doc, err := client.Collection(collection).Doc(docId).Get(ctx)
	if err != nil {
		return nil, err
	}
	return doc.Data(), nil
}
