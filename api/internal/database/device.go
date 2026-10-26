package database

import (
	"context"
	"time"
)

const onlineWindow = 60 * time.Second

type Device struct {
	ID            string `json:"device_id"`
	Name          string `json:"name"`
	LEDEnabled    bool   `json:"led_enabled"`
	BuzzerEnabled bool   `json:"buzzer_enabled"`
	Online        bool   `json:"online"`
	LastHeartbeat *int64 `json:"last_heartbeat"`
}

type DeviceConfig struct {
	LEDEnabled    bool `json:"led_enabled"`
	BuzzerEnabled bool `json:"buzzer_enabled"`
}

type DeviceStatus struct {
	DeviceConfig
	Online        bool   `json:"online"`
	LastHeartbeat *int64 `json:"last_heartbeat"`
}

type Detection struct {
	ID string `json:"id"`
	TS int64  `json:"ts"`
}

func (p *Postgres) ListDevices(ctx context.Context) ([]Device, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, name, led_enabled, buzzer_enabled, last_heartbeat
		FROM devices
		ORDER BY created_at, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.LEDEnabled, &d.BuzzerEnabled, &d.LastHeartbeat); err != nil {
			return nil, err
		}
		d.Online = isOnline(d.LastHeartbeat)
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (p *Postgres) RenameDevice(ctx context.Context, id, name string) (Device, error) {
	if err := p.ensureDevice(ctx, id); err != nil {
		return Device{}, err
	}

	var d Device
	err := p.pool.QueryRow(ctx, `
		UPDATE devices SET name = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, name, last_heartbeat`, id, name).Scan(&d.ID, &d.Name, &d.LastHeartbeat)
	if err != nil {
		return Device{}, err
	}
	d.Online = isOnline(d.LastHeartbeat)
	return d, nil
}

func (p *Postgres) GetConfig(ctx context.Context, id string) (DeviceConfig, error) {
	if err := p.ensureDevice(ctx, id); err != nil {
		return DeviceConfig{}, err
	}

	var c DeviceConfig
	err := p.pool.QueryRow(ctx, `SELECT led_enabled, buzzer_enabled FROM devices WHERE id = $1`, id).
		Scan(&c.LEDEnabled, &c.BuzzerEnabled)
	if err != nil {
		return DeviceConfig{}, err
	}
	return c, nil
}

func (p *Postgres) UpdateConfig(ctx context.Context, id string, cfg DeviceConfig) (DeviceConfig, error) {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO devices (id, led_enabled, buzzer_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET led_enabled = EXCLUDED.led_enabled,
		    buzzer_enabled = EXCLUDED.buzzer_enabled,
		    updated_at = now()`,
		id, cfg.LEDEnabled, cfg.BuzzerEnabled)
	if err != nil {
		return DeviceConfig{}, err
	}
	return cfg, nil
}

func (p *Postgres) RecordHeartbeat(ctx context.Context, id string) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO devices (id, last_heartbeat)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET last_heartbeat = EXCLUDED.last_heartbeat, updated_at = now()`,
		id, time.Now().Unix())
	return err
}

func (p *Postgres) GetStatus(ctx context.Context, id string) (DeviceStatus, error) {
	if err := p.ensureDevice(ctx, id); err != nil {
		return DeviceStatus{}, err
	}

	var s DeviceStatus
	err := p.pool.QueryRow(ctx, `SELECT led_enabled, buzzer_enabled, last_heartbeat FROM devices WHERE id = $1`, id).
		Scan(&s.LEDEnabled, &s.BuzzerEnabled, &s.LastHeartbeat)
	if err != nil {
		return DeviceStatus{}, err
	}
	s.Online = isOnline(s.LastHeartbeat)
	return s, nil
}

func (p *Postgres) ListDetections(ctx context.Context, id string, limit, offset int) ([]Detection, error) {
	if err := p.ensureDevice(ctx, id); err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, ts
		FROM detections
		WHERE device_id = $1
		ORDER BY ts DESC, id
		LIMIT $2 OFFSET $3`, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	detections := []Detection{}
	for rows.Next() {
		var d Detection
		if err := rows.Scan(&d.ID, &d.TS); err != nil {
			return nil, err
		}
		detections = append(detections, d)
	}
	return detections, rows.Err()
}

func (p *Postgres) CreateDetection(ctx context.Context, id string, ts int64) (Detection, error) {
	if err := p.ensureDevice(ctx, id); err != nil {
		return Detection{}, err
	}

	var d Detection
	err := p.pool.QueryRow(ctx, `
		INSERT INTO detections (device_id, ts)
		VALUES ($1, $2)
		RETURNING id::text, ts`, id, ts).Scan(&d.ID, &d.TS)
	if err != nil {
		return Detection{}, err
	}
	return d, nil
}

func (p *Postgres) ensureDevice(ctx context.Context, id string) error {
	_, err := p.pool.Exec(ctx, `INSERT INTO devices (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`, id)
	return err
}

func isOnline(last *int64) bool {
	if last == nil {
		return false
	}
	return time.Now().Unix()-*last <= int64(onlineWindow/time.Second)
}
