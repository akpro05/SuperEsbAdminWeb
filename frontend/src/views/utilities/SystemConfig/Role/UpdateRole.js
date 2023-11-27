import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import {
  Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, Table, TableBody, TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Checkbox
} from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_SUBSCRIBER_URL, SEARCH_SUBSCRIBER_URL } from '../../../../UrlPath.js';

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================|| Breadcrumb ||============================= //
function Breadcrumb({ label, onClick }) {
  return (
    <IconButton onClick={onClick}>
      {label}
    </IconButton>
  );
}

import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

const UpdateRole = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);
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
    <Typography color="text.primary">{t('updateRole')}</Typography>
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

  const initialValues = {
    name: '',
    input_status: '',

  };

  const validationSchema = Yup.object().shape({
    name: Yup.string()
      .matches(/^[A-Za-z\s]+$/, 'Role Name must contain only characters')
      .max(30, 'Role Name should be a maximum of 30 characters')
      .required('Role Name is required'),
    input_status: Yup.string().required(t('validationStatus')),
  });

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/Role/UpdateRole/${id}`);
        setUserData(response.data);
        setLoading(false);
      } catch (error) {
        setError('Failed to fetch data: ' + error.message);
        setLoading(false);
      }
    };

    fetchUserData();
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  // Parse Privilage into an object
  const Privilage = userData && userData.Privilage ? JSON.parse(userData.Privilage) : { Menus: [], Submenus: [] };

  const menuItems = Privilage.Menus;
const submenuItems = Privilage.Submenus;

//console.log("Menus:", menuItems);
//console.log("Submenus:", submenuItems);



  const handleSubmit = async (values, { setValues }) => {
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

      const response = await axios.post(`/Role/UpdateRole/${id}`, values);
      if (response.data.message) {
        setPopupMessage(response.data.message);
        setOpenPopup(true); // Open the popup dialog
      } else {
        openPopupWithMessage('An error occurred');
        setOpenPopup(true); // Open the popup dialog
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error updating user:', error);
      openPopupWithMessage('An error occurred');
    }
  };

  // -----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };

  const handleDialogClose = () => {
    setOpenPopup(false);
    navigate('/Role/SearchRole'); // Redirect to /SysUser/SearchsysUser
  };



  const initialData = [
    [
      {
        "menu_label": "System Configuration",
        "menu_value": "sysusermgmt",
        "submenu_array": [
          {
            "label": "Search System User",
            "value": "searchuser"
          }
        ]
      },
      {
        "menu_label": "Producer Management",
        "menu_value": "prodmgmt",
        "submenu_array": [
          {
            "label": "Search Producer",
            "value": "searchproducer"
          }
        ]
      },
      {
        "menu_label": "Consumer Management",
        "menu_value": "consumers",
        "submenu_array": [
          {
            "label": "Search Consumer",
            "value": "searchConsumers"
          }
        ]
      },
      {
        "menu_label": "System Configuration",
        "menu_value": "systconfig",
        "submenu_array": [
          {
            "label": "Set Role",
            "value": "setrole"
          },
          {
            "label": "Search Producer To Consumer",
            "value": "searchproducertoconsumer"
          }
        ]
      },
      {
        "menu_label": "Reports",
        "menu_value": "reports",
        "submenu_array": [
          {
            "label": "ESB Logs Report",
            "value": "esblogsreport"
          },
          {
            "label": "Audit Report",
            "value": "auditreport"
          }
        ]
      }
    ]
  ];


  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('updateRole')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: userData?.Name || '', // Populate the value from userData
              input_status: userData?.Status || '', // Populate the value from userData
              menuchecked: menuItems,
              submenuchecked: submenuItems,
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Field
                      name="name"
                      as={TextField}
                      label="Role Name"
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
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
                <br />
                <Grid container spacing={2}>
                  <Grid item xs={12}>
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
                                    // checked={Privilage.Menus.includes(row.menu_value)}
                                  />
                                  {row.menu_label}
                                </label>
                              </TableCell>
                              <TableCell>
                                {row.submenu_array.map((submenuItem, subIndex) => (
                                  <div key={submenuItem.value}>
                                    <label>
                                      <Field
                                        type="checkbox"
                                        name="submenuchecked"
                                        value={submenuItem.value}
                                        style={{marginRight:'10px'}}
                                        // checked={Privilage.Submenus.includes(submenuItem.value)}
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
      <Dialog open={openPopup} onClose={handleDialogClose}>
        <DialogTitle>{t('superEsbAdminAlert')}</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDialogClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default UpdateRole;