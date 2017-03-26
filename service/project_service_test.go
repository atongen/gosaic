package service

import (
	"testing"

	"github.com/atongen/gosaic/model"
)

func TestProjectServiceInsert(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	p1 := model.Project{
		Name: "test1",
	}

	err := projectService.Insert(&p1)
	if err != nil {
		t.Fatalf("Error inserting project: %s\n", err.Error())
	}

	if p1.Id == int64(0) {
		t.Fatalf("Inserted project id not set")
	}

	p2, err := projectService.Get(p1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted project: %s\n", err.Error())
	} else if p2 == nil {
		t.Fatalf("Project not inserted\n")
	}

	if p1.Id != p2.Id ||
		p1.Name != p2.Name {
		t.Fatalf("Inserted project (%+v) does not match: %+v\n", p2, p1)
	}

	if p1.CreatedAt.IsZero() ||
		p2.CreatedAt.IsZero() {
		t.Fatal("Project created at not set")
	} else if p1.CreatedAt.Unix() != p2.CreatedAt.Unix() {
		t.Fatal("Inserted project created at does not match")
	}
}

func TestProjectServiceGetOneBy(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	c1 := model.Project{
		Name:    "testme1",
		MacroId: int64(234),
	}

	err := projectService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting project: %s\n", err.Error())
	}

	c2, err := projectService.GetOneBy("macro_id = ? and name = ?", c1.MacroId, c1.Name)
	if err != nil {
		t.Fatalf("Error getting inserted project: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Project not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.MacroId != c2.MacroId ||
		c1.Name != c2.Name ||
		c1.CreatedAt.Unix() != c2.CreatedAt.Unix() {
		t.Fatalf("Inserted project (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestProjectServiceGetOneByNot(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	c, err := projectService.GetOneBy("macro_id = ? and name = ?", int64(123), "not a valid name")
	if err != nil {
		t.Fatalf("Error getting inserted project: %s\n", err.Error())
	}

	if c != nil {
		t.Fatal("Project found when should not exist")
	}
}

func TestProjectServiceExistsBy(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	c1 := model.Project{
		MacroId: int64(321),
		Name:    "testme1",
	}

	err := projectService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting project: %s\n", err.Error())
	}

	found, err := projectService.ExistsBy("macro_id = ? and name = ?", c1.MacroId, "testme1")
	if err != nil {
		t.Fatalf("Error getting inserted project: %s\n", err.Error())
	} else if !found {
		t.Fatalf("Project not inserted\n")
	}
}

func TestProjectServiceExistsByNot(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	found, err := projectService.ExistsBy("macro_id = ? and name = ?", int64(123), "not a valid name")
	if err != nil {
		t.Fatalf("Error getting inserted project: %s\n", err.Error())
	} else if found {
		t.Fatal("Project found when should not exist")
	}
}

func TestProjectServiceUpdate(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	updateProject := model.Project{
		Name: "testme1",
	}
	err := projectService.Insert(&updateProject)
	if err != nil {
		t.Fatalf("Error inserting project: %s\n", err.Error())
	}

	newName := "testme2"
	updateProject.Name = newName

	num, err := projectService.Update(&updateProject)
	if err != nil {
		t.Error("Error updating project", err)
	}

	if num == 0 {
		t.Error("Nothing was updated")
	}

	project2, err := projectService.Get(updateProject.Id)
	if err != nil {
		t.Error("Error finding update project", err)
	}

	if project2.Name != newName {
		t.Error("project was not updated")
	}
}

func TestProjectServiceFindAll(t *testing.T) {
	setTestServiceFactory()
	projectService := serviceFactory.MustProjectService()
	defer projectService.Close()

	c1 := model.Project{
		Name: "test1",
	}

	err := projectService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting project: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted project id not set")
	}

	projects, err := projectService.FindAll("id asc")
	if err != nil {
		t.Fatalf("Error finding all projects: %s\n", err.Error())
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d\n", len(projects))
	}

	c2 := projects[0]

	if c1.Id != c2.Id ||
		c1.Name != c2.Name {
		t.Fatalf("Inserted project (%+v) does not match: %+v\n", c2, c1)
	}

	if c1.CreatedAt.IsZero() ||
		c2.CreatedAt.IsZero() {
		t.Fatal("Project created at not set")
	} else if c1.CreatedAt.Unix() != c2.CreatedAt.Unix() {
		t.Fatal("Inserted project created at does not match")
	}
}
