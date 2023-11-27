import React from 'react';
import dashboard from './dashboard';
import pages from './pages';
import utilities from './utilities';
import other from './other';

// ==============================|| MENU ITEMS ||============================== //

const menu = sessionStorage.getItem('menu');
const menuItemsToInclude = JSON.parse(menu) || []; // Use an empty array if menu is null or undefined
const submenu = sessionStorage.getItem('submenu');
const submenuItemsToInclude = JSON.parse(submenu) || []; // Use an empty array if submenu is null or undefined

// Filter the initial_utilities object to include only the specified IDs
const filtered_utilities = {
  ...utilities,
  children: utilities.children.filter(menuItem => {
    return menuItemsToInclude.includes(menuItem.id);
  })
};

//console.log("Original filtered_utilities:", filtered_utilities);

function filterUtilities(data, submenuItems) {
  for (const child of data.children) {
    child.children = child.children.filter((nestedChild) =>
      submenuItems.includes(nestedChild.id)
    );
  }
  return data;
}

const filtered_utilities_filtered = filterUtilities(filtered_utilities, submenuItemsToInclude);
//console.log("Filtered filtered_utilities:", filtered_utilities_filtered);

const menuItems = {
  items: [dashboard,  filtered_utilities_filtered]
};

export default menuItems;
