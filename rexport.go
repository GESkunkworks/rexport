package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	etypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

var version string
var conf Config
var ASCII_ART string
var wgv sync.WaitGroup

func setArt() {
	ASCII_ART = fmt.Sprintf(`
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡷⡷⣿⣿⣿⣿⣿⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠓⢉⣌⣬⣮⣌⢜⠹⣷⣿⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⣀⣾⣿⠁⣌⣌⠌⣳⣿⢎⡱⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀⣿⣿⡿⠀⡷⡷⠇⡰⣿⣿⠎⣱⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢏⠰⣿⣿⠀⣮⠦⠢⣮⠀⣿⣿⠏⣰⣿⣿⣿
⣿⣿⣿⣿⣿⣿⡷⠳⢙⣉⣌⣌⢌⢙⠱⡳⣿⣿⣿⣿⣿⡿⠷⣉⠌⣳⣿⠀⣽⣌⣈⣿⠀⣿⡿⠁⣾⣿⣿⣿
⣿⣿⣿⣿⠷⣉⣮⣿⣿⡷⠷⡳⡷⣷⣿⣮⢌⠱⣿⠷⢓⣌⣾⣿⣿⢌⠱⣇⣌⣌⣌⣌⡬⠷⣉⣾⣿⣿⣿⣿     REXPORT
⣿⣿⡿⠃⣼⣿⣿⣿⠁⣮⣮⣮⣮⠈⣳⣿⣿⣯⠈⣤⣿⣿⣿⣿⣿⣿⣯⣎⢌⢙⢙⣈⣌⣾⣿⣿⣿⣿⣿⣿
⣿⣿⠃⣼⣿⣿⣿⣿⠀⣿⣿⣿⣿⠏⣰⣿⣿⣿⣯⠐⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿     Re-encrypt and export snapshots
⣿⣿⠀⣿⣿⣿⣿⠓⢈⢙⢙⢙⢙⢉⠰⣿⣿⣿⣿⠎⣱⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿     to another AWS account using
⣿⣿⠀⣿⣿⣿⣿⠀⣿⣿⠳⠳⣷⣿⠀⣿⣿⣿⣿⠇⣸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿     a new customer managed KMS key
⣿⣿⠌⣳⣿⣿⣿⠀⣿⣯⠘⢁⣼⣿⠀⣿⣿⣿⡿⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
⣿⣿⣿⠌⡳⣿⣿⠀⠳⠳⠳⠳⠳⠳⠀⣿⣿⡿⠁⡲⣿⣿⣿⣿⣿⣿⡿⠷⠓⢙⢙⠹⠳⣷⣿⣿⣿⣿⣿⣿     Version %s
⣿⣿⣿⣿⣎⠹⡳⣿⣿⣿⣿⣿⣿⣿⣿⡷⢓⣈⣿⣯⢌⠹⡳⣿⣿⠓⣈⣾⠷⠑⢑⠳⣧⣎⠸⣷⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣮⣌⢙⠙⠳⠳⢓⢙⣈⣬⣿⣿⣿⣿⣿⣿⣮⢌⠁⣼⣿⣿⠀⣿⣿⠏⣰⣿⣯⠈⣷⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⣀⣿⣿⠓⢈⢙⢙⢉⠐⣿⣿⠏⣰⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⠀⣿⣿⠀⣻⠉⠀⣿⠀⣿⣿⠇⣸⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣏⠰⣷⢌⠱⠑⠳⠱⢀⣿⠗⣨⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣌⠙⠳⠳⠷⠳⢓⣈⣾⣿⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣮⣮⣿⣿⣿⣿⣿⣿⣿⣿⣿
`, version)
}

func sliceCsvInput(csvstring string) (rstrings []string, err error) {
	rstrings = strings.Split(csvstring, ",")
	if len(rstrings) < 1 {
		err = errors.New("error processing csv user string no values in slice")
	}
	return rstrings, err
}

