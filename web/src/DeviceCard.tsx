import { useState } from 'react';
import type { Detection, Device, DeviceConfig } from './types';

interface DeviceCardProps {
  device: Device;
  detections: Detection[];
  selected: boolean;
  onSelect: (deviceId: string) => void;
  onRename: (deviceId: string, name: string) => Promise<void>;
  onUpdateConfig: (deviceId: string, config: DeviceConfig) => void;
}

function formatTimestamp(ts: number | null): string {
  if (ts === null) return 'never';
  return new Date(ts * 1000).toLocaleTimeString();
}

export default function DeviceCard({
  device,
  detections,
  selected,
  onSelect,
  onRename,
  onUpdateConfig,
}: DeviceCardProps) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(device.name);

  const displayName = device.name || device.device_id;

  const submitRename = async () => {
    setEditing(false);
    const next = draft.trim();
    if (next === displayName || next === '') return;
    await onRename(device.device_id, next);
  };

  const toggle = (key: 'led_enabled' | 'buzzer_enabled') => {
    onUpdateConfig(device.device_id, {
      led_enabled: device.led_enabled,
      buzzer_enabled: device.buzzer_enabled,
      [key]: !device[key],
    });
  };

  return (
    <article className={`device-card${selected ? ' is-selected' : ''}`}>
      <div
        className="device-main"
        role="button"
        tabIndex={0}
        onClick={() => onSelect(device.device_id)}
        onKeyDown={(event) => {
          if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            onSelect(device.device_id);
          }
        }}
      >
        <span className={`status-dot ${device.online ? 'status-online' : 'status-offline'}`} />
        <span className="device-identity">
          {editing ? (
            <input
              className="device-name-input"
              autoFocus
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              onClick={(event) => event.stopPropagation()}
              onBlur={() => void submitRename()}
              onKeyDown={(event) => {
                if (event.key === 'Enter') {
                  event.preventDefault();
                  void submitRename();
                } else if (event.key === 'Escape') {
                  setEditing(false);
                }
              }}
            />
          ) : (
            <>
              <button
                type="button"
                className="device-name"
                onClick={(event) => {
                  event.stopPropagation();
                  setDraft(displayName);
                  setEditing(true);
                }}
              >
                {displayName}
              </button>
              <span className="device-id">{device.device_id}</span>
            </>
          )}
        </span>
        <span className={`device-state ${device.online ? 'is-online' : 'is-offline'}`}>
          {device.online ? 'online' : 'offline'}
        </span>
      </div>

      {selected && (
        <div className="device-details">
          <div className="device-controls">
            <label className="toggle">
              <input
                type="checkbox"
                checked={device.led_enabled}
                onChange={() => toggle('led_enabled')}
              />
              <span>LED</span>
            </label>
            <label className="toggle">
              <input
                type="checkbox"
                checked={device.buzzer_enabled}
                onChange={() => toggle('buzzer_enabled')}
              />
              <span>Buzzer</span>
            </label>
            <span className="device-heartbeat">
              Last heartbeat: {formatTimestamp(device.last_heartbeat)}
            </span>
          </div>

          <div className="detection-list">
            <h3>Recent detections</h3>
            {detections.length === 0 ? (
              <p className="empty">No detections yet.</p>
            ) : (
              <ul>
                {detections.map((detection) => (
                  <li key={detection.id}>{formatTimestamp(detection.ts)}</li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
    </article>
  );
}
