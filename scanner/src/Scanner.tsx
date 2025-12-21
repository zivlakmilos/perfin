import { createSignal, type Component } from 'solid-js';
import QrScanner from './components/QrScanner';
import { useNavigate } from '@solidjs/router';

const Scanner: Component = () => {
  const navigate = useNavigate();

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
