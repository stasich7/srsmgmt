package srsconfig

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"text/template"
)

type SRSstream struct {
	ID       string
	Password string
}

type SRSConfig struct {
	Streams    []SRSstream
	ConfigPath string
	SRSClient  ConfigReloader
	tplStorage embed.FS
}

type ConfigReloader interface {
	ConfigReload() error
}

func New(srsConfigPath string, tplStorage embed.FS, srsClient ConfigReloader) *SRSConfig {
	return &SRSConfig{
		Streams:    make([]SRSstream, 0),
		ConfigPath: srsConfigPath,
		SRSClient:  srsClient,
		tplStorage: tplStorage,
	}
}

func (s *SRSConfig) Init(streams []SRSstream) {
	for i := 0; i < len(streams); i++ {
		s.Streams = append(s.Streams, SRSstream{
			ID:       streams[i].ID,
			Password: streams[i].Password,
		})
	}
	err := s.configReload()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (s *SRSConfig) AddRTC(id, pass string) {
	for i := 0; i < len(s.Streams); i++ {
		if s.Streams[i].ID == id {
			return
		}
	}
	s.Streams = append(s.Streams, SRSstream{
		ID:       id,
		Password: pass,
	})
	err := s.configReload()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

}

func (s *SRSConfig) RemoveRTC(id string) {
	for i := 0; i < len(s.Streams); i++ {
		if s.Streams[i].ID == id {
			s.Streams = append(s.Streams[:i], s.Streams[i+1:]...)
			err := s.configReload()
			if err != nil {
				fmt.Printf("error: %v", err)
				return
			}
		}
	}
}

func (s *SRSConfig) configReload() error {
	var newConfig []byte
	tData, err := s.tplStorage.ReadFile("srs_custom.tpl")
	if err != nil {
		return err
	}
	t, err := template.New("custom").Parse(string(tData))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(newConfig)

	for _, v := range s.Streams {
		item := struct {
			RTCStreamUUID     string
			RTCStreamPassword string
		}{
			RTCStreamUUID:     v.ID,
			RTCStreamPassword: v.Password,
		}

		if err := t.Execute(buf, item); err != nil {
			return err
		}
	}

	tBaseData, err := s.tplStorage.ReadFile("srs_base.tpl")
	if err != nil {
		return err
	}
	tBase, err := template.New("base").Parse(string(tBaseData))
	if err != nil {
		return err
	}

	if err := tBase.Execute(buf, nil); err != nil {
		return err
	}
	s.SRSClient.ConfigReload()

	// fmt.Printf("newConfig: %s", buf.String())
	err = ioutil.WriteFile(s.ConfigPath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
