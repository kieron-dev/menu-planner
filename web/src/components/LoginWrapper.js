import React from 'react';
import {useAuth} from './AuthProvider';
import Home from './Home';
import Welcome from './Welcome';

const LoginWrapper = () => {
  const {auth} = useAuth();

  return auth.isAuthed ? <Home /> : <Welcome />;
}

export default LoginWrapper;