// Config is an internal struct for storing
// configuration needed to run this application
// such as the snapshot id's to encrypt
type Config struct {
	Profile        string   `yaml:"aws_creds_profile"`
	Region         string   `yaml:"region"`
	ShareAccountId string   `yaml:"share_account_id"`
	KmsKeyId       string   `yaml:"kms_key_id"`
	Snapshots      []string `yaml:"snapshots"`
}

// ParseConfigFile takes a yaml filename as input and
// attempts to parse it into a config object.
func (rc *Config) ParseConfigFile(filename string) (err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, rc)
	return err
}

func addSessionTagToSpec(f *folio, sp etypes.TagSpecification) etypes.TagSpecification {
	sessionTag := etypes.Tag{
		Key:   f.SessionTagKey,
		Value: f.SessionTagValue,
	}
	sp.Tags = append(sp.Tags, sessionTag)
	return sp
}

func shareFolio(f *folio, svc *ec2.Client) {
	var cvp etypes.CreateVolumePermission
	cvp.UserId = f.ShareAccountId
	msai := ec2.ModifySnapshotAttributeInput{
		SnapshotId: f.NewSnapshotId,
		CreateVolumePermission: &etypes.CreateVolumePermissionModifications{
			Add: []etypes.CreateVolumePermission{cvp},
		},
		UserIds: []string{*f.ShareAccountId},
	}
	_, err := svc.ModifySnapshotAttribute(context.TODO(), &msai)
	if err != nil {
		log.Fatal(err)
	}
	Loggo.Info("successfully shared snapshot", "snapshotId", *f.NewSnapshotId, "accountNumber", *f.ShareAccountId)
}

