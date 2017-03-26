package service

import (
	"testing"

	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"

	_ "github.com/mattn/go-sqlite3"
)

func setupAspectServiceTest() {
	setTestServiceFactory()
	aspectService := serviceFactory.MustAspectService()

	aspect = model.Aspect{Columns: 1, Rows: 1}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}
}

func TestAspectServiceGet(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	aspect2, err := aspectService.Get(aspect.Id)
	if err != nil {
		t.Fatal("Error finding aspect by id", err)
	}

	if aspect.Id != aspect2.Id ||
		aspect.Columns != aspect2.Columns ||
		aspect.Rows != aspect2.Rows {
		t.Fatal("Found aspect does not match data")
	}
}

func TestAspectServiceGetMissing(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	aspect2, err := aspectService.Get(1234)
	if err != nil {
		t.Fatal("Error finding aspect by id", err)
	}

	if aspect2 != nil {
		t.Fatal("Found non-existent aspect")
	}
}

func TestAspectServiceFind(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	a1, err := aspectService.Find(333, 333)
	if err != nil {
		t.Fatal("Error finding 1:1 aspect: %s\n", err.Error())
	}

	if a1 == nil {
		t.Fatal("Expected to find 1:1 aspect")
	}

	a2, err := aspectService.Find(13, 27)
	if err != nil {
		t.Fatal("Error finding 13:27 aspect: %s\n", err.Error())
	}

	if a2 != nil {
		t.Fatal("Did not expect to find 13:27 aspect")
	}
}

func TestAspectServiceCreate(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	a1, err := aspectService.Create(123, 234)
	if err != nil {
		t.Fatal("Error creating 123:234 aspect: %s\n", err.Error())
	}

	if a1 == nil {
		t.Fatal("Expected to create 123:234 aspect")
	}

	a2, err := aspectService.Find(123, 234)
	if err != nil {
		t.Fatal("Error finding 123:234 aspect: %s\n", err.Error())
	}

	if a2 == nil {
		t.Fatal("Expect to find 123:234 aspect")
	}

	a3, err := aspectService.Create(123, 234)
	if err == nil {
		t.Fatal("Expected error creating duplicate 123:234 aspect")
	}

	if a3 != nil {
		t.Fatal("Expected no result when creating duplicate aspect, got: %+v\n", a3)
	}
}

func TestAspectServiceFindOrCreate(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	n1, err := aspectService.Count()
	if err != nil {
		t.Fatal("Unable to count aspects")
	}

	a1, err := aspectService.FindOrCreate(100, 100)
	if err != nil {
		t.Fatal("Unable to find or create 100x100 aspect")
	}

	a2, err := aspectService.FindOrCreate(200, 200)
	if err != nil {
		t.Fatal("Unable to find or create 200x200 aspect")
	}

	a3, err := aspectService.FindOrCreate(300, 300)
	if err != nil {
		t.Fatal("Unable to find or create 300x300 aspect")
	}

	_, err = aspectService.FindOrCreate(400, 600)
	if err != nil {
		t.Fatal("Unable to find or create 400x600 aspect")
	}

	n2, err := aspectService.Count()
	if err != nil {
		t.Fatal("Unable to re-count aspects")
	}

	if n1 != n2-1 {
		t.Fatal("Created aspect when shouldn't have")
	}

	if a1.Id != a2.Id {
		t.Fatal("Aspects not equal")
	}

	if a2.Id != a3.Id {
		t.Fatal("Aspects not equal")
	}
}

func TestAspectServiceFindIn(t *testing.T) {
	setupAspectServiceTest()
	aspectService := serviceFactory.MustAspectService()
	defer aspectService.Close()

	a1, err := aspectService.FindOrCreate(20, 30)
	if err != nil {
		t.Fatal("Unable to find or create 20x30 aspect")
	}

	a2, err := aspectService.FindOrCreate(30, 40)
	if err != nil {
		t.Fatal("Unable to find or create 30x40 aspect")
	}

	a3, err := aspectService.FindOrCreate(40, 50)
	if err != nil {
		t.Fatal("Unable to find or create 40x50 aspect")
	}

	ids := make([]int64, 4)
	ids[0] = aspect.Id
	ids[1] = a1.Id
	ids[2] = a2.Id
	ids[3] = a3.Id

	aspects, err := aspectService.FindIn(ids)
	if err != nil {
		t.Fatalf("Error finding in aspect ids: %s\n", err.Error())
	}

	if len(aspects) != 4 {
		t.Fatalf("Expected 4 aspects, got %d\n", len(aspects))
	}

	for _, aspect := range []*model.Aspect{&aspect, a1, a2, a3} {
		if !util.SliceContainsInt64(ids, aspect.Id) {
			t.Fatalf("Expected %d to be in ids slice", aspect.Id)
		}
	}

	for _, aspect := range aspects {
		if aspect == nil {
			t.Fatal("Received nil aspect from FindIn")
		}
	}
}
