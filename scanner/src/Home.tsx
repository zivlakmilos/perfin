import { useNavigate } from '@solidjs/router';
import type { Component } from 'solid-js';
import { loadToken, token } from './utils/token';

const Home: Component = () => {
  const navigate = useNavigate();

  loadToken();
  if (!token()) {
    navigate("/auth/login", { replace: true });
  }

  return (
    <div class="flex items-center justify-center h-screen">
      <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
        <button class="btn btn-neutral" onClick={() => navigate("/receipt/scan")}>Scan QR</button>
        <button class="btn btn-neutral mt-4">Maunal Input</button>
      </fieldset>
    </div>
  );
};

export default Home;
