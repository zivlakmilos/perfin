import { Html5QrcodeScanner, QrcodeErrorCallback, QrcodeSuccessCallback } from 'html5-qrcode';
import { Html5QrcodeScannerConfig } from 'html5-qrcode/esm/html5-qrcode-scanner';
import { onCleanup, onMount, type Component } from 'solid-js';

const qrcodeRegionId = "html5qr-code-full-region";

export interface TQtScannerProps {
  fps?: number | undefined;
  qrbox?: number | undefined;
  aspectRatio?: number | undefined;
  disableFlip?: boolean | undefined;
  qrCodeSuccessCallback?: QrcodeSuccessCallback | undefined;
  qrCodeErrorCallback?: QrcodeErrorCallback | undefined;
  verbose?: boolean | undefined;
}

const createConfig = (props: TQtScannerProps) => {
  const config: Html5QrcodeScannerConfig = {
    fps: 0,
    showTorchButtonIfSupported: true,
  };
  if (props.fps) {
    config.fps = props.fps;
  }
  if (props.qrbox) {
    config.qrbox = props.qrbox;
  }
  if (props.aspectRatio) {
    config.aspectRatio = props.aspectRatio;
  }
  if (props.disableFlip !== undefined) {
    config.disableFlip = props.disableFlip;
  }
  return config;
};

const QrScanner: Component = (props: TQtScannerProps) => {
  let html5QrcodeScanner: Html5QrcodeScanner | undefined = undefined;

  onMount(() => {
    const config = createConfig(props);
    const verbose = props.verbose === true;
    if (!(props.qrCodeSuccessCallback)) {
      throw "qrCodeSuccessCallback is required callback.";
    }
    html5QrcodeScanner = new Html5QrcodeScanner(qrcodeRegionId, config, verbose);
    html5QrcodeScanner.render(props.qrCodeSuccessCallback, props.qrCodeErrorCallback);
  });

  onCleanup(() => {
    if (!html5QrcodeScanner) {
      return;
    }

    html5QrcodeScanner.clear().catch(error => {
      console.error("Failed to clear html5QrcodeScanner. ", error);
    });
  });

  return (
    <div id={qrcodeRegionId} />
  );
};

export default QrScanner;
