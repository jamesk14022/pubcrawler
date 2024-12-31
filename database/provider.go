package dbprovider

import (
	"time"

	"github.com/jamesk14022/barcrawler/types"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"errors"
	"log"
	"os"
)

type Manager interface {
	AddPlace(place *types.Place) interface{}
	AddRoute(route *types.Route) interface{}
	FindRouteBetweenPlaces(start_placeID string, end_placeID string) types.Route
	FindCachedRouteBetweenPlaces(start_placeID string, end_placeID string) types.Route
	FindRoutesByCity(city string) []types.Route
	FindPlacesByCity(city string) []types.Place
	FindPlaceByID(placeID string) types.Place
	FindUniqueCities() []string
	buildDistanceCache()
}

type manager struct {
	client        *mongo.Client
	cache         *cache.Cache
	DistanceCache map[string]map[string]types.Route
}

var Mgr Manager

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	Mgr = &manager{client: client, cache: cache.New(10*time.Minute, 15*time.Minute)}
	Mgr.buildDistanceCache()
}

func (mgr *manager) AddPlace(place *types.Place) (InsertedID interface{}) {

	collection := mgr.client.Database("dev").Collection("places")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	doc, err := bson.Marshal(place)
	if err != nil {
		log.Fatal(err)
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	return res.InsertedID
}

func (mgr *manager) AddRoute(route *types.Route) (InsertedID interface{}) {

	collection := mgr.client.Database("dev").Collection("routes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	doc, err := bson.Marshal(route)
	if err != nil {
		log.Fatal(err)
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	return res.InsertedID
}

func (mgr *manager) FindRouteBetweenPlaces(start_placeID string, end_placeID string) types.Route {

	collection := mgr.client.Database("dev").Collection("routes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result types.Route

	var filter bson.D
	if start_placeID <= end_placeID {
		filter = bson.D{{"point1", start_placeID}, {"point2", end_placeID}}
	} else {
		filter = bson.D{{"point2", start_placeID}, {"point1", end_placeID}}
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return types.Route{}
	} else if err != nil {
		log.Fatal(err)
	}

	return result
}

func (mgr *manager) FindRoutesByCity(city string) []types.Route {

	if cachedRoutes, found := mgr.cache.Get(city); found {
		return cachedRoutes.([]types.Route)
	}

	collection := mgr.client.Database("dev").Collection("routes")
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	filter := bson.D{{"city", city}}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err) // Log and exit on a genuine error
	}
	defer cur.Close(ctx)

	var results []types.Route // Ensure results is initialized as an empty slice

	for cur.Next(ctx) {
		var route_result types.Route
		if err := cur.Decode(&route_result); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}
		results = append(results, route_result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err) // Catch any cursor errors
	}

	mgr.cache.Set(city, results, cache.DefaultExpiration)
	return results
}

func (mgr *manager) buildDistanceCache() {
	distanceMap := make(map[string]map[string]types.Route)
	uniqueCities := Mgr.FindUniqueCities()
	// uniqueCities := []string{"dublin"}

	for _, city := range uniqueCities {

		routes := Mgr.FindRoutesByCity(city)

		for _, r := range routes {

			if r.Point1 <= r.Point2 {
				if _, ok := distanceMap[r.Point1]; !ok {
					distanceMap[r.Point1] = make(map[string]types.Route)
				}
				distanceMap[r.Point1][r.Point2] = r
			} else {
				if _, ok := distanceMap[r.Point2]; !ok {
					distanceMap[r.Point2] = make(map[string]types.Route)
				}
				distanceMap[r.Point2][r.Point1] = r
			}

		}

	}
	mgr.DistanceCache = distanceMap
	log.Printf("Distance cache built with %d cities", len(uniqueCities))
}

func (mgr *manager) FindCachedRouteBetweenPlaces(start_placeID string, end_placeID string) types.Route {
	var route types.Route
	if route = mgr.DistanceCache[start_placeID][end_placeID]; route.Distance != 0 {
		return route
	} else {
		return mgr.DistanceCache[end_placeID][start_placeID]
	}
}

func (mgr *manager) FindPlaceByID(placeID string) types.Place {

	collection := mgr.client.Database("dev").Collection("places")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result types.Place

	filter := bson.D{{"placeid", placeID}}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return types.Place{}
	} else if err != nil {
		log.Fatal(err)
	}

	return result
}

func (mgr *manager) FindPlacesByCity(city string) []types.Place {

	collection := mgr.client.Database("dev").Collection("places")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []types.Place
	filter := bson.D{{"city", city}}
	cur, err := collection.Find(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []types.Place{}
	} else if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var place_result types.Place
		if err := cur.Decode(&place_result); err != nil {
			log.Fatal(err)
		}
		results = append(results, place_result)
	}

	return results
}

func (mgr *manager) FindUniqueCities() []string {

	collection := mgr.client.Database("dev").Collection("places")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cities, err := collection.Distinct(ctx, "city", bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	listedCities := make([]string, len(cities))
	for i, v := range cities {
		listedCities[i] = v.(string)
	}
	return listedCities
}
