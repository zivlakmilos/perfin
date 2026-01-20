import { createSignal } from "solid-js";

export const [token, setToken] = createSignal<string | null>(null);

export const storeToken = (token: string): void => {
  localStorage.setItem("token", token);
}

export const loadToken = () => {
  setToken(localStorage.getItem("token"));
}
