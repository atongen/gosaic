package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupGidxPartialServiceTest() {
	setTestServiceFactory()
	aspectService := serviceFactory.MustAspectService()
	gidxService := serviceFactory.MustGidxService()

	aspect = model.Aspect{Columns: 87, Rows: 128}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       696,
		Height:      1024,
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		panic(err)
	}
}

func TestGidxPartialServiceInsert(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	if mp.Id == int64(0) {
		t.Fatalf("Inserted gidx partial id not set")
	}

	mp2, err := gidxPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting inserted gidx partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Gidx partial not inserted\n")
	}

	if mp.Id != mp2.Id ||
		mp.GidxId != mp2.GidxId ||
		mp.AspectId != mp2.AspectId {
		t.Fatal("Inserted gidx partial data does not match")
	}

	if len(mp2.Pixels) != 1 {
		t.Fatal("Gidx partial pixels not serialized correctly")
	}

	plab := mp2.Pixels[0]

	if plab.L != 0.4 &&
		plab.A != 0.5 &&
		plab.B != 0.6 &&
		plab.Alpha != 0.0 {
		t.Fatal("Gidx partial pixel data is not correct")
	}
}

func TestGidxPartialServiceBulkInsert(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	aspectService := serviceFactory.MustAspectService()
	defer gidxPartialService.Close()

	aspect2, err := aspectService.FindOrCreate(5, 7)
	if err != nil {
		t.Fatalf("Failed to create 2nd aspect: %s\n", err.Error())
	}

	gidxPartials := []*model.GidxPartial{
		&model.GidxPartial{
			GidxId:   gidx.Id,
			AspectId: aspect.Id,
			Pixels: []*model.Lab{
				&model.Lab{
					L:     0.4,
					A:     0.5,
					B:     0.6,
					Alpha: 0.0,
				},
			},
		},
		&model.GidxPartial{
			GidxId:   gidx.Id,
			AspectId: aspect2.Id,
			Pixels: []*model.Lab{
				&model.Lab{
					L:     0.7,
					A:     0.8,
					B:     0.9,
					Alpha: 0.1,
				},
			},
		},
	}

	for _, gp := range gidxPartials {
		err = gp.EncodePixels()
		if err != nil {
			t.Fatalf("Error encoding gidx partial pixels: %s\n", err.Error())
		}
	}

	num, err := gidxPartialService.BulkInsert(gidxPartials)
	if err != nil {
		t.Fatalf("Error bulk inserting gidx partial: %s\n", err.Error())
	}

	if num != 2 {
		t.Fatalf("Expected 2 bulk insert gidx partial, but got: %d\n", num)
	}

	found := make([]*model.GidxPartial, 2)
	f0, err := gidxPartialService.Find(&gidx, &aspect)
	if err != nil {
		t.Fatalf("Error finding bulk inserted gidx partial: %s\n", err.Error())
	}

	f1, err := gidxPartialService.Find(&gidx, aspect2)
	if err != nil {
		t.Fatalf("Error finding bulk inserted gidx partial: %s\n", err.Error())
	}
	found[0] = f0
	found[1] = f1

	for i := 0; i < 2; i++ {
		f := found[i]
		// gp id is still zero
		gp := gidxPartials[i]

		if f.GidxId != gp.GidxId ||
			f.AspectId != gp.AspectId {
			t.Fatal("Bulk inserted gidx partial data does not match")
		}

		if len(f.Pixels) == 0 {
			t.Fatal("Expected found bulk inserted gidx partial pixel data to have length greater than zero")
		}

		if len(f.Pixels) != len(gp.Pixels) {
			t.Fatal("Found bulk inserted pixel data does not match")
		}

		lab1 := f.Pixels[0]
		lab2 := gp.Pixels[0]

		if lab1.L == 0.0 &&
			lab1.A == 0.0 &&
			lab1.B == 0.0 &&
			lab1.Alpha == 0.0 {
			t.Fatal("Bulk inserted gidx partial pixel data is empty")
		}

		if lab1.L != lab2.L &&
			lab2.A != lab2.A &&
			lab2.B != lab2.B &&
			lab2.Alpha != lab2.Alpha {
			t.Fatal("Bulk inserted gidx partial pixel data not match")
		}
	}
}

func TestGidxPartialServiceUpdate(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	mp.Pixels[0].L = 0.75
	err = gidxPartialService.Update(&mp)
	if err != nil {
		t.Fatalf("Error updating gidx partial: %s\n", err.Error())
	}

	mp2, err := gidxPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting updated gidx partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Gidx partial not inserted\n")
	}

	if mp2.Pixels[0].L != 0.75 {
		t.Fatal("Updated gidx partial data does not match")
	}
}

func TestGidxPartialServiceDelete(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	err = gidxPartialService.Delete(&mp)
	if err != nil {
		t.Fatalf("Error deleting gidx partial: %s\n", err.Error())
	}

	mp2, err := gidxPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting deleted gidx partial: %s\n", err.Error())
	} else if mp2 != nil {
		t.Fatalf("Gidx partial not deleted\n")
	}
}

