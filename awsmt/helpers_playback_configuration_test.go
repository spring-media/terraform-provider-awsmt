package awsmt

import (
	"testing"

	"terraform-provider-mediatailor/awsmt/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
)

func TestAddDashConfigurationToModel(t *testing.T) {
	t.Run("maps ManifestEndpointPrefix from its own field, not from MpdLocation", func(t *testing.T) {
		endpointPrefix := "https://endpoint.prefix/"
		mpdLocation := "https://mpd.location/manifest.mpd"

		builder := &putPlaybackConfigurationModelbuilder{
			model: &models.PlaybackConfigurationModel{},
			output: mediatailor.PutPlaybackConfigurationOutput{
				DashConfiguration: &awsTypes.DashConfiguration{
					ManifestEndpointPrefix: aws.String(endpointPrefix),
					MpdLocation:            aws.String(mpdLocation),
				},
			},
			isResource: false,
		}

		builder.addDashConfigurationToModel()

		if builder.model.DashConfiguration == nil {
			t.Fatal("expected DashConfiguration to be set on the model")
		}
		// Regression: the field used to be populated from MpdLocation.
		if got := builder.model.DashConfiguration.ManifestEndpointPrefix.ValueString(); got != endpointPrefix {
			t.Errorf("ManifestEndpointPrefix = %q, want %q (must not be sourced from MpdLocation)", got, endpointPrefix)
		}
		if builder.model.DashConfiguration.MpdLocation == nil || *builder.model.DashConfiguration.MpdLocation != mpdLocation {
			t.Errorf("MpdLocation = %v, want %q", builder.model.DashConfiguration.MpdLocation, mpdLocation)
		}
	})

	t.Run("leaves ManifestEndpointPrefix null and does not panic when the output field is nil", func(t *testing.T) {
		builder := &putPlaybackConfigurationModelbuilder{
			model: &models.PlaybackConfigurationModel{},
			output: mediatailor.PutPlaybackConfigurationOutput{
				DashConfiguration: &awsTypes.DashConfiguration{
					ManifestEndpointPrefix: nil,
				},
			},
			isResource: false,
		}

		builder.addDashConfigurationToModel()

		if builder.model.DashConfiguration == nil {
			t.Fatal("expected DashConfiguration to be set on the model")
		}
		if !builder.model.DashConfiguration.ManifestEndpointPrefix.IsNull() {
			t.Errorf("ManifestEndpointPrefix = %q, want null", builder.model.DashConfiguration.ManifestEndpointPrefix.ValueString())
		}
	})
}
