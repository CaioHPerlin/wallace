import { useEffect, useState } from 'react';

type ApiStatus = 'checking' | 'online' | 'offline';

export default function App() {
  const [apiStatus, setApiStatus] = useState<ApiStatus>('checking');

  useEffect(() => {
    let active = true;

    async function refreshStatus() {
      try {
        const response = await fetch('/api/health/ready');
        if (active) setApiStatus(response.ok ? 'online' : 'offline');
      } catch {
        if (active) setApiStatus('offline');
      }
    }

    void refreshStatus();
    const interval = window.setInterval(refreshStatus, 5000);

    return () => {
      active = false;
      window.clearInterval(interval);
    };
  }, []);

  const statusLabel = {
    checking: 'Checking connection',
    online: 'System online',
    offline: 'System unavailable',
  }[apiStatus];

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

      <section className="hero" aria-labelledby="page-title">
        <div className="hero-copy">
          <p className="eyebrow">YOUR SPACE, IN SIGHT</p>
          <h1 id="page-title">
            A quieter kind
            <br />
            of peace of mind.
          </h1>
          <p className="intro">
            WallAce connects your motion sensor to a clear, simple view of your space.
          </p>
          <div className={`system-status status-${apiStatus}`} role="status" aria-live="polite">
            <span className="status-dot" />
            <span>{statusLabel}</span>
          </div>
        </div>

        <div className="radar" aria-hidden="true">
          <div className="radar-ring ring-outer" />
          <div className="radar-ring ring-middle" />
          <div className="radar-ring ring-inner" />
          <div className="radar-sweep" />
          <span className="radar-center" />
          <span className="radar-point point-one" />
          <span className="radar-point point-two" />
          <span className="radar-point point-three" />
          <span className="radar-label label-top">LIVE</span>
          <span className="radar-label label-side">PIR SENSOR</span>
        </div>
      </section>

      <section className="setup-card" aria-labelledby="setup-title">
        <div className="setup-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none">
            <path d="M12 3v3m0 12v3M3 12h3m12 0h3M5.64 5.64l2.12 2.12m8.48 8.48 2.12 2.12m0-12.72-2.12 2.12m-8.48 8.48-2.12 2.12" />
            <circle cx="12" cy="12" r="4" />
          </svg>
        </div>
        <div>
          <p className="eyebrow">GETTING READY</p>
          <h2 id="setup-title">Your dashboard is taking shape.</h2>
          <p className="setup-description">
            Motion history and device controls will appear here when your sensor is connected.
          </p>
        </div>
        <span className="setup-tag">SETUP IN PROGRESS</span>
      </section>

      <footer className="footer">
        <span>WallAce Motion Security</span>
        <span>Built for a little more peace of mind.</span>
      </footer>
    </main>
  );
}
