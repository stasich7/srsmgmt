package srsmgmtrepo

import (
	"database/sql"
	"srsmgmt/internal/srsmgmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stream struct {
	StreamID  uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	App       string
	Password  string
	Status    int
	ClientId  string
	StartedAt *sql.NullTime
	StopedAt  *sql.NullTime
	RTC       bool
}

func (repo Repo) CreateStream(s srsmgmt.Stream) (*srsmgmt.Stream, error) {
	stream := Stream{
		StreamID: s.StreamID,
		App:      s.App,
		Password: s.Password,
		Status:   srsmgmt.StreamStatusWaitPublish,
		RTC:      s.RTC,
	}
	result := repo.Db.Create(&stream)
	if result == nil || result.Error != nil {
		return &srsmgmt.Stream{}, result.Error
	}

	repo.Db.First(&stream, stream.StreamID)

	resp := srsmgmt.Stream{
		StreamID:  stream.StreamID,
		App:       stream.App,
		Password:  stream.Password,
		Status:    stream.Status,
		CreatedAt: stream.CreatedAt,
		UpdatedAt: stream.UpdatedAt,
		ClientId:  stream.ClientId,
		RTC:       stream.RTC,
	}

	return &resp, nil
}

func (repo Repo) GetStream(streamID uuid.UUID) (*srsmgmt.Stream, error) {
	stream := Stream{}
	result := repo.Db.First(&stream, streamID)

	if result.Error != nil {
		return &srsmgmt.Stream{}, srsmgmt.ErrNotFound
	}

	var startedT *time.Time
	if stream.StartedAt != nil && stream.StartedAt.Valid {
		startedT = &stream.StartedAt.Time
	}

	var stopedT *time.Time
	if stream.StopedAt != nil && stream.StopedAt.Valid {
		stopedT = &stream.StopedAt.Time
	}

	resp := srsmgmt.Stream{
		StreamID:  stream.StreamID,
		App:       stream.App,
		Password:  stream.Password,
		Status:    stream.Status,
		CreatedAt: stream.CreatedAt,
		UpdatedAt: stream.UpdatedAt,
		ClientId:  stream.ClientId,
		StartedAt: startedT,
		StopedAt:  stopedT,
		RTC:       stream.RTC,
	}

	return &resp, nil
}

func (repo Repo) GetRTCStreams() (*[]srsmgmt.Stream, error) {
	streams := []Stream{}
	result := repo.Db.Model(&Stream{}).Where("rtc = true").Scan(&streams)

	if result.Error != nil {
		return &([]srsmgmt.Stream{}), srsmgmt.ErrNotFound
	}

	resp := []srsmgmt.Stream{}
	for _, v := range streams {
		var startedT *time.Time
		if v.StartedAt != nil && v.StartedAt.Valid {
			startedT = &v.StartedAt.Time
		}

		var stopedT *time.Time
		if v.StopedAt != nil && v.StopedAt.Valid {
			stopedT = &v.StopedAt.Time
		}

		resp = append(resp, srsmgmt.Stream{
			StreamID:  v.StreamID,
			App:       v.App,
			Password:  v.Password,
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			ClientId:  v.ClientId,
			StartedAt: startedT,
			StopedAt:  stopedT,
			RTC:       v.RTC,
		})
	}

	return &resp, nil
}

func (repo Repo) DeleteStream(streamID uuid.UUID) (uuid.UUID, error) {
	result := repo.Db.Unscoped().Delete(&Stream{}, streamID)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return streamID, nil
}

func (repo Repo) UpdateStream(s srsmgmt.Stream) (*srsmgmt.Stream, error) {
	stream := map[string]interface{}{
		"StreamID":  s.StreamID,
		"App":       s.App,
		"Password":  s.Password,
		"Status":    s.Status,
		"CreatedAt": s.CreatedAt,
		"UpdatedAt": s.UpdatedAt,
		"ClientId":  s.ClientId,
		"StartedAt": s.StartedAt,
		"StopedAt":  s.StopedAt,
		"RTC":       s.RTC,
	}

	result := repo.Db.Model(&Stream{}).Clauses(clause.Returning{}).Where("stream_id = ?", s.StreamID).Updates(&stream)
	if result.Error != nil {
		return &srsmgmt.Stream{}, result.Error
	}

	updated := Stream{}
	result.Scan(&updated)

	var startedT *time.Time
	if updated.StartedAt != nil && updated.StartedAt.Valid {
		startedT = &updated.StartedAt.Time
	}

	var stopedT *time.Time
	if updated.StopedAt != nil && updated.StopedAt.Valid {
		stopedT = &updated.StopedAt.Time
	}

	resp := srsmgmt.Stream{
		StreamID:  updated.StreamID,
		App:       updated.App,
		Password:  updated.Password,
		Status:    updated.Status,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
		ClientId:  updated.ClientId,
		StartedAt: startedT,
		StopedAt:  stopedT,
		RTC:       updated.RTC,
	}

	return &resp, nil
}
