import type { Component } from 'solid-js';
import { Route, Router } from '@solidjs/router';
import Home from './Home';
import Login from './Login';

const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/auth/login" component={Login} />
    </Router >
  );
};

export default App;
