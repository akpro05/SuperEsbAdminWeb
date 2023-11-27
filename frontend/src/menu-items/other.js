import React from 'react';
// assets
import { IconReportSearch, IconLogout } from '@tabler/icons';

// constant
const icons = { IconReportSearch, IconLogout };

// ==============================|| SAMPLE PAGE & DOCUMENTATION MENU ITEMS ||============================== //

const other = {
  id: 'sample-docs-roadmap',
  type: 'group',
  children: [
    {
      id: 'logout',
      title: 'Logout',
      type: 'item',
      url: '/Login',
      icon: icons.IconLogout
      
    }
  ]
};

export default other;
