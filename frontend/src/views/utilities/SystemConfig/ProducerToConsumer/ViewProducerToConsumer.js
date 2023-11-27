import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, FormGroup, FormControlLabel, Checkbox } from '@mui/material';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { VIEW_USER_URL } from '../../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
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

const ViewProducerToConsumer = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);
  const [service_list, setService_list] = useState([]);
  const [consumerDomain, setDomainValue] = useState(null);
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
      <Typography color="text.primary">{t('viewProducerToConsumer')}</Typography>
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
        const response = await axios.get(`/ProducerToConsumer/ViewProducerToConsumer/${id}`);
        setUserData(response.data);
        setService_list(response.data.service_list);
        setLoading(false);
        setDomainValue(response.data.ConsumerDomainAddress);
      } catch (error) {
        setError('Failed to fetch data: ' + error.message);
        setLoading(false);
      }
    };
    fetchUserData();
  }, [id]);

  // Add a console.log to check the serviceList
  console.log("service_list", service_list);





  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  const ProdName = userData && userData.ProducerName;
  const ConsName = userData && userData.ConsumerName;
  const ProdServices = userData && userData.ProducerServices;
  const Status = userData && userData.Status;

  const producerServicesData = userData && userData.ProducerServices;
  let producerServices = null;

  try {
    producerServices = producerServicesData ? JSON.parse(producerServicesData) : null;
  } catch (error) {
    console.error("Error parsing JSON data:", error);
  }

  console.log(producerServices);
  console.log("fhgdsfjdsf");


  //-----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };
  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('viewProducerToConsumer')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: '',
              email: '',
              clientcode: '',
              validateparameter: '',
              endpointsconfig: '',
              input_status: '',
              isChecked: false,
            }}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={4}>
                    <Field
                      name="prodname"
                      as={TextField}
                      label={t('producer')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={ProdName}
                      InputProps={{
					   readOnly: true,
					 }}
					   inputProps={{
					      style: {
					      backgroundColor: 'rgb(208 208 214)',
					      },
					 }}
                    />

                  </Grid>
                  <Grid item xs={4}>
                    <Field
                      name="consname"
                      as={TextField}
                      label={t('consumer')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={ConsName}
                      InputProps={{
					   readOnly: true,
					 }}
					   inputProps={{
					      style: {
					      backgroundColor: 'rgb(208 208 214)',
					      },
					 }}
                    />

                  </Grid>

                  <Grid item xs={4}>
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

                  {/* <Grid item xs={6}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel htmlFor="amf_metadata"></InputLabel>
                      <Field
                        as={TextField}
                        id="prodservices"
                        name="prodservices"
                        label="Producer Services Configuration"
                        variant="outlined"
                        multiline
                        value={ProdServices}
                        InputProps={{
                          readOnly: true,
                          style: { backgroundColor: '#f0f0f0' },
                        }}
                      />

                    </FormControl>

                    {touched.endpointsconfig && errors.endpointsconfig && (
                      <FormHelperText error>{errors.endpointsconfig}</FormHelperText>
                    )}
                  </Grid> */}

                </Grid>

                <Grid container spacing={2} style={{justifyContent:'center',display:'flex',alignItems:'center',marginTop:'10px'}}>
                  <FormControl sx={{ m: 3 }} component="fieldset" variant="standard">
                    <span style={{ fontWeight: 'bold' }}> {t('listSubscribed')}</span>
                    <br />
                    <TableContainer component={Paper} sx={{ border: '1px solid #000' }}>
                      <Table>
                        <TableHead>
                          <TableRow>
                            <TableCell
                              sx={{
                                fontWeight: 'bold',
                                background: 'rgb(12, 53, 106)',
                                color: 'white',
                              }}
                            >
                              {t('serviceName')}
                            </TableCell>
                            <TableCell
                              sx={{
                                fontWeight: 'bold',
                                background: 'rgb(12, 53, 106)',
                                color: 'white',
                              }}
                            >
                              {t('serviceURL')}
                            </TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {producerServices && producerServices.subsrcibed_services ? (
                            producerServices.subsrcibed_services.map((service, index) => (
                              <TableRow key={index}>
                                <TableCell>{service.service_name}</TableCell>
                                <TableCell>{consumerDomain}{service.service_url}</TableCell>
                              </TableRow>
                            ))
                          ) : (
                            <TableRow>
                              <TableCell colSpan={2}>{t('noSubscribed')}</TableCell>
                            </TableRow>
                          )}
                        </TableBody>
                      </Table>
                    </TableContainer>

                  </FormControl>
                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" type="submit" style={{ backgroundColor: '#1478c8' }} onClick={() => navigate(`/ProducerToConsumer/UpdateProducerToConsumer/${id}`)}>
                    {t('updateButton')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/ProducerToConsumer/SearchProducerToConsumer')}>
                    {t('back')}
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Card>
      </MainCard>
      <Dialog open={openPopup} onClose={() => setOpenPopup(false)}>
        <DialogTitle>SuperESB Admin Alert</DialogTitle>
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

export default ViewProducerToConsumer;