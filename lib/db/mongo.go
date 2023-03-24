package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// WarningLog *log.Logger
	// InfoLog   *log.Logger
	ErrorLog *log.Logger
)

func init() {
	// file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}

// Filter used to filter bson document
type Filter map[string]interface{}

// Change used to define change elements
type Change map[string]interface{}

// Database interface
type Database interface {
	Ping() error
	Find(collectionName string, filter Filter, options *options.FindOptions) ([]interface{}, error)
	// FindInterface(name string, filter interface{}, opt *options.FindOptions) ([]interface{}, error)
	FindOne(collectionName string, filter Filter, options *options.FindOneOptions, result interface{}) error
	FindOneAndDelete(collectionName string, filter Filter, opt *options.FindOneAndDeleteOptions, result interface{}) error
	// FindOneAndReplace(name string, filter Filter, replace interface{}, result interface{}) error
	FindOneAndUpdate(collectionName string, filter Filter, update Change) (*mongo.SingleResult, error)
	InsertOne(collectionName string, options *options.InsertOneOptions, document interface{}) (*mongo.InsertOneResult, error)
	InsertMany(collectionName string, opt *options.InsertManyOptions, obj []interface{}) (*mongo.InsertManyResult, error)
	UpdateOne(collectionName string, filter Filter, update Change, opt *options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(collectionName string, filter Filter, update Change, opt *options.UpdateOptions) (*mongo.UpdateResult, error)
	// ReplaceOne(name string, filter Filter, replacement Change, opt *options.ReplaceOptions) (*mongo.UpdateResult, error)
	DeleteOne(collectionName string, opt *options.DeleteOptions, obj interface{}) (*mongo.DeleteResult, error)
	// DeleteMany(name string, opt *options.DeleteOptions, filter Filter) (*mongo.DeleteResult, error)
	Exists(collectionName string, filter Filter) bool
	// CountDocuments(name string, filter Filter, opt *options.CountOptions) (int64, error)
	// CreateIndexesOne(name string, keysDoc bsonx.Doc) error
	// ReturnDB() *mongo.Database
	NewClient() *MongoDBClient
	// ListCollectionNames() ([]string, error)
	FindOneAndUpdateWithOptions(collectionName string, filter Filter, update Change, opt options.FindOneAndUpdateOptions) (*mongo.SingleResult, error)
}

type Options struct {
	MongoURI                  string
	DatabaseName              string
	MongoDBTLS                bool
	MongoDBClientKey          string
	MongoDBClientCert         string
	MongoDBRootCA             string
	MongoDBInsecureSkipVerify bool
	DialTimeout               time.Duration
}

func New(options *Options) *MongoDB {
	db := &MongoDB{}
	db.options = options

	_ = db.newClient() // establish connection to Mongo instance

	return db
}

type MongoDB struct {
	currentDb *mongo.Database
	options   *Options
}

type MongoDBClient struct {
	*mongo.Client
	*Options
}

func (db *MongoDB) newClient() *MongoDBClient {

	if db.options.MongoDBTLS {
		var tlsConfig *tls.Config
		pemClientCA, err := os.ReadFile(db.options.MongoDBRootCA)
		if err != nil {
			ErrorLog.Printf("Ca read error: %s, file: %s", err.Error(), db.options.MongoDBRootCA)
			os.Exit(1)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemClientCA) {
			ErrorLog.Println("AppendCertsFromPEM error")
			os.Exit(1)
		}
		cert, err := tls.LoadX509KeyPair(db.options.MongoDBClientCert, db.options.MongoDBClientKey)
		if err != nil {
			ErrorLog.Printf("X509KeyPair error: %s", err.Error())
			os.Exit(1)
		}
		// Create the credentials and return it
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			//ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:          certPool,
			InsecureSkipVerify: db.options.MongoDBInsecureSkipVerify,
		}
		client, err := mongo.NewClient(options.Client().SetTLSConfig(tlsConfig), options.Client().ApplyURI(db.options.MongoURI))
		if err != nil {
			ErrorLog.Printf("NewClient error: %s", err.Error())
			os.Exit(1)
		}

		// err = client.Connect(context.Background())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			os.Exit(1)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			ErrorLog.Printf("Ping error: %s", err.Error())
			os.Exit(1)
		}

		db.currentDb = client.Database(db.options.DatabaseName)
		return &MongoDBClient{client, db.options}

	}

	client, err := mongo.NewClient(options.Client().ApplyURI(db.options.MongoURI))
	if err != nil {
		ErrorLog.Printf("NewClient error: %s", err.Error())
		os.Exit(1)
	}

	err = client.Connect(context.Background())
	if err != nil {
		os.Exit(1)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		ErrorLog.Printf("Ping error: %s", err.Error())
		os.Exit(1)
	}
	client.Database(db.options.DatabaseName)
	db.currentDb = client.Database(db.options.DatabaseName)

	return &MongoDBClient{client, db.options}
}

