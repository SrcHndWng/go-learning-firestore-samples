package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Room struct {
	Name string `firestore:"name,omitempty"`
}

type Message struct {
	From string `firestore:"from,omitempty"`
	Msg  string `firestore:"msg,omitempty"`
}

func main() {
	fmt.Println("----- main start -----")

	jsonPath := os.Getenv("FIREBASE_JSON_PATH")

	opt := option.WithCredentialsFile(jsonPath)
	ctx := context.Background()
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	err = create(ctx, client)
	if err != nil {
		log.Fatalln(err)
	}

	err = reference(ctx, client)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("----- main end -----")
}

func reference(ctx context.Context, client *firestore.Client) error {
	docs := client.Collection("rooms").Documents(ctx)
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("document id = %s\n", doc.Ref.ID)
		msgDocs := doc.Ref.Collection("messages").Documents(ctx)
		for {
			msgDoc, err := msgDocs.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			fmt.Println(msgDoc.Data())
		}

	}

	return nil
}

func create(ctx context.Context, client *firestore.Client) error {
	roomA := Room{
		Name: "my chat room",
	}
	message1 := Message{
		From: "alex",
		Msg:  "Hello World!",
	}
	_, err := client.Collection("rooms").Doc("roomA").Set(ctx, roomA)
	if err != nil {
		return err
	}
	_, err = client.Collection("rooms").Doc("roomA").Collection("messages").Doc("message1").Set(ctx, message1)
	if err != nil {
		return err
	}

	message2 := Message{
		From: "tom",
		Msg:  "Hello Tom!",
	}
	_, err = client.Collection("rooms").Doc("roomA").Collection("messages").Doc("message2").Set(ctx, message2)
	if err != nil {
		return err
	}

	roomB := Room{
		Name: "beta room",
	}
	message3 := Message{
		From: "ken",
		Msg:  "Hello Ken!",
	}
	_, err = client.Collection("rooms").Doc("roomB").Set(ctx, roomB)
	if err != nil {
		return err
	}
	_, err = client.Collection("rooms").Doc("roomB").Collection("messages").Doc("message3").Set(ctx, message3)
	if err != nil {
		return err
	}

	return nil
}
