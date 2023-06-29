package mysql

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type location struct {
	Point wkb.Point
}

func (loc location) GormDataType() string {
	return "point"
}

func (loc location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%v %v)", loc.Point.X(), loc.Point.Y())},
	}
}

// Scan implements the sql.Scanner interface
func (loc *location) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	mysqlEncoding, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", src)
	}

	var srid = binary.LittleEndian.Uint32(mysqlEncoding[0:4])

	err := loc.Point.Scan(mysqlEncoding[4:])

	loc.Point.SetSRID(int(srid))

	return err
}

type polygon struct {
	MPoly wkb.MultiPolygon
}

func (pol polygon) GormDataType() string {
	return "multipolygon"
}

func (pol polygon) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_MultiPolygonFromText(?)",
		Vars: []interface{}{generateMultiPolygon(convertCoords(pol.MPoly.MultiPolygon.Coords()))},
	}
}

// Scan implements the sql.Scanner interface
func (pol *polygon) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	mysqlEncoding, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", src)
	}

	var srid = binary.LittleEndian.Uint32(mysqlEncoding[0:4])

	err := pol.MPoly.Scan(mysqlEncoding[4:])
	pol.MPoly.SetSRID(int(srid))

	return err
}

func convertCoords(coords [][][]geom.Coord) [][][][]float64 {
	converted := make([][][][]float64, len(coords))
	for i, level1 := range coords {
		converted[i] = make([][][]float64, len(level1))
		for j, level2 := range level1 {
			converted[i][j] = make([][]float64, len(level2))
			for k, level3 := range level2 {
				converted[i][j][k] = make([]float64, len(level3))
				copy(converted[i][j][k], level3)
			}
		}
	}
	return converted
}

func convertFloats(coords [][][][]float64) [][][]geom.Coord {
	converted := make([][][]geom.Coord, len(coords))
	for i, level1 := range coords {
		converted[i] = make([][]geom.Coord, len(level1))
		for j, level2 := range level1 {
			converted[i][j] = make([]geom.Coord, len(level2))
			for k, level3 := range level2 {
				converted[i][j][k] = make(geom.Coord, len(level3))
				copy(converted[i][j][k], level3)
			}
		}
	}
	return converted
}
func generateMultiPolygon(coordinates [][][][]float64) string {
	result := "MULTIPOLYGON("
	for _, polygon := range coordinates {
		result += "("
		for _, ring := range polygon {
			result += "("
			for _, point := range ring {
				result += fmt.Sprintf("%.6f %.6f,", point[0], point[1])
			}
			result = strings.TrimSuffix(result, ",") + "),"
		}
		result = strings.TrimSuffix(result, ",") + "),"
	}
	result = strings.TrimSuffix(result, ",") + ")"

	return result
}
