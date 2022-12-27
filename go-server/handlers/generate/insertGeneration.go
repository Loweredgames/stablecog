package generate

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/yekta/stablecog/go-server/loggers"
	"github.com/yekta/stablecog/go-server/shared"
)

var red = color.New(color.FgHiRed).SprintFunc()

func InsertGenerationInitial(g SInsertGenerationProps) string {
	start := time.Now().UTC().UnixMilli()
	insertBody := SDBGenerationInsertInitial{
		Width:             g.Width,
		Height:            g.Height,
		GuidanceScale:     g.GuidanceScale,
		NumInferenceSteps: g.NumInferenceSteps,
		Seed:              g.Seed,
		ModelId:           g.ModelId,
		SchedulerId:       g.SchedulerId,
		UserId:            g.UserId,
		UserTier:          g.UserTier,
		ServerUrl:         g.ServerUrl,
		Status:            "started",
		DeviceType:        g.DeviceType,
		DeviceBrowser:     g.DeviceBrowser,
		DeviceOS:          g.DeviceOS,
		UserAgent:         g.UserAgent,
		CountryCode:       g.CountryCode,
	}
	var res SDBGenerationInsertInitialRes
	_, err := shared.SupabaseDb.From("generation").Insert(insertBody, false, "", "", "").Single().ExecuteTo(&res)
	if err != nil {
		log.Printf(red("-- DB - Error inserting generation: %v --"), err)
		return ""
	}
	end := time.Now().UTC().UnixMilli()
	log.Printf("-- DB - Generation initial insert in: %s --", green(end-start, "ms"))
	g.GenerationIdChan <- res.Id
	return res.Id
}

type SInsertGenerationProps struct {
	Width             int
	Height            int
	GuidanceScale     int
	NumInferenceSteps int
	Seed              int
	ModelId           string
	SchedulerId       string
	Status            string
	UserId            string
	UserTier          string
	ServerUrl         string
	UserAgent         string
	DeviceType        string
	DeviceBrowser     string
	DeviceOS          string
	CountryCode       string
	LogObject         loggers.SGenerationLogObject
	GenerationIdChan  chan string
}

type SDBGenerationInsertInitial struct {
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	GuidanceScale     int    `json:"guidance_scale"`
	NumInferenceSteps int    `json:"num_inference_steps"`
	Seed              int    `json:"seed"`
	ModelId           string `json:"model_id"`
	SchedulerId       string `json:"scheduler_id"`
	UserId            string `json:"user_id,omitempty"`
	UserTier          string `json:"user_tier"`
	Status            string `json:"status"`
	ServerUrl         string `json:"server_url"`
	UserAgent         string `json:"user_agent,omitempty"`
	CountryCode       string `json:"country_code,omitempty"`
	DeviceType        string `json:"device_type,omitempty"`
	DeviceBrowser     string `json:"device_browser,omitempty"`
	DeviceOS          string `json:"device_os,omitempty"`
}

type SDBGenerationInsertInitialRes struct {
	Id                string `json:"id"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	GuidanceScale     int    `json:"guidance_scale"`
	NumInferenceSteps int    `json:"num_inference_steps"`
	Seed              int    `json:"seed"`
	ModelId           string `json:"model_id"`
	SchedulerId       string `json:"scheduler_id"`
	UserId            string `json:"user_id,omitempty"`
	UserTier          string `json:"user_tier"`
	Status            string `json:"status"`
	ServerUrl         string `json:"server_url"`
	UserAgent         string `json:"user_agent,omitempty"`
	CountryCode       string `json:"country_code,omitempty"`
	DeviceType        string `json:"device_type,omitempty"`
	DeviceBrowser     string `json:"device_browser,omitempty"`
	DeviceOS          string `json:"device_os,omitempty"`
}
