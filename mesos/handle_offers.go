package mesos

import (
	"github.com/sirupsen/logrus"

	mesosproto "../proto"
)

func defaultResources() []*mesosproto.Resource {
	CPU := "cpus"
	MEM := "mem"
	cpu := config.ResCPU
	mem := config.ResMEM

	return []*mesosproto.Resource{
		{
			Name:   &CPU,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &cpu},
		},
		{
			Name:   &MEM,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &mem},
		},
	}
}

// HandleOffers will handle the offers event of mesos
func HandleOffers(offers *mesosproto.Event_Offers) error {
	offerIds := []*mesosproto.OfferID{}
	for _, offer := range offers.Offers {
		offerIds = append(offerIds, offer.Id)
	}

	select {
	case cmd := <-config.CommandChan:

		firstOffer := offers.Offers[0]
		agentID := offerIds[0].Value

		var taskInfo []*mesosproto.TaskInfo
		RefuseSeconds := 10.0

		switch cmd.ContainerType {
		case "NONE":
			taskInfo, _ = prepareTaskInfoExecuteCommand(firstOffer.AgentId, cmd)
		case "MESOS":
			taskInfo, _ = prepareTaskInfoExecuteContainer(firstOffer.AgentId, cmd)
		case "DOCKER":
			taskInfo, _ = prepareTaskInfoExecuteContainer(firstOffer.AgentId, cmd)
		}

		logrus.Debug("HandleOffers cmd: ", taskInfo)

		accept := &mesosproto.Call{
			Type: mesosproto.Call_ACCEPT.Enum(),
			Accept: &mesosproto.Call_Accept{
				OfferIds: []*mesosproto.OfferID{{
					Value: agentID,
				}},
				Filters: &mesosproto.Filters{
					RefuseSeconds: &RefuseSeconds,
				},
				Operations: []*mesosproto.Offer_Operation{{
					Type: mesosproto.Offer_Operation_LAUNCH.Enum(),
					Launch: &mesosproto.Offer_Operation_Launch{
						TaskInfos: taskInfo,
					}}}}}

		logrus.Debug("Offer Accept: ", offerIds)
		return Call(accept)
	default:
		logrus.Debug("Offer Decline: ", offerIds)
		decline := &mesosproto.Call{
			Type:    mesosproto.Call_DECLINE.Enum(),
			Decline: &mesosproto.Call_Decline{OfferIds: offerIds},
		}
		return Call(decline)
	}
}
