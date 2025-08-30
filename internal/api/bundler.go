package groupieapi

type artistDataBundle struct {
	Artist
	Locations ArtistLocations
	Dates     ArtistDates
	Relations Relation
}

func BundleArtistData(id int) (artistDataBundle, chan error) {
	artistChan := make(chan Artist)
	locationsChan := make(chan ArtistLocations)
	datesChan := make(chan ArtistDates)
	relationChan := make(chan Relation)
	errChan := make(chan error, 4)
	go bundleHelper(id, GetArtist, artistChan, errChan)
	go bundleHelper(id, GetLocations, locationsChan, errChan)
	go bundleHelper(id, GetDates, datesChan, errChan)
	go bundleHelper(id, GetRelation, relationChan, errChan)
	return artistDataBundle{<-artistChan, <-locationsChan, <-datesChan, <-relationChan}, errChan
}

func bundleHelper[T HasAPI](id int, f func(int) (T, error), dataChan chan T, errChan chan error) {
	data, err := f(id)
	if err != nil {
		errChan <- err
	}
	dataChan <- data
}
