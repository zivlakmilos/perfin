import { type Component } from 'solid-js';
import QrScanner from './components/QrScanner';
import { useNavigate } from '@solidjs/router';
import { loadToken, token } from './utils/token';

const Scanner: Component = () => {
  const navigate = useNavigate();

  loadToken();
  if (!token()) {
    navigate("/auth/login", { replace: true });
  }
  const onScanSuccess = (txt: string) => {
    navigate(`/receipt/add?value=${encodeURIComponent(txt)}`, {
      replace: true,
    });
  }

  return (
    <QrScanner
      qrCodeSuccessCallback={onScanSuccess}
      fps={60}
      qrbox={300}
      aspectRatio={1.777778}
      disableFlip={false}
    />
  );
};

export default Scanner;
