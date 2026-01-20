import { createSignal, type Component } from 'solid-js';
import { login } from './utils/api';
import { storeToken } from './utils/token';
import { useNavigate } from '@solidjs/router';

const Login: Component = () => {
  const navigate = useNavigate();

  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");

  const handleLogin = async (e: any) => {
    e.preventDefault();
    try {
      const token = await login(username(), password());
      storeToken(token);
      navigate("/", { replace: true });
      console.log(token);
    } catch (err) {
      console.log(err);
      alert(err);
    }
  }

  return (
    <div class="flex items-center justify-center h-screen">
      <form onSubmit={handleLogin}>
        <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
          <legend class="fieldset-legend">Login</legend>

          <label class="label">Username</label>
          <input type="text" class="input" placeholder="Email" value={username()} onChange={e => setUsername(e.target.value)} />

          <label class="label">Password</label>
          <input type="password" class="input" placeholder="Password" value={password()} onChange={e => setPassword(e.target.value)} />

          <button class="btn btn-neutral mt-4" type="submit">Login</button>
        </fieldset>
      </form>
    </div>
  );
};

export default Login;
