package main

import (
    "context"
    "errors"
    "fmt"

    pb "github.com/dmanev/shippy-service-vessel/proto/vessel"
    "github.com/micro/go-micro"
)

type Repository interface {
    FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
    vessels []*pb.Vessel
}

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
    for _, vessel := range repo.vessels {
        if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
            return vessel, nil
        }
    }
    return nil, errors.New("No vessel found by this spec")
}

type service struct {
    repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

    //Find the next avialable vessel
    vessel, err := s.repo.FindAvailable(req)
    if err != nil {
        return err
    }

    // Set the vessel as part of the response message type
    res.Vessel = vessel
    return nil
}

func main() {
    vessel := []*pb.Vessel{
        &pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
    }
    repo := &VesselRepository{vessel}

    srv := micro.NewService(
        micro.Name("shippy.service.vessel"),
    )

    srv.Init()

    pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

    if err := srv.Run(); err != nil {
        fmt.Println(err)
    }
}
