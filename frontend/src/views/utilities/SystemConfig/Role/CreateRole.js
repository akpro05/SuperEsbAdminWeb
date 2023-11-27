import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import {
  Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, Table, TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Checkbox
} from '@mui/material';
import FormControlLabel from '@mui/material/FormControlLabel';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';
import Typography from '@mui/material/Typography';
// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_USER_URL, SEARCH_USER_URL } from '../../../../UrlPath.js';
import Link from '@mui/material/Link';
// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================||Breadcrumb ||============================= //
function Breadcrumb({ label, onClick }) {
  return (
    <IconButton onClick={onClick}>
      {label}
    </IconButton>
  );
}
import Breadcrumbs from '@mui/material/Breadcrumbs';
const CreateRole = () => {
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const { t, i18n } = useTranslation();

  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link underline="hover" color="inherit" href="#">
        {t('System Configuration')}
      </Link>
      <Typography color="text.primary">{t('createRole')}</Typography>
    </Breadcrumbs>
  ];

  const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
      border: 0,
    },
  }));


  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const [formData, setFormData] = useState({
    name: '',
    input_status: '',
    menuchecked: [], // Initially an empty array
    submenuchecked: [], // Initially an empty array
    selectAll: false, // Add a state for the "Select All" checkbox
  });

  const handleSubmit = async (values, { resetForm }) => {
    try {
      console.log('Sending data to backend:', values);

      if (values.menuchecked.length <= 0) {
        openPopupWithMessage(t('emptyMenuMsg'));
        return; 
      }
  
      if (values.submenuchecked.length <= 0) {
        openPopupWithMessage(t('emptySubMenuMsg'));
        return; 
      }

      const response = await axios.post('/Role/CreateRole', values);
      if (response.data.success) {
        openPopupWithMessage('Role created successfully');
        resetForm({});
      } else {
        setPopupMessage(response.data.message);
        setOpenPopup(true); // Open the popup dialog
        resetForm({});
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error creating user:', error);
      openPopupWithMessage('An error occurred');
      resetForm({});
    }
  };

  // --------For Validation Message------//
  const validationSchema = Yup.object().shape({
    name: Yup.string()
      .matches(/^[A-Za-z\s]+$/, t('validationRolename1'))
      .max(30, t('validationRolename2'))
      .required(t('validationRolename')),
    input_status: Yup.string().required(t('validationStatus')),
  });

  //-----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };


  const initialData = [
    [
      {
        "menu_label": t('System User Management'),
        "menu_value": "sysusermgmt",
        "submenu_array": [
          {
            "label": t('Search System User'),
            "value": "searchuser"
          }
        ]
      },
      {
        "menu_label": t('Producer Management'),
        "menu_value": "prodmgmt",
        "submenu_array": [
          {
            "label": t('Search Producer'),
            "value": "searchproducer"
          }
        ]
      },
      {
        "menu_label": t('Consumer Management'),
        "menu_value": "consumers",
        "submenu_array": [
          {
            "label": t('Search Consumer'),
            "value": "searchConsumers"
          }
        ]
      },
      {
        "menu_label": t('System Configuration'),
        "menu_value": "systconfig",
        "submenu_array": [
          {
            "label": t('Set Role'),
            "value": "setrole"
          },
          {
            "label": t('Search Producer To Consumer'),
            "value": "searchproducertoconsumer"
          }
        ]
      },
      {
        "menu_label": t('Reports'),
        "menu_value": "reports",
        "submenu_array": [
          {
            "label": t('ESB Logs Report'),
            "value": "esblogsreport"
          },
          {
            "label": t('Audit Report'),
            "value": "auditreport"
          }
        ]
      }
    ]
  ];



  const handleMenuCheckboxChange = (event, setFieldValue) => {
    const menuValue = event.target.value;
    let updatedMenuchecked = [...formData.menuchecked];
    let updatedSubmenuChecked = [...formData.submenuchecked];

    // Check if the menuValue is already in menuchecked
    const isMenuChecked = updatedMenuchecked.includes(menuValue);

    if (isMenuChecked) {
      // If the menu is already checked, uncheck it and remove its submenus
      updatedMenuchecked = updatedMenuchecked.filter(value => value !== menuValue);
      const selectedMenu = initialData[0].find(menu => menu.menu_value === menuValue);

      if (selectedMenu) {
        const submenuValues = selectedMenu.submenu_array.map(submenu => submenu.value);
        updatedSubmenuChecked = updatedSubmenuChecked.filter(value => !submenuValues.includes(value));
      }
    } else {
      // If the menu is not checked, check it and add its submenus
      updatedMenuchecked = [...updatedMenuchecked, menuValue];
      const selectedMenu = initialData[0].find(menu => menu.menu_value === menuValue);

      if (selectedMenu) {
        const submenuValues = selectedMenu.submenu_array.map(submenu => submenu.value);
        updatedSubmenuChecked = [...updatedSubmenuChecked, ...submenuValues];
      }
    }

    // Update the state with the updated menuchecked and submenuchecked arrays
    setFormData({ ...formData, menuchecked: updatedMenuchecked, submenuchecked: updatedSubmenuChecked });

    // Update the Formik field value
    setFieldValue('menuchecked', updatedMenuchecked);
    setFieldValue('submenuchecked', updatedSubmenuChecked);
  };

  const handleSelectAll = (setFieldValue) => {
    let updatedMenuchecked = [...formData.menuchecked];
    let updatedSubmenuChecked = [...formData.submenuchecked];

    // Check if any menu is checked to determine the desired action
    const isAllChecked = updatedMenuchecked.length === initialData[0].length;

    if (isAllChecked) {
      // If all are checked, uncheck them all
      updatedMenuchecked = [];
      updatedSubmenuChecked = [];
    } else {
      // If not all are checked, check them all
      initialData[0].forEach((menu) => {
        if (!updatedMenuchecked.includes(menu.menu_value)) {
          updatedMenuchecked.push(menu.menu_value);
          updatedSubmenuChecked.push(...menu.submenu_array.map(submenu => submenu.value));
        }
      });
    }

    // Update the state and Formik field values
    setFormData({ ...formData, menuchecked: updatedMenuchecked, submenuchecked: updatedSubmenuChecked });
    setFieldValue('menuchecked', updatedMenuchecked);
    setFieldValue('submenuchecked', updatedSubmenuChecked);
    setFieldValue('selectAll', !isAllChecked);
  };

// Inside your component
const handleSubmenuCheckboxChange = (event, setFieldValue) => {
  const submenuValue = event.target.value;
  const menuValue = getMenuValueBySubmenuValue(submenuValue); // Function to find the corresponding menu

  let updatedMenuchecked = [...formData.menuchecked];
  let updatedSubmenuChecked = [...formData.submenuchecked];

  // Check if the submenu is checked
  if (event.target.checked) {
    updatedSubmenuChecked.push(submenuValue);
  } else {
    updatedSubmenuChecked = updatedSubmenuChecked.filter((value) => value !== submenuValue);
  }

  // Check if the menu is checked
  const isMenuChecked = updatedMenuchecked.includes(menuValue);

  // If the submenu is checked and the menu is not, check the menu
  if (event.target.checked && !isMenuChecked) {
    updatedMenuchecked.push(menuValue);
  }

  // Update the state and Formik field values
  setFormData({ ...formData, menuchecked: updatedMenuchecked, submenuchecked: updatedSubmenuChecked });
  setFieldValue('menuchecked', updatedMenuchecked);
  setFieldValue('submenuchecked', updatedSubmenuChecked);
};

// Function to find the corresponding menu value based on submenu value
const getMenuValueBySubmenuValue = (submenuValue) => {
  for (const menu of initialData[0]) {
    for (const submenu of menu.submenu_array) {
      if (submenu.value === submenuValue) {
        return menu.menu_value;
      }
    }
  }
  return ''; // Return an empty string if not found
};

  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('createRole')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: '',
              input_status: '',
              menuchecked: [], // Initially an empty array
              submenuchecked: [], // Initially an empty array
              selectAll: false, // Initialize "Select All" checkbox to false
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ values, setFieldValue, errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Field
                      name="name"
                      as={TextField}
                      label={t('roleNameLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    />
                    {touched.name && errors.name && (
                      <FormHelperText error>{errors.name}</FormHelperText>
                    )}
                  </Grid>




                  <Grid item xs={3}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel id="status-label">{t('statusLabel')}</InputLabel>
                      <Field
                        name="input_status"
                        as={Select}
                        labelId="status-label"
                        label="Status"
                      >
                        <MenuItem value="">{t('selectStatus')}</MenuItem>
                        {hardcodedStatuses.map((status) => (
                          <MenuItem key={status.id} value={status.id}>
                            {status.name}
                          </MenuItem>
                        ))}
                      </Field>

                    </FormControl>

                    {touched.input_status && errors.input_status && (
                      <FormHelperText error>{errors.input_status}</FormHelperText>
                    )}
                  </Grid>
                </Grid>
                <br/>
                <Grid container spacing={2}>
				
                  <Grid item xs={12}>
                    <div>
                      <label style={{ fontWeight: 'bold' }}>
                        <input
                          type="checkbox"
                          name="selectAll"
                          style={{marginRight:'10px'}}
                          checked={values.selectAll}
                          onChange={() => handleSelectAll(setFieldValue)}
                        />
                         {t('SelectAll')}
                      </label>
                    </div>
                    <TableContainer component={Paper} sx={{ border: '1px solid #000' }}>
                      <Table>
                        <TableHead>
                          <TableRow>
                            <TableCell sx={{ fontWeight: 'bold', background: 'rgb(12, 53, 106)', color: 'white' }}>{t('srno')}</TableCell>
                            <TableCell sx={{ fontWeight: 'bold', background: 'rgb(12, 53, 106)', color: 'white' }}>{t('menu')}</TableCell>
                            {/* <TableCell sx={{ fontWeight: 'bold', background: '#814141', color: 'white' }}>Action</TableCell> */}
                            <TableCell sx={{ fontWeight: 'bold', background: 'rgb(12, 53, 106)', color: 'white' }}>{t('subMenu')}</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {initialData[0].map((row, index) => (
                            <StyledTableRow key={row.menu_value}>
                              <TableCell>{index + 1}</TableCell>
                              <TableCell>
                                <label>
                                  <Field
                                    type="checkbox"
                                    name="menuchecked"
                                    value={row.menu_value}
                                    style={{marginRight:'10px'}}
                                    checked={values.menuchecked.includes(row.menu_value)}
                                    onChange={(event) => handleMenuCheckboxChange(event, setFieldValue)}
                                  />
                                  {row.menu_label}
                                </label>
                              </TableCell>
                              {/* <TableCell>
                                <Button variant="outlined" color="primary">
                                  Show
                                </Button>
                              </TableCell> */}
                              <TableCell>
                                {row.submenu_array.map((submenuItem) => (
                                  <div key={submenuItem.value}>
                                    <label>
                                      <Field
                                        type="checkbox"
                                        name="submenuchecked"
                                        value={submenuItem.value}
                                        style={{marginRight:'10px'}}
                                        onChange={(event) => handleSubmenuCheckboxChange(event, setFieldValue)}
                                      />
                                      {submenuItem.label}
                                    </label>
                                  </div>
                                ))}

                              </TableCell>
                            </StyledTableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Grid>

                </Grid>

                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/Role/SearchRole')}>
                    {t('back')}
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Card>
      </MainCard>
      <Dialog open={openPopup} onClose={() => setOpenPopup(false)}>
        <DialogTitle>{t('superEsbAdminAlert')}</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenPopup(false)} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default CreateRole;