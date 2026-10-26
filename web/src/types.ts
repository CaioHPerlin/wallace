export interface Device {
  device_id: string;
  name: string;
  led_enabled: boolean;
  buzzer_enabled: boolean;
  online: boolean;
  last_heartbeat: number | null;
}

export interface DeviceConfig {
  led_enabled: boolean;
  buzzer_enabled: boolean;
}

export interface Detection {
  id: string;
  ts: number;
}

export interface DetectionEvent {
  type: 'detection';
  device_id: string;
  ts: number;
}
