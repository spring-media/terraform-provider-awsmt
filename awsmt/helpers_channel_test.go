package awsmt

import (
	"testing"

	"terraform-provider-mediatailor/awsmt/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
)

// TestWriteChannelToPlan_Audiences verifies the audiences assignment logic in writeChannelToPlan.
//
// The rule: only assign channel.Audiences when the plan had audiences configured OR the API
// returns a non-empty list. This handles two distinct cases:
//  1. No audiences in config (plan nil) + API returns [] → keep nil (avoid null→[] inconsistency)
//  2. Audiences in config → API returns [] (all removed) → write [] to state (reflect removal)
func TestWriteChannelToPlan_Audiences(t *testing.T) {
	t.Run("sets audiences when API returns non-empty slice", func(t *testing.T) {
		model := models.ChannelModel{}
		output := mediatailor.CreateChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{"audience1", "audience2"},
		}

		result := writeChannelToPlan(model, output)

		if len(result.Audiences) != 2 {
			t.Fatalf("expected 2 audiences, got %d", len(result.Audiences))
		}
		if result.Audiences[0] != "audience1" || result.Audiences[1] != "audience2" {
			t.Errorf("unexpected audiences: %v", result.Audiences)
		}
	})

	t.Run("keeps nil when plan has no audiences and API returns nil", func(t *testing.T) {
		model := models.ChannelModel{} // plan: audiences not set
		output := mediatailor.CreateChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   nil,
		}

		result := writeChannelToPlan(model, output)

		if result.Audiences != nil {
			t.Errorf("expected Audiences to remain nil, got %v", result.Audiences)
		}
	})

	t.Run("keeps nil when plan has no audiences and API returns empty slice", func(t *testing.T) {
		model := models.ChannelModel{} // plan: audiences not set
		output := mediatailor.CreateChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{},
		}

		result := writeChannelToPlan(model, output)

		if result.Audiences != nil {
			t.Errorf("expected Audiences to remain nil when config has no audiences, got %v", result.Audiences)
		}
	})

	t.Run("clears audiences when plan had audiences but API returns empty (all removed)", func(t *testing.T) {
		model := models.ChannelModel{Audiences: []string{"old-audience"}}
		output := mediatailor.CreateChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{},
		}

		result := writeChannelToPlan(model, output)

		if len(result.Audiences) != 0 {
			t.Errorf("expected Audiences to be empty after all removed, got %v", result.Audiences)
		}
	})
}

// TestWriteChannelToState_Audiences tests the same behaviour for writeChannelToState,
// which is called on Read and populates Terraform state from a DescribeChannel response.
func TestWriteChannelToState_Audiences(t *testing.T) {
	t.Run("sets audiences when API returns non-empty slice", func(t *testing.T) {
		model := models.ChannelModel{}
		output := mediatailor.DescribeChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{"audience1", "audience2"},
		}

		result := writeChannelToState(model, output)

		if len(result.Audiences) != 2 {
			t.Fatalf("expected 2 audiences, got %d", len(result.Audiences))
		}
		if result.Audiences[0] != "audience1" || result.Audiences[1] != "audience2" {
			t.Errorf("unexpected audiences: %v", result.Audiences)
		}
	})

	t.Run("keeps nil when plan has no audiences and API returns nil", func(t *testing.T) {
		model := models.ChannelModel{}
		output := mediatailor.DescribeChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   nil,
		}

		result := writeChannelToState(model, output)

		if result.Audiences != nil {
			t.Errorf("expected Audiences to remain nil, got %v", result.Audiences)
		}
	})

	t.Run("keeps nil when plan has no audiences and API returns empty slice", func(t *testing.T) {
		model := models.ChannelModel{}
		output := mediatailor.DescribeChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{},
		}

		result := writeChannelToState(model, output)

		if result.Audiences != nil {
			t.Errorf("expected Audiences to remain nil when config has no audiences, got %v", result.Audiences)
		}
	})

	t.Run("clears audiences when plan had audiences but API returns empty (all removed)", func(t *testing.T) {
		model := models.ChannelModel{Audiences: []string{"old-audience"}}
		output := mediatailor.DescribeChannelOutput{
			ChannelName: aws.String("test"),
			Audiences:   []string{},
		}

		result := writeChannelToState(model, output)

		if len(result.Audiences) != 0 {
			t.Errorf("expected Audiences to be empty after all removed, got %v", result.Audiences)
		}
	})
}