func checkFolio(f *folio, svc *ec2.Client, share bool) {
	defer wgv.Done()
	var done bool
	for i := 0; i < 1; i++ {
		dsi := ec2.DescribeSnapshotsInput{
			SnapshotIds: []string{*f.NewSnapshotId},
		}
		dso, err := svc.DescribeSnapshots(context.TODO(), &dsi)
		if err != nil {
			log.Fatal(err)
		}
		for _, snap := range dso.Snapshots {
			state := snap.State
			Loggo.Info("got snapshot status",
				"newSnapshotState", state,
				"newVolumeId", *f.NewVolumeId,
				"sourceVolumeId", *f.SourceVolumeId,
				"newSnapshotId", *f.NewSnapshotId,
				"newSnapshotProgress", *snap.Progress)
			switch state {
			case "completed":
				done = true
				break
			default:
				done = false
			}
		}
		if done {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if done && share {
		shareFolio(f, svc)
	}
}

func createAndWatchVolume(f *folio, svc *ec2.Client) {
	defer wgv.Done()
	encr := true
	var tagSpecVolume etypes.TagSpecification
	tagSpecVolume.Tags = f.SourceTags
	tagSpecVolume.ResourceType = "volume"
	tagSpecVolume = addSessionTagToSpec(f, tagSpecVolume)
	cvi := ec2.CreateVolumeInput{
		AvailabilityZone:  f.NewVolumeAz,
		Encrypted:         &encr,
		KmsKeyId:          &conf.KmsKeyId,
		SnapshotId:        f.SourceSnapshotId,
		TagSpecifications: []etypes.TagSpecification{tagSpecVolume},
	}
	cvo, err := svc.CreateVolume(context.TODO(), &cvi)
	if err != nil {
		log.Fatal(err)
	}
	f.NewVolumeId = cvo.VolumeId
	Loggo.Info("created new volume",
		"id", *f.NewVolumeId,
		"sourceVolume", *f.SourceVolumeId,
		"sourceSnapshot", *f.SourceSnapshotId)
	for i := 0; i < 300; i++ {
		var done bool
		dvi := ec2.DescribeVolumesInput{
			VolumeIds: []string{*f.NewVolumeId},
		}
		dvo, err := svc.DescribeVolumes(context.TODO(), &dvi)
		if err != nil {
			log.Fatal(err)
		}
		for _, vol := range dvo.Volumes {
			state := vol.State
			Loggo.Info("got volume status", "state", state, "volumeId", *vol.VolumeId)
			switch state {
			case "available":
				done = true
				break
			default:
				done = false
			}
		}
		if done {
			break
		}
		time.Sleep(5 * time.Second)
	}
	// Now create snapshot from new volume
	var tagSpecSnapshot etypes.TagSpecification
	tagSpecSnapshot.Tags = f.SourceTags
	tagSpecSnapshot.ResourceType = "snapshot"
	tagSpecSnapshot = addSessionTagToSpec(f, tagSpecSnapshot)
	description := fmt.Sprintf("Created for export to account %s in %s:%s", *f.ShareAccountId, *f.SessionTagKey, *f.SessionTagValue)
	csi := ec2.CreateSnapshotInput{
		VolumeId:          f.NewVolumeId,
		TagSpecifications: []etypes.TagSpecification{tagSpecSnapshot},
		Description:       &description,
	}
	cso, err := svc.CreateSnapshot(context.TODO(), &csi)
	if err != nil {
		log.Fatal(err)
	}
	f.NewSnapshotId = cso.SnapshotId
	for i := 0; i < 2; i++ {
		var done bool
		dsi := ec2.DescribeSnapshotsInput{
			SnapshotIds: []string{*f.NewSnapshotId},
		}
		dso, err := svc.DescribeSnapshots(context.TODO(), &dsi)
		if err != nil {
			log.Fatal(err)
		}
		for _, snap := range dso.Snapshots {
			state := snap.State
			Loggo.Info("got snapshot status",
				"newSnapshotState", state,
				"newVolumeId", *f.NewVolumeId,
				"sourceVolumeId", *f.SourceVolumeId,
				"newSnapshotId", *f.NewSnapshotId,
				"newSnapshotProgress", *snap.Progress)
			switch state {
			case "completed":
				done = true
				break
			default:
				done = false
			}
		}
		if done {
			break
		}
		time.Sleep(10 * time.Second)
	}
	share := false
	// since checkFolio decrements waitGroup for resumption support
	// we'll need to artificially add a wait
	wgv.Add(1)
	go checkFolio(f, svc, share)
}

func exportFolios(folios []*folio, filename string) error {
	b, err := json.Marshal(folios)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

type folio struct {
	SourceTags         []etypes.Tag `json:"sourceTags"`
	SourceVolumeId     *string      `json:"sourceVolumeId"`
	SourceSnapshotId   *string      `json:"sourceSnapshotId"`
	SourceMapping      *string      `json:"sourceMapping"`
	SourceInstanceId   *string      `json:"sourceInstanceId"`
	SourceInstanceName *string      `json:"sourceInstanceName"`
	NewVolumeId        *string      `json:"newVolumeId"`
	NewVolumeAz        *string      `json:"newVolumeAz"`
	NewSnapshotId      *string      `json:"newSnapshotId"`
	SessionTagKey      *string      `json:"sessionTagKey"`
	SessionTagValue    *string      `json:"sessionTagValue"`
	ShareAccountId     *string      `json:"shareAccountId"`
}

func newRandomString() string {
	rand.Seed(time.Now().Unix())

	//Only lowercase
	charSet := "abcdedfghijklmnopqrst"
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func resumeSession(folioFile string, svc *ec2.Client) (err error) {
	jsonFile, err := os.Open(folioFile)
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var folios []*folio

	json.Unmarshal(byteValue, &folios)
	share := true
	for _, f := range folios {
		wgv.Add(1)
		go checkFolio(f, svc, share)
	}
	wgv.Wait()
	return err
}

func main() {
	var logLevel string
	var logFile string
	var folioFile string
	var noLogFile bool
	var resume bool
	var daemonFlag bool

	var configFile string

	flag.StringVar(&configFile, "config", "", "Filename of YAML configuration file.")
	flag.StringVar(&logFile, "logfile", "_rexport.log.json", "JSON logfile location")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level (info or debug)")
	flag.StringVar(&folioFile, "folios", "", "json folio file to import for monitoring/sharing")
	flag.BoolVar(&resume, "resume", false, "resume previous session, requires folioFile parameter for json file")
	flag.Parse()
	SetLogger(daemonFlag, noLogFile, logFile, logLevel)
	Loggo.Info("Starting rexport")
	setArt()
	color.HiBlue(ASCII_ART)
	err := conf.ParseConfigFile(configFile)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(conf.Profile),
	)
	if err != nil {
		log.Fatal(err)
	}
	svc := ec2.NewFromConfig(cfg)
	if resume {
		err = resumeSession(folioFile, svc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		si := ec2.DescribeSnapshotsInput{
			SnapshotIds: conf.Snapshots,
		}
		out, err := svc.DescribeSnapshots(context.TODO(), &si)
		if err != nil {
			log.Fatal(err)
		}
		tagKey := "rexport-session"
		tagValue := newRandomString()
		Loggo.Info("starting rexport session and all objects created will be tagged with", "tag-key", tagKey, "tag-value", tagValue)
		var folios []*folio
		na := "n/a"
		for _, snap := range out.Snapshots {
			var f folio
			f.SessionTagKey = &tagKey
			f.SessionTagValue = &tagValue
			f.SourceSnapshotId = snap.SnapshotId
			f.SourceVolumeId = snap.VolumeId
			f.SourceTags = snap.Tags
			f.ShareAccountId = &conf.ShareAccountId
			Loggo.Info("pulling info", "volumeId", *f.SourceVolumeId, "snapshotId", *f.SourceSnapshotId)
			dvi := ec2.DescribeVolumesInput{
				VolumeIds: []string{*f.SourceVolumeId},
			}
			dvo, err := svc.DescribeVolumes(context.TODO(), &dvi)
			if err != nil {
				log.Fatal(err)
			}
			attached := false
			for _, vol := range dvo.Volumes {
				f.NewVolumeAz = vol.AvailabilityZone
				if len(vol.Attachments) > 0 {
					attached = true
					for _, attachment := range vol.Attachments {
						f.SourceMapping = attachment.Device
						f.SourceInstanceId = attachment.InstanceId
					}
				} else {
					f.SourceMapping = &na
					f.SourceInstanceId = &na
				}
			}
			if attached {
				fmt.Println("here")
				dii := ec2.DescribeInstancesInput{
					InstanceIds: []string{*f.SourceInstanceId},
				}
				dio, err := svc.DescribeInstances(context.TODO(), &dii)
				if err != nil {
					log.Fatal(err)
				}
				for _, res := range dio.Reservations {
					for _, i := range res.Instances {
						for _, tag := range i.Tags {
							if *tag.Key == "Name" {
								f.SourceInstanceName = tag.Value
							}
						}
					}
				}
			} else {
				f.SourceInstanceName = &na
			}
			folios = append(folios, &f)
			Loggo.Info("gathered info on snapshot",
				"id", *f.SourceSnapshotId,
				"createdFrom", *f.SourceVolumeId,
				"sourceMapping", *f.SourceMapping,
				"sourceInstanceId", *f.SourceInstanceId,
				"sourceInstanceName", *f.SourceInstanceName)
		}

		for _, f := range folios {
			wgv.Add(1)
			go createAndWatchVolume(f, svc)
		}
		wgv.Wait()
		exportFilename := fmt.Sprintf("%s.json", tagValue)
		err = exportFolios(folios, exportFilename)
		if err != nil {
			log.Fatal(err)
		}
		msg := fmt.Sprintf("Please re-run with the '-resume' and '-foliofile %s' argument to resume session, check for snapshot completion, and final sharing", exportFilename)
		Loggo.Info("completed volume creation with re-encryption and re-snapshotting", "msg", msg)
	}

}
