import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions, TableContainer, Table, TableHead, TableRow, TableCell, TableBody, Paper } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, FormGroup, FormControlLabel, Checkbox } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_USER_URL, SEARCH_USER_URL } from '../../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

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

const UpdateProducerToConsumer = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);
  const [checkboxes, setCheckboxes] = useState([]); // Maintain the state for checkboxes
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
      <Typography color="text.primary">{t('updateProducerToConsumer')}</Typography>
    </Breadcrumbs>
  ];


  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const initialValues = {
    subcname: '',
    prodname: '',
    subscriberservices: '',
    input_status: '',
  };

  const validationSchema = Yup.object().shape({
    // producerservices: Yup.string().required('Producer Services is required'),
    input_status: Yup.string().required(t('validationStatus')),
  });

 

  const handleCheckboxChange = (e) => {
    const serviceUrl = e.target.value;
    const serviceName = e.target.dataset.serviceName;
  
    // Now you have the corresponding service name based on the checkbox.
    console.log('Service Name:', serviceName);
    // You can use the `serviceName` as needed in your `handleSubmit` function.
  };


  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/ProducerToConsumer/UpdateProducerToConsumer/${id}`);
        setUserData(response.data);
        console.log(response.data);

        const servicesData = response.data.services_list;
        const ConsumerDomainAddress = response.data.ConsumerDomainAddress;
        const producerServicesData = response.data.ProducerServices;
        let producerServices = null;
        const serviceUrlToNameMapping = {};

        servicesData.forEach((service) => {
          serviceUrlToNameMapping[service.service_url] = service.service_name;
        });


        try {
          producerServices = producerServicesData ? JSON.parse(producerServicesData) : null;
        } catch (error) {
          console.error("Error parsing JSON data:", error);
        }

        const subscribedServiceUrls = producerServices.subsrcibed_services.map(service => service.service_url); // Extract subscribed service URLs

        const newCheckboxes = servicesData.map((service, index) => {
          const isChecked = subscribedServiceUrls.includes(service.service_url); // Check if the service is subscribed
          return (
            <TableRow key={index}>
              <TableCell>
                <label>
                  <input
                    type="checkbox"
                    name={`service_${index}`}
                    value={service.service_url}
                    defaultChecked={isChecked} 
                    data-service-name={service.service_name} // Include service_name as a data attribute
                    onChange={handleCheckboxChange}
                  />
                </label>
              </TableCell>
              <TableCell>{service.service_name}</TableCell>
              <TableCell>{ConsumerDomainAddress}{service.service_url}</TableCell>
            </TableRow>
          );
        });

        setCheckboxes(newCheckboxes);
        setLoading(false);
      } catch (error) {
        setError('Failed to fetch data: ' + error.message);
        setLoading(false);
      }
    };

    fetchUserData();
  }, [id]);

  // Add a console.log to check the serviceList



  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  const handleSubmit = async (values, { setValues }) => {
    try {

      // Get all the checkboxes that are checked
      const selectedCheckboxes = document.querySelectorAll('input[type="checkbox"]:checked');

      // Create the JSON struct array from the selected checkboxes
      const subscribedServices = [...selectedCheckboxes].map((checkbox) => {
        const label = checkbox.parentElement.textContent; // Get the label text
        const serviceName = checkbox.dataset.serviceName;
        const serviceUrl = checkbox.value; // Get the checkbox value

        return {
          service_name: serviceName,
          service_url: serviceUrl,
        };
      });

      if (subscribedServices.length <= 0) {
        openPopupWithMessage(t('emptyServiceMsg'));
        return;
      }

      // Include the selected services in the request data
      const requestData = {
        ...values,
        subscribed_services: subscribedServices, // 
      };

      console.log('Sending data to backend:', requestData);

      const response = await axios.post(`/ProducerToConsumer/UpdateProducerToConsumer/${id}`, requestData);
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
    navigate('/ProducerToConsumer/SearchProducerToConsumer'); // Redirect to /SysUser/SearchsysUser
  };





  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('updateProducerToConsumer')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              prodname: userData?.ProducerName || '', // Populate the value from userData
              consname: userData?.ConsumerName || '', // Populate the value from userData
              // producerservices: userData?.ProducerServices || '', // Populate the value from userData
              input_status: userData?.Status || '', // Populate the value from userData
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
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
                      {touched.input_status && errors.input_status && (
                        <FormHelperText error>{errors.input_status}</FormHelperText>
                      )}
                    </FormControl>
                  </Grid>

                  {/* <Grid item xs={6}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel htmlFor="amf_metadata"></InputLabel>
                      <Field
                        as={TextField}
                        id="producerservices"
                        name="producerservices"
                        label="Producer Services"
                        variant="outlined"
                        multiline
                      />

                    </FormControl>

                    {touched.producerservices && errors.producerservices && (
                      <FormHelperText error>{errors.producerservices}</FormHelperText>
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
                              {t('action')}
                            </TableCell>
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
                          {checkboxes}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </FormControl>

                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
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

export default UpdateProducerToConsumer;