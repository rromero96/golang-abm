package abm

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

func createDocument() (*firestore.CollectionRef, context.Context, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "golang-abm"}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Printf("error initializing app: %v", err)
		return nil, nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Error when creating client")
		log.Fatalln(err)
	}

	return client.Collection(ABM_COLLECTION), ctx, nil
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	var dto DTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	dto.CreatedAt = now
	dto.UpdatedAt = now

	collection, ctx, _ := createDocument()
	ref := collection.NewDoc()
	_, err := ref.Set(ctx, map[string]interface{}{
		"name":       dto.Name,
		"role":       dto.Role,
		"created_at": dto.CreatedAt,
		"updated_at": dto.UpdatedAt,
	})
	if err != nil {
		log.Printf("an error ocurred: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto)
}

func ListRecords(w http.ResponseWriter, r *http.Request) {
	var dtos []DTO

	collection, ctx, _ := createDocument()
	records := collection.Documents(ctx)

	for {
		doc, err := records.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("an error ocurred: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dtos = append(dtos, DTO{
			ID:        doc.Ref.ID,
			Name:      doc.Data()["name"].(string),
			Role:      doc.Data()["role"].(string),
			CreatedAt: doc.Data()["created_at"].(time.Time),
			UpdatedAt: doc.Data()["updated_at"].(time.Time),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dtos)
}

func UpdateRecordByID(w http.ResponseWriter, r *http.Request) {
	var dto DTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.UpdatedAt = time.Now()
	collection, ctx, _ := createDocument()
	_, err := collection.Doc(dto.ID).Set(ctx, map[string]interface{}{
		"name":       dto.Name,
		"role":       dto.Role,
		"updated_at": dto.UpdatedAt,
	}, firestore.MergeAll)
	if err != nil {
		log.Printf("an error ocurred: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto)
}
