import type { Component } from 'solid-js';

const Home: Component = () => {
  return (
    <div class="flex items-center justify-center h-screen">
      <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
        <button class="btn btn-neutral">Scan QR</button>
        <button class="btn btn-neutral mt-4">Maunal Input</button>
      </fieldset>
    </div>
  );
};

export default Home;
