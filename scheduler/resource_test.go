package scheduler

import (
	"testing"

	"github.com/stevedonovan/luar"
)

func TestWeight(t *testing.T) {
	lua = luar.Init()
	defer lua.Close()

	plug := `
function GetWeight(resource)
    return 1
end
`

	if err := loadPlugin(plug); err != nil {
		t.Fatal(err)
	}

	weight, err := Weight(&Resource{})
	if err != nil {
		t.Fatal(err)
	}
	if weight != 1 {
		t.Fatalf("expected a weight of 1 got %d", weight)
	}
}

func TestResourceValues(t *testing.T) {
	lua = luar.Init()
	defer lua.Close()

	plug := `
function GetWeight(resource)
    if resource.Cpus ~= 4 then
        error("cpus not equal to 4")
    end

    if resource.Memory ~= 1024 then
        error("memory not equal to 1024")
    end

    if resource.CpuProfile ~= "high" then 
        error("cpu profile not equal to high" .. resource.CpuProfile)
    end

    if resource.Disk ~= 1000 then 
        error("disk not equal to 1000")
    end
    return 1
end
`

	if err := loadPlugin(plug); err != nil {
		t.Fatal(err)
	}

	weight, err := Weight(&Resource{Cpus: 4, Memory: 1024, CpuProfile: "high", Disk: 1000})
	if err != nil {
		t.Fatal(err)
	}
	if weight != 1 {
		t.Fatalf("expected a weight of 1 got %d", weight)
	}
}
