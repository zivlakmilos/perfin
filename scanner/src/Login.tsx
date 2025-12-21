import type { Component } from 'solid-js';

const Login: Component = () => {
  return (
    <div class="flex items-center justify-center h-screen">
      <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
        <legend class="fieldset-legend">Login</legend>

        <label class="label">Email</label>
        <input type="email" class="input" placeholder="Email" />

        <label class="label">Password</label>
        <input type="password" class="input" placeholder="Password" />

        <button class="btn btn-neutral mt-4">Login</button>
      </fieldset>
    </div>
  );
};

export default Login;
