import React from 'react';
// assets
import { IconUser } from '@tabler/icons';

// constant
const icons = {
  IconUser
};

// ==============================|| EXTRA PAGES MENU ITEMS ||============================== //

const pages = {
  id: 'pages',
  title: 'Pages',
  caption: 'Pages Caption',
  type: 'group',
  children: [
    {
      id: 'authentication',
      title: 'Authentication',
      type: 'collapse',
      icon: icons.IconUser,

      children: [
        {
          id: 'login',
          title: 'Login',
          type: 'item',
          url: '/Login',
          target: true
        },
        {
          id: 'forgotpassword',
          title: 'ForgotPassword',
          type: 'item',
          url: '/ForgotPassword',
          target: true
        }
      ]
    }
  ]
};

export default pages;