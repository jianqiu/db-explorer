package db

import (
	"fmt"

	"github.com/vm-pool-server/db/models"
	"code.cloudfoundry.org/lager"

	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

//go:generate counterfeiter -o fake_bbs/fake_service_client.go . ServiceClient

type ServiceClient interface {
	VirtualGuestById(logger lager.Logger, cid int32) (*models.VirtualGuest, error)
}

type serviceClient struct {
	softlayerClient session
}

func NewServiceClient(username string, apikey string) ServiceClient {
	return &serviceClient{
		softlayerClient: session.New(username,apikey),
	}
}

func (slc *serviceClient) VirtualGuestById(logger lager.Logger, cid int32) (*models.VirtualGuest, error) {
	service := services.GetVirtualGuestService(slc.softlayerClient)
	mask := "mask[id, globalIdentifier, hostname, domain, fullyQualifiedDomainName, status.name, " +
	"powerState.name, activeTransaction, datacenter.name, " +
	"operatingSystem[softwareLicense[softwareDescription[name,version]],passwords[username,password]], " +
	" maxCpu, maxMemory, primaryIpAddress, primaryBackendIpAddress, " +
	"privateNetworkOnlyFlag, dedicatedAccountHostOnlyFlag, createDate, modifyDate, " +
	"billingItem[nextInvoiceTotalRecurringAmount, children[nextInvoiceTotalRecurringAmount]], notes, tagReferences.tag.name, networkVlans[id,vlanNumber,networkSpace]]"

	sl_virtual_guest, err := service.Id(cid).Mask(mask).GetObject()
	if err != nil {
		// Note: type assertion is only necessary for inspecting individual fields
		apiErr := err.(sl.Error)
		return nil, convertSoftlayerError(apiErr)
	}

	virtual_guest := models.VirtualGuest{}

	virtual_guest.Cid = *sl_virtual_guest.Id
	virtual_guest.Hostname = *sl_virtual_guest.Hostname
	virtual_guest.Ip = *sl_virtual_guest.PrimaryBackendIpAddress
	virtual_guest.Cpu = *sl_virtual_guest.MaxCpu
	virtual_guest.Memory = *sl_virtual_guest.MaxMemory

	if vlans := sl_virtual_guest.NetworkVlans; len(vlans) > 0 {
		for _, vlan := range vlans {
			if vlan.NetworkSpace != nil && vlan.VlanNumber != nil && vlan.Id != nil {
				switch *vlan.NetworkSpace {
				case "PRIVATE":
					virtual_guest.PrivateVlan = *vlan.Id
				case "PUBLIC":
					virtual_guest.PublicVlan = *vlan.Id
				default:
					return nil, models.NewError(models.Error_SoftLayerAPIError,"invalid vlan.Networkspace")
				}
			}
		}
	}

	return &virtual_guest, nil
}

func convertSoftlayerError(apiErr sl.Error) error {
	return models.NewError(models.Error_SoftLayerAPIError, fmt.Sprintf("HTTP Status Code: %d\n API Code: %s\n API Error: %s\n",
		apiErr.StatusCode,
		apiErr.Exception,
		apiErr.Message))
}
