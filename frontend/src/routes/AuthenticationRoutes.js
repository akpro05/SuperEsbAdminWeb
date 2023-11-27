import React from 'react';
import { lazy } from 'react';

// project imports
import Loadable from '../ui-component/Loadable.js';
import MinimalLayout from '../layout/MinimalLayout';

// login option 3 routing
const AuthLogin3 = Loadable(lazy(() => import('../views/pages/authentication/authentication3/Login3')));
const AuthForgotPassword = Loadable(lazy(() => import('../views/pages/authentication/authentication3/ForgotPassword')));

// ==============================|| AUTHENTICATION ROUTING ||============================== //

const AuthenticationRoutes = {
  path: '/',
  element: <MinimalLayout />,
  children: [
    {
      path: '',
      element: <AuthLogin3 />
    },
    {
      path: '/Login',
      element: <AuthLogin3 />
    },
    {
      path: '/ForgotPassword',
      element: <AuthForgotPassword />
    }
  ]
};

export default AuthenticationRoutes;
