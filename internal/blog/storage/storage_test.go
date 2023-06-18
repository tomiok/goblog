package storage

import "testing"

func TestStorage_GeoSearch(t *testing.T) {
	type fields struct {
		clientDB    *db.Database
		clientRedis *redis.Client
	}
	type args struct {
		lon float64
		lat float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				clientDB:    tt.fields.clientDB,
				clientRedis: tt.fields.clientRedis,
			}
			s.GeoSearch(tt.args.lon, tt.args.lat)
		})
	}
}