func (db *MongoDB) Ping() error {
	return db.currentDb.RunCommand(context.Background(), bson.D{bson.E{Key: "ping", Value: 1}}).Err()
}

// func (db *MongoDB) Find(collectionName string, filter Filter, options *options.FindOptions) ([]interface{}, error) {
// var results []interface{}
func (db *MongoDB) Find(collectionName string, filter Filter, options *options.FindOptions) ([]bson.M, error) {

	cur, err := db.currentDb.Collection(collectionName).Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	// for cur.Next(context.Background()) {
	// 	results = append(results, cur.Current)
	// }
	// if err := cur.Err(); err != nil {
	// 	return nil, err
	// }
	var results []bson.M
	if err := cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, err
}

func (db *MongoDB) FindOne(collectionName string, filter Filter, options *options.FindOneOptions, result interface{}) error {
	if options == nil {
		return db.currentDb.Collection(collectionName).FindOne(context.Background(), filter).Decode(result)
	}
	return db.currentDb.Collection(collectionName).FindOne(context.Background(), filter, options).Decode(result)
}

func (db *MongoDB) FindOneAndUpdate(collectionName string, filter Filter, update Change) (*mongo.SingleResult, error) {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := db.currentDb.Collection(collectionName).FindOneAndUpdate(context.Background(), filter, update, &opt)
	return result, result.Err()
}

func (db *MongoDB) FindOneAndDelete(collectionName string, filter Filter, opt *options.FindOneAndDeleteOptions, result interface{}) error {
	return db.currentDb.Collection(collectionName).FindOneAndDelete(context.Background(), filter, opt).Decode(result)
}

func (db *MongoDB) FindOneAndUpdateWithOptions(collectionName string, filter Filter, update Change, opt options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	result := db.currentDb.Collection(collectionName).FindOneAndUpdate(context.Background(), filter, update, &opt)
	return result, result.Err()
}

func (db *MongoDB) InsertOne(collectionName string, options *options.InsertOneOptions, document interface{}) (*mongo.InsertOneResult, error) {
	return db.currentDb.Collection(collectionName).InsertOne(context.Background(), document, options)
}

func (db *MongoDB) UpdateOne(collectionName string, filter Filter, update Change, opt *options.UpdateOptions) (*mongo.UpdateResult, error) {
	return db.currentDb.Collection(collectionName).UpdateOne(context.Background(), filter, update, opt)
}

func (db *MongoDB) Exists(collectionName string, filter Filter) bool {
	var result interface{}
	err := db.currentDb.Collection(collectionName).FindOne(context.Background(), filter).Decode(&result)
	return (err == nil)
}

func (db *MongoDB) InsertMany(collectionName string, opt *options.InsertManyOptions, obj []interface{}) (*mongo.InsertManyResult, error) {
	return db.currentDb.Collection(collectionName).InsertMany(context.Background(), obj, opt)
}

func (db *MongoDB) UpdateMany(collectionName string, filter Filter, update Change, opt *options.UpdateOptions) (*mongo.UpdateResult, error) {
	return db.currentDb.Collection(collectionName).UpdateMany(context.Background(), filter, update, opt)
}

// func (db *MongoDB) ListCollectionNames() ([]string, error) {
// 	return db.currentDb.ListCollectionNames(context.Background(), bson.D{})
// }

func (db *MongoDB) DeleteOne(collectionName string, opt *options.DeleteOptions, obj interface{}) (*mongo.DeleteResult, error) {
	return db.currentDb.Collection(collectionName).DeleteOne(context.Background(), obj, opt)
}
