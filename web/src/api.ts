import type { Detection, Device, DeviceConfig } from './types';

const token = import.meta.env.VITE_API_TOKEN ?? '';

const headers: Record<string, string> = {
  'Content-Type': 'application/json',
};
if (token) {
  headers.Authorization = `Bearer ${token}`;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, { ...init, headers });
  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export function listDevices(): Promise<Device[]> {
  return request<Device[]>('/api/devices');
}

export function renameDevice(deviceId: string, name: string): Promise<Device> {
  return request<Device>(`/api/devices/${deviceId}`, {
    method: 'PATCH',
    body: JSON.stringify({ name }),
  });
}

export function updateConfig(deviceId: string, config: DeviceConfig): Promise<DeviceConfig> {
  return request<DeviceConfig>(`/api/devices/${deviceId}/config`, {
    method: 'PUT',
    body: JSON.stringify(config),
  });
}

export function listDetections(deviceId: string): Promise<Detection[]> {
  return request<Detection[]>(`/api/devices/${deviceId}/detections`);
}

export function dashboardSocketURL(): string {
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  return `${protocol}://${window.location.host}/ws/dashboard`;
}

export function makeId(): string {
  return typeof crypto !== 'undefined' && 'randomUUID' in crypto
    ? crypto.randomUUID()
    : `${Date.now()}-${Math.random().toString(36).slice(2)}`;
}
