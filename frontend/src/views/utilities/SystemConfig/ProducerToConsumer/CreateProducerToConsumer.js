import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions, FormLabel ,TableContainer, Table, TableHead, TableRow, TableCell, TableBody, Paper} from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, FormGroup, FormControlLabel, Checkbox } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_USER_URL, CREATE_PROD2CONS_URL, GET_CONSUMER_SERVICES_URL } from '../../../../UrlPath.js';
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

const CreateProducerToConsumer = () => {
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [data, setData] = useState(null);
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
      <Typography color="text.primary">{t('createProducerToConsumer')}</Typography>
    </Breadcrumbs>
  ];

  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    input_status: '',
  });


  const structDefinition = `type ProducerServices struct {
    SubsrcibedServices []struct {
      ServiceName string \`json:"service_name"\`
      ServiceURL  string \`json:"service_url"\`
    } \`json:"subsrcibed_services"\`
  }`;

  const SamplestructDefinition = `{
    "subsrcibed_services":[
       {
          "service_name":"Cnps PGS Payment Page",
          "service_url":"/PG/txnsummery"
       }
    ]
 }`;


  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(CREATE_PROD2CONS_URL);
        // console.log(response.data);
        setData(response.data);
      } catch (error) {
        setError('Failed to fetch data: ' + error.message);
      }
    };

    fetchData();
  }, []);


  const handleCheckboxChange = (e) => {
    const serviceUrl = e.target.value;
    const serviceName = e.target.dataset.serviceName;
  
    // Now you have the corresponding service name based on the checkbox.
    console.log('Service Name:', serviceName);
    // You can use the `serviceName` as needed in your `handleSubmit` function.
  };


  const handleSubmit = async (values, { resetForm }) => {
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

      // Combine the checkbox values with the form values
      const dataToSend = {
        ...values,
        subscribed_services: subscribedServices, // 
      };

      //console.log('Sending data to backend:', dataToSend);

      const response = await axios.post(
        '/ProducerToConsumer/CreateProducerToConsumer',
        dataToSend
      );

      // Clear the checked state of all checkboxes
      selectedCheckboxes.forEach((checkbox) => {
        checkbox.checked = false;
      });

      if (response.data.success) {
        openPopupWithMessage('Producer To Consumer mapping created successfully');
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


  // styles
  const CustomCheckbox = styled(Checkbox)(({ theme }) => ({
    '&.Mui-checked': {
      color: 'red', // Change the color to red when checked
    },
  }));

  // --------For Validation Message------//
  const validationSchema = Yup.object().shape({
    input_producer: Yup.string().required(t('validationProducer')),
    input_consumer: Yup.string().required(t('validationConsumer')),
    // producerservices: Yup.string().required('Producer Services is required'),
    input_status: Yup.string().required(t('validationStatus')),
  });

  //-----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };



  const handleConsumerChange = async (event, setFieldValue) => {
    const selectedConsumerValue = event.target.value;

    // Set the selected consumer value in the Formik values
    setFieldValue('input_consumer', selectedConsumerValue);

    try {
      // Create the request body
      const requestBody = {
        consumer_id: selectedConsumerValue,
      };

      // Send a POST request with the request body
      const response = await axios.post(GET_CONSUMER_SERVICES_URL, requestBody);
      const servicesData = response.data.services_list;
      const ConsumerDomainAddress = response.data.ConsumerDomainAddress;


      const newCheckboxes = servicesData.map((service, index) => (
        <TableRow key={index}>
          <TableCell>
            <label>
              <input
                type="checkbox"
                name={`service_${index}`}
                value={service.service_url}
                data-service-name={service.service_name} // Include service_name as a data attribute
                onChange={handleCheckboxChange}
              />
            </label>
          </TableCell>
          <TableCell>{service.service_name}</TableCell>
          <TableCell>{ConsumerDomainAddress}{service.service_url}</TableCell>
        </TableRow>
      ));
  
      setCheckboxes(
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
                  Action
                </TableCell>
                <TableCell
                  sx={{
                    fontWeight: 'bold',
                    background: 'rgb(12, 53, 106)',
                    color: 'white',
                  }}
                >
                  Service Name
                </TableCell>
                <TableCell
                  sx={{
                    fontWeight: 'bold',
                    background: 'rgb(12, 53, 106)',
                    color: 'white',
                  }}
                >
                  Service URL
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {newCheckboxes}
            </TableBody>
          </Table>
        </TableContainer>
      );

    } catch (error) {
      console.error('Error fetching services:', error);
      // Handle error here
    }
  };




  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('createProducerToConsumer')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              input_producer: '',
              input_status: '',
              input_consumer: '',
              // producerservices: '',
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ values, errors, touched, handleSubmit, setFieldValue }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>


                  {data && (
                    <>
                      <Grid item xs={4}>
                        <FormControl fullWidth variant="outlined" margin="normal">
                          <InputLabel id="producer-label">{t('producer')}</InputLabel>
                          <Field
                            name="input_producer"
                            as={Select}
                            labelId="producer-label"
                            label="Producer"
                          >
                            {data.ProducerData.Fields1
                              ? data.ProducerData.Fields1.map((producer) => (
                                <MenuItem key={producer.Id} value={producer.Id}>
                                  {producer.Name}
                                </MenuItem>
                              ))
                              : <MenuItem disabled>No data available</MenuItem>
                            }
                          </Field>
                        </FormControl>
                        {touched.input_producer && errors.input_producer && (
                          <FormHelperText error>{errors.input_producer}</FormHelperText>
                        )}
                      </Grid>
                      <Grid item xs={4}>
                        <FormControl fullWidth variant="outlined" margin="normal">
                          <InputLabel id="consumer-label">{t('consumer')}</InputLabel>
                          <Field
                            name="input_consumer"
                            as={Select}
                            labelId="consumer-label"
                            label="Consumer"
                            onChange={(event) => handleConsumerChange(event, setFieldValue)}
                          >
                            {data.ConsumerData.Fields1
                              ? data.ConsumerData.Fields1.map((consumer) => (
                                <MenuItem key={consumer.Id} value={consumer.Id}>
                                  {consumer.Name}
                                </MenuItem>
                              ))
                              : <MenuItem disabled>No data available</MenuItem>
                            }
                          </Field>
                        </FormControl>
                        {touched.input_consumer && errors.input_consumer && (
                          <FormHelperText error>{errors.input_consumer}</FormHelperText>
                        )}
                      </Grid>
                    </>
                  )}



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
                    </FormControl>
                    {touched.input_status && errors.input_status && (
                      <FormHelperText error>{errors.input_status}</FormHelperText>
                    )}
                  </Grid>

                  {/* <Grid item xs={4}>
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
                {/* <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <br />
                    <span style={{ fontWeight: 'bold' }}>Producer Services Input  should be a json struct  in the form of following type</span>
                    <br />
                    <pre>
                      <code>{structDefinition}</code>
                    </pre>

                  </Grid>
                  <Grid item xs={6}>
                    <br />
                    <span style={{ fontWeight: 'bold' }}>Sample json for reference</span>
                    <br />
                    <pre>
                      <code>{SamplestructDefinition}</code>
                    </pre>

                  </Grid>
                </Grid> */}
                <Grid container spacing={2} style={{justifyContent:'center',display:'flex',alignItems:'center',marginTop:'10px'}}>
                  <FormControl sx={{ m: 3 }} component="fieldset" variant="standard">
                    <span style={{ fontWeight: 'bold' }}>{t('allAssociated')}</span>
                    <br />
                    {checkboxes} {/* Render the dynamic checkboxes here */}
                  </FormControl>

                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  {/* <Button
                    variant="contained"
                    color="primary"
                    onClick={() => fetchServicesData(values)}
                  >
                    Get Services
                  </Button> */}
                  <Box sx={{ mx: 2 }} />
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
      <Dialog open={openPopup} onClose={() => setOpenPopup(false)}>
        <DialogTitle>{t('superEsbAdminAlert')}</DialogTitle>
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

export default CreateProducerToConsumer;