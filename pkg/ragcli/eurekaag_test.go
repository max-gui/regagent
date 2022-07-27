package ragcli

import (
	"fmt"
	"testing"
)

func TestEurekaInstance_String(t *testing.T) {
	tests := []struct {
		name   string
		fields EurekaInstance
		want   string
	}{
		// TODO: Add test cases.
		{name: "a", fields: EurekaInstance{App: "test", IpAddr: "1.2.3.4"}, want: "name:test,x-baggage-AF-env:,x-baggage-AF-region:,dc:,Address:1.2.3.4,Port:0,extaddress:,extport:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			euinst := EurekaInstance{
				InstanceId:                    tt.fields.InstanceId,
				HostName:                      tt.fields.HostName,
				App:                           tt.fields.App,
				Status:                        tt.fields.Status,
				Overriddenstatus:              tt.fields.Overriddenstatus,
				IpAddr:                        tt.fields.IpAddr,
				Port:                          tt.fields.Port,
				SecurePort:                    tt.fields.SecurePort,
				CountryId:                     tt.fields.CountryId,
				DataCenterInfo:                tt.fields.DataCenterInfo,
				LeaseInfo:                     tt.fields.LeaseInfo,
				Metadata:                      tt.fields.Metadata,
				HomePageUrl:                   tt.fields.HomePageUrl,
				StatusPageUrl:                 tt.fields.StatusPageUrl,
				HealthCheckUrl:                tt.fields.HealthCheckUrl,
				VipAddress:                    tt.fields.VipAddress,
				SecureVipAddress:              tt.fields.SecureVipAddress,
				IsCoordinatingDiscoveryServer: tt.fields.IsCoordinatingDiscoveryServer,
				LastUpdatedTimestamp:          tt.fields.LastUpdatedTimestamp,
				LastDirtyTimestamp:            tt.fields.LastDirtyTimestamp,
				ActionType:                    tt.fields.ActionType,
			}
			if got := fmt.Sprint(euinst); got != tt.want {
				t.Errorf("EurekaInstance.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
