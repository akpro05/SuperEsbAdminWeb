import React from 'react';
// assets
import { IconUsers, IconSettings, IconTool, IconServerCog, IconFiles, IconWindmill, IconUserCheck, IconReport, IconServerBolt } from '@tabler/icons';

import * as TablerIcons from '@tabler/icons';
for (const iconName in TablerIcons) {
  // console.log(`Icon Name: ${iconName}`);
}


// constant
const icons = {
  IconUsers,
  IconSettings,
  IconTool,
  IconServerCog,
  IconFiles,
  IconWindmill,
  IconUserCheck,
  IconServerBolt,
  IconReport
};

// ==============================|| UTILITIES MENU ITEMS ||============================== //

const utilities = {
  id: 'utilities',
  title: 'Menu',
  type: 'group',
  children: [
    {
      id: 'sysusermgmt',
      title: 'System User Management',
      type: 'collapse',
      icon: icons.IconUsers,
      children: [
        {
          id: 'searchuser',
          title: 'Search System User',
          type: 'item',
          url: '/SysUser/SearchsysUser',
          breadcrumbs: false
        }
      ]
    },
    {
      id: 'prodmgmt',
      title: 'Producer Management',
      type: 'collapse',
      icon: icons.IconUserCheck,
      children: [
        {
          id: 'searchproducer',
          title: 'Search Producer',
          type: 'item',
          url: '/Producers/SearchProducers',
          breadcrumbs: false
        }
      ]
    },
    {
      id: 'consumers',
      title: 'Consumer Management',
      type: 'collapse',
      icon: icons.IconServerBolt,
      children: [
        {
          id: 'searchConsumers',
          title: 'Search Consumer',
          type: 'item',
          url: '/Consumers/SearchConsumers',
          breadcrumbs: false
        }
      ]
    },
    {
      id: 'systconfig',
      title: 'System Configuration',
      type: 'collapse',
      icon: icons.IconSettings,
      children: [
        {
          id: 'searchproducertoconsumer',
          title: 'Search Producer To Consumer',
          type: 'item',
          url: '/ProducerToConsumer/SearchProducerToConsumer',
          breadcrumbs: false
        },
        {
          id: 'setrole',
          title: 'Set Role',
          type: 'item',
          url: '/Role/SearchRole',
          breadcrumbs: false
        }
      ]
    },
    {
      id: 'reports',
      title: 'Reports',
      type: 'collapse',
      icon: icons.IconReport,
      children: [
        
        {
          id: 'esblogsreport',
          title: 'ESB Logs Report',
          type: 'item',
          url: '/Reports/ESBLogsReport',
          breadcrumbs: false
        },
        {
          id: 'producerreport',
          title: 'Producer Report',
          type: 'item',
          url: '/Reports/ProducerReport',
          breadcrumbs: false
        },
        {
          id: 'consumerreport',
          title: 'Consumer Report',
          type: 'item',
          url: '/Reports/ConsumerReport',
          breadcrumbs: false
        },
        {
          id: 'auditreport',
          title: 'Audit Report',
          type: 'item',
          url: '/Reports/AuditReport',
          breadcrumbs: false
        },
      ]
    }
  ]
};

export default utilities;