func TestGidxPartialServiceGetOneBy(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	mp2, err := gidxPartialService.GetOneBy("gidx_id", mp.GidxId)
	if err != nil {
		t.Fatalf("Error getting one by gidx partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Gidx partial not found by\n")
	}

	if mp2.GidxId != mp.GidxId {
		t.Fatal("Gidx partial gidx id does not match")
	}
}

func TestGidxPartialServiceExistsBy(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	found, err := gidxPartialService.ExistsBy("gidx_id", mp.GidxId)
	if err != nil {
		t.Fatalf("Error getting one by gidx partial: %s\n", err.Error())
	}

	if !found {
		t.Fatalf("Gidx partial not exists by\n")
	}
}

func TestGidxPartialServiceCount(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	num, err := gidxPartialService.Count()
	if err != nil {
		t.Fatalf("Error counting gidx partial: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Gidx partial count incorrect\n")
	}
}

func TestGidxPartialServiceCountBy(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	mp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	num, err := gidxPartialService.CountBy("gidx_id = ?", gidx.Id)
	if err != nil {
		t.Fatalf("Error counting by gidx partial: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Gidx partial count incorrect\n")
	}
}

func TestGidxPartialServiceCountForMacro(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	gp := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := gidxPartialService.Insert(&gp)
	if err != nil {
		t.Fatalf("Error inserting gidx partial: %s\n", err.Error())
	}

	num, err := gidxPartialService.CountForMacro(&macro)
	if err != nil {
		t.Fatalf("Error counting gidx partial for macro: %s\n", err.Error())
	}

	// we have no marco partials at this point
	if num != int64(0) {
		t.Fatalf("Expected 0 gidx partial count for macro, but got %d\n", num)
	}
}

func TestGidxPartialServiceFindOrCreate(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	gidxPartial, err := gidxPartialService.FindOrCreate(&gidx, &aspect)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate gidxPartial: %s\n", err.Error())
	}

	if gidxPartial.GidxId != gidx.Id {
		t.Errorf("gidxPartial.GidxId was %d, expected %d\n", gidxPartial.GidxId, gidx.Id)
	}

	if gidxPartial.AspectId != aspect.Id {
		t.Errorf("gidxPartial.AspectId was %d, expected %d\n", gidxPartial.AspectId, aspect.Id)
	}

	if len(gidxPartial.Data) == 0 {
		t.Error("gidxPartial.Data was empty")
	}

	numPixels := len(gidxPartial.Pixels)
	if numPixels != 100 {
		t.Errorf("gidxPartial.Pixels len was %d, expected %d\n", numPixels, 100)
	}

	for i, pix := range gidxPartial.Pixels {
		if pix.L == 0.0 && pix.A == 0.0 && pix.B == 0.0 && pix.Alpha == 0.0 {
			t.Errorf("pixel %d was empty\n", i)
		}
	}
}

func TestGidxPartialServiceFindMissing(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	gidxs, err := gidxPartialService.FindMissing(&aspect, "gidx.id", 100, 0)
	if err != nil {
		t.Fatalf("Failed to FindMissing gidxPartial: %s\n", err.Error())
	}

	if len(gidxs) != 1 {
		t.Fatalf("Expected 1 Missing gidxPartial, got %d\n", len(gidxs))
	}

	if gidxs[0].Id != gidx.Id {
		t.Errorf("Expected missing gidx id %d, got %d\n", gidx.Id, gidxs[0].Id)
	}

	_, err = gidxPartialService.Create(&gidx, &aspect)
	if err != nil {
		t.Fatalf("Failed to Create gidxPartial: %s\n", err.Error())
	}

	gidxs, err = gidxPartialService.FindMissing(&aspect, "gidx.id", 100, 0)
	if err != nil {
		t.Fatalf("Failed to FindMissing gidxPartial: %s\n", err.Error())
	}

	if len(gidxs) != 0 {
		t.Fatalf("Expected 0 Missing gidxPartial, got %d\n", len(gidxs))
	}
}

func TestGidxPartialServiceCountMissing(t *testing.T) {
	setupGidxPartialServiceTest()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	defer gidxPartialService.Close()

	num, err := gidxPartialService.CountMissing([]*model.Aspect{&aspect})
	if err != nil {
		t.Fatalf("Failed to CountMissing gidxPartial: %s\n", err.Error())
	}

	if num != 1 {
		t.Fatalf("Expected 1 Missing gidxPartial, got %d\n", num)
	}

	_, err = gidxPartialService.Create(&gidx, &aspect)
	if err != nil {
		t.Fatalf("Failed to Create gidxPartial: %s\n", err.Error())
	}

	num, err = gidxPartialService.CountMissing([]*model.Aspect{&aspect})
	if err != nil {
		t.Fatalf("Failed to CountMissing gidxPartial: %s\n", err.Error())
	}

	if num != 0 {
		t.Fatalf("Expected 0 Missing gidxPartial, got %d\n", num)
	}
}
