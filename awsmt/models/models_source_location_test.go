package models

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccessConfigurationModel", func() {
	var (
		config1 *AccessConfigurationModel
		config2 *AccessConfigurationModel
	)

	BeforeEach(func() {
		config1 = &AccessConfigurationModel{
			AccessType: aws.String("token"),
			SecretsManagerAccessTokenConfiguration: &SecretsManagerAccessTokenConfigurationModel{
				HeaderName:      aws.String("Authorization"),
				SecretArn:       aws.String("arn:aws:secretsmanager:region:123456789012:secret:my-secret"),
				SecretStringKey: aws.String("api-key"),
			},
		}
	})

	Context("Equal method", func() {
		It("should return true when comparing identical configurations", func() {
			config2 = &AccessConfigurationModel{
				AccessType: aws.String("token"),
				SecretsManagerAccessTokenConfiguration: &SecretsManagerAccessTokenConfigurationModel{
					HeaderName:      aws.String("Authorization"),
					SecretArn:       aws.String("arn:aws:secretsmanager:region:123456789012:secret:my-secret"),
					SecretStringKey: aws.String("api-key"),
				},
			}
			Expect(config1.Equal(config2)).To(BeTrue())
		})

		It("should return true when both are nil", func() {
			var nilConfig1, nilConfig2 *AccessConfigurationModel
			Expect(nilConfig1.Equal(nilConfig2)).To(BeTrue())
		})

		It("should return false when one is nil", func() {
			var nilConfig *AccessConfigurationModel
			Expect(config1.Equal(nilConfig)).To(BeFalse())
			Expect(nilConfig.Equal(config1)).To(BeFalse())
		})

		It("should return false when AccessType differs", func() {
			config2 = config1.copy()
			config2.AccessType = aws.String("different")
			Expect(config1.Equal(config2)).To(BeFalse())
		})

		It("should return false when SecretsManagerAccessTokenConfiguration differs", func() {
			config2 = config1.copy()
			config2.SecretsManagerAccessTokenConfiguration.HeaderName = aws.String("different")
			Expect(config1.Equal(config2)).To(BeFalse())
		})

		It("should handle nil SecretsManagerAccessTokenConfiguration", func() {
			config2 = config1.copy()
			config2.SecretsManagerAccessTokenConfiguration = nil
			Expect(config1.Equal(config2)).To(BeFalse())

			config1.SecretsManagerAccessTokenConfiguration = nil
			Expect(config1.Equal(config2)).To(BeTrue())
		})
	})
})

var _ = Describe("SecretsManagerAccessTokenConfigurationModel", func() {
	var (
		config1 *SecretsManagerAccessTokenConfigurationModel
		config2 *SecretsManagerAccessTokenConfigurationModel
	)

	BeforeEach(func() {
		config1 = &SecretsManagerAccessTokenConfigurationModel{
			HeaderName:      aws.String("Authorization"),
			SecretArn:       aws.String("arn:aws:secretsmanager:region:123456789012:secret:my-secret"),
			SecretStringKey: aws.String("api-key"),
		}
	})

	Context("Equal method", func() {
		It("should return true when comparing identical configurations", func() {
			config2 = &SecretsManagerAccessTokenConfigurationModel{
				HeaderName:      aws.String("Authorization"),
				SecretArn:       aws.String("arn:aws:secretsmanager:region:123456789012:secret:my-secret"),
				SecretStringKey: aws.String("api-key"),
			}
			Expect(config1.Equal(config2)).To(BeTrue())
		})

		It("should return true when both are nil", func() {
			var nilConfig1, nilConfig2 *SecretsManagerAccessTokenConfigurationModel
			Expect(nilConfig1.Equal(nilConfig2)).To(BeTrue())
		})

		It("should return false when one is nil", func() {
			var nilConfig *SecretsManagerAccessTokenConfigurationModel
			Expect(config1.Equal(nilConfig)).To(BeFalse())
			Expect(nilConfig.Equal(config1)).To(BeFalse())
		})

		It("should return false when HeaderName differs", func() {
			config2 = config1.copy()
			config2.HeaderName = aws.String("different")
			Expect(config1.Equal(config2)).To(BeFalse())
		})

		It("should return false when SecretArn differs", func() {
			config2 = config1.copy()
			config2.SecretArn = aws.String("different")
			Expect(config1.Equal(config2)).To(BeFalse())
		})

		It("should return false when SecretStringKey differs", func() {
			config2 = config1.copy()
			config2.SecretStringKey = aws.String("different")
			Expect(config1.Equal(config2)).To(BeFalse())
		})
	})
})

// Helper method to create a deep copy for testing
func (a *AccessConfigurationModel) copy() *AccessConfigurationModel {
	if a == nil {
		return nil
	}
	tmp := &AccessConfigurationModel{
		AccessType: aws.String(*a.AccessType),
	}
	if a.SecretsManagerAccessTokenConfiguration != nil {
		tmp.SecretsManagerAccessTokenConfiguration = a.SecretsManagerAccessTokenConfiguration.copy()
	}
	return tmp
}

func (s *SecretsManagerAccessTokenConfigurationModel) copy() *SecretsManagerAccessTokenConfigurationModel {
	if s == nil {
		return nil
	}
	return &SecretsManagerAccessTokenConfigurationModel{
		HeaderName:      aws.String(*s.HeaderName),
		SecretArn:       aws.String(*s.SecretArn),
		SecretStringKey: aws.String(*s.SecretStringKey),
	}
}
