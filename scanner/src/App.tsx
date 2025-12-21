import type { Component } from 'solid-js';
import { Route, Router } from '@solidjs/router';
import Home from './Home';
import Login from './Login';
import Scanner from './Scanner';
import AddReceipt from './AddReceipt';

const App: Component = () => {
  return (
    <Router>
      <Route path="/auth/login" component={Login} />
      <Route path="/receipt/scan" component={Scanner} />
      <Route path="/receipt/add" component={AddReceipt} />
      <Route path="/" component={Home} />
    </Router >
  );
};

export default App;
