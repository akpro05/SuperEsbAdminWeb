import React, { useEffect } from 'react';
// material-ui
import { Typography } from '@mui/material';

// project imports
import NavGroup from './NavGroup';
import menuItem from '../../../../menu-items';

import axios from 'axios';
// ==============================|| SIDEBAR MENU LIST ||============================== //

const MenuList = () => {

  // useEffect(() => {
  //   const fetchData = async () => {
  //     try {
  //       const response = await axios.get('/Dashboard'); 
  //      // console.log('Data Received:', response.data);
  //     } catch (error) {
  //       console.error('Error:', error);
  //     }
  //   };
  
  //   fetchData();
  // }, []);

  

  const navItems = menuItem.items.map((item) => {
    switch (item.type) {
      case 'group':
        return <NavGroup key={item.id} item={item} />;
      default:
        return (
          <Typography key={item.id} variant="h6" color="error" align="center">
            Menu Items Error
          </Typography>
        );
    }
  });

  return <>{navItems}</>;
};

export default MenuList;
