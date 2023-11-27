import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, Table, TableBody, TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Checkbox} from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { VIEW_SUBSCRIBER_URL } from '../../../../UrlPath.js';

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
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';




const ViewRole = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);

  
  const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
      border: 0,
    },
  }));
  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link underline="hover" color="inherit" href="#">
        {t('System Configuration')}
      </Link>
      <Typography color="text.primary">{t('viewRole')}</Typography>
    </Breadcrumbs>
  ];
  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/Role/ViewRole/${id}`);
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
  const Name = userData && userData.Name;
  const Status = userData && userData.Status;

  
  // Parse Privilage into an object
const Privilage = userData && userData.Privilage ? JSON.parse(userData.Privilage) : { Menus: [], Submenus: [] };

  console.log(Privilage);


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


  return (
    <>
      <MainCard    title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('viewRole')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: '',
              input_status: '',
            }}
          >
            {({ errors, touched, handleSubmit }) => (
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
                      value={Name}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                    />
                    {touched.fullName && errors.fullName && (
                      <FormHelperText error>{errors.fullName}</FormHelperText>
                    )}
                  </Grid>
                  <Grid item xs={3}>
                    <Field
                      name="input_status"
                      as={TextField}
                      label={t('statusLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={Status}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                    />
                    {touched.address && errors.address && (
                      <FormHelperText error>{errors.address}</FormHelperText>
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
                              <Checkbox
                                value={row.menu_value}
                                label={row.menu_label}
                                checked={Privilage.Menus.includes(row.menu_value)}
                                readOnly={true}
                              />
                              {row.menu_label}
                            </TableCell>
                            <TableCell>
                              {row.submenu_array.map((submenuItem) => (
                                <div key={submenuItem.value}>
                                  <Checkbox
                                    value={submenuItem.value}
                                    label={submenuItem.label}
                                    checked={Privilage.Submenus.includes(submenuItem.value)}
                                    readOnly={true}
                                  />
                                  {submenuItem.label}
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
                  <Button variant="contained" color="primary" type="submit" style={{ backgroundColor: '#1478c8' }} onClick={() => navigate(`/Role/UpdateRole/${id}`)}>
                    {t('updateButton')}
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
        <DialogTitle>SuperEsb Admin Alert</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenPopup(false)} color="primary">
            close
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ViewRole;
