package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"testing"
)

func TestDeleteVodSources(t *testing.T) {
	// arrange: create source location, add vod sources
	conn := testAccProvider.Meta().(*mediatailor.MediaTailor)
	sourceLocationName := aws.String("source_location_test_vod_deletion")
	httpConfiguration := &mediatailor.HttpConfiguration{BaseUrl: aws.String("https://www.example.com")}
	vodSourceName := aws.String("vod_source_test_vod_deletion")
	if _, err := conn.CreateSourceLocation(&mediatailor.CreateSourceLocationInput{SourceLocationName: sourceLocationName, HttpConfiguration: httpConfiguration}); err != nil {
		t.Fatalf(`Error creating source location: %v`, err)
	}
	httpPackageConfiguration := &mediatailor.HttpPackageConfiguration{Path: aws.String("/"), SourceGroup: aws.String("default"), Type: aws.String("HLS")}
	if _, err := conn.CreateVodSource(&mediatailor.CreateVodSourceInput{SourceLocationName: sourceLocationName, VodSourceName: vodSourceName, HttpPackageConfigurations: []*mediatailor.HttpPackageConfiguration{httpPackageConfiguration}}); err != nil {
		t.Fatalf(`Error creating vod source: %v`, err)
	}
	// act: delete vod sources
	if err := deleteVodSources(sourceLocationName, conn); err != nil {
		t.Fatalf(`Error deleting vod sources: %v`, err)
	}
	// assert: vod source has been deleted
	if _, err := conn.DescribeVodSource(&mediatailor.DescribeVodSourceInput{SourceLocationName: sourceLocationName, VodSourceName: vodSourceName}); err == nil {
		t.Fatalf(`VodSource was not deleted`)
	}
	// cleanup
	if _, err := conn.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: sourceLocationName}); err != nil {
		t.Fatalf(`Error cleaning up: %v`, err)
	}
}

func TestDeleteLiveSources(t *testing.T) {
	// arrange: create source location, add vod sources
	conn := testAccProvider.Meta().(*mediatailor.MediaTailor)
	sourceLocationName := aws.String("source_location_test_vod_deletion")
	httpConfiguration := &mediatailor.HttpConfiguration{BaseUrl: aws.String("https://www.example.com")}
	liveSourceName := aws.String("vod_source_test_vod_deletion")
	if _, err := conn.CreateSourceLocation(&mediatailor.CreateSourceLocationInput{SourceLocationName: sourceLocationName, HttpConfiguration: httpConfiguration}); err != nil {
		t.Fatalf(`Error creating source location: %v`, err)
	}
	httpPackageConfiguration := &mediatailor.HttpPackageConfiguration{Path: aws.String("/"), SourceGroup: aws.String("default"), Type: aws.String("HLS")}
	if _, err := conn.CreateLiveSource(&mediatailor.CreateLiveSourceInput{SourceLocationName: sourceLocationName, LiveSourceName: liveSourceName, HttpPackageConfigurations: []*mediatailor.HttpPackageConfiguration{httpPackageConfiguration}}); err != nil {
		t.Fatalf(`Error creating live source: %v`, err)
	}
	// act: delete vod sources
	if err := deleteLiveSources(sourceLocationName, conn); err != nil {
		t.Fatalf(`Error deleting live sources: %v`, err)
	}
	// assert: vod source has been deleted
	if _, err := conn.DescribeLiveSource(&mediatailor.DescribeLiveSourceInput{SourceLocationName: sourceLocationName, LiveSourceName: liveSourceName}); err == nil {
		t.Fatalf(`LiveSource was not deleted`)
	}
	// cleanup
	if _, err := conn.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: sourceLocationName}); err != nil {
		t.Fatalf(`Error cleaning up: %v`, err)
	}
}
