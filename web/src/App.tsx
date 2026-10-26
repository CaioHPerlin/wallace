import { useCallback, useEffect, useState } from 'react';
import type { Detection, DetectionEvent, Device, DeviceConfig } from './types';
import {
  dashboardSocketURL,
  listDetections,
  listDevices,
  makeId,
  renameDevice,
  updateConfig,
} from './api';
import DeviceCard from './DeviceCard';

export default function App() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detections, setDetections] = useState<Record<string, Detection[]>>({});
  const [error, setError] = useState<string | null>(null);

  const refreshDevices = useCallback(async () => {
    try {
      setDevices(await listDevices());
      setError(null);
    } catch {
      setError('Unable to reach the API');
    }
  }, []);

  useEffect(() => {
    void refreshDevices();
    const interval = window.setInterval(refreshDevices, 5000);
    return () => window.clearInterval(interval);
  }, [refreshDevices]);

  useEffect(() => {
    const socket = new WebSocket(dashboardSocketURL());
    socket.onmessage = (message) => {
      const event = JSON.parse(message.data) as DetectionEvent;
      if (event.type !== 'detection') return;
      const detection: Detection = { id: makeId(), ts: event.ts };
      setDetections((prev) => ({
        ...prev,
        [event.device_id]: [detection, ...(prev[event.device_id] ?? [])],
      }));
    };
    return () => socket.close();
  }, []);

  const selectDevice = useCallback(
    async (id: string) => {
      setSelectedId((current) => (current === id ? null : id));
      if (!detections[id]) {
        try {
          const history = await listDetections(id);
          setDetections((prev) => ({ ...prev, [id]: history }));
        } catch {
          // keep the feed empty until the next successful load
        }
      }
    },
    [detections],
  );

  const handleRename = useCallback(async (id: string, name: string) => {
    const updated = await renameDevice(id, name);
    setDevices((prev) =>
      prev.map((device) => (device.device_id === id ? { ...device, name: updated.name } : device)),
    );
  }, []);

  const handleUpdateConfig = useCallback(
    (id: string, config: DeviceConfig) => {
      setDevices((prev) =>
        prev.map((device) => (device.device_id === id ? { ...device, ...config } : device)),
      );
      void updateConfig(id, config).catch(() => refreshDevices());
    },
    [refreshDevices],
  );

  return (
    <main className="page-shell">
      <header className="topbar">
        <a className="brand" href="/" aria-label="WallAce home">
          <span className="brand-mark" aria-hidden="true">
            W
          </span>
          <span>WallAce</span>
        </a>
        <span className="topbar-caption">MOTION SECURITY</span>
      </header>

      <section className="devices" aria-labelledby="devices-title">
        <div className="section-heading">
          <p className="eyebrow">YOUR SENSORS</p>
          <h1 id="devices-title">Devices</h1>
          <p className="intro">Motion history and actuator controls for every connected sensor.</p>
        </div>

        {error && <p className="error-banner">{error}</p>}

        {devices.length === 0 ? (
          <div className="setup-card">
            <div className="setup-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M12 3v3m0 12v3M3 12h3m12 0h3M5.64 5.64l2.12 2.12m8.48 8.48 2.12 2.12m0-12.72-2.12 2.12m-8.48 8.48-2.12 2.12" />
                <circle cx="12" cy="12" r="4" />
              </svg>
            </div>
            <div>
              <p className="eyebrow">GETTING READY</p>
              <h2>No sensors connected yet.</h2>
              <p className="setup-description">
                Devices appear here after their first heartbeat. Flash the ESP32 firmware to get
                started.
              </p>
            </div>
          </div>
        ) : (
          <div className="device-list">
            {devices.map((device) => (
              <DeviceCard
                key={device.device_id}
                device={device}
                detections={detections[device.device_id] ?? []}
                selected={selectedId === device.device_id}
                onSelect={(id) => void selectDevice(id)}
                onRename={handleRename}
                onUpdateConfig={handleUpdateConfig}
              />
            ))}
          </div>
        )}
      </section>

      <footer className="footer">
        <span>WallAce Motion Security</span>
        <span>Built for a little more peace of mind.</span>
      </footer>
    </main>
  );
}
